package v1

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/openebs/openebs-e2e/common"
)

func (cp CPv1) GetSnapshots() ([]common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshots")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.SnapshotSchema
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshots)", "string", string(jsonInput))
		return []common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func (cp CPv1) GetSnapshot(snapshotId string) (common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	var response []common.SnapshotSchema
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshots", "--snapshot", snapshotId)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return common.SnapshotSchema{}, err
	}

	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshot by snapshot id)", "string", string(jsonInput))
		return common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	if len(response) == 0 {
		return common.SnapshotSchema{}, fmt.Errorf("snapshot not found, snapshot: %s", snapshotId)
	} else if len(response) > 1 {
		return common.SnapshotSchema{}, fmt.Errorf("multiple snapshot found, snapshot: %s, snapshots in cluster: %v", snapshotId, response)
	}
	return response[0], nil
}

func (cp CPv1) GetVolumeSnapshot(volUuid string, snapshotId string) (common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	var response []common.SnapshotSchema
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshots", "--volume", volUuid, "--snapshot", snapshotId)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return common.SnapshotSchema{}, err
	}
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshot by snapshot id for volume)", "string", string(jsonInput))
		return common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	if len(response) == 0 {
		return common.SnapshotSchema{}, fmt.Errorf("snapshot not found, snapshot: %s, volume: %s", snapshotId, volUuid)
	} else if len(response) > 1 {
		return common.SnapshotSchema{}, fmt.Errorf("multiple snapshot found, snapshot: %s, volume: %s, snapshots in cluster: %v", snapshotId, volUuid, response)
	}
	return response[0], nil
}

func (cp CPv1) GetVolumeSnapshots(volUuid string) ([]common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshots", "--volume", volUuid)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.SnapshotSchema
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshots for volume)", "string", string(jsonInput))
		return []common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func (cp CPv1) GetVolumeSnapshotTopology() ([]common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshot-topology")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.SnapshotSchema
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshots)", "string", string(jsonInput))
		return []common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func (cp CPv1) GetPerSnapshotVolumeSnapshotTopology(snapshotId string) (common.SnapshotSchema, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	var response []common.SnapshotSchema
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume-snapshot-topology", "--snapshot", snapshotId)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return common.SnapshotSchema{}, err
	}

	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		// logf.Log.Info("Failed to unmarshal (get snapshot by snapshot id)", "string", string(jsonInput))
		return common.SnapshotSchema{}, fmt.Errorf("%s", errMsg)
	}
	if len(response) == 0 {
		return common.SnapshotSchema{}, fmt.Errorf("snapshot not found, snapshot: %s", snapshotId)
	} else if len(response) > 1 {
		return common.SnapshotSchema{}, fmt.Errorf("multiple snapshot found, snapshot: %s, snapshots in cluster: %v", snapshotId, response)
	}
	return response[0], nil
}
