package event

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/k8stest"
)

// Version agnostic structure for obtaining the version string
type EventSparse struct {
	Metadata struct {
		Version int `json:"version"`
	} `json:"metadata"`
}

// Event action.
type Action int

// Event category.
type Category int

// Event component.
type Component int

// Outcome of rebuild.
type RebuildStatus int

type EventMessage struct {
	// Event Category.
	Category Category `json:"category"`
	// Event Action.
	Action Action `json:"action"`
	// Target id for the category against which action is performed.
	Target string `json:"target"`
	// Event meta data.
	Metadata EventMeta `json:"metadata"`
}

const (
	LIST_ALL     = "events."
	LIST_POOL    = "events.1."
	LIST_VOLUME  = "events.2."
	LIST_NEXUS   = "events.3."
	LIST_REPLICA = "events.4."
	LIST_NODE    = "events.5."
	LIST_HA      = "events.6."
	LIST_NVME    = "events.7."

	SUBSCRIBE_ALL     = "events.>"
	SUBSCRIBE_POOL    = "events.1.>"
	SUBSCRIBE_VOLUME  = "events.2.>"
	SUBSCRIBE_NEXUS   = "events.3.>"
	SUBSCRIBE_REPLICA = "events.4.>"
	SUBSCRIBE_NODE    = "events.5.>"
	SUBSCRIBE_HA      = "events.6.>"
	SUBSCRIBE_NVME    = "events.7.>"
)

//////////////////////// rebuild events

const (
	RebuildStatusUnknown   RebuildStatus = 0
	RebuildStatusStarted   RebuildStatus = 1
	RebuildStatusCompleted RebuildStatus = 2
	RebuildStatusStopped   RebuildStatus = 3
	RebuildStatusFailed    RebuildStatus = 4
)

type RebuildDetails struct {
	Error              string        `json:"error"`
	SourceReplica      string        `json:"source_replica"`
	DestinationReplica string        `json:"destination_replica"`
	RebuildStatus      RebuildStatus `json:"rebuild_status"`
}

//////////////////// replica events

// Replica event details
type ReplicaEventDetails struct {
	// Pool name
	PoolName string `json:"pool_name"`
	// Pool uuid
	PoolUuid string `json:"pool_uuid"`
	// Replica name
	ReplicaName string `json:"replica_name"`
}

///////////////////// switchover events

type SwitchOverStatus int

// Switch over status
const (
	// Unknown
	UnknownSwitchOverStatus SwitchOverStatus = 0
	// Switch over started
	SwitchOverStarted SwitchOverStatus = 1
	// Switch over is completed successfully
	SwitchOverCompleted SwitchOverStatus = 2
	// Switch over failed
	SwitchOverFailed SwitchOverStatus = 3
)

// HA switch over event details
type SwitchOverEventDetails struct {
	// Switch over status
	SwitchOverStatus SwitchOverStatus `json:"switch_over_status"`
	// Timestamp when switchover request was initialized
	StartTime string `json:"start_time"`
	// Failed nexus path of the volume
	ExistingNqn string `json:"existing_nqn"`
	// New nexus path of the volume
	NewPath string `json:"new_path"`
	// Number of failed attempts in the current Stage
	RetryCount uint64 `json:"retry_count"`
}

/////////////////////// nexus child events

// Nexus child event details
type NexusChildEventDetails struct {
	// Nexus child uri
	Uri string `json:"uri"`
}

/////////////////////// nvme path events

// Nvme path event details
type NvmePathEventDetails struct {
	Nqn  string `json:"nqn"`
	Path string `json:"path"`
}

// ///////////////////// Host initiator event details
type HostInitiatorEventDetails struct {
	HostNqn      string `json:"host_nqn"`
	SubsystemNqn string `json:"subsystem_nqn"`
	// The target on which the host is connected to the subsystem - Nexus/Replica
	Target string `json:"target"`
	// Target uuid
	Uuid string `json:"uuid"`
}

// ///////////////////// State change event details
type StateChangeEventDetails struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

// ///////////////////// Reactor event details
type ReactorEventDetails struct {
	// The logical core this reactor is created on
	Lcore uint64 `json:"lcore"`
	// Reactor state
	State string `json:"state"`
}

// ///////////////////// Snapshot event details
type SnapshotEventDetails struct {
	// Parent id of the snapshot
	ReplicaId string `json:"replica_id"`
	// Snapshot creation time
	CreateTime string `json:"create_time"`
	// Entity id of the snapshot
	VolumeId string `json:"volume_id"`
}

// ///////////////////// Clone event details
type CloneEventDetails struct {
	// Source uuid from which clone is created
	SourceUuid string `json:"source_uuid"`
	// Clone creation time
	CreateTime string `json:"create_time"`
}

// ///////////////////// Error details
type ErrorDetails struct {
	Error string `json:"error"`
}

// ///////////////////// Subsystem pause details
type SubsystemPauseDetails struct {
	// Nexus pause state
	NexusPauseState string `json:"nexus_pause_state"`
}

/////////////////////// generic details

type EventDetails struct {
	CloneEventDetails         *CloneEventDetails         `json:"clone_details"`
	ErrorDetails              *ErrorDetails              `json:"error_details"`
	HostInitiatorEventDetails *HostInitiatorEventDetails `json:"host_initiator_details"`
	NexusChildEventDetails    *NexusChildEventDetails    `json:"nexus_child_details"`
	NvmePathEventDetails      *NvmePathEventDetails      `json:"nvme_path_details"`
	ReactorEventDetails       *ReactorEventDetails       `json:"reactor_details"`
	RebuildDetails            *RebuildDetails            `json:"rebuild_details"`
	ReplicaEventDetails       *ReplicaEventDetails       `json:"replica_details"`
	SnapshotEventDetails      *SnapshotEventDetails      `json:"snapshot_details"`
	StateChangeEventDetails   *StateChangeEventDetails   `json:"state_change_details"`
	SubsystemPauseDetails     *SubsystemPauseDetails     `json:"subsystem_pause_details"`
	SwitchOverEventDetails    *SwitchOverEventDetails    `json:"switch_over_details"`
}

type EventMeta struct {
	// Something that uniquely identifies events.
	// UUIDv4.
	// GUID.
	Id     string      `json:"id"`
	Source EventSource `json:"source"`
	// Event timestamp.
	EventTimestamp string `json:"timestamp"`
	// Version of the event message.
	Version int `json:"version"`
}

// Event source.
type EventSource struct {
	// Io-engine or core-agent.
	Component Component `json:"component"`
	// Node name
	Node         string        `json:"node"`
	EventDetails *EventDetails `json:"event_details"`
}

type EventContext struct {
	E2eAgentAddress string
	NatsServer      string
}

const (
	UnknownAction              Action = 0
	ActionCreate               Action = 1
	ActionDelete               Action = 2
	ActionStateChange          Action = 3
	ActionRebuildBegin         Action = 4
	ActionRebuildEnd           Action = 5
	ActionSwitchOver           Action = 6
	ActionAddChild             Action = 7
	ActionRemoveChild          Action = 8
	ActionNvmePathSuspect      Action = 9
	ActionNvmePathFail         Action = 10
	ActionNvmePathFix          Action = 11
	ActionOnlineChild          Action = 12
	ActionNvmeConnect          Action = 13
	ActionNvmeDisconnect       Action = 14
	ActionNvmeKeepAliveTimeout Action = 15
	ActionReactorFreeze        Action = 16
	ActionReactorUnfreeze      Action = 17
	ActionShutdown             Action = 18
	ActionStart                Action = 19
	ActionStop                 Action = 20
	ActionSubsystemPause       Action = 21
	ActionSubsystemResume      Action = 22
	ActionInit                 Action = 23
	ActionReconfiguring        Action = 24
)

const (
	UnknownCategory          Category = 0
	CategoryPool             Category = 1
	CategoryVolume           Category = 2
	CategoryNexus            Category = 3
	CategoryReplica          Category = 4
	CategoryNode             Category = 5
	CategoryHighAvailability Category = 6
	CategoryNvmePath         Category = 7
	CategoryHostInitiator    Category = 8
	CategoryIoEngine         Category = 9
	CategorySnapshot         Category = 10
	CategoryClone            Category = 11
)

const (
	UnknownComponent        Component = 0
	ComponentCoreAgent      Component = 1
	ComponentIoEngine       Component = 2
	ComponentHaClusterAgent Component = 3
	ComponentHaNodeAgent    Component = 4
)

const Version = 1

func NewEventContext(natsSts string, namespace string, natsPort string) (EventContext, error) {
	var err error
	var context EventContext
	nodeIPs := k8stest.GetMayastorNodeIPAddresses()

	if len(nodeIPs) < 1 {
		return context, fmt.Errorf("no mayastor nodes")
	}
	context.E2eAgentAddress = nodeIPs[0]

	eventsServer, err := k8stest.GetPodAddress(natsSts+"-0", namespace)
	if err != nil {
		return context, fmt.Errorf("failed to get nats pod, error %s", err.Error())
	}
	eventsServer += ":" + natsPort
	context.NatsServer = eventsServer
	return context, err
}

func (context *EventContext) Subscribe(subject string) error {
	out, err := e2e_agent.EventSubscribe(context.E2eAgentAddress, context.NatsServer, subject)
	if err != nil {
		return fmt.Errorf("failed to subscribe, output: %s, error: %s", out, err.Error())
	}
	return err
}

// Used only for test development purposes. data is a json-encoded event
func (context *EventContext) PublishRaw(subject string, data string) error {
	out, err := e2e_agent.EventPublish(context.E2eAgentAddress, context.NatsServer, subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish, output: %s, error: %s", out, err.Error())
	}
	return err
}

func (context *EventContext) UnsubscribeAll() error {
	out, err := e2e_agent.EventUnsubscribeAll(context.E2eAgentAddress)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe all, output: %s, error: %s", out, err.Error())
	}
	return err
}

func (context *EventContext) GetAllEvents() ([]EventMessage, error) {
	return context.GetEvents("")
}

func (context *EventContext) GetEvents(subject_pattern string) ([]EventMessage, error) {
	var events []EventMessage
	var messages []string
	out, err := e2e_agent.EventList(context.E2eAgentAddress, subject_pattern)
	if err != nil {
		return events, fmt.Errorf("failed to get all events, output: %s, error: %s", out, err.Error())
	}
	// The output is a json-encoded array of strings.
	err = json.Unmarshal([]byte(out), &messages)
	if err != nil {
		return events, fmt.Errorf("failed to unmarshall, data %s, error %s", out, err.Error())
	}
	// Each string is a json-encoded event. Decode it and append to the array to be returned.
	for _, s := range messages {
		// get just the version which should (hopefully) work with all variants
		var es EventSparse
		err = json.Unmarshal([]byte(s), &es)
		if err != nil {
			return events, fmt.Errorf("failed to unmarshall version info, data %s, error %s", s, err.Error())
		}
		if es.Metadata.Version != Version {
			return events, fmt.Errorf("unsupported event format, version %d, expected %d", es.Metadata.Version, Version)
		}
		var e EventMessage
		err = json.Unmarshal([]byte(s), &e)
		if err != nil {
			return events, fmt.Errorf("failed to unmarshall message, data %s, error %s", s, err.Error())
		}
		events = append(events, e)
	}
	return events, err
}

func CheckRebuildEvent(
	rebuildEvent *EventMessage,
	source_uri string,
	dest_uri string,
	has_err bool,
	rebuild_status RebuildStatus) error {

	if has_err {
		if rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.Error == "" {
			return fmt.Errorf("event is missing error message")
		}
	} else {
		if rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.Error != "" {
			return fmt.Errorf("event has error message: %s", rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.Error)
		}
	}

	if rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.SourceReplica != source_uri {
		return fmt.Errorf(
			"event has invalid source uri, expected %s, got: %s",
			source_uri, rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.SourceReplica,
		)
	}

	if rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.DestinationReplica != dest_uri {
		return fmt.Errorf(
			"event has invalid destination uri, expected %s, got: %s",
			dest_uri, rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.DestinationReplica,
		)
	}

	if rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.RebuildStatus != rebuild_status {
		return fmt.Errorf(
			"event has invalid rebuild status, expected %d, got: %d",
			int(rebuild_status), int(rebuildEvent.Metadata.Source.EventDetails.RebuildDetails.RebuildStatus),
		)
	}
	return nil
}

func VerifyReplicaEventMetadata(replicaEvent *EventMessage, poolName string, replicaName string) error {

	if replicaEvent.Metadata.Source.EventDetails.ReplicaEventDetails.PoolName != poolName {
		return fmt.Errorf(
			"replica pool name does not match with replica events metadata, expected: %s, got: %s",
			poolName, replicaEvent.Metadata.Source.EventDetails.ReplicaEventDetails.PoolName,
		)
	}

	if replicaEvent.Metadata.Source.EventDetails.ReplicaEventDetails.ReplicaName != replicaName {
		return fmt.Errorf(
			"replica name does not match with replica events metadata, expected: %s, got: %s",
			poolName, replicaEvent.Metadata.Source.EventDetails.ReplicaEventDetails.ReplicaName,
		)
	}

	return nil
}

func VerifyNexusChildEventMetadata(nexusChildEvent *EventMessage, uri string) error {
	if nexusChildEvent.Metadata.Source.EventDetails.NexusChildEventDetails.Uri != uri {
		return fmt.Errorf(
			"nexus child uri does not match with nexus child event metadata, expected: %s, got: %s",
			uri, nexusChildEvent.Metadata.Source.EventDetails.NexusChildEventDetails.Uri,
		)
	}
	return nil
}

func VerifyHASwitchoverEventMetadata(switchoverEvent *EventMessage, switchoverStatus SwitchOverStatus, newPath string) error {

	if switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.SwitchOverStatus != switchoverStatus {
		return fmt.Errorf(
			"switchover status does not match with switchover events metadata, expected: %v, got: %v",
			switchoverStatus, switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.SwitchOverStatus,
		)
	}

	if switchoverStatus == SwitchOverStarted {
		if switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.NewPath != "" {
			return fmt.Errorf(
				"new path should be a empty string when switchover started event generated, but got: %s",
				switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.NewPath,
			)
		}
	}

	if switchoverStatus == SwitchOverCompleted {
		if switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.NewPath != newPath {
			return fmt.Errorf(
				"new path does not match with switchover events metadata, expected: %s, got: %s",
				newPath, switchoverEvent.Metadata.Source.EventDetails.SwitchOverEventDetails.NewPath,
			)
		}
	}

	return nil
}

func VerifyNvmePathEventMetadata(nvmePathEvent *EventMessage, uuid string) error {

	if !strings.Contains(nvmePathEvent.Metadata.Source.EventDetails.NvmePathEventDetails.Nqn, uuid) {
		return fmt.Errorf(
			"nvmePath nqn does not match with nvmePath events metadata, expected: %v, got: %v",
			uuid, nvmePathEvent.Metadata.Source.EventDetails.NvmePathEventDetails.Nqn,
		)
	}

	return nil
}
