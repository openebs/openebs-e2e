package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// TODO: `drain pool` / `get drain pool` don't exist in the plugin yet.
// Calling these against a live cluster fails with "unrecognized subcommand".

// DrainPool requests a drain on a pool. unsafeRebuildOtherwiseEvict force-evicts an
// unplaceable replica after that grace period; pass nil to disable forced eviction.
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
		case common.DrainIgnoreSnapshots:
			args = append(args, "--ignore-snapshots")
		case common.DrainAcceptSnapshotLoss:
			args = append(args, "--accept-snapshot-loss")
		case common.DrainUnsafeEvict:
			args = append(args, "--unsafe-evict")
		case common.DrainDryRun:
			args = append(args, "--dry-run")
		}
	}
	if unsafeRebuildOtherwiseEvict != nil {
		args = append(args, "--unsafe-rebuild-otherwise-evict", unsafeRebuildOtherwiseEvict.String())
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

// PoolDrainUsage is a pool usage snapshot.
type PoolDrainUsage struct {
	ReplicaCount  uint64  `json:"replicaCount"`
	SnapshotCount uint64  `json:"snapshotCount"`
	Used          uint64  `json:"used"`
	Committed     *uint64 `json:"committed,omitempty"`
}

// PoolDrainStatistics is the baseline usage captured entering Draining, versus live usage.
type PoolDrainStatistics struct {
	Initial *PoolDrainUsage `json:"initial,omitempty"`
	Current *PoolDrainUsage `json:"current,omitempty"`
}

// PoolDrainPolicy is the snapshot/eviction policy a drain was requested with.
type PoolDrainPolicy struct {
	SnapshotPolicy              string  `json:"snapshotPolicy"`
	UnsafeRebuildOtherwiseEvict *uint64 `json:"unsafeRebuildOtherwiseEvict,omitempty"`
	UnsafeEvict                 bool    `json:"unsafeEvict"`
}

// PoolDrainSpec is the requested drain: when it was requested, its policy, and the
// user's own cordon if the pool was cordoned before the drain began.
type PoolDrainSpec struct {
	RequestTimestamp time.Time          `json:"requestTimestamp"`
	Policy           PoolDrainPolicy    `json:"policy"`
	UserCordon       *PoolCordonedState `json:"userCordon,omitempty"`
}

// PoolSpareReplica is the over-replicated spare of an in-flight move.
type PoolSpareReplica struct {
	ReplicaID *string `json:"replicaId,omitempty"`
}

// PoolReplicaMove is a single in-flight replica move off a draining pool.
type PoolReplicaMove struct {
	Volume             string            `json:"volume"`
	PlacementStartedAt *time.Time        `json:"placementStartedAt,omitempty"`
	MovingReplica      *string           `json:"movingReplica,omitempty"`
	SpareReplica       *PoolSpareReplica `json:"spareReplica,omitempty"`
	// Unwind is one of Cancelled, Respare; empty in steady state.
	Unwind string `json:"unwind,omitempty"`
}

// PoolDrainDetail is the full picture of one pool's drain, surfaced by
// `kubectl mayastor get drain pool <pool-id>`.
// TODO: verify field names/casing once the command exists for real.
type PoolDrainDetail struct {
	Spec         PoolDrainSpec       `json:"spec"`
	Phase        string              `json:"phase"`
	PhaseReason  string              `json:"phaseReason,omitempty"`
	Statistics   PoolDrainStatistics `json:"statistics"`
	ReplicaMoves []PoolReplicaMove   `json:"replicaMoves"`
}

// GetPoolDrainProgress fetches drain progress for a pool.
// TODO: verify the field names above once the command exists for real.
func (cp CPv1) GetPoolDrainProgress(poolID string) (*PoolDrainDetail, error) {
	args := []string{"-n", common.NSMayastor(), "-ojson", "get", "drain", "pool", poolID}

	logf.Log.Info("Executing command", "command", "kubectl mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get drain progress for pool %s, error %v, output: %s", poolID, err, out.String())
	}

	var detail PoolDrainDetail
	if err := json.Unmarshal(out.Bytes(), &detail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drain progress for pool %s, output %s, error %v", poolID, out.String(), err)
	}
	return &detail, nil
}
