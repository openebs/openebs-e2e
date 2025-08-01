package k8stest

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
)

const (
	// ZFSVolumesResource is the Kubernetes resource type for ZFS volumes
	ZFSVolumesResource = "zfsvolumes.zfs.openebs.io"
)

// GetZFSVolume returns a ZFSVolume CR by name
func GetZFSVolume(name string) (*common.ZFSVolume, error) {
	return controlplane.GetZFSVolume(name)
}

// ListZFSVolumes returns all ZFSVolume CRs in the mayastor namespace
func ListZFSVolumes() ([]common.ZFSVolume, error) {
	return controlplane.ListZFSVolumes()
}

// GetZFSVolumeStatus returns the status.state of a ZFSVolume by name
func GetZFSVolumeStatus(name string) (string, error) {
	return controlplane.GetZFSVolumeStatus(name)
}

// GetZfsVolumeCr gets a ZFS volume using direct kubectl command
func GetZfsVolumeCr(name, namespace string) (*common.ZFSVolume, error) {
	var jsonInput []byte
	var err error
	cmd := exec.Command("kubectl", "get", ZFSVolumesResource, name, "-n", namespace, "-o", "json")
	jsonInput, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var response common.ZFSVolume
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		msg := string(jsonInput)
		return nil, fmt.Errorf("%s", msg)
	}
	return &response, nil
}

// ListZFSVolumeCr returns all ZFS volumes using direct kubectl
func ListZFSVolumeCr(namespace string) ([]string, error) {
	cmd := exec.Command("kubectl", "get", ZFSVolumesResource, "-n", namespace, "-o", "jsonpath={.items[*].metadata.name}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list ZFS volumes: %v", err)
	}

	if strings.TrimSpace(string(output)) == "" {
		return []string{}, nil
	}

	names := strings.Split(strings.TrimSpace(string(output)), " ")
	return names, nil
}

// GetZFSVolumeCrStatus returns the status of a ZFS volume using direct kubectl
func GetZFSVolumeCrStatus(name, namespace string) (string, error) {
	cmd := exec.Command("kubectl", "get", ZFSVolumesResource, name, "-n", namespace, "-o", "jsonpath={.status.state}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get ZFS volume status for %s: %v", name, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// DeleteZFSVolumeCr deletes a ZFS volume using direct kubectl
func DeleteZFSVolumeCr(name, namespace string) error {
	cmd := exec.Command("kubectl", "delete", ZFSVolumesResource, name, "-n", namespace)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete ZFS volume %s: %v, output: %s", name, err, string(output))
	}
	return nil
}
