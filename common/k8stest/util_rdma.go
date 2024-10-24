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

const RdmaDeviceName = "rxe0"

type RdmaDeviceNetworkInterface struct {
	IfIndex       int    `json:"ifindex"`
	IfName        string `json:"ifname"`
	Port          int    `json:"port"`
	State         string `json:"state"`
	PhysicalState string `json:"physical_state"`
	NetDev        string `json:"netdev"`
	NetDevIndex   int    `json:"netdev_index"`
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

func CreateRdmaDeviceOnAllWorkerNodes() error {
	workerNodes, err := ListWorkerNode()
	if err != nil {
		return err
	}
	logf.Log.Info("Worker", "Nodes", workerNodes)
	for _, node := range workerNodes {
		err := CreateRdmaDeviceOnNode(node.NodeName)
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
	// deviceUri: nvmf+tcp://<some-random-string>
	// Parse the device URI
	u, err := url.Parse(deviceUri)
	if err != nil {
		return "", fmt.Errorf("error parsing URI: %s, error: %v", deviceUri, err)
	}
	return u.Scheme, nil
}

// IsVolumeAccessibleOverRdma return true if volume device uri scheme contains rdma
// if volume is accessible over rdma then device uri will be like nvmf+tcp+rdma://<some-random-string>
func IsVolumeAccessibleOverRdma(volUuid string) (bool, error) {
	protocol, err := GetVolumeProtocol(volUuid)
	if err != nil {
		return false, err
	}
	if strings.Contains(protocol, "rdma") {
		return true, nil
	}
	return false, nil
}

// IsVolumeAccessibleOverTcp return true if volume device uri scheme contains tcp and not rdma
// if volume is accessible over rdma then device uri will be like nvmf+tcp://<some-random-string>
func IsVolumeAccessibleOverTcp(volUuid string) (bool, error) {
	protocol, err := GetVolumeProtocol(volUuid)
	if err != nil {
		return false, err
	}
	if !strings.Contains(protocol, "rdma") && strings.Contains(protocol, "tcp") {
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

func DisableRdmaOnNode(node string) error {
	//disable rdma on the io-engine node
	platformName := e2e_config.GetConfig().Platform.Name
	if platformName == "Maas" {
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}
		// get interface name
		iface := e2e_config.GetConfig().NetworkInterface
		out, err := e2e_agent.DisableNetworkInterface(*nodeIp, iface)
		if err != nil {
			logf.Log.Info("failed to disable network interface", "platform", platformName, "node", node, "iface", iface, "output", out)
			return err
		}
	} else if platformName == "Hetzner" {
		err := RemoveRdmaDeviceOnNode(node)
		if err != nil {
			logf.Log.Info("failed to remove rdma device", "platform", platformName, "node", node, "device", RdmaDeviceName)
			return err
		}
	} else {
		return fmt.Errorf("unsupported platform")
	}

	return nil
}

func RemoveRdmaDeviceOnAllWorkerNodes() error {
	workerNodes, err := ListWorkerNode()
	if err != nil {
		return err
	}
	logf.Log.Info("Worker", "Nodes", workerNodes)
	for _, node := range workerNodes {
		err := DisableRdmaOnNode(node.NodeName)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableRdmaOnNode(node string) error {
	//enable rdma on the io-engine node
	platformName := e2e_config.GetConfig().Platform.Name
	if platformName == "Maas" {
		nodeIp, err := GetNodeIPAddress(node)
		if err != nil {
			return fmt.Errorf("failed to get node %s ip, error: %v", node, err)
		}
		// get interface name
		iface := e2e_config.GetConfig().NetworkInterface
		out, err := e2e_agent.EnableNetworkInterface(*nodeIp, iface)
		if err != nil {
			logf.Log.Info("failed to enable network interface", "platform", platformName, "node", node, "iface", iface, "output", out)
			return err
		}
	} else if platformName == "Hetzner" {
		err := CreateRdmaDeviceOnNode(node)
		if err != nil {
			logf.Log.Info("failed to create rdma device", "platform", platformName, "node", node, "device", RdmaDeviceName)
			return err
		}
	} else {
		return fmt.Errorf("unsupported platform")
	}

	return nil
}
