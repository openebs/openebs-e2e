package k8stest

import (
	"context"
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"gopkg.in/yaml.v3"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GetConfiguredClusterNodePoolDevices read the diskpool configmap
// and return the contents as a map of pool device lists keyed on the node name.
func GetConfiguredClusterNodePoolDevices() (map[string][]string, error) {
	configMap, err := gTestEnv.KubeInt.CoreV1().ConfigMaps(common.NSDefault).Get(context.TODO(), "e2e-diskpools-fqn", metaV1.GetOptions{})
	nodesPoolDevices := make(map[string][]string)
	if err == nil {
		for k, v := range configMap.Data {
			var devices []string
			if err = yaml.Unmarshal([]byte(v), &devices); err == nil {
				nodesPoolDevices[k] = devices
			} else {
				log.Log.Info("GetConfiguredClusterNodePoolDevices", k, v, "err", err)
				nodesPoolDevices = make(map[string][]string)
				break
			}
		}
	}
	return nodesPoolDevices, err
}

// GetConfiguredNodePoolDevices given a nodename,
// return the array of diskpool devices configured for that node
// returns error if no pools are configured for the node
func GetConfiguredNodePoolDevices(nodeName string) ([]string, error) {
	var devices []string
	var ok bool
	nodesPoolDevices, err := GetConfiguredClusterNodePoolDevices()
	if err == nil {
		if devices, ok = nodesPoolDevices[nodeName]; !ok {
			log.Log.Info("Configure cluster node pool devices", nodesPoolDevices)
			err = fmt.Errorf("no pool devices configured for node %s", nodeName)
		}
	}
	return devices, err
}

// CreateConfiguredPools (re)create pools as defined by the configuration.
// No check is made on the status of pools
func CreateConfiguredPools() error {
	log.Log.Info("CreateConfiguredPools")
	var nodes []IOEngineNodeLocation
	nodesPoolDevices, err := GetConfiguredClusterNodePoolDevices()
	if err != nil {
		return err
	}
	nodes, err = GetIOEngineNodes()
	if err != nil {
		return fmt.Errorf("failed to get list of nodes, error: %v", err)
	}
	var errs common.ErrorAccumulator
	for _, node := range nodes {
		if poolDevices, ok := nodesPoolDevices[node.NodeName]; ok {
			for ix, device := range poolDevices {
				poolName := fmt.Sprintf("pool-%d-on-%s", ix+1, node.NodeName)
				pool, err := custom_resources.CreateMsPool(poolName, node.NodeName, []string{device})
				if err != nil {
					errs.Accumulate(fmt.Errorf("failed to create pool on %v , disks: %s, error: %v", node, device, err))
				}
				log.Log.Info("Created", "pool", pool)
			}
		}
	}
	return errs.GetError()
}

// GetCapacityAllReplicasOnPool given a poolName,
// returns totalCapacity of all the replicas on that pool
func GetCapacityAllReplicasOnPool(poolName string) (int64, error) {
	var totalCapacity int64
	totalCapacity = 0
	volumes, err := ListMsvs()
	if err == nil {
		for _, volume := range volumes {
			replicas, err := GetMsvReplicas(volume.Spec.Uuid)
			if err != nil {
				return totalCapacity, err
			}
			for _, replica := range replicas {
				if replica.Replica.Pool == poolName {
					totalCapacity += replica.Replica.Usage.Capacity
				}
			}
		}
	}
	return totalCapacity, err
}

// GetCapacityAndAllocatedAllSnapshotsAllReplicasOnPool given a poolName,
// returns totalCapacity which is sum of capacity of all the replicas on that pool
// and AllocatedAllSnapshots of all replicas on that pool
func GetCapacityAndAllocatedAllSnapshotsAllReplicasOnPool(poolName string) (int64, error) {
	var totalCapacity int64
	totalCapacity = 0
	volumes, err := ListMsvs()
	if err == nil {
		for _, volume := range volumes {
			replicas, err := GetMsvReplicas(volume.Spec.Uuid)
			if err != nil {
				return totalCapacity, err
			}
			for _, replica := range replicas {
				if replica.Replica.Pool == poolName {
					totalCapacity += replica.Replica.Usage.Capacity + replica.Replica.Usage.AllocatedAllSnapshots
				}
			}
		}
	}
	return totalCapacity, err
}
