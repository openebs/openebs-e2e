package k8stest

import (
	"encoding/json"
	"fmt"

	agent "github.com/openebs/openebs-e2e/common/e2e_agent"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
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
