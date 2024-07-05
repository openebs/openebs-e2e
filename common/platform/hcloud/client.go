package client

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/platform/types"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type hcloud struct {
}

func New() types.Platform {
	return &hcloud{}
}

func (h *hcloud) PowerOffNode(node string) error {
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
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud server poweroff %s", node))
	_, err = cmd.Output()
	return err
}

func (h *hcloud) PowerOnNode(node string) error {
	logf.Log.Info("Power on", "node", node)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud server poweron %s", node))
	_, err := cmd.Output()
	return err
}

func (h *hcloud) RebootNode(node string) error {
	logf.Log.Info("flush container engine cache before rebooting", "node", node)

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
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud server reboot %s", node))
	_, err = cmd.Output()
	return err
}

func (h *hcloud) DetachVolume(volName string, node string) error {
	logf.Log.Info("Detach Volume ", "volName", volName)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud volume detach %s", volName))
	_, err := cmd.Output()
	return err
}

func (h *hcloud) AttachVolume(volName, node string) error {
	logf.Log.Info("Attach Volume to node", "volName", volName, "node", node)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud volume attach %s --server %s", volName, node))
	_, err := cmd.Output()
	return err
}

func (h *hcloud) GetNodeStatus(node string) (string, error) {
	logf.Log.Info("Get status", "node", node)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("hcloud  server list | grep %s", node))
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(stdout), "running") {
		return "running", nil
	}
	return "off", nil
}
