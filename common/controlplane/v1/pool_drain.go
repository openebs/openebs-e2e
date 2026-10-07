package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// DrainPool requests a drain on a pool. unsafeRebuildOtherwiseEvict force-evicts an
// unplaceable replica after that grace period; pass nil to disable forced eviction.
// Re-issuing the drain on a pool is only accepted while it's Queued or PartiallyDrained.
func (cp CPv1) DrainPool(poolID string, unsafeRebuildOtherwiseEvict *time.Duration, opts ...common.PoolDrainOption) error {
	optStrings := make([]string, 0, len(opts))
	for _, opt := range opts {
		if flagName := opt.String(); flagName != "" {
			optStrings = append(optStrings, flagName)
		}
	}
	logf.Log.Info("Executing drain pool command", "pool", poolID, "options", optStrings, "unsafeRebuildOtherwiseEvict", unsafeRebuildOtherwiseEvict)

	args := []string{"-n", common.NSMayastor(), "drain", "pool"}

	for _, opt := range opts {
		switch opt {
		case common.DrainSnapshotPolicyAcceptLoss:
			args = append(args, "--snapshot-policy", "accept-loss")
		case common.DrainUnsafeEvict:
			args = append(args, "--unsafe-evict")
		}
	}
	if unsafeRebuildOtherwiseEvict != nil {
		args = append(args, "--unsafe-rebuild-otherwise-evict", fmt.Sprintf("%.0f", unsafeRebuildOtherwiseEvict.Seconds()))
	}

	args = append(args, poolID)

	logf.Log.Info("Executing command", "command", "kubectl mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to drain pool %s with options %v, error %v, output: %s", poolID, optStrings, err, out.String())
	}

	logf.Log.Info("Successfully requested drain on pool", "pool", poolID, "options", optStrings, "output", out.String())
	return nil
}

// AbortPoolDrain cancels a queued or in-progress drain: there is no dedicated abort
// command, it's just uncordoning the pool with the drain scope.
func (cp CPv1) AbortPoolDrain(poolID string) error {
	logf.Log.Info("Executing abort drain (uncordon pool --drain) command", "pool", poolID)

	args := []string{"-n", common.NSMayastor(), "uncordon", "pool", "--drain", poolID}

	logf.Log.Info("Executing command", "command", "kubectl mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to abort drain on pool %s, error %v, output: %s", poolID, err, out.String())
	}

	logf.Log.Info("Successfully aborted drain", "pool", poolID, "output", out.String())
	return nil
}

// PoolDrainUsage mirrors the REST PoolDrainUsage schema: a pool usage snapshot.
type PoolDrainUsage struct {
	ReplicaCount  uint64  `json:"replicaCount"`
	SnapshotCount uint64  `json:"snapshotCount"`
	Used          uint64  `json:"used"`
	Committed     *uint64 `json:"committed,omitempty"`
}

// PoolDrainRecord mirrors the REST PoolDrainRecord schema, surfaced at
// `get pool <id>`'s .meta.drain - there is no separate `get drain pool` command.
// Persists after the drain reaches a terminal phase, as the pool's last-drain
// record; cleared only when the drain is cancelled. Current usage isn't part of
// this record - it's the pool's own live .meta.replicaCount/.meta.snapshotCount
// and .state.used/.state.committed, diffed against Initial by the caller.
type PoolDrainRecord struct {
	Phase  string  `json:"phase"`
	Reason *string `json:"reason,omitempty"`
	// Initial is absent while the drain is still Queued.
	Initial *PoolDrainUsage `json:"initial,omitempty"`
	// MovingReplicas are the replica ids this drain is currently moving off the pool.
	MovingReplicas []string `json:"movingReplicas"`
}

// Pool drain phase constants (PoolDrainPhase in the REST schema).
const (
	PoolDrainPhaseUnknown          = "Unknown"
	PoolDrainPhaseQueued           = "Queued"
	PoolDrainPhaseDraining         = "Draining"
	PoolDrainPhaseAwaitingCleanup  = "AwaitingCleanup"
	PoolDrainPhasePartiallyDrained = "PartiallyDrained"
	PoolDrainPhaseDrained          = "Drained"
	PoolDrainPhaseCancelled        = "Cancelled"
)

// Pool drain phase-reason constants (PoolDrainPhaseReason in the REST schema).
const (
	PoolDrainPhaseReasonUnknown              = "Unknown"
	PoolDrainPhaseReasonWaitingForSlot       = "WaitingForSlot"
	PoolDrainPhaseReasonOfflinePool          = "OfflinePool"
	PoolDrainPhaseReasonSingleReplicaEvction = "SingleReplicaEviction"
	PoolDrainPhaseReasonImportCordoned       = "ImportCordoned"
	PoolDrainPhaseReasonSnapshotsRetained    = "SnapshotsRetained"
)

// Pool drain snapshot-policy constants (PoolDrainSnapshotPolicy in the REST schema).
const (
	PoolDrainSnapshotPolicyIgnore     = "Ignore"
	PoolDrainSnapshotPolicyAcceptLoss = "AcceptLoss"
)

// PoolDrainPolicy mirrors the REST PoolDrainPolicy schema.
type PoolDrainPolicy struct {
	SnapshotPolicy              string  `json:"snapshotPolicy"`
	UnsafeRebuildOtherwiseEvict *uint64 `json:"unsafeRebuildOtherwiseEvict,omitempty"`
	UnsafeEvict                 bool    `json:"unsafeEvict"`
}

// PoolDrainSpec mirrors the REST PoolDrainSpec schema, at .spec.cordonDrain.drain.
type PoolDrainSpec struct {
	RequestTimestamp string             `json:"requestTimestamp"`
	Policy           PoolDrainPolicy    `json:"policy"`
	UserCordon       *PoolCordonedState `json:"userCordon,omitempty"`
}

// GetPoolDrainProgress fetches a pool's drain record via `get pool <id>` -
// `get drain pool` does not exist, drain progress rides on the regular pool GET.
func (cp CPv1) GetPoolDrainProgress(poolID string) (*PoolDrainRecord, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its drain progress, error %v", poolID, err)
	}
	if pool.Meta == nil || pool.Meta.Drain == nil {
		return nil, fmt.Errorf("pool %s has no drain record - has a drain ever been requested on it?", poolID)
	}
	return pool.Meta.Drain, nil
}

// GetPoolDrainProgressOrNil fetches a pool's drain record like
// GetPoolDrainProgress, but returns (nil, nil) instead of an error when no
// drain record exists yet - for callers checking whether a record was cleared
// (e.g. after abort/uncordon), where "no record" is an expected outcome, not a
// failure.
func (cp CPv1) GetPoolDrainProgressOrNil(poolID string) (*PoolDrainRecord, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its drain progress, error %v", poolID, err)
	}
	if pool.Meta == nil {
		return nil, nil
	}
	return pool.Meta.Drain, nil
}

// GetPoolDrainSpec fetches a pool's requested drain policy.
func (cp CPv1) GetPoolDrainSpec(poolID string) (*PoolDrainSpec, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its drain spec, error %v", poolID, err)
	}
	if pool.Spec.CordonDrain == nil || pool.Spec.CordonDrain.Drain == nil {
		return nil, fmt.Errorf("pool %s has no drain spec - has a drain ever been requested on it?", poolID)
	}
	return pool.Spec.CordonDrain.Drain, nil
}

// PoolSpareReplica speculatively mirrors a proposed per-move REST schema - not real yet.
type PoolSpareReplica struct {
	ReplicaId *string `json:"replicaId,omitempty"`
}

// Unwind-reason constants for a not-yet-real per-move REST schema.
const (
	UnwindAbort   = "Abort"
	UnwindRespare = "Respare"
)

// PoolReplicaMove speculatively mirrors a proposed per-move REST schema - not real yet.
type PoolReplicaMove struct {
	Volume string `json:"volume"`
	// PlacementStartedAt is absent until the first spare-placement attempt.
	PlacementStartedAt *string `json:"placementStartedAt,omitempty"`
	// MovingReplica is the replica on the draining pool; absent once removed.
	MovingReplica *string `json:"movingReplica,omitempty"`
	// SpareReplica's presence/absence is the real signal for the BDD's "need_spare".
	SpareReplica *PoolSpareReplica `json:"spareReplica,omitempty"`
	// Unwind is set while this move's spare is being torn down (UnwindAbort/UnwindRespare).
	Unwind *string `json:"unwind,omitempty"`
}

// mayastorCpVolumeMeta is a guessed envelope for a `get volume` field that doesn't exist yet.
type mayastorCpVolumeMeta struct {
	Meta struct {
		ReplicaMove *PoolReplicaMove `json:"replicaMove,omitempty"`
	} `json:"meta"`
}

// GetVolumeReplicaMove speculatively fetches a volume's replica-move marker - TODO: guessed field.
func (cp CPv1) GetVolumeReplicaMove(volUuid string) (*PoolReplicaMove, error) {
	args := []string{"-n", common.NSMayastor(), "-ojson", "get", "volume", volUuid}
	cmd := GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("plugin failed to get volume %s, error %v, output: %s", volUuid, err, out.String())
	}

	var vol mayastorCpVolumeMeta
	if err := json.Unmarshal(out.Bytes(), &vol); err != nil {
		return nil, fmt.Errorf("failed to unmarshal volume %s output, error %v", volUuid, err)
	}
	return vol.Meta.ReplicaMove, nil
}

// GetPoolLiveUsage fetches a pool's current (live) replica/snapshot/usage tallies,
// to diff against a PoolDrainRecord's Initial - the record itself only ever
// carries the initial snapshot, never a persisted "current".
func (cp CPv1) GetPoolLiveUsage(poolID string) (*PoolDrainUsage, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its live usage, error %v", poolID, err)
	}
	if pool.Meta == nil {
		return nil, fmt.Errorf("pool %s has no meta in its GET response", poolID)
	}
	committed := pool.State.Committed
	return &PoolDrainUsage{
		ReplicaCount:  pool.Meta.ReplicaCount,
		SnapshotCount: pool.Meta.SnapshotCount,
		Used:          pool.State.Used,
		Committed:     &committed,
	}, nil
}
