package events

// Package events provides CLI filter constants, helper functions, and
// reusable utilities for the eventing aggregator plugin.
// Enum types and event struct types live in common/types.go.
// Plugin execution functions live in controlplane/v1.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"gopkg.in/yaml.v3"
)

// ── Default timeouts ──

const (
	DefaultTimeout      = 60 * time.Second
	DefaultPollInterval = 3 * time.Second
)

// ── FilterOption builds CLI flag pairs for plugin queries ──

type FilterOption func() []string

// WithCategory adds --category flag to filter events by category (e.g. "pool", "volume", "nexus").
func WithCategory(filter string) FilterOption {
	return func() []string { return []string{FlagCategory, filter} }
}

// WithAction adds --action flag to filter events by action (e.g. "create", "delete", "state-change").
func WithAction(filter string) FilterOption {
	return func() []string { return []string{FlagAction, filter} }
}

// WithComponent adds --component flag to filter events by source component (e.g. "core-agent", "io-engine").
func WithComponent(filter string) FilterOption {
	return func() []string { return []string{FlagComponent, filter} }
}

// WithNode adds --node flag to filter events by the node name that generated them.
func WithNode(node string) FilterOption {
	return func() []string { return []string{FlagNode, node} }
}

// WithTarget adds --target flag to filter events by target UUID (volume, nexus, replica, etc.).
func WithTarget(target string) FilterOption {
	return func() []string { return []string{FlagTarget, target} }
}

// WithPool adds --pool flag to filter events by pool name or UUID.
func WithPool(pool string) FilterOption {
	return func() []string { return []string{FlagPool, pool} }
}

// WithVolume adds --volume flag to filter events by volume UUID.
func WithVolume(volume string) FilterOption {
	return func() []string { return []string{FlagVolume, volume} }
}

// WithReplica adds --replica flag to filter events by replica UUID.
func WithReplica(replica string) FilterOption {
	return func() []string { return []string{FlagReplica, replica} }
}

// WithRebuildStatusFilter adds --rebuild-status flag to filter rebuild events by status (e.g. "started", "completed").
func WithRebuildStatusFilter(status string) FilterOption {
	return func() []string { return []string{FlagRebuildStatus, status} }
}

// WithState adds --state flag to filter events by state value.
func WithState(state string) FilterOption {
	return func() []string { return []string{FlagState, state} }
}

// WithFilter adds --filter flag for free-form text filtering on event records.
func WithFilter(filter string) FilterOption {
	return func() []string { return []string{FlagFilter, filter} }
}

// WithSince adds --since flag to filter events from a relative duration (e.g. "5m", "1h").
func WithSince(duration string) FilterOption {
	return func() []string { return []string{FlagSince, duration} }
}

// WithLimit adds --limit flag to cap the number of returned events.
func WithLimit(limit int) FilterOption {
	return func() []string { return []string{FlagLimit, strconv.Itoa(limit)} }
}

// WithLokiEndpoint adds --loki-endpoint flag to force a specific Loki URL (hard error if unreachable).
func WithLokiEndpoint(endpoint string) FilterOption {
	return func() []string { return []string{FlagLokiEndpoint, endpoint} }
}

// WithFromFile adds --from-file flag to read events from a local file instead of querying Loki/aggregator.
func WithFromFile(filePath string) FilterOption {
	return func() []string { return []string{FlagFromFile, filePath} }
}

// WithNatsEndpoint adds --nats-endpoint flag with auto-discovery (no value) to read events directly from NATS JetStream.
func WithNatsEndpoint() FilterOption {
	return func() []string { return []string{FlagNatsEndpoint} }
}

// BuildFlags collects all FilterOption functions into a flat slice of CLI flag pairs.
func BuildFlags(opts ...FilterOption) []string {
	var flags []string
	for _, opt := range opts {
		flags = append(flags, opt()...)
	}
	return flags
}

// ── Query helpers ──

// GetEvents queries the plugin for events matching the given filters and returns them newest-first.
func GetEvents(opts ...FilterOption) ([]common.EventRecord, error) {
	return controlplane.GetEvents(BuildFlags(opts...)...)
}

// GetEventsWithStderr queries the plugin and returns events plus the raw stderr output (warnings).
func GetEventsWithStderr(opts ...FilterOption) ([]common.EventRecord, string, error) {
	return controlplane.GetEventsWithStderr(BuildFlags(opts...)...)
}

// GetRawEventsOutput runs the plugin with the specified output format and returns raw bytes.
func GetRawEventsOutput(format string, opts ...FilterOption) ([]byte, error) {
	return controlplane.GetRawEventsOutput(format, BuildFlags(opts...)...)
}

// EventsCount returns the number of events matching the given filters.
func EventsCount(opts ...FilterOption) int {
	return controlplane.CountEvents(BuildFlags(opts...)...)
}

// IsEventFound returns true if at least one event matches the given filters.
func IsEventFound(opts ...FilterOption) bool {
	return EventsCount(opts...) > 0
}

// ── Wait helpers ──

// WaitForEvents polls until at least 1 event matches the filters, using DefaultTimeout.
func WaitForEvents(opts ...FilterOption) ([]common.EventRecord, error) {
	return WaitForEventCountWithTimeout(1, DefaultTimeout, opts...)
}

// WaitForEventCount polls until at least count events match the filters, using DefaultTimeout.
func WaitForEventCount(count int, opts ...FilterOption) ([]common.EventRecord, error) {
	return WaitForEventCountWithTimeout(count, DefaultTimeout, opts...)
}

// WaitForEventCountWithTimeout polls until at least count events match or the timeout expires.
func WaitForEventCountWithTimeout(count int, timeout time.Duration, opts ...FilterOption) ([]common.EventRecord, error) {
	var records []common.EventRecord
	deadline := time.Now().Add(timeout)
	flags := BuildFlags(opts...)
	for {
		var err error
		records, err = controlplane.GetEvents(flags...)
		if err == nil && len(records) >= count {
			return records, nil
		}
		if time.Now().After(deadline) {
			break
		}
		if err != nil {
			logf.Log.Info("WaitForEventCountWithTimeout", "poll error", err)
		}
		time.Sleep(DefaultPollInterval)
	}
	return records, fmt.Errorf("timed out waiting for %d events (got %d) after %v",
		count, len(records), timeout)
}

// EventCount returns a closure suitable for Gomega Eventually that polls the matching event count.
func EventCount(opts ...FilterOption) func() int {
	flags := BuildFlags(opts...)
	return func() int {
		return controlplane.CountEvents(flags...)
	}
}

// ── Client-side filters (post-retrieval) ──

// FilterEventByCategory filters a pre-fetched event slice to only records matching the given category.
func FilterEventByCategory(records []common.EventRecord, category common.EventCategory) []common.EventRecord {
	var result []common.EventRecord
	for _, r := range records {
		if r.Category == category {
			result = append(result, r)
		}
	}
	return result
}

// FilterEventByAction filters a pre-fetched event slice to only records matching the given action.
func FilterEventByAction(records []common.EventRecord, action common.EventAction) []common.EventRecord {
	var result []common.EventRecord
	for _, r := range records {
		if r.Action == action {
			result = append(result, r)
		}
	}
	return result
}

// FilterEventByTarget filters a pre-fetched event slice to only records matching the given target UUID.
func FilterEventByTarget(records []common.EventRecord, target string) []common.EventRecord {
	var result []common.EventRecord
	for _, r := range records {
		if r.Target == target {
			result = append(result, r)
		}
	}
	return result
}

// FilterEventByComponent filters a pre-fetched event slice to only records from the given source component.
func FilterEventByComponent(records []common.EventRecord, component common.EventComponent) []common.EventRecord {
	var result []common.EventRecord
	for _, r := range records {
		if r.Metadata.Source.Component == component {
			result = append(result, r)
		}
	}
	return result
}

// GetRebuildEventsDetails extracts RebuildDetails from an event record, returns nil if absent.
func GetRebuildEventsDetails(record common.EventRecord) *common.RebuildDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.RebuildDetails
}

// GetSwitchOverEventDetails extracts SwitchOverDetails from an event record, returns nil if absent.
func GetSwitchOverEventDetails(record common.EventRecord) *common.SwitchOverDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.SwitchOverDetails
}

// GetStateChangeEventDetails extracts StateChangeDetails from an event record, returns nil if absent.
func GetStateChangeEventDetails(record common.EventRecord) *common.StateChangeDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.StateChangeDetails
}

// GetSnapshotEventDetails extracts SnapshotDetails from an event record, returns nil if absent.
func GetSnapshotEventDetails(record common.EventRecord) *common.SnapshotDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.SnapshotDetails
}

// GetNvmePathEventDetails extracts NvmePathDetails from an event record, returns nil if absent.
func GetNvmePathEventDetails(record common.EventRecord) *common.NvmePathDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.NvmePathDetails
}

// GetReplicaEventDetails extracts ReplicaDetails from an event record, returns nil if absent.
func GetReplicaEventDetails(record common.EventRecord) *common.ReplicaDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.ReplicaDetails
}

// GetReactorEventDetails extracts ReactorDetails from an event record, returns nil if absent.
func GetReactorEventDetails(record common.EventRecord) *common.ReactorDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.ReactorDetails
}

// GetErrorEventDetails extracts ErrorDetails from an event record, returns nil if absent.
func GetErrorEventDetails(record common.EventRecord) *common.ErrorDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.ErrorDetails
}

// GetNexusChildEventDetails extracts NexusChildDetails from an event record, returns nil if absent.
func GetNexusChildEventDetails(record common.EventRecord) *common.NexusChildDetails {
	if record.Metadata.Source.EventDetails == nil {
		return nil
	}
	return record.Metadata.Source.EventDetails.NexusChildDetails
}

// ── Event metadata verification helpers ──

// CheckRebuildEvent validates rebuild event metadata: source/dest URIs, error presence, and rebuild status.
func CheckRebuildEvent(
	record *common.EventRecord,
	sourceUri string,
	destUri string,
	hasErr bool,
	rebuildStatus common.RebuildStatus) error {

	details := GetRebuildEventsDetails(*record)
	if details == nil {
		return fmt.Errorf("event is missing rebuild details")
	}

	if hasErr {
		if details.Error == "" {
			return fmt.Errorf("event is missing error message")
		}
	} else {
		if details.Error != "" {
			return fmt.Errorf("event has error message: %s", details.Error)
		}
	}

	if details.SourceReplica != sourceUri {
		return fmt.Errorf(
			"event has invalid source uri, expected %s, got: %s",
			sourceUri, details.SourceReplica,
		)
	}

	if details.DestinationReplica != destUri {
		return fmt.Errorf(
			"event has invalid destination uri, expected %s, got: %s",
			destUri, details.DestinationReplica,
		)
	}

	if details.RebuildStatus != rebuildStatus {
		return fmt.Errorf(
			"event has invalid rebuild status, expected %s, got: %s",
			string(rebuildStatus), string(details.RebuildStatus),
		)
	}
	return nil
}

// VerifyHASwitchoverEventMetadata validates switchover event status and new path (empty for Started, device URI for Completed).
func VerifyHASwitchoverEventMetadata(record *common.EventRecord, switchoverStatus common.SwitchOverStatus, newPath string) error {

	details := GetSwitchOverEventDetails(*record)
	if details == nil {
		return fmt.Errorf("event is missing switchover details")
	}

	if details.SwitchOverStatus != switchoverStatus {
		return fmt.Errorf(
			"switchover status does not match with switchover events metadata, expected: %v, got: %v",
			switchoverStatus, details.SwitchOverStatus,
		)
	}

	if switchoverStatus == common.SwitchOverStarted {
		if details.NewPath != "" {
			return fmt.Errorf(
				"new path should be a empty string when switchover started event generated, but got: %s",
				details.NewPath,
			)
		}
	}

	if switchoverStatus == common.SwitchOverCompleted {
		if details.NewPath != newPath {
			return fmt.Errorf(
				"new path does not match with switchover events metadata, expected: %s, got: %s",
				newPath, details.NewPath,
			)
		}
	}

	return nil
}

// VerifyNexusChildEventMetadata validates that the nexus child event contains the expected URI.
func VerifyNexusChildEventMetadata(record *common.EventRecord, uri string) error {

	details := GetNexusChildEventDetails(*record)
	if details == nil {
		return fmt.Errorf("event is missing nexus child details")
	}

	if details.Uri != uri {
		return fmt.Errorf(
			"nexus child uri does not match with nexus child event metadata, expected: %s, got: %s",
			uri, details.Uri,
		)
	}
	return nil
}

// VerifyReplicaEventMetadata validates that the replica event contains the expected pool name and replica name.
func VerifyReplicaEventMetadata(record *common.EventRecord, poolName string, replicaName string) error {

	details := GetReplicaEventDetails(*record)
	if details == nil {
		return fmt.Errorf("event is missing replica details")
	}

	if details.PoolName != poolName {
		return fmt.Errorf(
			"replica pool name does not match with replica events metadata, expected: %s, got: %s",
			poolName, details.PoolName,
		)
	}

	if details.ReplicaName != replicaName {
		return fmt.Errorf(
			"replica name does not match with replica events metadata, expected: %s, got: %s",
			replicaName, details.ReplicaName,
		)
	}

	return nil
}

// VerifyNvmePathEventMetadata validates that the NVMe path event's NQN contains the given volume UUID.
func VerifyNvmePathEventMetadata(record *common.EventRecord, uuid string) error {

	details := GetNvmePathEventDetails(*record)
	if details == nil {
		return fmt.Errorf("event is missing nvme path details")
	}

	if !strings.Contains(details.Nqn, uuid) {
		return fmt.Errorf(
			"nvmePath nqn does not match with nvmePath events metadata, expected: %v, got: %v",
			uuid, details.Nqn,
		)
	}

	return nil
}

// ── Output format parsing helpers ──

// ParseJSONEventsOutput unmarshals raw JSON plugin output into EventRecord slice and validates payload fields.
func ParseJSONEventsOutput(rawOutput []byte) ([]common.EventRecord, error) {
	if len(rawOutput) == 0 {
		return nil, fmt.Errorf("empty JSON output")
	}
	var records []common.EventRecord
	if err := json.Unmarshal(rawOutput, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON output: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("JSON output contains no events")
	}
	return records, nil
}

// ParseYAMLEventsOutput unmarshals raw YAML plugin output into a generic list and validates structure.
func ParseYAMLEventsOutput(rawOutput []byte) ([]map[string]interface{}, error) {
	if len(rawOutput) == 0 {
		return nil, fmt.Errorf("empty YAML output")
	}
	var yamlData []map[string]interface{}
	if err := yaml.Unmarshal(rawOutput, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML output: %w", err)
	}
	if len(yamlData) == 0 {
		return nil, fmt.Errorf("YAML output contains no events")
	}
	return yamlData, nil
}

// VerifyEventRecordPayload checks that core EventMessage fields are populated on a record.
func VerifyEventRecordPayload(record common.EventRecord) error {
	if string(record.Category) == "" {
		return fmt.Errorf("event category is empty")
	}
	if string(record.Action) == "" {
		return fmt.Errorf("event action is empty")
	}
	if record.Metadata.Id == "" {
		return fmt.Errorf("event metadata id is empty")
	}
	if string(record.Metadata.Source.Component) == "" {
		return fmt.Errorf("event source component is empty")
	}
	if record.Metadata.EventTimestamp == "" {
		return fmt.Errorf("event timestamp is empty")
	}
	if string(record.Metadata.Version) == "" {
		return fmt.Errorf("event version is empty")
	}
	return nil
}

// VerifyYAMLEventPayload checks that core EventMessage fields are present in a parsed YAML event map.
func VerifyYAMLEventPayload(event map[string]interface{}) error {
	if event["category"] == nil {
		return fmt.Errorf("event category is missing in YAML")
	}
	if event["action"] == nil {
		return fmt.Errorf("event action is missing in YAML")
	}
	if event["metadata"] == nil {
		return fmt.Errorf("event metadata is missing in YAML")
	}
	metadata, ok := event["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("event metadata is not a map")
	}
	if metadata["id"] == nil {
		return fmt.Errorf("event metadata id is missing in YAML")
	}
	if metadata["source"] == nil {
		return fmt.Errorf("event metadata source is missing in YAML")
	}
	if metadata["timestamp"] == nil {
		return fmt.Errorf("event timestamp is missing in YAML")
	}
	if metadata["version"] == nil {
		return fmt.Errorf("event version is missing in YAML")
	}
	return nil
}

// ── CLI filter flag names ──

const (
	FlagCategory      = "--category"
	FlagAction        = "--action"
	FlagNode          = "--node"
	FlagTarget        = "--target"
	FlagComponent     = "--component"
	FlagPool          = "--pool"
	FlagVolume        = "--volume"
	FlagReplica       = "--replica"
	FlagRebuildStatus = "--rebuild-status"
	FlagState         = "--state"
	FlagFilter        = "--filter"
	FlagSince         = "--since"
	FlagLimit         = "--limit"
	FlagLokiEndpoint  = "--loki-endpoint"
	FlagNatsEndpoint  = "--nats-endpoint"
	FlagFromFile      = "--from-file"
	FlagOutputFormat  = "-o"
)

// ── CLI filter values (kebab-case — passed to --category, --action, etc.) ──

const (
	FilterCategoryPool             = "pool"
	FilterCategoryVolume           = "volume"
	FilterCategoryNexus            = "nexus"
	FilterCategoryReplica          = "replica"
	FilterCategoryNode             = "node"
	FilterCategoryHighAvailability = "high-availability"
	FilterCategoryNvmePath         = "nvme-path"
	FilterCategoryHostInitiator    = "host-initiator"
	FilterCategoryIoEngine         = "io-engine-category"
	FilterCategorySnapshot         = "snapshot"
	FilterCategoryClone            = "clone"
)

const (
	FilterActionCreate               = "create"
	FilterActionDelete               = "delete"
	FilterActionStateChange          = "state-change"
	FilterActionRebuildBegin         = "rebuild-begin"
	FilterActionRebuildEnd           = "rebuild-end"
	FilterActionSwitchOver           = "switch-over"
	FilterActionAddChild             = "add-child"
	FilterActionRemoveChild          = "remove-child"
	FilterActionOnlineChild          = "online-child"
	FilterActionNvmePathSuspect      = "nvme-path-suspect"
	FilterActionNvmePathFail         = "nvme-path-fail"
	FilterActionNvmePathFix          = "nvme-path-fix"
	FilterActionNvmeConnect          = "nvme-connect"
	FilterActionNvmeDisconnect       = "nvme-disconnect"
	FilterActionNvmeKeepAliveTimeout = "nvme-keep-alive-timeout"
	FilterActionReactorFreeze        = "reactor-freeze"
	FilterActionReactorUnfreeze      = "reactor-unfreeze"
	FilterActionShutdown             = "shutdown"
	FilterActionStart                = "start"
	FilterActionStop                 = "stop"
	FilterActionSubsystemPause       = "subsystem-pause"
	FilterActionSubsystemResume      = "subsystem-resume"
	FilterActionInit                 = "init"
	FilterActionReconfiguring        = "reconfiguring"
	FilterActionNvmePathDeleting     = "nvme-path-deleting"
)

const (
	FilterComponentCoreAgent      = "core-agent"
	FilterComponentIoEngine       = "io-engine"
	FilterComponentHaClusterAgent = "ha-cluster-agent"
	FilterComponentHaNodeAgent    = "ha-node-agent"
)

const (
	FilterRebuildStatusStarted   = "started"
	FilterRebuildStatusCompleted = "completed"
	FilterRebuildStatusStopped   = "stopped"
	FilterRebuildStatusFailed    = "failed"
)

// ── Output formats ──

const (
	OutputJSON  = "json"
	OutputYAML  = "yaml"
	OutputTable = "table"
)
