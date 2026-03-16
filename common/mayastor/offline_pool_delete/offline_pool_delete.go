package offline_pool_delete

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// FixME : error messages related to offline pool delete

var (
	PoolOnlineState                         = "purge not allowed"
	DeleteWithoutPurgeFlag                  = "purge flag is required to delete pool, pool is not in online state"
	DeleteOfflinePoolWithoutCordon          = "cannot delete pool that is offline without cordoning it first"
	DeleteOfflinePoolWithReplica            = "pool has replicas, cannot delete without confirm flag"
	DeleteOfflinePoolWithSnapshots          = "pool has snapshots, cannot delete without confirm snapshot loss flag"
	DeleteOfflinePoolWithOnlyReplicaCordon  = "pool must be cordoned with snapshot and replica flag to delete pool"
	DeleteOfflinePoolWithOnlySnapshotCordon = "pool must be cordoned with snapshot and replica flag to delete pool"
	DeleteOfflinePoolWithData               = "pool has data, cannot delete without confirm data loss flag" // last healthy replica is scheduled on this pool
)

// DeleteOfflinePool deletes the offline pool via control plane plugin
func DeleteOfflinePool(poolName string, flags ...common.OfflinePoolDelete) error {
	logf.Log.Info("Deleting offline pool via plugin", "pool", poolName, "flags", flags)
	err := controlplane.DeleteOfflinePoolViaPlugin(poolName, flags...)
	if err != nil {
		logf.Log.Error(err, "Failed to delete offline pool via plugin", "pool", poolName, "flags", flags)
		return err
	}
	return nil
}

// GetNodeNameFromPool returns the node name for a given pool.
// It is a small helper used by multiple tests that need to
// map a DiskPool/MSP name to the corresponding node.
func GetNodeNameFromPool(poolName string) (string, error) {
	poolsInCluster, err := controlplane.ListMsPools()
	if err != nil {
		return "", fmt.Errorf("failed to list disk pools: %v", err)
	}

	for _, pool := range poolsInCluster {
		if pool.Name == poolName {
			return pool.Spec.Node, nil
		}
	}

	return "", fmt.Errorf("pool %s not found in cluster", poolName)
}
