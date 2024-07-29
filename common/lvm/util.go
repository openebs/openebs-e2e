package lvm

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ListLvmNode list all nodes where lvm daemonset pods are scheduled
func ListLvmNode(namespace string) ([]string, error) {
	lvmDaemonSet, err := k8stest.GetDaemonSet(e2e_config.GetConfig().Product.LvmEngineDaemonSetName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get lvm DaemonSet %s, error: %v", e2e_config.GetConfig().Product.LvmEngineDaemonSetName, err)
	}
	readyCount := lvmDaemonSet.Status.NumberReady
	labelKey := e2e_config.GetConfig().Product.LocalEngineComponentPodLabelKey
	labelValue := e2e_config.GetConfig().Product.LvmEngineComponentDsPodLabelValue
	label := map[string]string{
		labelKey: labelValue,
	}

	listLvmDaemonSetPodList, err := k8stest.ListPodsWithLabel(namespace, label)
	if err != nil {
		return nil, fmt.Errorf("failed to list lvm daemonset pod with label%v, error: %v", label, err)
	}
	if int(readyCount) != len(listLvmDaemonSetPodList.Items) {
		logf.Log.Info("LVM daemonset pod count does not match with lvm daemonset ready count",
			"pod count:", len(listLvmDaemonSetPodList.Items),
			"lvm daemonset ready count:", readyCount,
		)
		return nil, fmt.Errorf("LVM daemonset pod count %d does not match with lvm daemonset ready count %d",
			len(listLvmDaemonSetPodList.Items),
			readyCount,
		)
	}
	lvmNode := make([]string, readyCount)
	for _, item := range listLvmDaemonSetPodList.Items {
		lvmNode = append(lvmNode, item.Spec.NodeName)
	}
	return lvmNode, err
}

// func ConfigureLoopDeviceOnLvmNodes(namespace string, size int64, imageDir string) (map[string]e2e_agent.LoopDevice, error) {
// 	lvmNodes, err := ListLvmNode(namespace)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return k8stest.ConfigureLoopDeviceOnNodes(lvmNodes, size, imageDir)
// }

type LvmNodesDevicePvVgConfig struct {
	VgName        string
	NodeDeviceMap map[string]e2e_agent.LoopDevice
}

func (lvmNodesDevicePvVgConfig *LvmNodesDevicePvVgConfig) RemoveConfiguredLvmNodesWithDeviceAndVg() error {
	logf.Log.Info("Deleting  ZFS pool on nodes")
	for node, device := range lvmNodesDevicePvVgConfig.NodeDeviceMap {
		logf.Log.Info("Deleting lvm vg and pv", "node", node, "vg", lvmNodesDevicePvVgConfig.VgName, "disk path", device.DiskPath)
		if lvmNodesDevicePvVgConfig.VgName == "" {
			return fmt.Errorf("lvm vg name not found for node %s", lvmNodesDevicePvVgConfig.VgName)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			_, lvmVgErr := e2e_agent.LvmRemoveVg(*nodeIp, lvmNodesDevicePvVgConfig.VgName)
			if err != nil {
				return fmt.Errorf("failed to delete lvm vg %s on node %s, error: %v", lvmNodesDevicePvVgConfig.VgName, node, lvmVgErr)
			}
			if device.DiskPath == "" {
				return fmt.Errorf("device disk path not found for node %s", device.DiskPath)
			} else {
				_, lvmPvErr := e2e_agent.LvmRemovePv(*nodeIp, device.DiskPath)
				if err != nil {
					return fmt.Errorf("failed to delete lvm pv %s on node %s, error: %v", device.DiskPath, node, lvmPvErr)
				}

			}
		}
	}
	logf.Log.Info("Verifying and deleting loop device on nodes if required")
	return k8stest.RemoveConfiguredLoopDeviceOnNodes(lvmNodesDevicePvVgConfig.NodeDeviceMap)
}

func (lvmNodesDevicePvVgConfig *LvmNodesDevicePvVgConfig) ConfigureLvmNodesWithDeviceAndVg() error {
	var err error
	if lvmNodesDevicePvVgConfig.VgName == "" {
		return fmt.Errorf("lvm vg name not found")
	}
	logf.Log.Info("Verifying and creating loop device on nodes if required")
	lvmNodesDevicePvVgConfig.NodeDeviceMap, err = k8stest.ConfigureLoopDeviceOnNodes(lvmNodesDevicePvVgConfig.NodeDeviceMap)
	if err != nil {
		return fmt.Errorf("failed to create loop device, error: %v", err)
	}

	logf.Log.Info("Creating LVM pv and vg on nodes")
	for node, device := range lvmNodesDevicePvVgConfig.NodeDeviceMap {
		logf.Log.Info("Creating lvm pv and vg", "node", node, "device path", device.DiskPath)
		if device.DiskPath == "" {
			return fmt.Errorf("device path not found for node %s", node)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			logf.Log.Info("Creating lvm pv", "node", node, "device path", device.DiskPath)
			_, lvmPvErr := e2e_agent.LvmCreatePv(*nodeIp, device.DiskPath)
			if err != nil {
				return fmt.Errorf("failed to create lvm pv on node %s with disk path %s, error: %v", node, device.DiskPath, lvmPvErr)
			}
			logf.Log.Info("Creating lvm vg", "node", node, "device path", device.DiskPath)
			_, lvmVgErr := e2e_agent.LvmCreateVg(*nodeIp, device.DiskPath, lvmNodesDevicePvVgConfig.VgName)
			if err != nil {
				return fmt.Errorf("failed to create lvm vg %s on node %s with disk path %s, error: %v",
					lvmNodesDevicePvVgConfig.VgName,
					node,
					device.DiskPath,
					lvmVgErr,
				)
			}
		}
	}
	return nil
}
