package v1_rest_api

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
)

// ZFSVolume interface stubs for REST API (not implemented)
func (cp CPv1RestApi) GetZFSVolume(name string) (*common.ZFSVolume, error) {
	return nil, fmt.Errorf("GetZFSVolume not implemented for REST API")
}

// ListZFSVolumes returns all ZFSVolume CRs
func (cp CPv1RestApi) ListZFSVolumes() ([]common.ZFSVolume, error) {
	return nil, fmt.Errorf("ListZFSVolumes not implemented for REST API")
}

// GetZFSVolumeStatus returns the status.state of a ZFSVolume by name
func (cp CPv1RestApi) GetZFSVolumeStatus(name string) (string, error) {
	return "", fmt.Errorf("GetZFSVolumeStatus not implemented for REST API")
}
