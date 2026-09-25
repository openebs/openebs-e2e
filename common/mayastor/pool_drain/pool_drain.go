package pool_drain

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Phase name constants for a pool drain.
// TODO: `drain pool` / `get drain pool` don't exist in the plugin yet, so every
// function here fails with "unrecognized subcommand" against a live cluster.
const (
	PhaseQueued           = "Queued"
	PhaseDraining         = "Draining"
	PhaseAwaitingCleanup  = "AwaitingCleanup"
	PhasePartiallyDrained = "PartiallyDrained"
	PhaseDrained          = "Drained"
	PhaseCancelled        = "Cancelled"
)

// Phase-reason constants: why a pool's drain sits in its current phase.
const (
	PhaseReasonWaitingForSlot              = "WaitingForSlot"
	PhaseReasonOfflinePool                 = "OfflinePool"
	PhaseReasonSingleReplicaUnsafeEviction = "SingleReplicaUnsafeEviction"
	PhaseReasonImportCordoned              = "ImportCordoned"
	PhaseReasonSnapshotsRetained           = "SnapshotsRetained"
)

// FixMe : error messages below are drafted from the design docs/BDD, not yet
// observed live - update with the real plugin output once it ships.
var (
	// DestroyPoolWithoutDrain: deleting a pool that still has replica/snapshot
	// allocation on it (drain not yet run, or not run to completion).
	DestroyPoolWithoutDrain = "please drain the diskpool before attempting destroy"
	// DrainPoolConflictingSnapshotFlags: --ignore-snapshots and
	// --accept-snapshot-loss passed together.
	DrainPoolConflictingSnapshotFlags = "cannot use --ignore-snapshots and --accept-snapshot-loss together"
)

// VerifyPoolDrainPhase verifies a pool's drain has reached the given phase.
func VerifyPoolDrainPhase(poolID string, expectedPhase string) (bool, error) {
	logf.Log.Info("Verifying pool drain phase", "pool", poolID, "expected", expectedPhase)

	progress, err := controlplane.GetPoolDrainProgress(poolID)
	if err != nil {
		return false, err
	}

	logf.Log.Info("Pool drain progress", "pool", poolID, "phase", progress.Phase)
	return progress.Phase == expectedPhase, nil
}

// VerifyPoolDrainPhaseReason verifies a pool's drain phase reason, e.g.
// distinguishing why a pool landed at PartiallyDrained (SnapshotsRetained vs
// SingleReplicaUnsafeEviction vs ImportCordoned).
func VerifyPoolDrainPhaseReason(poolID string, expectedReason string) (bool, error) {
	progress, err := controlplane.GetPoolDrainProgress(poolID)
	if err != nil {
		return false, err
	}
	return progress.PhaseReason == expectedReason, nil
}

// VerifyPoolDrained verifies a pool has fully drained: phase Drained with zero
// replicas and zero allocation remaining.
func VerifyPoolDrained(poolID string) (bool, error) {
	progress, err := controlplane.GetPoolDrainProgress(poolID)
	if err != nil {
		return false, err
	}
	if progress.Phase != PhaseDrained {
		logf.Log.Info("Pool not yet Drained", "pool", poolID, "phase", progress.Phase)
		return false, nil
	}
	current := progress.Statistics.Current
	if current == nil {
		return false, nil
	}
	return current.ReplicaCount == 0 && current.Used == 0, nil
}

// VerifyPoolPartiallyDrained verifies a pool reached PartiallyDrained. Call
// VerifyPoolDrainPhaseReason to distinguish SnapshotsRetained, SingleReplicaUnsafeEviction, or ImportCordoned.
func VerifyPoolPartiallyDrained(poolID string) (bool, error) {
	progress, err := controlplane.GetPoolDrainProgress(poolID)
	if err != nil {
		return false, err
	}
	if progress.Phase != PhasePartiallyDrained {
		logf.Log.Info("Pool not PartiallyDrained", "pool", poolID, "phase", progress.Phase)
		return false, nil
	}
	return true, nil
}

// GetMovingReplicas returns the ids of replicas currently being moved off a draining pool.
func GetMovingReplicas(poolID string) ([]string, error) {
	progress, err := controlplane.GetPoolDrainProgress(poolID)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, move := range progress.ReplicaMoves {
		if move.MovingReplica != nil {
			ids = append(ids, *move.MovingReplica)
		}
	}
	return ids, nil
}

// AbortAllPoolDrains is an AfterEach-style cleanup utility: aborts any drain still
// queued or in progress on every pool in the cluster.
func AbortAllPoolDrains() error {
	pools, err := controlplane.ListMsPools()
	if err != nil {
		return err
	}

	var lastErr error
	for _, pool := range pools {
		progress, err := controlplane.GetPoolDrainProgress(pool.Name)
		if err != nil {
			continue
		}
		// AwaitingCleanup is already unwinding to Drained on its own - nothing left to abort.
		if progress.Phase == "" || progress.Phase == PhaseCancelled || progress.Phase == PhaseAwaitingCleanup {
			continue
		}
		logf.Log.Info("pool still draining, aborting", "pool", pool.Name, "phase", progress.Phase)
		if err := controlplane.AbortPoolDrain(pool.Name); err != nil {
			lastErr = fmt.Errorf("AbortPoolDrain(%s) failed: %v", pool.Name, err)
		}
	}
	return lastErr
}
