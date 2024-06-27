package snapshot

import (
	"fmt"
	"strings"
	"time"

	v1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs                           = 330 //exponential retry has max timeout of 300 seconds, so DefTimeoutSecs is set to 330 seconds
	DefSnapsErrorTimeout                     = 60  // in seconds
	IoEnginePodName                          = e2e_config.GetConfig().Product.IOEnginePodName
	FioRunTime                               = 1200 // in seconds
	VolSizeMb                                = 8192 // in  MiB
	BadRequestErrorSubstring                 = "400 Bad Request"
	PreConditionFailedErrorSubstring         = "412 Precondition Failed"
	PreConditionFailedFsfreezeErrorSubstring = "Preflight check for fsfreeze failed, nvmf subsystem not in desired state"
	DefCpTimeoutSecs                         = 60 // in seconds, timeout value used for snapshot verification via control plane
)

// CreateVolumeSnapshot create snapshot class and snapshot
// it return snapshot object , snapshot content name and err
func CreateVolumeSnapshot(snapshotClassName string, snapshotName string, pvc string, namespace string) (*v1.VolumeSnapshot, string, error) {
	err := k8stest.CreateSnapshotClass(snapshotClassName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create snapshot class %s, error: %v", snapshotClassName, err)
	}

	// add sleep of two seconds before fetching snapshot class
	time.Sleep(5 * time.Second)
	// Get snapshot class before creating snapshot to make sure that snapshot class created successfully
	vsc, err := k8stest.GetSnapshotClass(snapshotClassName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get snapshot class %s, error: %v", snapshotClassName, err)
	} else if vsc == nil {
		return nil, "", fmt.Errorf("snapshot class %s not found, vsc: %v", snapshotClassName, *vsc)
	}
	snapshot, snapshotContentName, err := k8stest.MakeSnapshot(pvc, snapshotClassName, namespace, snapshotName)
	return snapshot, snapshotContentName, err
}

// DeleteVolumeSnapshot delete snapshot class and snapshot and returns an err
func DeleteVolumeSnapshot(snapshotClassName string, snapshotName string, namespace string) error {
	err := k8stest.RemoveSnapshot(snapshotName, namespace)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot %s, error: %v", snapshotName, err)
	}
	return k8stest.DeleteSnapshotClass(snapshotClassName)
}

// Wait for the Snapshot content to be ready.
func WaitForSnapshotContentReadyStatus(snapshotContentName string, skipSnapError bool) (bool, error) {
	var ready bool
	var err error
	var snapErr *v1.VolumeSnapshotError
	const timeSleepSecs = 1
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		ready, snapErr, err = k8stest.GetSnapshotContentReadyStatus(snapshotContentName)
		if snapErr != nil {
			if !skipSnapError {
				break
			}
		} else if err == nil && ready {
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return ready, err
}

// Wait for the Snapshot to be ready.
func WaitForSnapshotReadyStatus(snapshotName string, namespace string, skipSnapError bool) (bool, error) {
	const timeSleepSecs = 1
	var ready bool
	var err error
	var snapErr *v1.VolumeSnapshotError
	// Wait for the Snapshot to be ready
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		ready, snapErr, err = k8stest.GetSnapshotReadyStatus(snapshotName, namespace)
		if snapErr != nil {
			if !skipSnapError {
				break
			}
		} else if err == nil && ready {
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return ready, err
}

// VerifySuccessfulSnapshotCreation verify:
// 1. snapshot and snapshot content ready status
// 2. snapshot restore size should not be zero
// 3. snapshot creation time should not be empty
func VerifySuccessfulSnapshotCreation(snapshotName string, snapshotContentName string, namespace string, skipSnapError bool) (bool, error) {
	logf.Log.Info("Verify snapshot content ready status")
	contentReady, err := WaitForSnapshotContentReadyStatus(snapshotContentName, skipSnapError)
	if err != nil {
		return contentReady, err
	} else if !contentReady {
		logf.Log.Info("Snapshot content not ready", "VolumeSnapshotContent.status.readyToUse", contentReady)
		return contentReady, err
	}
	logf.Log.Info("Verify snapshot ready status")
	snapshotReady, err := WaitForSnapshotReadyStatus(snapshotName, namespace, skipSnapError)
	if err != nil {
		return snapshotReady, err
	} else if !snapshotReady {
		logf.Log.Info("Snapshot not ready", "VolumeSnapshot.status.readyToUse", snapshotReady)
		return snapshotReady, err
	}

	restoreSize, err := k8stest.GetSnapshotRestoreSize(snapshotName, namespace)
	if err != nil {
		return false, err
	}

	restoreSizeInt, conversionStatus := restoreSize.AsInt64()
	if !conversionStatus {
		return false, fmt.Errorf("failed to convert snapshot restore size into int:, restore size: %v", restoreSize)
	} else if restoreSizeInt == 0 {
		return false, fmt.Errorf("snapshot restore size is 0")
	}

	creationTime, err := k8stest.GetSnapshotCreationTime(snapshotName, namespace)
	if err != nil {
		return false, err
	}

	if !isValidTimestamp(creationTime) {
		return false, fmt.Errorf("snapshot creation time %s not in RFC3339 time layout", creationTime)
	}

	return true, err

}

// VerifySuccessfulCpSnapshotCreation verify control plane using kubectl mayastor plugin :
// 1. snapshot and snapshot content ready status
// 2. snapshot restore size should not be zero
// 3. snapshot creation time should not be empty
func VerifySuccessfulCpSnapshotCreation(snapshotName string, volUuid string, namespace string) (bool, error) {
	//get snapshot Uid
	snapUid, err := k8stest.GetSnapshotUid(snapshotName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get snapshot Uid for snapshot %s, error: %v", snapshotName, err)
	}
	logf.Log.Info("Snapshot", "Uid", snapUid)

	snapshotReady, err := WaitForCpSnapshotReadyStatus(snapUid, volUuid)
	if err != nil {
		return snapshotReady, err
	} else if !snapshotReady {
		logf.Log.Info("Snapshot not ready", "Snapshot.State.ReplicaSnapshot[0].online", snapshotReady)
		return snapshotReady, err
	}

	restoreSize, err := k8stest.GetSnapshotCpRestoreSize(snapUid, volUuid)
	if err != nil {
		return false, err
	}
	if restoreSize == 0 {
		return false, fmt.Errorf("snapshot restore size is 0")
	}

	creationTime, err := k8stest.GetSnapshotCpTimestamp(snapUid, volUuid)
	if err != nil {
		return false, err
	} else if creationTime == "" {
		return false, fmt.Errorf("snapshot creation time should not empty string")
	}

	msv, err := k8stest.GetMSV(volUuid)
	if err != nil {
		return false, err
	}
	replicaCountsMatch, err := VerifySnapshotReplicas(snapUid, msv.Spec.Num_replicas, volUuid)
	if err != nil {
		return false, err
	}
	if !replicaCountsMatch {
		return false, fmt.Errorf("snapshot replica count does not match volume replica count")
	}

	return true, err

}

// Wait for the Snapshot to be ready as per control plane.
func WaitForCpSnapshotReadyStatus(snapshotUid string, volUuid string) (bool, error) {
	const timeSleepSecs = 1
	var ready bool
	var err error
	// Wait for the Snapshot to be ready
	for ix := 0; ix < DefCpTimeoutSecs/timeSleepSecs; ix++ {
		ready, err = k8stest.GetSnapshotCpCloneReadyStatus(snapshotUid, volUuid)
		if err == nil && ready {
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return ready, err
}

func isValidTimestamp(timestamp string) bool {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err == nil
}

// VerifySnapshotReplicas
// to check if number of SNAPSHOT_REPLICAS for a volume
// matches with the number of replicas for a volume
func VerifySnapshotReplicas(snapUuid string, replicaCount int, volUuid string) (bool, error) {
	snap, err := controlplane.GetSnapshot(snapUuid)
	if err != nil {
		return false, fmt.Errorf("failed to get snapshot object with uuid %s, error: %v", snapUuid, err)
	}
	logf.Log.Info("Snapshot Replicas", "num_snapshot_replicas", snap.Definition.Metadata.NumSnapshotReplicas)
	logf.Log.Info("Volume Replicas as per Storage Class", "Volume Replica Count", replicaCount)
	if snap.Definition.Metadata.NumSnapshotReplicas != replicaCount {
		return false, fmt.Errorf("number of snapshot_replicas for volUuid %s does not match with the number of replicas %d", volUuid, replicaCount)
	}
	return true, err
}

// DeleteVolumeSnapshot delete snapshot class and snapshot and returns an err
func DeleteFailedVolumeSnapshot(snapshotClassName string, snapshotName string, snapshotContentName string, namespace string) error {
	err := RemoveFailedSnapshot(snapshotName, snapshotContentName, namespace)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot %s, error: %v", snapshotName, err)
	}
	return k8stest.DeleteSnapshotClass(snapshotClassName)
}

// RemoveFailedSnapshot delete a snapshot and remove annotation from snapshot content and verify that
//  1. The Volume Snapshot is deleted
//  2. The associated Volume Snapshot ContentPV is deleted
func RemoveFailedSnapshot(snapshotName string, snapshotContentName string, namespace string) error {
	const timoSleepSecs = 1
	logf.Log.Info("Removing Volume Snapshot", "snapshot", snapshotName, "namespace ", namespace)
	var isDeleted bool
	t0 := time.Now()

	// Confirm Volume snapshot exist before deleting
	snapshot, getErr := k8stest.GetSnapshot(snapshotName, namespace)

	if k8serrors.IsNotFound(getErr) {
		return fmt.Errorf("Volume Snapshot %s not found, namespace: %s, error: %v", snapshotName, namespace, getErr)
	} else if getErr != nil {
		return fmt.Errorf("failed to get volume snapshot %s, namespace: %s, error: %v", snapshotName, namespace, getErr)
	} else if snapshot == nil {
		return fmt.Errorf("Volume Snapshot %s not found, namespace: %s", snapshotName, namespace)
	}

	// Delete Volume Snapshot
	logf.Log.Info("Delete volume snapshot")
	deleteErr := k8stest.DeleteSnapshot(snapshotName, namespace)
	if deleteErr != nil {
		return fmt.Errorf("failed to delete Volume Snapshot %s, namespace: %s, error: %v", snapshotName, namespace, deleteErr)
	}

	logf.Log.Info("Waiting for Snapshot to be deleted", "snapshot", snapshotName)
	// Wait for the Volume Snapshot content annotations
	logf.Log.Info("Verify volume snapshot content annotation")
	var err error
	var annotations map[string]string
	for ix := 0; ix < DefTimeoutSecs/timoSleepSecs; ix++ {
		annotations, err = k8stest.GetSnapshotContentAnnotation(snapshotContentName)
		if len(annotations) > 1 {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
		logf.Log.Info("RemoveFailedSnapshot", "annotations", annotations)
	}
	if err != nil {
		return err
	} else if len(annotations) != 2 {
		return fmt.Errorf("more then two annotations exist on snapshot content , snapshot content: %s, annotations: %s", snapshotContentName, annotations)
	}

	// remove volume snapshot content annotations
	logf.Log.Info("Delete volume snapshot content annotation")
	err = k8stest.RemoveSnapshotContentAnnotation(snapshotContentName)

	// Wait for the Volume Snapshot to be deleted.
	logf.Log.Info("Wait for volume snapshot to be deleted")
	for ix := 0; ix < DefTimeoutSecs/timoSleepSecs; ix++ {
		isDeleted, err = k8stest.IsSnapshotDeleted(snapshotName, namespace)
		if isDeleted {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return err
	} else if !isDeleted {
		return fmt.Errorf("volume snapshot not deleted, snapshot: %s, namespace: %s", snapshotName, namespace)
	}

	// Wait for the Volume Snapshot Content to be deleted.
	logf.Log.Info("Waiting for Volume Snapshot Content to be deleted", "snapshot content", snapshotContentName)

	for ix := 0; ix < DefTimeoutSecs/timoSleepSecs; ix++ {
		isDeleted, err = k8stest.IsSnapshotContentDeleted(snapshotContentName)
		if isDeleted {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return err
	} else if !isDeleted {
		return fmt.Errorf("Volume Snapshot content not deleted, snapshot content: %s", snapshotContentName)
	}
	logf.Log.Info("Deleted", "Snapshot", snapshotName, "Snapshot Content", snapshotContentName, "elapsed time", time.Since(t0))
	return nil
}

// Wait for the Snapshot error to be updated.
func WaitForSnapshotError(snapshot string, namespace string, errorMessageSubstring string) (bool, error) {
	var err error
	var isPresent bool
	const timeSleepSecs = 1
	for ix := 0; ix < DefSnapsErrorTimeout/timeSleepSecs; ix++ {
		var snap *v1.VolumeSnapshot
		snap, err = k8stest.GetSnapshot(snapshot, namespace)
		if err == nil &&
			snap != nil &&
			snap.Status.Error != nil &&
			snap.Status.Error.Message != nil &&
			strings.Contains(*snap.Status.Error.Message, errorMessageSubstring) {
			isPresent = true
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return isPresent, err
}

// func GetSnapshotEvents(snapshotName string, namespace string) (*coreV1.EventList, error) {
// 	options := metaV1.ListOptions{
// 		TypeMeta:      metaV1.TypeMeta{Kind: "VolumeSnapshot"},
// 		FieldSelector: fmt.Sprintf("involvedObject.name=%s", snapshotName),
// 	}
// 	return k8stest.GetEvents(namespace, options)
// }

// func IsWarningSnapshotEventPresent(snapshotName string, namespace string, errorSubstring string) (bool, error) {
// 	events, err := k8stest.GetPvcEvents(snapshotName, namespace)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to get snapshot events in namespace %s, error: %v", namespace, err)
// 	}
// 	for _, event := range events.Items {
// 		logf.Log.Info("Found Snapshot event", "message", event.Message)
// 		if event.Type == "Warning" && strings.Contains(event.Message, errorSubstring) {
// 			return true, err
// 		}
// 	}
// 	return false, err
// }
