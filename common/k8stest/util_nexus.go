package k8stest

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/controlplane"
)

// GetNexusNode return nexus node given the volume uuid
func GetNexusNode(vol_uuid string) (string, error) {
	return controlplane.GetMsvTargetNode(vol_uuid)
}

// GetNexusNodeIp return nexus node IP given the volume uuid
func GetNexusNodeIp(vol_uuid string) (string, error) {
	msv, err := GetMSV(vol_uuid)
	if err != nil {
		return "", fmt.Errorf("failed to get msv %s, error %v", vol_uuid, err)
	}

	nexusNodeIp, err := GetNodeIPAddress(msv.State.Target.Node)
	if err != nil {
		return "", fmt.Errorf("failed to get node %s ip address, error %v", msv.State.Target.Node, err)
	}
	return *nexusNodeIp, err
}

// GetNexusUuid return the nexus uuid given the volume uuid
func GetNexusUuid(vol_uuid string) (string, error) {
	return controlplane.GetMsvTargetUuid(vol_uuid)
}
