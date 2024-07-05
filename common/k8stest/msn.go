package k8stest

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
)

// GetMSN Get pointer to a mayastor volume custom resource
// returns nil and no error if the msn is in pending state.
func GetMSN(nodeName string) (*common.MayastorNode, error) {
	return controlplane.GetMSN(nodeName)
}

func ListMsns() ([]common.MayastorNode, error) {
	return controlplane.ListMsns()
}

func GetMsNodeStatus(nodeName string) (string, error) {
	return controlplane.GetMsNodeStatus(nodeName)
}
