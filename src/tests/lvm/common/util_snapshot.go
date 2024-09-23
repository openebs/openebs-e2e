package common

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/mayastor/snapshot"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// LvmVolumeSnapshotVerify verify snapshot and content to be ready
// it also verify that snapshot and content restore size should be zero
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
	restoreSize, err := k8stest.GetSnapshotRestoreSize(snapshotName, namespace)
	if err != nil {
		return false, err
	}

	restoreSizeInt, conversionStatus := restoreSize.AsInt64()
	if !conversionStatus {
		return false, fmt.Errorf("failed to convert snapshot restore size into int:, restore size: %v", restoreSize)
	} else if restoreSizeInt != 0 {
		return false, fmt.Errorf("snapshot restore size is not 0")
	}
	return true, nil
}
