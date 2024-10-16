package lvm

import (
	"fmt"
	"strings"

	"github.com/openebs/openebs-e2e/common"
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
	lvmNode := make([]string, 0, readyCount)
	for _, item := range listLvmDaemonSetPodList.Items {
		lvmNode = append(lvmNode, item.Spec.NodeName)
	}
	logf.Log.Info("LVM", "nodes", lvmNode)
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
			if lvmPvErr != nil {
				return fmt.Errorf("failed to create lvm pv on node %s with disk path %s, error: %v", node, device.DiskPath, lvmPvErr)
			}
			logf.Log.Info("Creating lvm vg", "node", node, "device path", device.DiskPath)
			_, lvmVgErr := e2e_agent.LvmCreateVg(*nodeIp, device.DiskPath, lvmNodesDevicePvVgConfig.VgName)
			if lvmVgErr != nil {
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

func SetupLvmNodes(vgName string, size int64) (LvmNodesDevicePvVgConfig, error) {
	var lvmNodeConfig LvmNodesDevicePvVgConfig
	workerNodes, err := ListLvmNode(common.NSOpenEBS())
	if err != nil {
		return lvmNodeConfig, fmt.Errorf("failed to list lvm worker nodes, error: %v", err)
	}
	if len(workerNodes) == 0 {
		return lvmNodeConfig, fmt.Errorf("lvm worker nodes not found")
	}
	var imgDir string
	// Cluster setup on github runner will have worker name starts with kind- like kind-worker, kind-worker2
	// create disk image file at /mnt on host which will /host/host/mnt in e2e-agent container because
	// on github runner one device is mounted
	// and if worker name does not start kind- then create disk image file at /tmp directory of host
	if strings.Contains(workerNodes[0], "kind-") {
		imgDir = "/host/host/mnt"
	} else {
		imgDir = "/mnt"
	}

	loopDevice := e2e_agent.LoopDevice{
		Size:   size,
		ImgDir: imgDir,
	}

	lvmNodeConfig = LvmNodesDevicePvVgConfig{
		VgName:        vgName,
		NodeDeviceMap: make(map[string]e2e_agent.LoopDevice), // Properly initialize the map
	}
	for _, node := range workerNodes {
		lvmNodeConfig.NodeDeviceMap[node] = loopDevice
	}

	logf.Log.Info("setup node with loop device, lvm pv and vg", "lvm node config", lvmNodeConfig)
	err = lvmNodeConfig.ConfigureLvmNodesWithDeviceAndVg()
	return lvmNodeConfig, err
}

// EnableLvmThinPoolAutoExpansion enable auto extending of the Thin Pool (Configure Over-Provisioning protection)
func EnableLvmThinPoolAutoExpansion(thinPoolAutoExtendThreshold, thinPoolAutoExtendPercent int) error {

	workerNodes, err := ListLvmNode(common.NSOpenEBS())
	if err != nil {
		return fmt.Errorf("failed to list lvm worker nodes, error: %v", err)
	}
	if len(workerNodes) == 0 {
		return fmt.Errorf("lvm worker nodes not found")
	}
	/*
		Editing the settings in the /etc/lvm/lvm.conf can allow auto growth of the thin pool when required.
		By default, the threshold is 100% which means that the pool will not grow.
		If we set this to, 75%, the Thin Pool will autoextend when the pool is 75% full.
		It will increase by the default percentage of 20% if the value is not changed.
		We can see these settings using the command grep against the file.
		$ grep -E ‘^\s*thin_pool_auto’ /etc/lvm/lvm.conf
		thin_pool_autoextend_threshold = 100
		thin_pool_autoextend_percent = 20
	*/

	for _, node := range workerNodes {
		nodeIp, err := k8stest.GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s IP, error: %v", node, err)
		}

		out, err := e2e_agent.LvmThinPoolAutoExtendThreshold(*nodeIp, thinPoolAutoExtendThreshold)
		if err != nil {
			return fmt.Errorf("failed to set up thin_pool_autoextend_threshold value %d on node %s,output: %s error: %v",
				thinPoolAutoExtendThreshold, node, out, err)
		}

		out, err = e2e_agent.LvmThinPoolAutoExtendPercent(*nodeIp, thinPoolAutoExtendPercent)
		if err != nil {
			return fmt.Errorf("failed to set up thin_pool_autoextend_percent value %d on node %s,output: %s error: %v",
				thinPoolAutoExtendPercent, node, out, err)
		}

	}

	return nil
}
