package v1

import (
	"encoding/json"
	"fmt"

	"github.com/openebs/openebs-e2e/common"
)

// Interface methods for CPv1
// GetZFSVolume returns a ZFSVolume CR by name
func (cp CPv1) GetZFSVolume(name string) (*common.ZFSVolume, error) {
	return getZfsCpVolume(name)
}

// ListZFSVolumes returns all ZFSVolume CRs in the mayastor namespace
func (cp CPv1) ListZFSVolumes() ([]common.ZFSVolume, error) {
	return listZfsCpVolumes()
}

// GetZFSVolumeStatus returns the status.state of a ZFSVolume by name
func (cp CPv1) GetZFSVolumeStatus(name string) (string, error) {
	zfsVol, err := getZfsCpVolume(name)
	if err != nil {
		return "", err
	}
	return zfsVol.Status.State, nil
}

// Plugin-based functions for ZFS volumes
// getZfsCpVolume gets a ZFS volume using plugin command
func getZfsCpVolume(name string) (*common.ZFSVolume, error) {
	var jsonInput []byte
	var err error
	cmd := GetOpenebsPluginCmd("-n", common.NSMayastor(), "-ojson", "get", "zfsvolume", name)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
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

// listZfsCpVolumes returns all ZFS volumes using plugin command
func listZfsCpVolumes() ([]common.ZFSVolume, error) {
	var jsonInput []byte
	var err error
	cmd := GetOpenebsPluginCmd("-n", common.NSMayastor(), "-ojson", "get", "zfsvolumes")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.ZFSVolume
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		return []common.ZFSVolume{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}
