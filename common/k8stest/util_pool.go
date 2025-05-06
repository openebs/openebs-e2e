package k8stest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	mcpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"gopkg.in/yaml.v3"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const configMapName = "e2e-diskpools-fqn"

// GetConfiguredClusterNodePoolDevices read the diskpool configmap
// and return the contents as a map of pool device lists keyed on the node name.
func GetConfiguredClusterNodePoolDevices() (map[string][]string, error) {
	configMap, err := gTestEnv.KubeInt.CoreV1().ConfigMaps(common.NSDefault).Get(context.TODO(), configMapName, metaV1.GetOptions{})
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

// CreateDiskPoolsConfiguration query mayastor for available devices and create configmap
// e2e-diskpools-fqn with a map of nodes and pool devices
// This function is a NOP if the config map is already present
// to allow external agents to determine the set of pool devices
// by creating the configmap before this function is invoked
func CreateDiskPoolsConfiguration() error {
	// weights for sorting device links in order of preference
	var devLinkWeights = map[string]int{
		"by-uuid":     0,
		"by-id":       1,
		"by-partuuid": 2,
		"by-path":     3,
	}
	// struct to de-serialise plugin output for block devices
	type msBlockDev struct {
		Available bool     `json:"available"`
		DevLinks  []string `json:"devlinks"`
		DevName   string   `json:"devname"`
		DevPath   string   `json:"devpath"`
	}
	configMapApi := gTestEnv.KubeInt.CoreV1().ConfigMaps(common.NSDefault)
	cmLst, err := configMapApi.List(context.TODO(), metaV1.ListOptions{})

	if err == nil {
		// if the configmap already exists the use it as is
		for _, cm := range cmLst.Items {
			if cm.Name == configMapName {
				log.Log.Info("CreateDiskPoolsConfiguration: using existing config map", "name", cm.Name)
				return nil
			}
		}
		// compile a list of devices on each node.
		nodes, err := GetMayastorNodeNames()
		if err != nil {
			log.Log.Info("CreateDiskPoolsConfiguration: GetMayastorNodeName", "error", err)
			return err
		}
		cmapData := make(map[string]string)
		for _, node := range nodes {
			// retrieve block devices on a node in json format
			cmd := mcpV1.GetMayastorPluginCmd("-n", common.NSMayastor(), "get", "block-devices", node, "-o", "json")
			var out bytes.Buffer
			log.Log.Info("About to execute:", "cmd", cmd)
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				log.Log.Info("CreateDiskPoolsConfiguration:", "cmd", cmd, "error", err)
				return err
			}
			log.Log.Info("CreateDiskPoolsConfiguration:", "raw data", out.String())
			// de-serialise the json data
			var data []msBlockDev
			err = json.Unmarshal(out.Bytes(), &data)
			if err != nil {
				log.Log.Info("CreateDiskPoolsConfiguration: json Unmarshal", "error", err)
				return err
			}
			var disks []string
			for _, bd := range data {
				if !bd.Available {
					// disk is not available for use by mayastor
					continue
				}
				// ignore loop devices
				if strings.HasPrefix(bd.DevName, "/dev/loop") {
					continue
				}
				// default to device name
				device := bd.DevName
				if len(bd.DevLinks) != 0 {
					// if a list device links exists then use that as the device
					// sort in order of preference and use the first one.
					sort.Slice(bd.DevLinks, func(i int, j int) bool {
						keyI := strings.Split(bd.DevLinks[i], "/")[3]
						keyJ := strings.Split(bd.DevLinks[j], "/")[3]
						weightI := 10
						if val, ok := devLinkWeights[keyI]; ok {
							weightI = val
						}
						weightJ := 10
						if val, ok := devLinkWeights[keyJ]; ok {
							weightJ = val
						}
						if weightI == weightJ {
							// if weights are the same sort alphabetically for consitent
							// results
							return bd.DevLinks[i] < bd.DevLinks[j]
						}
						return weightI < weightJ
					})
					device = bd.DevLinks[0]
				}
				disks = append(disks, device)
			}
			if len(disks) > 0 {
				var yamlBytes []byte
				yamlBytes, err = yaml.Marshal(disks)
				if err != nil {
					log.Log.Info("CreateDiskPoolsConfiguration: yaml Marshal", "error", err)
					return err
				}
				cmapData[node] = string(yamlBytes)
			}
		}

		cmap := coreV1.ConfigMap{
			TypeMeta: metaV1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metaV1.ObjectMeta{
				Name:      configMapName,
				Namespace: common.NSDefault,
			},
			Immutable:  nil,
			Data:       cmapData,
			BinaryData: nil,
		}
		_, err = configMapApi.Create(context.TODO(), &cmap, metaV1.CreateOptions{})
		if err == nil {
			log.Log.Info("Created config map", "name", cmap.ObjectMeta.Name)
		} else {
			log.Log.Info("failed to create config map", "name", cmap.ObjectMeta.Name, "error", err)
		}
	}
	return err
}

// GetDiskPoolUsage returns pool used size
func GetDiskPoolUsage(poolName string) (uint64, error) {
	pool, err := custom_resources.GetMsPool(poolName)
	if err != nil {
		return 0, fmt.Errorf("failed to get pool %s, error: %v", poolName, err)
	}
	if pool == nil {
		return 0, fmt.Errorf("pool %s not found", poolName)
	}
	return pool.GetStatusUsed(), nil
}

// WaitForDiskPoolUsageToBeZero waits for the disk pool usage to be zero
func WaitForDiskPoolUsageToBeZero(poolName string, timeoutsecs int) error {
	const sleepTime = 5
	var err error
	for ix := 1; ix < timeoutsecs/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		var usedSize uint64
		usedSize, err := GetDiskPoolUsage(poolName)
		if err != nil {
			logf.Log.Info("Error in WaitForDiskPoolUsageToBeZero", "poolName", poolName, "error", err)
		}
		if usedSize == 0 {
			logf.Log.Info("WaitForDiskPoolUsageToBeZero", "poolName", poolName, "usedSize", usedSize)
			break
		}
	}

	return err
}
