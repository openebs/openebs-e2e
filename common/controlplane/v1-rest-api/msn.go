package v1_rest_api

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/generated/openapi"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func oaNodeToMsn(oaNode *openapi.Node) common.MayastorNode {
	nodeSpec := oaNode.GetSpec()
	nodeState := oaNode.GetState()
	msn := common.MayastorNode{
		Name: oaNode.GetId(),
		Spec: common.MayastorNodeSpec{
			GrpcEndpoint: nodeSpec.GrpcEndpoint,
			ID:           nodeSpec.Id,
		},
		State: common.MayastorNodeState{
			GrpcEndpoint: nodeState.GrpcEndpoint,
			ID:           nodeState.Id,
			Status:       string(nodeState.GetStatus()),
		},
	}

	return msn
}

func (cp CPv1RestApi) GetMSN(nodeName string) (*common.MayastorNode, error) {
	oaNode, err := cp.oa.getNode(nodeName)

	if err != nil {
		logf.Log.Info("getNode failed", "err", err)
		return nil, err
	}

	msn := oaNodeToMsn(&oaNode)
	return &msn, err
}

func (cp CPv1RestApi) ListMsns() ([]common.MayastorNode, error) {
	oaNodes, err := cp.oa.getNodes()
	if err != nil {
		logf.Log.Info("getNodes failed", "err", err)
		return nil, err
	}
	var msns []common.MayastorNode
	for _, oaNode := range oaNodes {
		msns = append(msns, oaNodeToMsn(&oaNode))
	}
	return msns, err
}

func (cp CPv1RestApi) GetMsNodeStatus(nodeName string) (string, error) {
	rNode, err := cp.oa.getNode(nodeName)

	if err == nil {
		return string(rNode.State.GetStatus()), err
	}
	return "", err
}

func (cp CPv1RestApi) UpdateNodeLabel(nodeName string, labelKey, labelValue string) error {
	panic(fmt.Errorf("not implemented REST api"))
}
