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
		case common.DrainSnapshotPolicyIgnore:
			args = append(args, "--snapshot-policy", "ignore")
		case common.DrainUnsafeEvict:
			args = append(args, "--unsafe-evict")
		}
	}
	if unsafeRebuildOtherwiseEvict != nil {
		args = append(args, "--unsafe-rebuild-otherwise-evict", fmt.Sprintf("%ds", int(unsafeRebuildOtherwiseEvict.Seconds())))
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

// PoolDrainStats mirrors the REST PoolDrainStats schema: a pool usage snapshot.
type PoolDrainStats struct {
	ReplicaCount  uint64  `json:"replicaCount"`
	SnapshotCount uint64  `json:"snapshotCount"`
	Used          uint64  `json:"used"`
	Committed     *uint64 `json:"committed,omitempty"`
}

// PoolDrainState mirrors the REST PoolDrainState schema, surfaced at `get pool <id>`'s
// .meta.drain. It only carries the number of moving replicas.
type PoolDrainState struct {
	Phase             string          `json:"phase"`
	Reason            *string         `json:"reason,omitempty"`
	Initial           *PoolDrainStats `json:"initial,omitempty"`
	NumMovingReplicas uint32          `json:"numMovingReplicas"`
}

// PoolReplicaMove mirrors the REST ReplicaMoveConfig schema: one in-flight replica move.
type PoolReplicaMove struct {
	Volume             string  `json:"volume"`
	PlacementStartedAt *string `json:"placement_started_at,omitempty"`
	// MovingReplica is absent once the replica is evacuated.
	MovingReplica *string `json:"moving_replica,omitempty"`
	SpareReplica  *string `json:"spare_replica,omitempty"`
	// UnwindSpare is set while the move's spare is being removed (UnwindCancelled/UnwindRespare).
	UnwindSpare *string `json:"unwind_spare,omitempty"`
}

// PoolDrainStateExt mirrors the REST PoolDrainStateExt schema, returned by `get drain pool <id>`.
type PoolDrainStateExt struct {
	Phase   string          `json:"phase"`
	Reason  *string         `json:"reason,omitempty"`
	Initial *PoolDrainStats `json:"initial,omitempty"`
	// Current is absent while the drain is Queued or the pool is Offline.
	Current        *PoolDrainStats   `json:"current,omitempty"`
	MovingReplicas []PoolReplicaMove `json:"movingReplicas"`
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
	UnsafeRebuildOtherwiseEvict *string `json:"unsafeRebuildOtherwiseEvict,omitempty"`
	UnsafeEvict                 bool    `json:"unsafeEvict"`
}

// PoolDrainSpec mirrors the REST PoolDrainSpec schema, at .spec.cordonDrain.drain.
type PoolDrainSpec struct {
	RequestTimestamp string             `json:"requestTimestamp"`
	Policy           PoolDrainPolicy    `json:"policy"`
	PoolCordon       *PoolCordonedState `json:"poolCordon,omitempty"`
}

// GetPoolDrainProgress fetches a pool's drain state via `get drain pool <id>`.
func (cp CPv1) GetPoolDrainProgress(poolID string) (*PoolDrainStateExt, error) {
	cmd := GetMayastorPluginCmd("-n", common.NSMayastor(), "-ojson", "get", "drain", "pool", poolID)
	out, err := cmd.CombinedOutput()
	if err = CheckPluginError(out, err); err != nil {
		return nil, fmt.Errorf("failed to get drain of pool %s, error %v, output: %s", poolID, err, string(out))
	}
	var state PoolDrainStateExt
	if err = json.Unmarshal(out, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drain of pool %s, error %v, output: %s", poolID, err, string(out))
	}
	return &state, nil
}

// GetPoolDrainProgressOrNil is like GetPoolDrainProgress but returns (nil, nil) when
// the pool has no drain state, e.g. after the drain was cancelled.
func (cp CPv1) GetPoolDrainProgressOrNil(poolID string) (*PoolDrainStateExt, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its drain state, error %v", poolID, err)
	}
	if pool.Meta == nil || pool.Meta.Drain == nil {
		return nil, nil
	}
	return cp.GetPoolDrainProgress(poolID)
}

// GetPoolDrainPhase returns a pool's drain phase from `get pool`, or "" when it has no drain state.
func (cp CPv1) GetPoolDrainPhase(poolID string) (string, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return "", fmt.Errorf("failed to get pool %s to read its drain phase, error %v", poolID, err)
	}
	if pool.Meta == nil || pool.Meta.Drain == nil {
		return "", nil
	}
	return pool.Meta.Drain.Phase, nil
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

// Unwind-reason constants of a replica move (UnwindSpare in the REST schema).
const (
	UnwindCancelled = "Cancelled"
	UnwindRespare   = "Respare"
)

// GetPoolLiveUsage fetches a pool's current (live) replica/snapshot/usage tallies,
// to diff against a PoolDrainStateExt's Initial - the record itself only ever
// carries the initial snapshot, never a persisted "current".
func (cp CPv1) GetPoolLiveUsage(poolID string) (*PoolDrainStats, error) {
	pool, err := GetMayastorCpPool(poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool %s to read its live usage, error %v", poolID, err)
	}
	if pool.Meta == nil {
		return nil, fmt.Errorf("pool %s has no meta in its GET response", poolID)
	}
	committed := pool.State.Committed
	return &PoolDrainStats{
		ReplicaCount:  pool.Meta.ReplicaCount,
		SnapshotCount: pool.Meta.SnapshotCount,
		Used:          pool.State.Used,
		Committed:     &committed,
	}, nil
}
