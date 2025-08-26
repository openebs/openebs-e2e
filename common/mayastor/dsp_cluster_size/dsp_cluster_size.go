package dsp_cluster_size

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/custom_resources/types"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var poolOnlineTimeoutSec = 120 // seconds

func SetMayastorDspClusterSize(size, helmChart, helmRelease, helmVersion string) error {
	var values map[string]interface{}
	prefix := "agents.core.poolClusterSize"
	if e2e_config.GetConfig().Product.UseUmbrellaOpenEBSChart {
		prefix = fmt.Sprintf("%s.%s", e2e_config.GetConfig().Product.ChartName, prefix)
	}

	logf.Log.Info("Setting DSP cluster size for Mayastor",
		"size", size,
		"helmChart", helmChart,
		"helmRelease", helmRelease,
		"helmVersion", helmVersion,
		"prefix", prefix)
	values = map[string]interface{}{
		prefix: size,
	}

	_, err := k8stest.UpgradeHelmChart(helmChart,
		common.NSMayastor(),
		helmRelease,
		helmVersion,
		values,
	)
	if err != nil {
		return fmt.Errorf("failed to set dsp cluster size via helm, error: %v", err)
	}

	return nil
}

// CreateClusterSizePools creates cluster size pools
// It takes a map of node, disk and cluster size as input and returns a map of pool name and pool object
func CreateClusterSizePools(nodeDiskMap map[string]string, clusterSize string) (map[string]types.DiskPool, error) {
	pools := make(map[string]types.DiskPool)
	for node, disk := range nodeDiskMap {
		poolName := "pool-" + node
		pool, err := custom_resources.CreateMsPoolWithClusterSize(poolName, node, []string{disk}, clusterSize)
		if err != nil {
			return pools, err
		}
		pools[poolName] = pool
	}
	return pools, nil
}

// CreateAndWaitForClusterSizePools creates cluster size pools and waits for them to be online
// It takes a map of node, disk and cluster size as input and returns a map of pool name and pool object
// It returns an error if the pools are not online within the given time
func CreateAndWaitForClusterSizePools(nodeDiskMap map[string]string, clusterSize string) (map[string]types.DiskPool, error) {
	pools, err := CreateClusterSizePools(nodeDiskMap, clusterSize)
	if err != nil {
		return pools, err
	}

	// wait for all pools to be online
	err = k8stest.WaitForPoolsToBeOnline(poolOnlineTimeoutSec)
	if err != nil {
		return pools, fmt.Errorf("all pools are not online, error: %v", err)
	}

	// refresh the pool list
	poolList, err := custom_resources.ListMsPools()
	if err != nil {
		return pools, fmt.Errorf("failed to list pools, error: %v", err)
	}

	for _, pool := range poolList {
		if pools[pool.GetName()] != nil {
			pools[pool.GetName()] = pool

		}
	}

	// check if all pools cluster size matches
	for _, pool := range pools {
		if pool.GetClusterSize() != clusterSize {
			return pools, fmt.Errorf("pool %s cluster size %s does not match with provide size %s",
				pool.GetName(),
				pool.GetClusterSize(),
				clusterSize,
			)
		}
	}
	return pools, nil
}
