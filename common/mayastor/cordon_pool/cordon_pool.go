package cordon_pool

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	InsufficientPoolToTakeVolume = "Not enough suitable pools available"
)

// VerifyPoolCordon verifies if a pool is cordoned with specific constraints
func VerifyPoolCordon(poolID string, expectedConstraints ...common.PoolCordonConstraint) (bool, error) {
	// Create a set of expected constraints for lookup
	expectedSet := make(map[string]bool)
	for _, c := range expectedConstraints {
		if s := c.String(); s != "" {
			expectedSet[s] = true
		}
	}

	logf.Log.Info("Verifying pool cordon", "pool", poolID, "expected constraints", expectedSet)

	status, err := controlplane.GetPoolCordonStatus(poolID)
	if err != nil {
		return false, err
	}

	if !status.IsCordoned {
		return false, nil
	}

	// Check if all expected constraints are present
	for _, actualConstraint := range status.Constraints {
		if !expectedSet[actualConstraint] {
			return false, nil
		}
	}

	return true, nil
}

// VerifyPoolUncordon verifies if a pool is completely uncordoned
func VerifyPoolUncordon(poolID string) (bool, error) {
	logf.Log.Info("Verifying pool uncordon", "pool", poolID)

	status, err := controlplane.GetPoolCordonStatus(poolID)
	if err != nil {
		return false, err
	}

	return !status.IsCordoned && len(status.Constraints) == 0, nil
}

// VerifyUncordonPoolsStatus verifies that all pools are uncordoned
func VerifyUncordonPoolsStatus() (bool, error) {
	logf.Log.Info("Verifying all pools uncordon status")

	pools, err := controlplane.ListMsPools()
	if err != nil {
		return false, fmt.Errorf("failed to list mayastor pools, error %v", err)
	}

	for _, pool := range pools {
		status, err := controlplane.GetPoolCordonStatus(pool.Name)
		if err != nil {
			return false, err
		}
		if status.IsCordoned {
			logf.Log.Info("pool still cordoned", "pool", pool.Name, "constraints", status.Constraints)
			return false, nil
		}
	}
	return true, nil
}

// CancelAllCordonsOnPool utility function to remove all cordon constraints from a pool
func CancelAllCordonsOnPool(poolID string) bool {
	status, err := controlplane.GetPoolCordonStatus(poolID)
	if err != nil {
		logf.Log.Info("get pool cordon status failed", "pool", poolID, "error", err)
		return false
	}

	if !status.IsCordoned {
		return true // already uncordoned
	}

	// Remove all constraints
	typed := make([]common.PoolCordonConstraint, 0, len(status.Constraints))
	for _, c := range status.Constraints {
		switch c {
		case "replicas":
			typed = append(typed, common.CordonReplicas)
		case "snapshots":
			typed = append(typed, common.CordonSnapshots)
		case "restores":
			typed = append(typed, common.CordonRestores)
		}
	}
	err = controlplane.UnCordonPool(poolID, typed...)
	if err != nil {
		logf.Log.Info("uncordon pool failed", "pool", poolID, "error", err)
		return false
	}

	return true
}

// UncordonAllPools utility function to uncordon all pools in the cluster
func UncordonAllPools() error {
	pools, err := controlplane.ListMsPools()
	if err != nil {
		return err
	}

	for _, pool := range pools {
		if !CancelAllCordonsOnPool(pool.Name) {
			err = fmt.Errorf("CancelAllCordonsOnPool(%s) failed", pool.Name)
		}
	}
	return err
}
