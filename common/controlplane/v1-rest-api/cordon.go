package v1_rest_api

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	cpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
)

func (cp CPv1RestApi) CordonNode(nodeName string, cordonLabel string) error {
	// #TODO implement REST api for node cordon
	panic(fmt.Errorf("not implemented REST api for node cordon"))
}

func (cp CPv1RestApi) GetCordonNodeLabels(nodeName string) ([]string, error) {
	// #TODO implement REST api for node cordon labels
	panic(fmt.Errorf("not implemented REST api for getting node cordon labels"))
}

func (cp CPv1RestApi) UnCordonNode(nodeName string, cordonLabel string) error {
	// #TODO implement REST api for node uncordon
	panic(fmt.Errorf("not implemented REST api for node uncordon"))
}

// pool cordon
func (cp CPv1RestApi) CordonPool(poolID string, constraints ...common.PoolCordonConstraint) error {
	// #TODO implement REST api for pool cordon
	panic(fmt.Errorf("not implemented REST api for pool cordon"))
}

func (cp CPv1RestApi) UnCordonPool(poolID string, constraints ...common.PoolCordonConstraint) error {
	// #TODO implement REST api for pool uncordon
	panic(fmt.Errorf("not implemented REST api for pool uncordon"))
}

func (cp CPv1RestApi) GetPoolCordonStatus(poolID string) (*cpV1.PoolCordonStatus, error) {
	// #TODO implement REST api for getting pool cordon status
	panic(fmt.Errorf("not implemented REST api for getting pool cordon status"))
}
