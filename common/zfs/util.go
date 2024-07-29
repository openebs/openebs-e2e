package zfs

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ListZfsNode list all nodes where zfs daemonset pods are scheduled
func ListZfsNode(namespace string) ([]string, error) {
	zfsDaemonSet, err := k8stest.GetDaemonSet(e2e_config.GetConfig().Product.ZfsEngineDaemonSetName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get lvm DaemonSet %s, error: %v", e2e_config.GetConfig().Product.LvmEngineDaemonSetName, err)
	}
	readyCount := zfsDaemonSet.Status.NumberReady
	labelKey := e2e_config.GetConfig().Product.LocalEngineComponentPodLabelKey
	labelValue := e2e_config.GetConfig().Product.ZfsEngineComponentDsPodLabelValue
	label := map[string]string{
		labelKey: labelValue,
	}

	listZfsDaemonSetPodList, err := k8stest.ListPodsWithLabel(namespace, label)
	if err != nil {
		return nil, fmt.Errorf("failed to list zfs daemonset pod with label%v, error: %v", label, err)
	}
	if int(readyCount) != len(listZfsDaemonSetPodList.Items) {
		logf.Log.Info("ZFS daemonset pod count does not match with zfs daemonset ready count",
			"pod count:", len(listZfsDaemonSetPodList.Items),
			"zfs daemonset ready count:", readyCount,
		)
		return nil, fmt.Errorf("ZFS daemonset pod count %d does not match with zfs daemonset ready count %d",
			len(listZfsDaemonSetPodList.Items),
			readyCount,
		)
	}
	zfsNode := make([]string, readyCount)
	for _, item := range listZfsDaemonSetPodList.Items {
		zfsNode = append(zfsNode, item.Spec.NodeName)
	}
	return zfsNode, err
}

// func ConfigureLoopDeviceOnZfsNodes(namespace string, size int64, imageDir string) (map[string]e2e_agent.LoopDevice, error) {
// 	zfsNodes, err := ListZfsNode(namespace)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return k8stest.ConfigureLoopDeviceOnNodes(zfsNodes, size, imageDir)
// }

type ZfsNodesDevicePoolConfig struct {
	PoolName      string
	NodeDeviceMap map[string]e2e_agent.LoopDevice
}

func (zfsDevicePoolConfig *ZfsNodesDevicePoolConfig) RemoveConfiguredDeviceZfsPool() error {
	logf.Log.Info("Deleting  ZFS pool on nodes")
	for node := range zfsDevicePoolConfig.NodeDeviceMap {
		logf.Log.Info("Deleting zfs pool", "node", node, "zfs pool name", zfsDevicePoolConfig.PoolName)
		if zfsDevicePoolConfig.PoolName == "" {
			return fmt.Errorf("device path not found for node %s", zfsDevicePoolConfig.PoolName)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			_, zfsPoolErr := e2e_agent.ZfsDestroyPool(*nodeIp, zfsDevicePoolConfig.PoolName)
			if err != nil {
				return fmt.Errorf("failed to delete zfs pool %s on node %s, error: %v", zfsDevicePoolConfig.PoolName, node, zfsPoolErr)
			}
		}
	}

	logf.Log.Info("Verifying and deleting loop device on nodes if required")
	return k8stest.RemoveConfiguredLoopDeviceOnNodes(zfsDevicePoolConfig.NodeDeviceMap)
}

func (zfsDevicePoolConfig *ZfsNodesDevicePoolConfig) ConfigureZfsNodesWithDeviceAndPool() error {
	var err error
	if zfsDevicePoolConfig.PoolName == "" {
		return fmt.Errorf("zfs pool name not found")
	}
	logf.Log.Info("Verifying and creating loop device on nodes if required")
	zfsDevicePoolConfig.NodeDeviceMap, err = k8stest.ConfigureLoopDeviceOnNodes(zfsDevicePoolConfig.NodeDeviceMap)
	if err != nil {
		return fmt.Errorf("failed to create loop device, error: %v", err)
	}

	logf.Log.Info("Creating ZFS pool on nodes")
	for node, device := range zfsDevicePoolConfig.NodeDeviceMap {
		logf.Log.Info("Creating zfs pool", "node", node, "device path", device.DiskPath)
		if device.DiskPath == "" {
			return fmt.Errorf("device path not found for node %s", node)
		} else {
			nodeIp, err := k8stest.GetNodeIPAddress(node)
			if err != nil {
				return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			_, zfsPoolErr := e2e_agent.ZfsCreatePool(*nodeIp, device.DiskPath, zfsDevicePoolConfig.PoolName)
			if err != nil {
				return fmt.Errorf("failed to create zfs pool on node %s, error: %v", node, zfsPoolErr)
			}
		}
	}
	return nil
}
