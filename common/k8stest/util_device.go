package k8stest

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func ConfigureLoopDeviceOnNode(node string, size int64, imageDir string) (e2e_agent.LoopDevice, error) {
	var loopDevice e2e_agent.LoopDevice
	nodeIP, err := GetNodeIPAddress(node)
	if err != nil {
		return loopDevice, fmt.Errorf("failed to get node %s IP address, error: %v", node, err)
	}
	loopDevice, err = e2e_agent.CreateLoopDevice(*nodeIP, size, imageDir)
	if err != nil {
		return loopDevice, fmt.Errorf("failed to create loop device on node %s, error: %v", node, err)
	}
	return loopDevice, nil
}

func RemoveConfiguredLoopDeviceOnNode(node string, device e2e_agent.LoopDevice) error {
	nodeIP, err := GetNodeIPAddress(node)
	if err != nil {
		return fmt.Errorf("failed to get node %s IP address, error: %v", node, err)
	}
	_, err = e2e_agent.DeleteLoopDevice(*nodeIP, device.DiskPath, device.ImageName)
	if err != nil {
		return fmt.Errorf("failed to delete loop device on node %s, error: %v", node, err)
	}
	return nil
}

func ConfigureLoopDeviceOnNodes(nodeDeviceMap map[string]e2e_agent.LoopDevice) (map[string]e2e_agent.LoopDevice, error) {
	logf.Log.Info("Verifying and creating loop device on nodes if required")
	logf.Log.Info("ConfigureLoopDeviceOnNodes", "GET", nodeDeviceMap)
	for node, device := range nodeDeviceMap {
		if device.DiskPath == "" {
			logf.Log.Info("Creating loop device", "node", node)
			if device.ImgDir == "" || device.Size == 0 {
				logf.Log.Info("Not sufficient values passed to create loop device", "node", node, "device", device)
				return nil, fmt.Errorf("missing values for image directory or device size")
			} else {
				nodeLoopDevice, err := ConfigureLoopDeviceOnNode(node,
					device.Size,
					device.ImgDir)
				if err != nil {
					logf.Log.Info("failed to create loop device", "node", node)
					return nil, fmt.Errorf("failed to create loop device on node %s, error: %v", node, err)
				}
				nodeDeviceMap[node] = nodeLoopDevice
			}
		}
	}
	logf.Log.Info("ConfigureLoopDeviceOnNodes", "POST", nodeDeviceMap)
	return nodeDeviceMap, nil
}

func RemoveConfiguredLoopDeviceOnNodes(nodeDeviceMap map[string]e2e_agent.LoopDevice) error {
	logf.Log.Info("Verifying and creating loop device on nodes if required")
	for node, device := range nodeDeviceMap {
		if device.DiskPath != "" {
			logf.Log.Info("Deleting loop device", "node", node)
			if device.ImageName == "" {
				logf.Log.Info("Device image name not found to delete loop device", "node", node)
				logf.Log.Info("No loop device created, skip loop device deletion", "node", node)
			} else {
				err := RemoveConfiguredLoopDeviceOnNode(node, device)
				if err != nil {
					logf.Log.Info("failed to delete loop device", "node", node)
					return fmt.Errorf("failed to delete loop device on node %s, error: %v", node, err)
				}
			}
		} else {
			return fmt.Errorf("disk path not found for node %s", node)
		}
	}
	return nil
}
