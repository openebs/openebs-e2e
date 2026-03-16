package client

import (
	"fmt"
	"os/exec"
	"regexp"
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

func (h *hcloud) DetachVolumeFromNode(node string) error {
	volName := fmt.Sprintf("mayastor-%s", node)

	logf.Log.Info("Detaching volume from node",
		"node", node,
		"volume", volName,
	)

	return h.DetachVolume(volName, node)
}

func (h *hcloud) AttachVolume(volName, node string) error {
	logf.Log.Info("Attach Volume to node", "volName", volName, "node", node)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud volume attach %s --server %s", volName, node))
	_, err := cmd.Output()
	return err
}

func (h *hcloud) AttachVolumeToNode(node string) error {
	volName := fmt.Sprintf("mayastor-%s", node)

	logf.Log.Info("Attaching volume to node",
		"node", node,
		"volume", volName,
	)

	return h.AttachVolume(volName, node)
}

func (h *hcloud) ResizeVolume(volName string, newSizeGB int) error {
	logf.Log.Info("Resize Volume", "volName", volName, "newSizeGB", newSizeGB)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("hcloud volume resize --size %d %s", newSizeGB, volName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("HCloud: Resize volume failed", "error", err.Error())
		logf.Log.Info("HCloud: Resize volume output", "output", string(output))
		return fmt.Errorf("hcloud resize volume failed: %v", err)
	}
	return nil
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

// ExtractVolumeIdFromDevicePath parses Hetzner by-id device path and extracts the numeric volume id.
// Example inputs:
//   - /dev/disk/by-id/scsi-0HC_Volume_12345678
//   - /dev/disk/by-id/scsi-0HC-VOLUME-12345678
//   - /dev/disk/by-id/scsi-0HC_Volume_12345678-part1
func (h *hcloud) ExtractVolumeIdFromDevicePath(dev string) (string, error) {
	re := regexp.MustCompile(`(?i)HC[_-]?VOLUME_(\d+)`)
	m := re.FindStringSubmatch(dev)
	if len(m) == 2 {
		return m[1], nil
	}
	// Fallback: try to find trailing numeric id possibly followed by -partX
	tail := regexp.MustCompile(`(\d+)(?:-?part\d+)?$`).FindStringSubmatch(strings.ToLower(dev))
	if len(tail) == 2 {
		return tail[1], nil
	}
	return "", fmt.Errorf("device path %s does not look like a Hetzner volume by-id path", dev)
}
