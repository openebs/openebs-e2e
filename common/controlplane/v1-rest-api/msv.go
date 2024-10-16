package v1_rest_api

// Utility functions for Mayastor CRDs
import (
	"fmt"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	openapiClient "github.com/openebs/openebs-e2e/common/generated/openapi"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// GetMSV Get pointer to a mayastor volume custom resource
// returns nil and no error if the msv is in pending state.
func (cp CPv1RestApi) GetMSV(uuid string) (*common.MayastorVolume, error) {
	vol, err, _ := cp.oa.getVolume(uuid)

	if err != nil {
		return nil, fmt.Errorf("GetMSV: %v", err)
	}

	// still being created
	volSpec, ok := vol.GetSpecOk()
	if !ok {
		logf.Log.Info("GetMSV !GetSpecOk()", "uuid", uuid)
		return nil, nil
	}

	volState, ok := vol.GetStateOk()
	if !ok {
		logf.Log.Info("GetMSV !vol.GetStateOk()", "uuid", uuid)
		return nil, nil
	}

	if volState.GetStatus() == openapiClient.VOLUMESTATUS_UNKNOWN {
		return nil, fmt.Errorf("GetMSV: state not defined, got msv.Status=\"%v\"", vol.GetState().Status)
	}

	if volSpec.GetNumReplicas() < 1 {
		return nil, fmt.Errorf("GetMsv msv.Spec.NumReplicas=\"%v\"", volSpec.GetNumReplicas())
	}

	nexus, ok := volState.GetTargetOk()
	if ok {
		if len(nexus.Children) < 1 {
			return nil, fmt.Errorf("GetMSV: nexus children =\"%v\"", nexus.Children)
		}
	}

	msv := cp.oa.volToMsv(vol)
	return &msv, nil
}

// GetMsvNodes Retrieve the nexus node hosting the Mayastor Volume,
// and the names of the replica nodes
// function asserts if the volume CR is not found.
func (cp CPv1RestApi) GetMsvNodes(uuid string) (string, []string) {
	var nexusNode string
	var replicaNodes []string

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		nexusNode = vol.State.Target.Node
		for _, replica := range vol.State.ReplicaTopology {
			replicaNodes = append(replicaNodes, replica.GetNode())
		}
	}
	return nexusNode, replicaNodes
}

func (cp CPv1RestApi) CanDeleteMsv() bool {
	return true
}

func (cp CPv1RestApi) DeleteMsv(uuid string) error {
	err, _ := cp.oa.deleteVolume(uuid)
	return err
}

func (cp CPv1RestApi) ListMsvs() ([]common.MayastorVolume, error) {
	vols, err, _ := cp.oa.getVolumes()
	if err != nil {
		return nil, fmt.Errorf("ListMsvs: %v", err)
	}

	var msvs []common.MayastorVolume
	if err == nil {
		for _, vol := range vols.Entries {
			msvs = append(msvs, cp.oa.volToMsv(vol))
		}
	}
	return msvs, err
}

func (cp CPv1RestApi) ListRestoredMsvs() ([]common.MayastorVolume, error) {
	return nil, fmt.Errorf("not implemented")
}

func (cp CPv1RestApi) SetMsvReplicaCount(uuid string, replicaCount int) error {
	err, _ := cp.oa.putReplicaCount(uuid, replicaCount)
	return err
}

func (cp CPv1RestApi) GetMsvState(uuid string) (string, error) {
	vol, err, _ := cp.oa.getVolume(uuid)
	var volState string

	if err == nil {
		volState = string(vol.State.GetStatus())
	}

	return volState, err
}

func (cp CPv1RestApi) GetMsvReplicas(uuid string) ([]common.MsvReplica, error) {
	var replicas []common.MsvReplica

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		for rUuid, topology := range vol.State.ReplicaTopology {
			var replicaUri string
			for _, child := range vol.State.Target.Children {
				if strings.Contains(child.Uri, uuid) {
					replicaUri = child.Uri
					break
				}
			}
			replica := common.Replica{
				Node:  topology.GetNode(),
				State: string(topology.GetState()),
				Pool:  topology.GetPool(),
				// FIXME: Update openapi client
				// Usage:       common.ReplicaUsage{},
				// ChildStatus: "",
			}
			replicas = append(replicas, common.MsvReplica{
				Uuid:    rUuid,
				Uri:     replicaUri,
				Replica: replica,
			})
		}
	}
	return replicas, err
}

func (cp CPv1RestApi) GetMsvReplicaTopology(volUuid string) (common.ReplicaTopology, error) {
	return common.ReplicaTopology{}, fmt.Errorf("not implemented")
}

func (cp CPv1RestApi) GetMsvNexusChildren(uuid string) ([]common.TargetChild, error) {
	var children []common.TargetChild

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		for _, inChild := range vol.State.GetTarget().Children {
			children = append(children, common.TargetChild{
				Uri:             inChild.Uri,
				State:           string(inChild.State),
				RebuildProgress: inChild.RebuildProgress,
			})
		}
	}
	return children, err
}

func (cp CPv1RestApi) GetMsvNexusState(uuid string) (string, error) {
	var nexusState string

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		nexusState = string(vol.State.Target.GetState())
	}
	return nexusState, err
}

func (cp CPv1RestApi) GetMsvTargetNode(uuid string) (string, error) {
	var targetNode string

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		targetNode = string(vol.State.Target.GetNode())
	}
	return targetNode, err
}

func (cp CPv1RestApi) GetMsvTargetUuid(uuid string) (string, error) {
	var targetUuid string
	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		targetUuid = string(vol.State.Target.GetUuid())
	}
	return targetUuid, err
}

func (cp CPv1RestApi) IsMsvPublished(uuid string) bool {
	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		return vol.Spec.Target.Node != ""
	}
	return false
}

func (cp CPv1RestApi) IsMsvDeleted(uuid string) bool {
	_, err, responseStatusCode := cp.oa.getVolume(uuid)
	return err != nil && responseStatusCode == 404
}

func (cp CPv1RestApi) CheckForMsvs() (bool, error) {
	vols, err, _ := cp.oa.getVolumes()
	return err == nil && len(vols.Entries) != 0, err
}

func (cp CPv1RestApi) CheckAllMsvsAreHealthy() error {
	allHealthy := true
	vols, err, _ := cp.oa.getVolumes()

	if err == nil {
		for _, vol := range vols.Entries {
			if vol.State.GetStatus() != openapiClient.VOLUMESTATUS_ONLINE {
				allHealthy = false
				logf.Log.Info("CheckAllMsvsAreHealthy", "vol", vol)
			}
		}
		if !allHealthy {
			err = fmt.Errorf("all MSVs were not healthy")
		}
	}

	return err
}

func (cp CPv1RestApi) GetMsvSize(uuid string) (int64, error) {
	var msvSize int64

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		msvSize = vol.State.GetSize()
	}
	return msvSize, err
}

func (cp CPv1RestApi) GetMsvDeviceUri(uuid string) (string, error) {
	vol, err, _ := cp.oa.getVolume(uuid)
	var deviceUri string

	if err == nil {
		deviceUri = string(vol.State.Target.DeviceUri)
	}

	return deviceUri, err
}

func (cp CPv1RestApi) SetVolumeMaxSnapshotCount(uuid string, maxSnapshotCount int32) error {
	err, _ := cp.oa.setVolumeMaxSnapshotCount(uuid, maxSnapshotCount)
	return err
}

func (cp CPv1RestApi) GetMsvMaxSnapshotCount(uuid string) (int32, error) {
	var maxSnapshotCount int32

	vol, err, _ := cp.oa.getVolume(uuid)
	if err == nil {
		maxSnapshotCount = *vol.GetSpec().MaxSnapshots
	}
	return maxSnapshotCount, err
}
