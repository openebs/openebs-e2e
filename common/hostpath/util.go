package hostpath

import (
	"encoding/json"
	"fmt"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type HostPathDeviceNodeConfig struct {
	NodeDeviceMap map[string]e2e_agent.LoopDevice
}

func (hostPathConfig *HostPathDeviceNodeConfig) RemoveConfiguredHostPathDevices() error {
	logf.Log.Info("Deleting  Hostpath device on nodes")
	for node, device := range hostPathConfig.NodeDeviceMap {
		logf.Log.Info("Deleting hostpath mount point", "node", node, "mount point", device.MountPoint)
		if hostPathConfig.NodeDeviceMap[node].MountPoint == "" {
			return fmt.Errorf("device mount point not found for node %s", hostPathConfig.NodeDeviceMap[node].MountPoint)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			err = e2e_agent.RemoveHostPathDisk(*nodeIp, device.DiskPath, device.MountPoint)
			if err != nil {
				return fmt.Errorf("failed to delete mount point %s on node %s, error: %v", device.MountPoint, node, err)
			}
		}
	}

	logf.Log.Info("Verifying and deleting loop device on nodes if required")
	return k8stest.RemoveConfiguredLoopDeviceOnNodes(hostPathConfig.NodeDeviceMap)
}

func (hostPathConfig *HostPathDeviceNodeConfig) ConfigureHostPathDevices() error {
	var err error
	nodeMountPoints := make(map[string]string)
	for node, device := range hostPathConfig.NodeDeviceMap {
		nodeMountPoints[node] = device.MountPoint
	}

	nodeDeviceMap, err := k8stest.ConfigureLoopDeviceOnNodes(hostPathConfig.NodeDeviceMap)
	if err != nil {
		return fmt.Errorf("failed to create loop device, error: %v", err)
	}

	for node, device := range nodeDeviceMap {
		mountPoint := nodeMountPoints[node]
		hostPathConfig.NodeDeviceMap[node] = e2e_agent.LoopDevice{
			Size:       device.Size,
			ImageName:  device.ImageName,
			ImgDir:     device.ImgDir,
			DiskPath:   device.DiskPath,
			MountPoint: mountPoint,
		}
	}

	logf.Log.Info("Creating Hostpath device on nodes")
	for node, device := range hostPathConfig.NodeDeviceMap {
		logf.Log.Info("Creating Hostpath disk", "node", node, "device path", device.DiskPath, "mount", device.MountPoint)
		if hostPathConfig.NodeDeviceMap[node].MountPoint == "" {
			return fmt.Errorf("device mount point not found for node %s", node)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			err = e2e_agent.CreateHostPathDisk(*nodeIp, device.DiskPath, hostPathConfig.NodeDeviceMap[node].MountPoint)
			if err != nil {
				return fmt.Errorf("failed to create hostpath device on node %s, error: %v", node, err)
			}
		}
	}
	return nil
}

func GetHostPathAnnotationConfig() (map[string]string, error) {
	type hostpathConfig struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	config := []hostpathConfig{
		{Name: "StorageType", Value: "hostpath"},
		{Name: "BasePath", Value: "/var/local-hostpath"},
	}
	// Marshal the configuration to a JSON string
	configBytes, err := json.Marshal(config)
	if err != nil {
		logf.Log.Info("Error marshalling configuration:", "Error", err)
		return nil, err

	}
	configString := string(configBytes)

	return map[string]string{
		"openebs.io/cas-type":   "local",
		"cas.openebs.io/config": configString,
	}, err
}
