package v1

// Utility functions for Mayastor control plane volume
import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const ErrOutput = "Error error"

func HasNotFoundRestJsonError(str string) bool {
	re := regexp.MustCompile(`Error error in response.*RestJsonError.*kind:\s*(\w+)`)
	frags := re.FindSubmatch([]byte(str))
	return len(frags) == 2 && string(frags[1]) == "NotFound"
}

func getMayastorCpVolume(uuid string) (*common.MayastorVolume, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volume", uuid)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response common.MayastorVolume
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		msg := string(jsonInput)
		if !HasNotFoundRestJsonError(msg) {
			logf.Log.Info("Failed to unmarshal (get volume)", "string", msg)
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return &response, nil
}

func listMayastorCpVolumes() ([]common.MayastorVolume, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volumes")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.MayastorVolume
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		logf.Log.Info("Failed to unmarshal (get volumes)", "string", string(jsonInput))
		return []common.MayastorVolume{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func scaleMayastorVolume(uuid string, replicaCount int) error {
	pluginpath := GetPluginPath()

	var err error
	var jsonInput []byte
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "scale", "volume", uuid, strconv.Itoa(replicaCount))
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return err
	}
	return nil
}

func GetMayastorVolumeState(volName string) (string, error) {
	msv, err := getMayastorCpVolume(volName)
	if err == nil {
		return msv.State.Status, nil
	}
	return "", err
}

func GetMayastorVolumeChildren(volName string) ([]common.TargetChild, error) {
	msv, err := getMayastorCpVolume(volName)
	if err != nil {
		return nil, err
	}
	return msv.State.Target.Children, nil
}

func GetMayastorVolumeChildState(uuid string) (string, error) {
	msv, err := getMayastorCpVolume(uuid)
	if err != nil {
		return "", err
	}
	return msv.State.Target.State, nil
}

func GetMayastorVolumeTargetNode(uuid string) (string, error) {
	msv, err := getMayastorCpVolume(uuid)
	if err != nil {
		return "", err
	}
	return msv.State.Target.Node, nil
}

func (cp CPv1) GetMsvTargetUuid(uuid string) (string, error) {
	msv, err := getMayastorCpVolume(uuid)
	if err != nil {
		return "", err
	}
	return msv.State.Target.Uuid, nil
}

func IsMayastorVolumePublished(uuid string) bool {
	msv, err := getMayastorCpVolume(uuid)
	if err == nil {
		return msv.Spec.Target.Node != ""
	}
	return false
}

func IsMayastorVolumeDeleted(uuid string) bool {
	msv, err := getMayastorCpVolume(uuid)
	if err != nil {
		if HasNotFoundRestJsonError(fmt.Sprintf("%v", err)) {
			return true
		}
		logf.Log.Error(err, "IsMayastorVolumeDeleted msv is nil")
		return false
	}
	if msv.Spec.Uuid == "" {
		return true
	}
	logf.Log.Info("IsMayastorVolumeDeleted", "msv", msv)
	return false
}

func CheckForMayastorVolumes() (bool, error) {
	logf.Log.Info("CheckForMayastorVolumes")
	foundResources := false

	msvs, err := listMayastorCpVolumes()
	if err == nil && msvs != nil && len(msvs) != 0 {
		logf.Log.Info("CheckForVolumeResources: found MayastorVolumes",
			"MayastorVolumes", msvs)
		foundResources = true
	}
	return foundResources, err
}

func CheckAllMayastorVolumesAreHealthy() error {
	allHealthy := true
	msvs, err := listMayastorCpVolumes()
	if err == nil && msvs != nil && len(msvs) != 0 {
		for _, msv := range msvs {
			if msv.State.Status != common.MsvStatusStateOnline {
				logf.Log.Info("CheckAllMayastorVolumesAreHealthy",
					"msv.State.Status", msv.State.Status,
					"msv.Spec", msv.Spec,
					"msv.State", msv.State,
				)
				allHealthy = false
			}
		}
	}

	if !allHealthy {
		return fmt.Errorf("CheckAllMayastorVolumesAreHealth: all MSVs were not healthy")
	}
	return err
}

// GetMSV Get pointer to a mayastor control plane volume
// returns nil and no error if the msv is in pending state.
func (cp CPv1) GetMSV(uuid string) (*common.MayastorVolume, error) {
	cpMsv, err := getMayastorCpVolume(uuid)
	if err != nil {
		return nil, fmt.Errorf("GetMSV: %v", err)
	}
	if cpMsv.Spec.Uuid == "" {
		logf.Log.Info("Msv not found", "uuid", uuid)
		return nil, nil
	}

	// pending means still being created
	if cpMsv.State.Status == "pending" {
		return nil, nil
	}

	//logf.Log.Info("GetMSV", "msv", msv)
	// Note: msVol.Node can be unassigned here if the volume is not mounted
	if cpMsv.State.Status == "" {
		return nil, fmt.Errorf("GetMSV: state not defined, got msv.Status=\"%v\"", cpMsv.State)
	}

	if cpMsv.Spec.Num_replicas < 1 {
		return nil, fmt.Errorf("GetMSV: msv.Spec.Num_replicas=\"%v\"", cpMsv.Spec.Num_replicas)
	}

	// msv := cpVolumeToMsv(cpMsv, cp.nodeIPAddresses)
	return cpMsv, nil
}

// GetMsvNodes Retrieve the nexus node hosting the Mayastor Volume,
// and the names of the replica nodes
func (cp CPv1) GetMsvNodes(uuid string) (string, []string) {
	msv, err := getMayastorCpVolume(uuid)
	if err != nil {
		logf.Log.Info("failed to get mayastor volume", "uuid", uuid)
		return "", nil
	}
	node := msv.State.Target.Node
	replicas := make([]string, msv.Spec.Num_replicas)

	msvReplicas, err := cp.GetMsvReplicas(uuid)
	if err != nil {
		logf.Log.Info("failed to get mayastor volume replica", "uuid", uuid)
		return node, nil
	}
	for ix, r := range msvReplicas {
		replicas[ix] = r.Replica.Node
	}
	return node, replicas
}

func (cp CPv1) ListMsvs() ([]common.MayastorVolume, error) {
	return listMayastorCpVolumes()
}

func (cp CPv1) ListRestoredMsvs() ([]common.MayastorVolume, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "volumes", "--source", "snapshot")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []common.MayastorVolume
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		logf.Log.Info("Failed to unmarshal (get volumes)", "string", string(jsonInput))
		return []common.MayastorVolume{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func (cp CPv1) SetMsvReplicaCount(uuid string, replicaCount int) error {
	err := scaleMayastorVolume(uuid, replicaCount)
	logf.Log.Info("ScaleMayastorVolume", "Num_replicas", replicaCount)
	return err
}

func (cp CPv1) GetMsvState(uuid string) (string, error) {
	return GetMayastorVolumeState(uuid)
}

func (cp CPv1) GetMsvReplicas(volName string) ([]common.MsvReplica, error) {
	vol, err := cp.GetMSV(volName)
	if err != nil {
		logf.Log.Info("Failed to get volume", "uuid", volName, "error", err)
		return nil, err
	}
	var msvReplicas []common.MsvReplica
	for uuid, replica := range vol.State.ReplicaTopology {
		var replicaUri string
		for _, child := range vol.State.Target.Children {
			if strings.Contains(child.Uri, uuid) {
				replicaUri = child.Uri
				break
			}
		}
		replica := common.MsvReplica{
			Uuid:    uuid,
			Uri:     replicaUri,
			Replica: replica,
		}
		msvReplicas = append(msvReplicas, replica)
	}

	return msvReplicas, nil
}

func (cp CPv1) GetMsvReplicaTopology(volUuid string) (common.ReplicaTopology, error) {
	vol, err := cp.GetMSV(volUuid)
	if err != nil {
		logf.Log.Info("Failed to get replicas", "uuid", volUuid, "error", err)
		return nil, err
	}
	return vol.State.ReplicaTopology, nil
}

func (cp CPv1) GetMsvNexusChildren(volName string) ([]common.TargetChild, error) {
	return GetMayastorVolumeChildren(volName)
}

func (cp CPv1) GetMsvNexusState(uuid string) (string, error) {
	return GetMayastorVolumeChildState(uuid)
}

func (cp CPv1) IsMsvPublished(uuid string) bool {
	return IsMayastorVolumePublished(uuid)
}

func (cp CPv1) IsMsvDeleted(uuid string) bool {
	return IsMayastorVolumeDeleted(uuid)
}

func (cp CPv1) CheckForMsvs() (bool, error) {
	return CheckForMayastorVolumes()
}

func (cp CPv1) CheckAllMsvsAreHealthy() error {
	return CheckAllMayastorVolumesAreHealthy()
}

func (cp CPv1) CanDeleteMsv() bool {
	return false
}

func (cp CPv1) DeleteMsv(volName string) error {
	return fmt.Errorf("delete of mayastor volume not supported %v", volName)
}

func (cp CPv1) GetMsvTargetNode(volName string) (string, error) {
	return GetMayastorVolumeTargetNode(volName)
}

func (cp CPv1) GetMsvSize(uuid string) (int64, error) {
	return getMsvSize(uuid)
}

func getMsvSize(volName string) (int64, error) {
	var size int64
	msv, err := getMayastorCpVolume(volName)
	if err == nil {
		size = msv.State.Size
	}
	return size, err
}

func (cp CPv1) SetVolumeMaxSnapshotCount(uuid string, maxSnapshotCount int32) error {
	pluginpath := GetPluginPath()

	var err error
	var jsonInput []byte
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "set", "volume", uuid, "max-snapshots", strconv.Itoa(int(maxSnapshotCount)))
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return err
	}
	return nil
}

func (cp CPv1) GetMsvMaxSnapshotCount(uuid string) (int32, error) {
	var maxSnapshotCount int32
	msv, err := getMayastorCpVolume(uuid)
	if err == nil {
		maxSnapshotCount = msv.Spec.MaxSnapshots
	}
	return maxSnapshotCount, err
}
