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

type gcp struct {
}

func New() types.Platform {
	return &gcp{}
}

func getGcpNodeZone(node string) (string, error) {
	logf.Log.Info("GCP: get node zone")
	command := fmt.Sprintf("gcloud compute instances list %s --format 'csv[no-heading](zone)'", node)
	logf.Log.Info("GCP", "command", command)
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		logf.Log.Info("GCP: Failed to get node zone", "node", node, "error", err.Error())
		return "", fmt.Errorf("GCP: Failed to get node %s zone, error:  %v", node, err)
	}
	nodeZone := strings.TrimSuffix(string(output), "\n")
	logf.Log.Info("GCP:", "node", node, "zone", nodeZone)
	return nodeZone, nil
}

func (h *gcp) PowerOffNode(node string) error {
	logf.Log.Info("GCP: Power off", "node", node)
	nodeZone, err := getGcpNodeZone(node)
	if err != nil {
		return fmt.Errorf("failed to get GCP node %s zone, error: %v", node, err)
	} else if nodeZone == "" {
		return fmt.Errorf("GCP node %s zone not found", node)
	}
	command := fmt.Sprintf("gcloud compute instances stop %s --zone %s", node, nodeZone)
	logf.Log.Info("GCP: Power off", "command", command)
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("GCP: Power off failed", "error", err.Error())
		logf.Log.Info("GCP: Power off command output", "output", string(output))
		return fmt.Errorf("GCP power off failed: %v", err)
	}
	return err
}

func (h *gcp) PowerOnNode(node string) error {
	logf.Log.Info("GCP: Power on", "node", node)
	nodeZone, err := getGcpNodeZone(node)
	if err != nil {
		return fmt.Errorf("failed to get GCP node %s zone, error: %v", node, err)
	} else if nodeZone == "" {
		return fmt.Errorf("GCP node %s zone not found", node)
	}
	command := fmt.Sprintf("gcloud compute instances start %s --zone %s", node, nodeZone)
	logf.Log.Info("GCP: Power on", "command", command)
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("GCP: Power on failed", "error", err.Error())
		logf.Log.Info("GCP: Power on command output", "output", string(output))
		return fmt.Errorf("GCP power on failed: %v", err)
	}
	return err
}

func (h *gcp) RebootNode(node string) error {
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

	logf.Log.Info("GCP: Reboot", "node", node)
	nodeZone, err := getGcpNodeZone(node)
	if err != nil {
		return fmt.Errorf("failed to get GCP node %s zone, error: %v", node, err)
	} else if nodeZone == "" {
		return fmt.Errorf("GCP node %s zone not found ", node)
	}
	command := fmt.Sprintf("gcloud compute instances reset %s --zone %s", node, nodeZone)
	logf.Log.Info("GCP: Reboot node", "command", command)
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("GCP: Reboot failed", "error", err.Error())
		logf.Log.Info("GCP: Reboot command output", "output", string(output))
		return fmt.Errorf("GCP: reboot node failed: %v", err)
	}
	return err
}

func (h *gcp) DetachVolume(volName string, node string) error {
	logf.Log.Info("GCP: Detach volume", "volume", volName, "node", node)
	nodeZone, err := getGcpNodeZone(node)
	if err != nil {
		return fmt.Errorf("failed to get GCP node %s zone, error: %v", node, err)
	} else if nodeZone == "" {
		return fmt.Errorf("GCP node %s zone not found ", node)
	}
	command := fmt.Sprintf("gcloud compute instances detach-disk %s --disk %s --zone %s", node, volName, nodeZone)
	logf.Log.Info("GCP: Detach volume", "command", command)
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("GCP: Detach volume failed", "error", err.Error())
		logf.Log.Info("GCP: Detach volume command output", "output", string(output))
		return fmt.Errorf("GCP Detach volume failed: %v", err)
	}
	return err
}

func (h *gcp) AttachVolume(volName, node string) error {
	deviceName := "mayastor-disk"
	logf.Log.Info("GCP: Attach volume", "volume", volName, "node", node, "device-name", deviceName)
	nodeZone, err := getGcpNodeZone(node)
	if err != nil {
		return fmt.Errorf("failed to get GCP node %s zone, error: %v", node, err)
	} else if nodeZone == "" {
		return fmt.Errorf("GCP node %s zone not found", node)
	}
	command := fmt.Sprintf("gcloud compute instances attach-disk %s --disk %s --device-name %s --zone %s", node, volName, deviceName, nodeZone)
	logf.Log.Info("GCP: Attach volume", "command", command)
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("GCP: Attach volume failed", "error", err.Error())
		logf.Log.Info("GCP: Attach volume command output", "output", string(output))
		return fmt.Errorf("GCP Attach volume failed: %v", err)
	}
	return err
}

func (h *gcp) GetNodeStatus(node string) (string, error) {
	logf.Log.Info("Get status", "node", node)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("gcloud compute instances list | grep %s", node))
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	logf.Log.Info("Get status", "node", node, "output", string(stdout))
	if strings.Contains(string(stdout), "RUNNING") {
		return "running", nil
	}
	return "off", nil
}
