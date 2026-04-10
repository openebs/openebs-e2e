package restore

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs                         = 180 // in seconds
	PreConditionFailedErrorSubstring       = "412 Precondition Failed"
	PoolNotReadyToTakeCloneOfSnapshot      = "Pool not ready to take clone of snapshot"
	NotReady                               = "is not Ready"
	RangeNotSatisfiable                    = "416 Range Not Satisfiable"
	CloneVolumeSizeErrorString             = "Cloned snapshot volume must match the snapshot size"
	BadRequest                             = "400 Bad Request"
	CloneVolumeThinErrorString             = "Cloned snapshot volumes must be thin provisioned"
	ErrorGettingHandleDataSource           = "error getting handle for DataSource"
	MultiReplicaRestore                    = "Cannot create a multi-replica volume from a snapshot of a single-replica volume"
	FailedToStageVolume                    = "Failed to stage volume"
	FailedToProvisionWithCrossVolumeMode   = "snapshot.storage.kubernetes.io/allow-volume-mode-change annotation is not present on snapshotcontent"
	VolumeModeConversionAnnotationKey      = "snapshot.storage.kubernetes.io/allow-volume-mode-change"
	VolumeModeConversionAnnotationValue    = "true"
	ClonedSnapshotVolumeThinErrorsubstring = "ClonedSnapshotVolumeThin: Cloned snapshot volumes must be thin provisioned"
)

func VerifyRestoredVolumeList(restVol map[string]common.Snapshot) (bool, error) {
	logf.Log.Info("Verify restored volumes list via control plane")
	restoredVolumes, err := controlplane.ListRestoredMsvs()
	if err != nil {
		return false, err
	} else if len(restoredVolumes) != len(restVol) {
		return false, fmt.Errorf("number of restored volumes (%d) do not  match with expected number of restored volume (%d)",
			len(restoredVolumes),
			len(restVol))
	}

	for volUuid, source := range restVol {
		var restoredVol common.MayastorVolume
		for _, vol := range restoredVolumes {
			if vol.Spec.Uuid == volUuid {
				restoredVol = vol
				break
			}
		}
		if restoredVol.Spec.Uuid == "" {
			return false, fmt.Errorf("volume %s not found in restored volume list", volUuid)
		} else if restoredVol.Spec.ContentSource.Snapshot != source {
			logf.Log.Info("volume source mismatch which source",
				"volume", volUuid,
				"actual source", restoredVol.Spec.ContentSource.Snapshot,
				"expected source", source)
			return false, nil
		}
	}
	return true, err

}
