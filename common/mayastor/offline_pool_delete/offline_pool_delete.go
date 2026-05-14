package offline_pool_delete

import (
	"bytes"
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	cpv1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/k8stest"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// FixME : error messages related to offline pool delete

var (
	// Purge is only permitted for pools in Offline/Unknown state.
	// Keep this as a substring match used by tests asserting plugin output.
	PoolOnlineState                         = "Only pools with Offline or Unknown state can be purged"
	DeleteWithoutPurgeFlag                  = "NodeNotOnline"           // Pool is not having resources in use, but is not in Online state
	DeleteWithoutPurgeFlagWithResources     = "InUse: Pool Resource id" // Pool has resources in use, and is not in Online state
	DeleteOfflinePoolWithoutCordon          = "Pool must be cordoned first. Use: cordon pool <id> --replicas --snapshots"
	DeleteOfflinePoolWithReplica            = "pool has replicas, cannot delete without confirm flag"
	DeleteOfflinePoolWithSnapshots          = "Volumes would lose their last healthy replica. Use --accept-volume-loss or --accept-data-loss to proceed"
	DeleteOfflinePoolWithOnlyReplicaCordon  = "Pool cordon must block both replicas and snapshots. Use: cordon pool <id> --replicas --snapshots"
	DeleteOfflinePoolWithOnlySnapshotCordon = "Pool cordon must block both replicas and snapshots. Use: cordon pool <id> --replicas --snapshots"
	DeleteOfflinePoolWithData               = "Volumes would lose their last healthy replica. Use --accept-volume-loss to proceed, or --accept-data-loss to also accept snapshot loss in a single flag" // last healthy replica is scheduled on this pool
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

// DeleteOfflinePoolWithOutput runs kubectl mayastor delete pool and returns combined stdout/stderr
// so tests can assert on plugin output (for example data loss / snapshot loss details).
func DeleteOfflinePoolWithOutput(poolName string, flags ...common.OfflinePoolDelete) (string, error) {
	args := []string{"-n", common.NSMayastor(), "delete", "pool", poolName}
	for _, flag := range flags {
		if s := flag.String(); s != "" {
			args = append(args, "--"+s)
		}
	}
	logf.Log.Info("Deleting offline pool via plugin", "pool", poolName, "args", args)
	cmd := cpv1.GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("plugin failed to delete pool %s, error %v, output: %s", poolName, err, out.String())
	}
	return out.String(), nil
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

// CreatePoolWithDeleteOpts creates a DiskPool CR on a given node/disk and
// immediately sets the requested delete options as annotations on the CR.
// From a caller's perspective this behaves as "create pool with delete options".
func CreatePoolWithDeleteOpts(poolName, node, disk string, opts ...string) error {
	logf.Log.Info("Creating DiskPool with delete options", "poolName", poolName, "node", node, "disk", disk, "opts", opts)

	// Create the DiskPool CR with delete options pre-configured via the v1beta3 helper.
	_, err := custom_resources.CreateMsPoolWithDeleteOpts(poolName, node, []string{disk}, opts...)
	if err != nil {
		return err
	}
	return nil
}

// CreatePoolsWithDeleteOptsOnAllNodes creates a DiskPool on each provided node and
// annotates it with the given delete-opts options (e.g. "purge", "accept",
// "accept_volume_loss", "accept_snapshot_loss", "accept_data_loss"). It returns the created pool names.

func CreatePoolsWithDeleteOptsOnAllNodes(allNodes []string, poolNamePrefix string, opts ...string) (createdPools []string, err error) {
	createdPools = make([]string, 0)

	for _, nodeName := range allNodes {
		// Get the first available disk on this node
		devices, deviceErr := k8stest.GetConfiguredNodePoolDevices(nodeName)
		if deviceErr != nil || len(devices) == 0 {
			logf.Log.Info("Skipping node without configured pool device", "node", nodeName, "err", deviceErr)
			continue
		}
		diskDevice := devices[0]

		poolName := fmt.Sprintf("%s-%s", poolNamePrefix, nodeName)
		logf.Log.Info("Creating pool for delete-opts tests", "poolName", poolName, "node", nodeName, "disk", diskDevice, "opts", opts)

		// Create the DiskPool CR and set delete options using the helper.
		if err = CreatePoolWithDeleteOpts(poolName, nodeName, diskDevice, opts...); err != nil {
			return
		}

		createdPools = append(createdPools, poolName)
	}

	if len(createdPools) == 0 {
		err = fmt.Errorf("no pools could be created on any node for delete-opts")
		return
	}
	return createdPools, nil
}

// WaitForDiskPoolCrState polls the DiskPool CR status.cr_state until it matches expected
// or the timeout elapses.
func WaitForDiskPoolCrState(poolName string, expected string, timeout, pollInterval time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for {
		st, err := custom_resources.GetDiskPoolCrStatus(poolName)
		if err == nil {
			last = st
			if st == expected {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for DiskPool CR %s cr_state=%q (last=%q)", poolName, expected, last)
		}
		time.Sleep(pollInterval)
	}
}
