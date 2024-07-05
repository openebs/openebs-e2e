package k8stest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	defGetSnapshotTimeout = 60
	DefCompletionTimeout  = 600
)

var csiDriver = e2e_config.GetConfig().Product.CsiProvisioner

// MakeSnapshot Create volume snapshot and verify that
//  1. The Snapshot created
//  2. The Snapshot class exist
//  3. Volume snapshot content provisioned
//  4. Bound Volume Snapshot content ready
//  5. Volume Snapshot ready
//
// MakeSnapshot function return snapshot object, corresponding snapshot content name and error
func MakeSnapshot(pvcName string, snapshotClassName string, nameSpace string, snapshotName string) (*snapshotv1.VolumeSnapshot, string, error) {
	const timoSleepSecs = 1
	var snapshotContentName string
	logf.Log.Info("Creating", "Snapshot", snapshotName, "Snapshot class", snapshotClassName, "pvc", pvcName, "namespace", nameSpace)
	t0 := time.Now()
	// Create the Snapshot
	snapshot, err := CreateSnapshot(pvcName, snapshotClassName, nameSpace, snapshotName)
	if err != nil {
		return snapshot, snapshotContentName, fmt.Errorf("failed to create snapshot: %s,  pvc: %s, error: %v", snapshotName, pvcName, err)
	}

	_, err = GetSnapshotClass(snapshotClassName)
	if err != nil {
		return snapshot, snapshotContentName, fmt.Errorf("failed to get snapshot class: %s, error: %v", snapshotClassName, err)
	}

	// Refresh the Snapshot , so that we can get the Snapshot content name.
	var snap *snapshotv1.VolumeSnapshot
	for ix := 0; ix < defGetSnapshotTimeout/timoSleepSecs; ix++ {
		snap, err = GetSnapshot(snapshotName, nameSpace)
		if err == nil && snap != nil && snap.Status != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	logf.Log.Info("GetSnapshot", "Snapshot", snap)
	if err != nil {
		return snap, snapshotContentName, fmt.Errorf("failed to get snapshot: %s, namespace: %s, error: %v", snapshotName, nameSpace, err)
	} else if snap == nil {
		return snap, snapshotContentName, fmt.Errorf("snapshot %s not found, namespace: %s", snapshotName, nameSpace)
	} else if snap.Status == nil {
		return snap, snapshotContentName, fmt.Errorf("snapshot status not found for snapshot: %s, namespace: %s", snapshotName, nameSpace)
	} else if snap.Status.BoundVolumeSnapshotContentName == nil {
		return snap, snapshotContentName, fmt.Errorf("snapshot content not found for snapshot: %s, namespace: %s", snapshotName, nameSpace)
	}
	snapshotContentName = *snap.Status.BoundVolumeSnapshotContentName

	// Wait for the volume snapshot content to be provisioned
	var snapContent *snapshotv1.VolumeSnapshotContent
	for ix := 0; ix < defGetSnapshotTimeout/timoSleepSecs; ix++ {
		snapContent, err = GetSnapshotContent(snapshotContentName)
		if err == nil && snapContent != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return snap, snapshotContentName, fmt.Errorf("failed to get snapshot content: %s, error: %v", snapshotContentName, err)
	} else if snapContent == nil {
		return snap, snapshotContentName, fmt.Errorf("volume snapshot content not found , snapshot: %s", snapshotName)
	}

	logf.Log.Info("Created", "Snapshot", snapshotName, "Snapshot Content", snapshotContentName, "Snapshot class", snapshotClassName, "pvc", pvcName, "elapsed time", time.Since(t0))
	return snap, snapshotContentName, nil
}

// Create the Volume Snapshot , this will create snapshot kubernetes object only. This will be called by MakeSnapshot function
func CreateSnapshot(pvcName string, snapshotClassName string, nameSpace string, snapshotName string) (*snapshotv1.VolumeSnapshot, error) {
	// Snapshot create options
	createOpts := &snapshotv1.VolumeSnapshot{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      snapshotName,
			Namespace: nameSpace,
		},
		Spec: snapshotv1.VolumeSnapshotSpec{
			VolumeSnapshotClassName: &snapshotClassName,
			Source: snapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &pvcName,
			},
		},
	}

	// Create the Snapshot
	snapshotAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshots
	snapshot, createErr := snapshotAPI(nameSpace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})

	if createErr != nil || snapshot == nil {
		return nil, fmt.Errorf("failed to create snapshot: %s,pvc: %s, error: %v", snapshotName, pvcName, createErr)
	}

	return snapshot, createErr
}

// Get the Volume Snapshot
func GetSnapshot(snapshotName string, namespace string) (*snapshotv1.VolumeSnapshot, error) {
	snapshotAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshots
	snapshot, getErr := snapshotAPI(namespace).Get(context.TODO(), snapshotName, metaV1.GetOptions{})
	if getErr != nil {
		return snapshot, fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, getErr)
	}
	return snapshot, getErr
}

// List the Volume Snapshot
func ListSnapshot(namespace string) (*snapshotv1.VolumeSnapshotList, error) {
	snapshotAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshots
	snapshotList, listErr := snapshotAPI(namespace).List(context.TODO(), metaV1.ListOptions{})
	if listErr != nil {
		return snapshotList, fmt.Errorf("failed to list snapshot in namespace %s, error: %v", namespace, listErr)
	}
	return snapshotList, listErr
}

// Delete the volume Snapshot
func DeleteSnapshot(snapshotName string, namespace string) error {
	snapshotAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshots
	delErr := snapshotAPI(namespace).Delete(context.TODO(), snapshotName, metaV1.DeleteOptions{})
	if delErr != nil {
		return fmt.Errorf("failed to delete snapshot: %s, error: %v", snapshotName, delErr)
	}
	return delErr
}

// Create the Snapshot class
func CreateSnapshotClass(snapshotClassName string) error {
	// Snapshotclass create options
	createOpts := &snapshotv1.VolumeSnapshotClass{
		ObjectMeta: metaV1.ObjectMeta{
			Name: snapshotClassName,
			Annotations: map[string]string{
				"snapshot.storage.kubernetes.io/is-default-class": "true",
			},
		},
		Driver:         csiDriver,
		DeletionPolicy: snapshotv1.VolumeSnapshotContentDelete,
	}
	snapshotClassAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotClasses
	_, createErr := snapshotClassAPI().Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	if createErr != nil {
		return fmt.Errorf("failed to create snapshot class: %s, error: %v", snapshotClassName, createErr)
	}
	return createErr
}

// Get the Volume Snapshot class
func GetSnapshotClass(snapshotClassName string) (*snapshotv1.VolumeSnapshotClass, error) {
	snapshotClassAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotClasses
	snapshotClass, getErr := snapshotClassAPI().Get(context.TODO(), snapshotClassName, metaV1.GetOptions{})
	return snapshotClass, getErr
}

// Delete the volume Snapshot class
func DeleteSnapshotClass(snapshotClassName string) error {
	snapshotClassAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotClasses
	delErr := snapshotClassAPI().Delete(context.TODO(), snapshotClassName, metaV1.DeleteOptions{})
	if delErr != nil {
		return fmt.Errorf("failed to delete snapshot class: %s, error: %v", snapshotClassName, delErr)
	}
	return delErr
}

// List the Volume Snapshot class
func ListSnapshotClass() (*snapshotv1.VolumeSnapshotClassList, error) {
	snapshotClassAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotClasses
	snapshotClassList, listErr := snapshotClassAPI().List(context.TODO(), metaV1.ListOptions{})
	if listErr != nil {
		return snapshotClassList, fmt.Errorf("failed to list snapshot class, error: %v", listErr)
	}
	return snapshotClassList, listErr
}

// Get the Volume Snapshot content
func GetSnapshotContent(snapshotContentName string) (*snapshotv1.VolumeSnapshotContent, error) {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	snapshotContent, getErr := snapshotContentAPI().Get(context.TODO(), snapshotContentName, metaV1.GetOptions{})
	if getErr != nil {
		return snapshotContent, fmt.Errorf("failed to get snapshot class: %s, error: %v", snapshotContentName, getErr)
	}
	return snapshotContent, getErr
}

// List the Volume Snapshot content
func ListSnapshotContent() (*snapshotv1.VolumeSnapshotContentList, error) {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	snapshotContentList, getErr := snapshotContentAPI().List(context.TODO(), metaV1.ListOptions{})
	if getErr != nil {
		return snapshotContentList, fmt.Errorf("failed to list snapshot content, error: %v", getErr)
	}
	return snapshotContentList, getErr
}

// Delete the volume Snapshot content
func DeleteSnapshotContent(snapshotContentName string) error {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	delErr := snapshotContentAPI().Delete(context.TODO(), snapshotContentName, metaV1.DeleteOptions{})
	if delErr != nil {
		return fmt.Errorf("failed to delete snapshot content: %s, error: %v", snapshotContentName, delErr)
	}
	return delErr
}

// Get the volume Snapshot content ready status
func GetSnapshotContentReadyStatus(snapshotContentName string) (bool, *snapshotv1.VolumeSnapshotError, error) {
	var snapErorr *snapshotv1.VolumeSnapshotError
	snapshotContent, err := GetSnapshotContent(snapshotContentName)
	if err != nil {
		return false, snapErorr, fmt.Errorf("failed to get snapshot content: %s, error: %v", snapshotContentName, err)
	} else if snapshotContent == nil {
		return false, snapErorr, fmt.Errorf("snapshot content %s not found", snapshotContentName)
	} else if snapshotContent.Status == nil {
		return false, snapErorr, fmt.Errorf("snapshot content %s status not found", snapshotContentName)
	} else if snapshotContent.Status.Error != nil {
		snapErorr = snapshotContent.Status.Error
		return false, snapErorr, nil
	} else if snapshotContent.Status.ReadyToUse == nil {
		return false, snapErorr, fmt.Errorf("snapshot content ready to use not found for snapshot content: %s", snapshotContentName)
	}
	return *snapshotContent.Status.ReadyToUse, snapErorr, err
}

// Get the volume Snapshot  ready status
func GetSnapshotReadyStatus(snapshotName string, namespace string) (bool, *snapshotv1.VolumeSnapshotError, error) {
	var snapErorr *snapshotv1.VolumeSnapshotError
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return false, snapErorr, fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return false, snapErorr, fmt.Errorf("snapshot %s not found", snapshotName)
	} else if snapshot.Status == nil {
		return false, snapErorr, fmt.Errorf("snapshot %s status not found", snapshotName)
	} else if snapshot.Status.Error != nil {
		snapErorr = snapshot.Status.Error
		return false, snapErorr, nil
	} else if snapshot.Status.ReadyToUse == nil {
		return false, snapErorr, fmt.Errorf("snapshot ready to use not found for snapshot: %s", snapshotName)
	}
	return *snapshot.Status.ReadyToUse, snapErorr, err
}

// Get the volume Snapshot bound content name
func GetSnapshotBoundContentName(snapshotName string, namespace string) (string, error) {
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return "", fmt.Errorf("snapshot %s bound snapshot content not found ", snapshotName)
	} else if snapshot.Status == nil {
		return "false", fmt.Errorf("snapshot %s status not found", snapshotName)
	} else if snapshot.Status.BoundVolumeSnapshotContentName == nil {
		return "", fmt.Errorf("snapshot bound snapshot content not found for snapshot: %s", snapshotName)
	}
	return *snapshot.Status.BoundVolumeSnapshotContentName, err
}

// Get the volume Snapshot restore size
func GetSnapshotRestoreSize(snapshotName string, namespace string) (*resource.Quantity, error) {
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return nil, fmt.Errorf("snapshot %s not found ", snapshotName)
	} else if snapshot.Status == nil {
		return nil, fmt.Errorf("snapshot %s status not found", snapshotName)
	} else if snapshot.Status.RestoreSize == nil {
		return nil, fmt.Errorf("snapshot restore size not found for snapshot: %s", snapshotName)
	}
	return snapshot.Status.RestoreSize, err
}

// Get the volumeSnapshot Creation Time
func GetSnapshotCreationTime(snapshotName string, namespace string) (string, error) {
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return "", fmt.Errorf("snapshot %s not found", snapshotName)
	} else if snapshot.Status == nil {
		return "false", fmt.Errorf("snapshot %s status not found", snapshotName)
	} else if snapshot.Status.CreationTime == nil {
		return "", fmt.Errorf("Volume Snapshot Creation Time not found for snapshot: %s", snapshotName)
	}

	if creationTimeStr, ok := snapshot.Status.CreationTime.ToUnstructured().(string); ok {
		logf.Log.Info(" Volume Snapshot", "Creation Time", creationTimeStr)
		return creationTimeStr, err
	} else {
		return "", fmt.Errorf("failed to convert Volume Snapshot Creation Time  %v to string", snapshot.Status.CreationTime.ToUnstructured())
	}
}

// Get the Kubernetes Snapshot "CreationTimestamp"
func GetSnapshotCreationTimeStamp(snapshotName string, namespace string) (string, error) {
	var emptyTime metaV1.Time
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return "", fmt.Errorf("snapshot %s not found", snapshotName)
	} else if snapshot.Status == nil {
		return "false", fmt.Errorf("snapshot %s status not found", snapshotName)
	} else if snapshot.CreationTimestamp == emptyTime {
		return "", fmt.Errorf("snapshot CreationTimeStamp not found for snapshot: %s", snapshotName)
	}

	creationTimestampKubeStr, ok := snapshot.CreationTimestamp.ToUnstructured().(string)
	if ok {
		logf.Log.Info("Snapshot Kubernetes", "creationTimestamp", creationTimestampKubeStr)
		return creationTimestampKubeStr, err
	} else {
		return "", fmt.Errorf("failed to convert snapshot creation time  %v to string", snapshot.CreationTimestamp.ToUnstructured())
	}
}

// Get the volume Snapshot Uid
func GetSnapshotUid(snapshotName string, namespace string) (string, error) {
	snapshot, err := GetSnapshot(snapshotName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot: %s, error: %v", snapshotName, err)
	} else if snapshot == nil {
		return "", fmt.Errorf("snapshot %s not found", snapshotName)
	} else if string(snapshot.ObjectMeta.UID) == "" {
		return "", fmt.Errorf("snapshot Uid not found for snapshot: %s", snapshotName)
	}
	return string(snapshot.ObjectMeta.UID), err
}

// RemoveSnapshot delete a snapshot in namespace and verify that
//  1. The Volume Snapshot is deleted
//  2. The associated Volume Snapshot ContentPV is deleted
func RemoveSnapshot(snapshotName string, namespace string) error {
	const timoSleepSecs = 1
	logf.Log.Info("Removing Volume Snapshot", "snapshot", snapshotName, "namespace ", namespace)
	var isDeleted bool
	t0 := time.Now()

	// Confirm Volume snapshot exist before deleting
	snapshot, getErr := GetSnapshot(snapshotName, namespace)

	if k8serrors.IsNotFound(getErr) {
		return fmt.Errorf("Volume Snapshot %s not found, namespace: %s, error: %v", snapshotName, namespace, getErr)
	} else if getErr != nil {
		return fmt.Errorf("failed to get volume snapshot %s, namespace: %s, error: %v", snapshotName, namespace, getErr)
	} else if snapshot == nil {
		return fmt.Errorf("Volume Snapshot %s not found, namespace: %s", snapshotName, namespace)
	}

	// Delete Volume Snapshot
	deleteErr := DeleteSnapshot(snapshotName, namespace)
	if deleteErr != nil {
		return fmt.Errorf("failed to delete Volume Snapshot %s, namespace: %s, error: %v", snapshotName, namespace, deleteErr)
	}

	logf.Log.Info("Waiting for Snapshot to be deleted", "snapshot", snapshotName)
	// Wait for the Volume Snapshot to be deleted.
	var err error
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		isDeleted, err = IsSnapshotDeleted(snapshotName, namespace)
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
	snapshotContentName := *snapshot.Status.BoundVolumeSnapshotContentName
	if snapshotContentName == "" {
		return fmt.Errorf("Volume Snapshot content not found,sna[shot: %s ,namespace: %s", snapshotName, namespace)
	}
	logf.Log.Info("Waiting for Volume Snapshot Content to be deleted", "snapshot content", snapshotContentName)

	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		isDeleted, err = IsSnapshotContentDeleted(snapshotContentName)
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

// IsSnapshotDeleted return true if Volume snapshot does not exist in cluster
func IsSnapshotDeleted(snapshotName string, namespace string) (bool, error) {
	// Confirm Volume snapshot exist before deleting
	snapshots, err := ListSnapshot(namespace)

	logf.Log.Info("IsSnapshotDeleted", "getErr", err)
	if err != nil {
		return false, fmt.Errorf("failed to list volume snapshots in %s namespace, error: %v", namespace, err)
	} else if len(snapshots.Items) == 0 {
		return true, nil
	}

	for _, snap := range snapshots.Items {
		if snap.Name == snapshotName {
			return false, nil
		}
	}
	return true, nil
}

// IsSnapshotContentDeleted return true if Volume snapshot content does not exist in cluster
func IsSnapshotContentDeleted(snapshotContentName string) (bool, error) {
	// Confirm Volume snapshot exist before deleting
	snapshotsContent, err := ListSnapshotContent()

	if err != nil {
		return false, fmt.Errorf("failed to list volume snapshots content, error: %v", err)
	} else if len(snapshotsContent.Items) == 0 {
		return true, nil
	}

	for _, snap := range snapshotsContent.Items {
		if snap.Name == snapshotContentName {
			return false, nil
		}
	}
	return true, nil
}

// Get the volume Snapshot  ready status from control plane
func GetSnapshotCpCloneReadyStatus(snapshotUid string, volUuid string) (bool, error) {
	snap, err := controlplane.GetVolumeSnapshot(volUuid, snapshotUid)

	if err != nil {
		return false, fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapshotUid, volUuid, err)
	}
	//if snapshot is online then timestamp exist otherwise timestamp will be empty string
	if len(snap.State.ReplicaSnapshots) == 0 {
		logf.Log.Info("Snapshot", "Snapshot.State.ReplicaSnapshots", snap.State.ReplicaSnapshots)
		return false, fmt.Errorf("failed to find replica snapshot, snapshot: %s, volume: %s", snapshotUid, volUuid)
	}

	for _, replicaSnapshot := range snap.State.ReplicaSnapshots {
		if replicaSnapshot.Online.Timestamp == "" {
			return false, err
		}
	}
	return true, err
}

// Get the volume Snapshot  restore size from control plane
func GetSnapshotCpRestoreSize(snapshotUid string, volUuid string) (int64, error) {
	snap, err := controlplane.GetVolumeSnapshot(volUuid, snapshotUid)

	if err != nil {
		return 0, fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapshotUid, volUuid, err)
	}
	return snap.State.AllocatedSize, err
}

// Get the volume Snapshot  timestamp from control plane
func GetSnapshotCpTimestamp(snapshotUid string, volUuid string) (string, error) {
	snap, err := controlplane.GetVolumeSnapshot(volUuid, snapshotUid)

	if err != nil {
		return "", fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapshotUid, volUuid, err)
	}
	return snap.State.Timestamp, err
}

// Get the volume Snapshot  replica snapshots from control plane
func GetSnapshotCpReplicaSnapshots(snapshotUid string, volUuid string) ([]common.ReplicaSnapshot, error) {
	snap, err := controlplane.GetVolumeSnapshot(volUuid, snapshotUid)

	if err != nil {
		return nil, fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapshotUid, volUuid, err)
	}
	return snap.State.ReplicaSnapshots, err
}

// Remove the Volume Snapshot content annotations
func RemoveSnapshotContentAnnotation(snapshotContentName string) error {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	snapshotContent, getErr := snapshotContentAPI().Get(context.TODO(), snapshotContentName, metaV1.GetOptions{})
	if getErr != nil {
		return fmt.Errorf("failed to get snapshot content: %s, error: %v", snapshotContentName, getErr)
	}
	if len(snapshotContent.Annotations) == 0 {
		return nil
	} else {
		snapshotContent.Annotations = map[string]string{}
		_, getErr = snapshotContentAPI().Update(context.TODO(), snapshotContent, metaV1.UpdateOptions{})
	}
	return getErr
}

// Add the Volume Snapshot content annotations
func AddSnapshotContentAnnotation(snapshotContentName, annotationKey, annotationValue string) error {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	snapshotContent, getErr := snapshotContentAPI().Get(context.TODO(), snapshotContentName, metaV1.GetOptions{})
	if getErr != nil {
		return fmt.Errorf("failed to get snapshot content: %s, error: %v", snapshotContentName, getErr)
	}

	snapshotContent.Annotations = map[string]string{annotationKey: annotationValue}
	_, getErr = snapshotContentAPI().Update(context.TODO(), snapshotContent, metaV1.UpdateOptions{})

	return getErr
}

// Get the Volume Snapshot content annotation
func GetSnapshotContentAnnotation(snapshotContentName string) (map[string]string, error) {
	snapshotContentAPI := gTestEnv.CsiInt.SnapshotV1().VolumeSnapshotContents
	snapshotContent, getErr := snapshotContentAPI().Get(context.TODO(), snapshotContentName, metaV1.GetOptions{})
	if getErr != nil {
		return nil, fmt.Errorf("failed to get snapshot content: %s, error: %v", snapshotContentName, getErr)
	}

	return snapshotContent.Annotations, getErr
}

func GetVolumeSnapshotChecksum(volUuid string, snapshotName string, ns string) ([]string, error) {
	var checksums []string

	//get snapshot Uid
	snapUuid, err := GetSnapshotUid(snapshotName, ns)
	if err != nil {
		return checksums, fmt.Errorf("failed to get snapshot Uid for snapshot %s, error: %v", snapshotName, err)
	}
	logf.Log.Info("Snapshot", "Uid", snapUuid)

	snapshot, err := controlplane.GetVolumeSnapshot(volUuid, snapUuid)

	if err != nil {
		return checksums, fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapUuid, volUuid, err)
	} else if snapshot.State.ReplicaSnapshots == nil {
		return checksums, fmt.Errorf("failed to get cp replica snapshot , uid: %s, volume uuid %s, error: %v", snapUuid, volUuid, err)
	}

	logf.Log.Info("Snapshot",
		"txn_id", snapshot.Definition.Metadata.TxnID,
		"uuid", snapshot.Definition.Spec.UUID,
	)

	dsp, err := ListMsPools()
	if err != nil {
		return checksums, fmt.Errorf("failed to list disk pool, error: %v", err)
	}

	// Get IP addresses of nodes
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return checksums, err
	}

	for _, replicaSnapshot := range snapshot.State.ReplicaSnapshots {
		var snapshotNode, snapshotNodeIp, snapshotPool string

		for _, pool := range dsp {
			if replicaSnapshot.Online.PoolID == pool.Name {
				snapshotNode = pool.Spec.Node
				snapshotPool = pool.Name
				break
			}
		}

		if snapshotNode == "" {
			return checksums, fmt.Errorf("failed to find node where the snapshot exists, error: %v", err)
		}

		for _, node := range nodeList {
			if node.NodeName == snapshotNode {
				snapshotNodeIp = node.IPAddress
				break
			}
		}

		if snapshotNodeIp == "" {
			return checksums, fmt.Errorf("failed to node %s IP where snapshot exist, error: %v", snapshotNode, err)
		}

		crc32, err := mayastorclient.ChecksumReplica(snapshotNodeIp, replicaSnapshot.Online.UUID, snapshotPool)
		if err != nil {
			return checksums, fmt.Errorf("failed to retrieve checksum : %v", err)
		}
		checksums = append(checksums, fmt.Sprintf("%v", crc32))
	}

	return checksums, err
}

/*
	func getSnapshotCheckSum(replica replicaInfo) (string, error) {
		var checksum string

		nqn, err := GetNodeNqn(replica.IP)
		if err == nil {
			snapNqn, err := getSnapshotNqn(replica.URI)
			if err != nil {
				checksum = fmt.Sprintf("%s; %v", replica, err)
			} else {
				checksum, err = ChecksumReplica(replica.IP, replica.IP, snapNqn, 10, nqn)
				if err != nil {
					logf.Log.Info("ChecksumSnapshot failed", "error", err)
					// do not return from here because we want to unshare if the original
					// replica URI was a bdev
					checksum = fmt.Sprintf("%s; %v", replica, err)
				}
			}

		} else {
			logf.Log.Info("GetNodeNqn failed", "IP", replica.IP, "error", err)
			checksum = fmt.Sprintf("%s; %v", replica, err)
		}

		// for now ignore unshare errors as we may have successfully retrieved a checksum
		unsErr := mayastorclient.UnshareBdev(replica.IP, replica.UUID)
		if unsErr != nil {
			logf.Log.Info("Unshare bdev failed", "bdev UUID", replica.UUID, "node IP", replica.IP, "error", unsErr)
		}

		return checksum, err
	}
*/
func GetVolumeSnapshotFsCheck(volUuid string, snapshotName string, ns string, fsType common.FileSystemType) ([]string, error) {
	var fschecks []string

	//get snapshot Uid
	snapUuid, err := GetSnapshotUid(snapshotName, ns)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot Uid for snapshot %s, error: %v", snapshotName, err)
	}
	logf.Log.Info("Snapshot", "Uid", snapUuid)

	snapshot, err := controlplane.GetVolumeSnapshot(volUuid, snapUuid)

	if err != nil {
		return nil, fmt.Errorf("failed to get cp snapshot, uid: %s, volume uuid %s, error: %v", snapUuid, volUuid, err)
	} else if snapshot.State.ReplicaSnapshots == nil {
		return nil, fmt.Errorf("failed to get cp replica snapshot , uid: %s, volume uuid %s, error: %v", snapUuid, volUuid, err)
	}

	logf.Log.Info("Snapshot",
		"txn_id", snapshot.Definition.Metadata.TxnID,
		"uuid", snapshot.Definition.Spec.UUID,
	)

	dsps, err := ListMsPools()
	if err != nil {
		return nil, fmt.Errorf("failed to list disk pool, error: %v", err)
	}

	// Get IP addresses of nodes
	nodeList, err := GetIOEngineNodes()
	if err != nil {
		return nil, err
	}

	for _, replicaSnapshot := range snapshot.State.ReplicaSnapshots {
		var snapshotNode, snapshotNodeIp string
		for _, pool := range dsps {
			if replicaSnapshot.Online.PoolID == pool.Name {
				snapshotNode = pool.Spec.Node
				break
			}
		}
		if snapshotNode == "" {
			return nil, fmt.Errorf("failed to find node where snapshot exist, error: %v", err)
		}

		for _, node := range nodeList {
			if node.NodeName == snapshotNode {
				snapshotNodeIp = node.IPAddress
				break
			}
		}
		if snapshotNodeIp == "" {
			return nil, fmt.Errorf("failed to node %s IP where snapshot exist, error: %v", snapshotNode, err)
		}

		// snapshot bdev uuid will snapshot_uuid/txn_id
		bdevUuid := fmt.Sprintf("%s/%s", snapshot.Definition.Spec.UUID, snapshot.Definition.Metadata.TxnID)
		// bdev share via gRPC
		bdevShareUri, err := mayastorclient.ShareBdev(snapshotNodeIp, bdevUuid)
		if err != nil {
			return nil, err
		} else if bdevShareUri == "" {
			return nil, fmt.Errorf("failed to get bdev share uri for bdev with uuid %s", bdevUuid)
		}

		fscheck, fsckErr := getSnapshotFsCheck(
			replicaInfo{
				IP:   snapshotNodeIp,
				URI:  bdevShareUri,
				UUID: bdevUuid,
			}, fsType)
		if fsckErr != nil {
			return nil, fmt.Errorf("%v; %v", fsckErr, err)
		}
		fschecks = append(fschecks, fscheck)
	}

	return fschecks, err
}

func getSnapshotFsCheck(replica replicaInfo, fsType common.FileSystemType) (string, error) {
	var fscheck string

	nqn, err := GetNodeNqn(replica.IP)
	if err == nil {
		snapNqn, err := getSnapshotNqn(replica.URI)
		if err != nil {
			fscheck = fmt.Sprintf("%s; %v", replica, err)
		} else {
			fscheck, err = FsConsistentReplica(replica.IP, replica.IP, snapNqn, 10, nqn, fsType)
			if err != nil {
				logf.Log.Info("SnapshotFsCheck failed", "error", err)
				// do not return from here because we want to unshare if the original
				// replica URI was a bdev
				fscheck = fmt.Sprintf("%s; %v", replica, err)
			}
		}

	} else {
		logf.Log.Info("GetNodeNqn failed", "IP", replica.IP, "error", err)
		fscheck = fmt.Sprintf("%s; %v", replica, err)
	}

	// for now ignore unshare errors as we may have successfully retrieved a checksum
	unsErr := mayastorclient.UnshareBdev(replica.IP, replica.UUID)
	if unsErr != nil {
		logf.Log.Info("Unshare bdev failed", "bdev UUID", replica.UUID, "node IP", replica.IP, "error", unsErr)
	}

	return fscheck, err
}

func getSnapshotNqn(snapshotUri string) (string, error) {
	nqnoffset := strings.Index(snapshotUri, "nqn.")
	if nqnoffset == -1 {
		logf.Log.Info("ChecksumReplica", "Invalid URI", snapshotUri)
		return "", fmt.Errorf("invalid nqn URI %v", snapshotUri)
	}
	return snapshotUri[nqnoffset:], nil
}

func GetMayastorSnapshotScMap() ([]string, error) {
	mayastorSnapshotStorageClasses := make([]string, 0)

	snapScs, err := ListSnapshotClass()
	if err == nil {
		for _, vsc := range snapScs.Items {
			if vsc.Driver == csiDriver {
				mayastorSnapshotStorageClasses = append(mayastorSnapshotStorageClasses, vsc.Name)
			}
		}
	}
	return mayastorSnapshotStorageClasses, err
}
