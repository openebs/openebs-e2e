package k8stest

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// GetMSV Get pointer to a mayastor volume custom resource
// returns nil and no error if the msv is in pending state.
func GetMSV(uuid string) (*common.MayastorVolume, error) {
	return controlplane.GetMSV(uuid)
}

// GetMsvNodes Retrieve the nexus node hosting the Mayastor Volume,
// and the names of the replica nodes
// function asserts if the volume CR is not found.
func GetMsvNodes(uuid string) (string, []string) {
	return controlplane.GetMsvNodes(uuid)
}

func DeleteMsv(volName string) error {
	return controlplane.DeleteMsv(volName)
}

func ListMsvs() ([]common.MayastorVolume, error) {
	return controlplane.ListMsvs()
}

func SetMsvReplicaCount(uuid string, replicaCount int) error {
	return controlplane.SetMsvReplicaCount(uuid, replicaCount)
}

func SetMsvEncryption(uuid string, encryption bool) error {
	return controlplane.SetMsvEncryption(uuid, encryption)
}

func GetMsvState(uuid string) (string, error) {
	return controlplane.GetMsvState(uuid)
}

func GetMsvReplicas(volName string) ([]common.MsvReplica, error) {
	return controlplane.GetMsvReplicas(volName)
}

func GetMsvReplicaTopology(volUuid string) (common.ReplicaTopology, error) {
	return controlplane.GetMsvReplicaTopology(volUuid)
}

func GetMsvNexusChildren(volName string) ([]common.TargetChild, error) {
	return controlplane.GetMsvNexusChildren(volName)
}

func GetMsvNexusState(uuid string) (string, error) {
	return controlplane.GetMsvNexusState(uuid)
}

func IsMsvPublished(uuid string) bool {
	return controlplane.IsMsvPublished(uuid)
}

func IsMsvEncrypted(uuid string) bool {
	return controlplane.IsMsvEncrypted(uuid)
}

func IsMsvDeleted(uuid string) bool {
	return controlplane.IsMsvDeleted(uuid)
}

func CheckForMsvs() (bool, error) {
	return controlplane.CheckForMsvs()
}

func CheckAllMsvsAreHealthy() error {
	return controlplane.CheckAllMsvsAreHealthy()
}

func GetMsvTargetNode(uuid string) (string, error) {
	return controlplane.GetMsvTargetNode(uuid)
}

func GetMsvSize(uuid string) (int64, error) {
	return controlplane.GetMsvSize(uuid)
}

func GetMsvDeviceUri(uuid string) (string, error) {
	return controlplane.GetMsvDeviceUri(uuid)
}

func GetMsvMaxSnapshotCount(uuid string) (int32, error) {
	return controlplane.GetMsvMaxSnapshotCount(uuid)
}

// GetMsvsForStatefulSet retrieves all MSVs associated with a StatefulSet.
func GetMsvsForStatefulSet(stsName, namespace string) ([]*common.MayastorVolume, error) {
	var msvs []*common.MayastorVolume

	// List PVCs in the namespace
	pvcList, err := ListPVCs(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list PVCs in namespace %s: %v", namespace, err)
	}

	for _, pvc := range pvcList.Items {
		// Check if the PVC belongs to the StatefulSet
		if !isPvcForStatefulSet(pvc.Name, stsName) {
			continue
		}

		// Get the PV associated with the PVC
		pv, err := GetPV(pvc.Spec.VolumeName)
		if err != nil {
			logf.Log.Info("Processing PVC", "PVC Name", pvc.Name)
			return nil, fmt.Errorf("failed to get PV for PVC %s: %v", pvc.Name, err)
		}

		// Get the MSV associated with the PV
		volUuid := pv.Spec.CSI.VolumeHandle
		msv, err := GetMSV(volUuid)
		if err != nil {
			logf.Log.Info("Retrieved MSV", "MSV UUID", msv.Spec.Uuid)
			return nil, fmt.Errorf("failed to get MSV for Volume UUID %s: %v", volUuid, err)
		}

		msvs = append(msvs, msv)
	}

	return msvs, nil
}
