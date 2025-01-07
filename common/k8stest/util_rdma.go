package k8stest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	agent "github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	RdmaDeviceName = "rxe0"
	TcpProtocol    = "tcp"
	RdmaProtocol   = "rdma"
)

type RdmaDeviceNetworkInterface struct {
	IfIndex       int    `json:"ifindex"`
	IfName        string `json:"ifname"`
	Port          int    `json:"port"`
	State         string `json:"state"`
	PhysicalState string `json:"physical_state"`
	NetDev        string `json:"netdev"`
	NetDevIndex   int    `json:"netdev_index"`
}

type PortInfo struct {
	PCI     string `json:"pci"`
	Type    string `json:"type"`
	Netdev  string `json:"netdev"`
	Flavour string `json:"flavour"`
	Port    int    `json:"port"`
}

type PortMap struct {
	Port map[string]PortInfo `json:"port"`
}

func ListRdmaDevice(node string) ([]RdmaDeviceNetworkInterface, error) {
	var rdmaDeiceList []RdmaDeviceNetworkInterface
	nodeIp, err := GetNodeIPAddress(node)
	if err != nil {
		return rdmaDeiceList, fmt.Errorf("failed to get node %s ip, error: %v", node, err)
	}

	rdmaDevice, err := agent.ListRdmaDevice(*nodeIp)
	if err != nil {
		return rdmaDeiceList, fmt.Errorf("failed to list RDMA device on node %s , error: %v", node, err)
	}
	if rdmaDevice == "" {
		logf.Log.Info("RDMA device list failed with empty string", "output", rdmaDevice)
		return rdmaDeiceList, fmt.Errorf("failed to list RDMA device on node %s", node)
	}
	output := trimForJson(rdmaDevice)
	if err = json.Unmarshal([]byte(output), &rdmaDeiceList); err != nil {
		logf.Log.Info("Failed to unmarshal rdma list", "output", output)
		return rdmaDeiceList, fmt.Errorf("failed to unmarshal rdma list on node %s , output: %s,error: %v", node, output, err)
	}
	logf.Log.Info("RDMA device", "node", node, "list", rdmaDeiceList)
	return rdmaDeiceList, nil
}

func CreateRdmaDeviceOnNode(node string) error {
	rdmaDeviceList, err := ListRdmaDevice(node)
	if err != nil {
		return err
	}
	if len(rdmaDeviceList) == 0 {
		logf.Log.Info("RDMA device not found", "node", node, "list", rdmaDeviceList)
		//create rdma device
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}

		// get interface name
		iface := e2e_config.GetConfig().NetworkInterface
		out, err := e2e_agent.CreateRdmaDevice(*nodeIp, RdmaDeviceName, iface)
		if err != nil {
			return err
		}
		logf.Log.Info("Device created", "node", node, "output", out, "interface", iface)

	}
	rdmaDeviceList, err = ListRdmaDevice(node)
	if err != nil {
		return err
	}
	logf.Log.Info("RDMA device", "node", node, "list", rdmaDeviceList)
	return nil
}

func CreateRdmaDeviceOnAllIoEngineNodes() error {
	mayastorNodes, err := ListIOEngineNodes()
	if err != nil {
		return err
	}
	for _, node := range mayastorNodes.Items {
		logf.Log.Info("Create rdma device", "Node", node.Name)
		err := CreateRdmaDeviceOnNode(node.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetVolumeProtocol(volUuid string) (string, error) {
	deviceUri, err := GetMsvDeviceUri(volUuid)
	if err != nil {
		return "", err
	}
	logf.Log.Info("Volume URI", "volume", volUuid, "URI", deviceUri)
	// deviceUri: nvmf://<some-random-string>
	// Parse the device URI
	u, err := url.Parse(deviceUri)
	if err != nil {
		return "", fmt.Errorf("error parsing URI: %s, error: %v", deviceUri, err)
	}
	return u.Scheme, nil
}

// IsVolumeAccessibleOverRdma return true if volume device uri scheme contains rdma
// if volume is accessible over rdma then device uri will be like nvmf+rdma://<some-random-string>
func IsVolumeAccessibleOverRdma(volUuid string) (bool, error) {
	protocol, err := GetVolumeProtocol(volUuid)
	if err != nil {
		return false, err
	}
	logf.Log.Info("Volume ", "name", volUuid, "URI protocol", protocol)
	if strings.Contains(protocol, "rdma") {
		return true, nil
	}
	return false, nil
}

// IsVolumeAccessibleOverTcp return true if volume device uri scheme contains tcp and not rdma
// if volume is accessible over rdma then device uri will be like nvmf://<some-random-string>
func IsVolumeAccessibleOverTcp(volUuid string) (bool, error) {
	protocol, err := GetVolumeProtocol(volUuid)
	if err != nil {
		return false, err
	}
	logf.Log.Info("Volume", "name", volUuid, "URI protocol", protocol)
	if !strings.Contains(protocol, "rdma") && strings.Contains(protocol, "nvmf") {
		return true, nil
	}
	return false, nil
}

func RemoveRdmaDeviceOnNode(node string) error {
	rdmaDeviceList, err := ListRdmaDevice(node)
	if err != nil {
		return err
	}
	if len(rdmaDeviceList) != 0 {
		logf.Log.Info("RDMA device found", "node", node, "list", rdmaDeviceList)
		//create rdma device
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}
		var iface, out string
		for _, device := range rdmaDeviceList {
			if device.IfName == RdmaDeviceName {
				// get interface name
				iface = e2e_config.GetConfig().NetworkInterface
				out, err = e2e_agent.DeleteRdmaDevice(*nodeIp, RdmaDeviceName)
				if err != nil {
					return err
				}
			}
		}
		if iface == "" {
			return fmt.Errorf("rdma device %s not found", RdmaDeviceName)
		}

		logf.Log.Info("Device deleted", "node", node, "output", out, "interface", iface, "device", RdmaDeviceName)

	}
	return nil
}

func DisableRdmaOnNode(node string, networkInterface string) error {
	logf.Log.Info("Disable rdma from IO engine node", "name", node)

	// get dev link port wrt to interface
	rdmaDevName, err := GetDevLinkName(node, networkInterface)
	if err != nil {
		return err
	}

	logf.Log.Info("rdma dev", "Name", rdmaDevName, "node", node)

	if rdmaDevName != "" {
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}

		out, err := e2e_agent.DisableDevLink(*nodeIp, rdmaDevName)
		if err != nil {
			logf.Log.Info("failed to disable rdma dev link", "node", node, "dev link", rdmaDevName, "output", out)
			return err
		}
	} else {
		err := RemoveRdmaDeviceOnNode(node)
		if err != nil {
			logf.Log.Info("failed to remove rdma device", "node", node, "device", RdmaDeviceName)
			return err
		}
	}

	// Restart csi node pod on the node
	return RestartCsiNodePodOnNode(node, 240, 120)

}

func RemoveRdmaDeviceOnAllWorkerNodes() error {
	workerNodes, err := ListIOEngineNodes()
	if err != nil {
		return err
	}
	logf.Log.Info("Remove rdma from IO engine node")
	for _, node := range workerNodes.Items {
		logf.Log.Info("IO engine", "Node", node.Name)
		err := DisableRdmaOnNode(node.Name, e2e_config.GetConfig().NetworkInterface)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableRdmaDeviceOnAllWorkerNodes() error {
	workerNodes, err := ListIOEngineNodes()
	if err != nil {
		return err
	}
	logf.Log.Info("Enable rdma from IO engine node")
	for _, node := range workerNodes.Items {
		logf.Log.Info("IO engine", "Node", node.Name)
		err := EnableRdmaOnNode(node.Name, e2e_config.GetConfig().NetworkInterface)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableRdmaOnNode(node string, networkInterface string) error {
	logf.Log.Info("Enable rdma on IO engine node", "name", node)

	// get dev link port wrt to interface
	rdmaDevPortName, err := GetDevLinkName(node, networkInterface)
	if err != nil {
		return err
	}

	if rdmaDevPortName != "" {
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}

		out, err := e2e_agent.EnableDevLink(*nodeIp, rdmaDevPortName)
		if err != nil {
			logf.Log.Info("failed to enable rdma dev link", "node", node, "dev link", rdmaDevPortName, "output", out)
			return err
		}
	} else {
		err := CreateRdmaDeviceOnNode(node)
		if err != nil {
			logf.Log.Info("failed to create rdma device", "node", node, "device", RdmaDeviceName)
			return err
		}
	}

	return nil
}

func ListDevLink(node string) (PortMap, error) {
	var devLink PortMap
	nodeIp, err := GetNodeIPAddress(node)
	if err != nil {
		return devLink, fmt.Errorf("failed to get node %s ip, error: %v", node, err)
	}

	devLinkOut, err := agent.ListDevLink(*nodeIp)
	if err != nil {
		return devLink, fmt.Errorf("failed to list dev link on node %s , error: %v", node, err)
	}
	if devLinkOut == "" {
		logf.Log.Info("Dev kink list failed with empty string", "output", devLinkOut)
		return devLink, fmt.Errorf("failed to list dev link on node %s", node)
	}
	output := trimForJson(devLinkOut)
	if err = json.Unmarshal([]byte(output), &devLink); err != nil {
		logf.Log.Info("Failed to unmarshal dev link list", "output", output)
		return devLink, fmt.Errorf("failed to unmarshal dev link list on node %s , output: %s,error: %v", node, output, err)
	}
	logf.Log.Info("Dev link", "node", node, "list", devLink)
	return devLink, nil
}

func GetDevLinkName(node, iface string) (string, error) {
	devLinkList, err := ListDevLink(node)
	if err != nil {
		return "", err
	}
	for key, val := range devLinkList.Port {
		if val.Netdev == iface {
			// dev link port will be like pci/0000:3b:00.0/65535
			// dev link will be pci/0000:3b:00.0 by removing port
			return key[:strings.LastIndex(key, "/")], nil
		}
	}
	return "", nil
}

// DisableConfiguredDisabledRdmaDevicesOnAllMayastorNodes disable rdma devices which are configured in platform config
// on all mayastor nodes. If no rdma devices are configured then it will return without doing anything.
// In some cases, multiple rdma devices are configured in cluster, so it should be disabled all the devices on nodes
func DisableConfiguredDisabledRdmaDevicesOnAllMayastorNodes() error {
	ifaceList := e2e_config.GetConfig().DisabledRdmaDevices
	if len(ifaceList) == 0 {
		logf.Log.Info("No rdma devices configured which needs to be disabled")
		return nil
	}
	workerNodes, err := ListIOEngineNodes()
	if err != nil {
		return err
	}
	for _, node := range workerNodes.Items {
		logf.Log.Info("Disable rdma device", "Node", node.Name)
		for _, iface := range ifaceList {
			err := DisableRdmaOnNode(node.Name, iface)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// RestoreConfiguredDisabledRdmaDevicesOnAllMayastorNodes enable rdma devices which are configured in platform config
// on all mayastor nodes. If no rdma devices are configured then it will return without doing anything.
// In some cases, multiple rdma devices are configured in cluster were disabled, so it should be enabled all the devices on nodes
func RestoreConfiguredDisabledRdmaDevicesOnAllMayastorNodes() error {
	ifaceList := e2e_config.GetConfig().DisabledRdmaDevices
	if len(ifaceList) == 0 {
		logf.Log.Info("No rdma devices configured which needs to be restored")
		return nil
	}
	workerNodes, err := ListIOEngineNodes()
	if err != nil {
		return err
	}
	for _, node := range workerNodes.Items {
		logf.Log.Info("Enable rdma device", "Node", node.Name)
		for _, iface := range ifaceList {
			err := EnableRdmaOnNode(node.Name, iface)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// IsVolumeConnectedOverSpecifiedProtocol return true if volume is connected over specified protocol
// by checking nvme list subsystem live path protocol for specified volume
func IsVolumeConnectedOverSpecifiedProtocol(volUuid, initiatorNode, protocol string) (bool, error) {
	// get node IP
	nodeIp, err := GetNodeIPAddress(initiatorNode)
	if err != nil {
		return false, fmt.Errorf("failed to get node IP for node %s, error: %v", initiatorNode, err)
	} else if *nodeIp == "" {
		return false, fmt.Errorf("node IP for node %s not found", initiatorNode)
	}
	logf.Log.Info("Node IP", "node", initiatorNode, "IP", *nodeIp)
	logf.Log.Info("volume connection entry in nvme sub system", "expected protocol", protocol)
	nvmeListSubsysProtocol, err := GetNvmeListSubsystemVolumeEntryLivePathProtocol(*nodeIp, volUuid)
	if err != nil {
		return false, fmt.Errorf("failed to get nvme list subsystem live path protocol for volume %s on node %s, error: %v", volUuid, initiatorNode, err)
	} else if nvmeListSubsysProtocol == "" {
		return false, fmt.Errorf("nvme list subsystem live path protocol for volume %s on node %s not found", volUuid, initiatorNode)
	}
	logf.Log.Info("Volume connection protocol", "volume", volUuid, "protocol", nvmeListSubsysProtocol)
	return true, nil
}
