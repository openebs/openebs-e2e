package client

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/platform/types"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type maas struct {
}

func New() types.Platform {
	return &maas{}
}

func (h *maas) PowerOffNode(node string) error {
	logf.Log.Info("flush container engine cache before powering off", "node", node)

	nodeIPAddr, err := k8stest.GetNodeIPAddress(node)

	if err != nil {
		return fmt.Errorf("failed to get node %s ip address %v", node, err)
	}
	if len(*nodeIPAddr) == 0 {
		return fmt.Errorf("node %s IP address is empty", node)
	}

	_, err = e2e_agent.FlushDiskWriteCache(*nodeIPAddr)
	if err != nil {
		return fmt.Errorf("failed to flush container engine cache on node %s with error %v", node, err)
	}

	logf.Log.Info("Power off", "node", node)
	nodeIp, err := k8stest.GetNodeIPAddress(node)
	if err != nil {
		return fmt.Errorf("failed to get node ip address for node %s, error %v", node, err)
	}
	systemId, err := getMaasMachineId(*nodeIp)
	if err != nil {
		return err
	}
	err = maasMachinePowerOnOff(powerOff, systemId)
	return err
}

func (h *maas) PowerOnNode(node string) error {
	logf.Log.Info("Power on", "node", node)
	nodeIp, err := k8stest.GetNodeIPAddress(node)
	if err != nil {
		return fmt.Errorf("failed to get node ip address for node %s, error %v", node, err)
	}
	systemId, err := getMaasMachineId(*nodeIp)
	if err != nil {
		return err
	}
	err = maasMachinePowerOnOff(powerOn, systemId)
	return err
}

func (h *maas) RebootNode(node string) error {
	logf.Log.Info("flush container engine cache before powering off", "node", node)

	nodeIPAddr, err := k8stest.GetNodeIPAddress(node)

	if err != nil {
		return fmt.Errorf("failed to get node %s ip address %v", node, err)
	}
	if len(*nodeIPAddr) == 0 {
		return fmt.Errorf("node %s IP address is empty", node)
	}

	_, err = e2e_agent.FlushDiskWriteCache(*nodeIPAddr)
	if err != nil {
		return fmt.Errorf("failed to flush container engine cache on node %s with error %v", node, err)
	}

	logf.Log.Info("Reboot", "node", node)
	nodeIp, err := k8stest.GetNodeIPAddress(node)
	if err != nil {
		return fmt.Errorf("failed to get node ip address for node %s, error %v", node, err)
	}
	systemId, err := getMaasMachineId(*nodeIp)
	if err != nil {
		return err
	}
	err = maasMachinePowerOnOff(powerOff, systemId)
	if err != nil {
		return err
	}
	const sleepTime = 10
	// Wait (with timeout) for maas node to be powerd off
	logf.Log.Info("Waiting for maas node to powered off", "timeoutSecs", timeoutSecs)
	for ix := 1; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		status, err := getMaasMachineStatus(node)
		if err != nil {
			return err
		} else if status == powerStateOff {
			break
		}
	}
	err = maasMachinePowerOnOff(powerOn, systemId)
	return err
}

// FIXME implement Maas API to detach disk
func (h *maas) DetachVolume(volName string, node string) error {
	logf.Log.Info("Detach Volume ", "volName", volName)
	panic("API to detach dsik from maas node not implemented")
}

// FIXME implement Maas API to attach disk
func (h *maas) AttachVolume(volName, node string) error {
	logf.Log.Info("Attach Volume to node", "volName", volName, "node", node)
	panic("API to attach dsik to maas node not implemented")
}

func (h *maas) GetNodeStatus(node string) (string, error) {
	logf.Log.Info("Get status", "node", node)
	nodeIp, err := k8stest.GetNodeIPAddress(node)
	if err != nil {
		return "", fmt.Errorf("failed to get node ip address for node %s, error %v", node, err)
	}
	status, err := getMaasMachineStatus(*nodeIp)
	if err != nil {
		return "", err
	}
	if status == powerStateOn {
		return "running", nil
	}
	return "off", nil
}
