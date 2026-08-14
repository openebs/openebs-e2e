package common

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/mayastor/snapshot"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// LvmVolumeSnapshotVerify verifies the snapshot and content are ready.
// Thick LVM snapshots have a zero restore size, represented as nil since external-snapshotter v8.
func LvmVolumeSnapshotVerify(snapshotName, snapshotContentName, namespace string, skipSnapError bool) (bool, error) {

	logf.Log.Info("Verify lvm snapshot content ready status")
	contentReady, err := snapshot.WaitForSnapshotContentReadyStatus(snapshotContentName, skipSnapError)
	if err != nil {
		return contentReady, err
	} else if !contentReady {
		logf.Log.Info("Snapshot content not ready", "VolumeSnapshotContent.status.readyToUse", contentReady)
		return contentReady, err
	}
	logf.Log.Info("Verify snapshot ready status")
	snapshotReady, err := snapshot.WaitForSnapshotReadyStatus(snapshotName, namespace, skipSnapError)
	if err != nil {
		return snapshotReady, err
	} else if !snapshotReady {
		logf.Log.Info("Snapshot not ready", "VolumeSnapshot.status.readyToUse", snapshotReady)
		return snapshotReady, err
	}

	logf.Log.Info("Verify snapshot restore size is zero")
	snapshotObj, err := k8stest.GetSnapshot(snapshotName, namespace)
	if err != nil {
		return false, err
	}
	if snapshotObj == nil {
		return false, fmt.Errorf("snapshot %s not found", snapshotName)
	}
	if snapshotObj.Status == nil {
		return false, fmt.Errorf("snapshot %s status not found", snapshotName)
	}
	if snapshotObj.Status.RestoreSize == nil {
		return true, nil
	}

	restoreSizeInt, conversionStatus := snapshotObj.Status.RestoreSize.AsInt64()
	if !conversionStatus {
		return false, fmt.Errorf("failed to convert snapshot restore size into int, restore size: %v", snapshotObj.Status.RestoreSize)
	} else if restoreSizeInt != 0 {
		return false, fmt.Errorf("snapshot restore size is not 0")
	}
	return true, nil
}
