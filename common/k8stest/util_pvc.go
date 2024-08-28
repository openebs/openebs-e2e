package k8stest

// Utility functions for Persistent Volume Claims and Persistent Volumes
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/mayastorclient"
	errors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	coreV1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	defTimeoutSecs           = 360
	insufficientStorageEvent = "error in response: status code '507 Insufficient Storage'"
)

// IsPVCDeleted Check for a deleted Persistent Volume Claim,
// either the object does not exist
// or the status phase is invalid.
func IsPVCDeleted(volName string, nameSpace string) (bool, error) {
	pvc, err := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if err != nil {
		// Unfortunately there is no associated error code, so we resort to string comparison
		if strings.HasPrefix(err.Error(), "persistentvolumeclaims") &&
			strings.HasSuffix(err.Error(), " not found") {
			return true, nil
		} else {
			return false, fmt.Errorf("failed to get pvc %s, namespace: %s, error: %v", volName, nameSpace, err)
		}
	}
	// After the PVC has been deleted it may still accessible, but status phase will be invalid
	if pvc != nil {
		switch pvc.Status.Phase {
		case
			coreV1.ClaimBound,
			coreV1.ClaimPending,
			coreV1.ClaimLost:
			return false, nil
		default:
			return true, nil
		}
	}
	return true, nil
}

// IsPVDeleted Check for a deleted Persistent Volume,
// either the object does not exist
// or the status phase is invalid.
func IsPVDeleted(volName string) (bool, error) {
	pv, err := gTestEnv.KubeInt.CoreV1().PersistentVolumes().Get(context.TODO(), volName, metaV1.GetOptions{})
	if err != nil {
		// Unfortunately there is no associated error code so we resort to string comparison
		if strings.HasPrefix(err.Error(), "persistentvolumes") &&
			strings.HasSuffix(err.Error(), " not found") {
			return true, nil
		} else {
			return false, fmt.Errorf("failed to get pv %s, error: %v", volName, err)
		}
	}

	if pv != nil {
		switch pv.Status.Phase {
		case
			coreV1.VolumeBound,
			coreV1.VolumeAvailable,
			coreV1.VolumeFailed,
			coreV1.VolumePending,
			coreV1.VolumeReleased:
			return false, nil
		default:
			return true, nil
		}
	}
	// After the PV has been deleted it may still accessible, but status phase will be invalid
	return true, nil
}

// IsPvcBound returns true if a PVC with the given name is bound otherwise false is returned.
func IsPvcBound(pvcName string, nameSpace string) (bool, error) {
	phase, err := GetPvcStatusPhase(pvcName, nameSpace)
	if err != nil {
		return false, err
	}
	return phase == coreV1.ClaimBound, nil
}

// GetPvcStatusPhase Retrieve status phase of a Persistent Volume Claim
func GetPvcStatusPhase(volname string, nameSpace string) (phase coreV1.PersistentVolumeClaimPhase, err error) {
	pvc, getPvcErr := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Get(context.TODO(), volname, metaV1.GetOptions{})
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v",
			volname,
			nameSpace,
			getPvcErr,
		)
	}
	if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s",
			volname,
			nameSpace,
		)
	}
	return pvc.Status.Phase, nil
}

// GetPvStatusPhase Retrieve status phase of a Persistent Volume
func GetPvStatusPhase(volname string) (phase coreV1.PersistentVolumePhase, err error) {
	pv, getPvErr := gTestEnv.KubeInt.CoreV1().PersistentVolumes().Get(context.TODO(), volname, metaV1.GetOptions{})
	if getPvErr != nil {
		return "", fmt.Errorf("failed to get pv: %s, error: %v",
			volname,
			getPvErr,
		)
	}
	if pv == nil {
		return "", fmt.Errorf("PV not found: %s", volname)
	}
	return pv.Status.Phase, nil
}

// MkPVC Create a PVC and verify that
//  1. The PVC status transitions to bound,
//  2. The associated PV is created and its status transitions bound
//  3. The associated MV is created and has a State "healthy"
// MkPVC is used by mayastor tests, so local arguments in MakePVC is set to false

// Deprecated:MkPVC is deprecated. Make use of MakePVC function
func MkPVC(volSizeMb int, volName string, scName string, volType common.VolumeType, nameSpace string) (string, error) {
	return MakePVC(volSizeMb, volName, scName, volType, nameSpace, false, false)
}

func VerifyMayastorPvcIsUsable(pvc *coreV1.PersistentVolumeClaim) error {
	const timoSleepSecs = 1
	var err error

	if pvc.Spec.VolumeName == "" {
		return fmt.Errorf("PVC spec.VolumeName is empty")
	}

	// Wait for the PV to be provisioned
	var pv *coreV1.PersistentVolume
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		pv, err = gTestEnv.KubeInt.CoreV1().PersistentVolumes().Get(context.TODO(), pvc.Spec.VolumeName, metaV1.GetOptions{})
		if err == nil && pv != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get pv, pvc.Spec.VolumeName: %s, error: %v", pvc.Spec.VolumeName, err)
	}

	// Wait for the PV to be bound.
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		var pvPhase coreV1.PersistentVolumePhase
		pvPhase, err = GetPvStatusPhase(pv.Name)
		if err == nil && pvPhase == coreV1.VolumeBound {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get pv status, pv: %s, error: %v", pvc.Spec.VolumeName, err)
	}

	// Wait for the PV to be provisioned
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		var msv *common.MayastorVolume
		msv, err = GetMSV(string(pvc.ObjectMeta.UID))
		if err == nil && msv != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get mayastor volume, uuid: %s, error: %v", pvc.ObjectMeta.UID, err)
	}

	err = MsvConsistencyCheck(string(pvc.ObjectMeta.UID))
	if err != nil {
		return fmt.Errorf("msv consistency check failed, msv uuid: %s, error: %v", string(pvc.ObjectMeta.UID), err)
	}

	return nil
}

func VerifyPvcCreateAndFail(
	volSizeMb int,
	volName string,
	scName string,
	volType common.VolumeType,
	nameSpace string,
	timeoutsecs int,
	errorMessage string,
) error {
	const timoSleepSecs = 1
	logf.Log.Info("Creating", "volume", volName, "storageClass", scName, "volume type", volType)
	volSizeMbStr := fmt.Sprintf("%dMi", volSizeMb)
	var err error

	// PVC create options
	createOpts := &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      volName,
			Namespace: nameSpace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes:      []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
			Resources: coreV1.ResourceRequirements{
				Requests: coreV1.ResourceList{
					coreV1.ResourceStorage: resource.MustParse(volSizeMbStr),
				},
			},
		},
	}

	switch volType {
	case common.VolFileSystem:
		var fileSystemVolumeMode = coreV1.PersistentVolumeFilesystem
		createOpts.Spec.VolumeMode = &fileSystemVolumeMode
	case common.VolRawBlock:
		var blockVolumeMode = coreV1.PersistentVolumeBlock
		createOpts.Spec.VolumeMode = &blockVolumeMode
	}

	// Create the PVC.
	PVCApi := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims
	_, createErr := PVCApi(nameSpace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	if createErr != nil {
		return fmt.Errorf("failed to create pvc: %s, error: %v", volName, createErr)
	}

	// Confirm the PVC has been created.
	pvc, getPvcErr := PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if getPvcErr != nil {
		return fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}
	var foundFailureEvent bool
	// Wait for the PVC to have a failure event
	for ix := 0; ix < timeoutsecs/timoSleepSecs; ix++ {

		eventList, listerr := GetPvcEvents(volName, nameSpace)
		if listerr != nil {
			err = fmt.Errorf("failed to get namespace events, pvc: %s, namespace:  %s, error: %v", volName, nameSpace, listerr)
			return err
		}
		for _, event := range eventList.Items {
			if event.Type == "Warning" && strings.Contains(event.Message, errorMessage) {
				logf.Log.Info("Found event", "message", event.Message)
				foundFailureEvent = true
				break
			}
		}
		if foundFailureEvent {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if !foundFailureEvent {
		return fmt.Errorf("timed out waiting for failure event")
	}
	return err
}

// Attempt to create a PV, but expect it to fail
func PVCCreateAndFailCordon(
	volSizeMb int,
	volName string,
	scName string,
	volType common.VolumeType,
	nameSpace string,
	timeoutsecs int,
) error {
	const timoSleepSecs = 1
	logf.Log.Info("Creating", "volume", volName, "storageClass", scName, "volume type", volType)
	volSizeMbStr := fmt.Sprintf("%dMi", volSizeMb)
	var err error

	// PVC create options
	createOpts := &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      volName,
			Namespace: nameSpace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes:      []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
			Resources: coreV1.ResourceRequirements{
				Requests: coreV1.ResourceList{
					coreV1.ResourceStorage: resource.MustParse(volSizeMbStr),
				},
			},
		},
	}

	switch volType {
	case common.VolFileSystem:
		var fileSystemVolumeMode = coreV1.PersistentVolumeFilesystem
		createOpts.Spec.VolumeMode = &fileSystemVolumeMode
	case common.VolRawBlock:
		var blockVolumeMode = coreV1.PersistentVolumeBlock
		createOpts.Spec.VolumeMode = &blockVolumeMode
	}

	// Create the PVC.
	PVCApi := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims
	_, createErr := PVCApi(nameSpace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	if createErr != nil {
		return fmt.Errorf("failed to create pvc: %s, error: %v", volName, createErr)
	}

	// Confirm the PVC has been created.
	pvc, getPvcErr := PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if getPvcErr != nil {
		return fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}
	// Verify that the PVC remains unbound.
	for ix := 0; ix < timeoutsecs/timoSleepSecs; ix++ {
		var pvcPhase coreV1.PersistentVolumeClaimPhase
		pvcPhase, err = GetPvcStatusPhase(volName, nameSpace)
		if err != nil {
			return fmt.Errorf("failed to get the PVC status, error: %s", err.Error())
		}
		if pvcPhase == coreV1.ClaimBound {
			return fmt.Errorf("the PVC %s became bound", volName)
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	// check that a volume hasn't been created
	msvs, msverr := controlplane.ListMsvs()
	_ = RmPVC(volName, scName, common.NSDefault)
	if msverr != nil {
		return fmt.Errorf("failed to list msvs, error: %s", msverr.Error())
	}
	if len(msvs) > 0 {
		err = fmt.Errorf("should fail to create a volume %s, but succeeded, found %d", volName, len(msvs))
		logf.Log.Info("Found volume", "name", volName)
	}
	return err
}

// MsvConsistencyCheck check consistency of  MSV Spec, Status, and associated objects returned by gRPC
func MsvConsistencyCheck(uuid string) error {
	msv, err := GetMSV(uuid)
	if msv == nil {
		return fmt.Errorf("MsvConsistencyCheck: GetMsv: %v, got nil pointer to msv", uuid)
	}
	if err != nil {
		return fmt.Errorf("MsvConsistencyCheck: GetMsv: %v", err)
	}

	replicas, err := GetMsvReplicas(uuid)
	if err != nil {
		return fmt.Errorf("MsvConsistencyCheck: GetMsvReplicas: %v", err)
	}

	if mayastorclient.CanConnect() {
		for _, replica := range replicas {
			gReplicas, err := mayastorclient.FindReplicas(replica.Uuid, GetMayastorNodeIPAddresses())
			if err != nil {
				return fmt.Errorf("failed to find replicas using gRPC %v", err)
			}
			if len(gReplicas) != 1 {
				return fmt.Errorf("invalid set of replicas returned for %v", replica.Uuid)
			}
			for _, gReplica := range gReplicas {
				if int64(gReplica.GetSize())-msv.State.Size < 0 {
					return fmt.Errorf("MsvConsistencyCheck: replica size  %d < msv status size %d", gReplica.GetSize(), msv.State.Size)
				}
			}
		}
		nexus := msv.State.Target
		// The nexus is only present when a volume is mounted by a pod.
		if nexus.Node != "" {
			if msv.Spec.Num_replicas != len(msv.State.Target.Children) {
				return fmt.Errorf("MsvConsistencyCheck: msv spec replica count %d != msv status nexus children %d", msv.Spec.Num_replicas, len(msv.State.Target.Children))
			}
			nexusNodeIp, err := GetNodeIPAddress(nexus.Node)
			if err != nil {
				return fmt.Errorf("MsvConsistencyCheck: failed to resolve nexus node IP address, %v", err)
			}
			grpcNexus, err := mayastorclient.FindNexus(nexus.Uuid, []string{*nexusNodeIp})
			if grpcNexus == nil {
				// if we failed to find the nexus, the error maybe significant
				if err != nil {
					return fmt.Errorf("MsvConsistencyCheck: failed to list nexuses gRPC, %v", err)
				}
				return fmt.Errorf("MsvConsistencyCheck: failed to find nexus gRPC %v", nexus.Uuid)
			} else {
				// if we found the nexus, the error is not relevant but noteworthy
				logf.Log.Info("MsvConsistencyCheck: nexus was found, ignoring", "error", err)
			}
			if (*grpcNexus).GetSize() != uint64(msv.State.Size) {
				return fmt.Errorf("MsvConsistencyCheck: nexus size mismatch msv and grpc")
			}
			if len((*grpcNexus).GetChildren()) != msv.Spec.Num_replicas {
				return fmt.Errorf("MsvConsistencyCheck: msv replica count != grpc nexus children")
			}
			if (*grpcNexus).GetStateString() != msv.State.Target.State {
				return fmt.Errorf("MsvConsistencyCheck: msv nexus state %v != grpc nexus state %v", msv.State.Target.State, (*grpcNexus).GetStateString())
			}
		} else {
			logf.Log.Info("MsvConsistencyCheck nexus unavailable")
		}
	} else {
		logf.Log.Info("MsvConsistencyCheck, gRPC calls to mayastor are not enabled, not checking MSVs using gRPC calls")
	}

	logf.Log.Info("MsvConsistencyCheck OK")
	return nil
}

// RmPVC Delete a PVC in the default namespace and verify that
//  1. The PVC is deleted
//  2. The associated PV is deleted
//  3. The associated MV is deleted
// RmPVC is used by mayastor tests, so local arguments in MakePVC is set to false

// Deprecated:RmPVC is deprecated. Make use of RemovePVC function
func RmPVC(volName string, scName string, nameSpace string) error {
	return RemovePVC(volName, scName, nameSpace, false)
}

// CreatePVC Create a PVC in default namespace, no options and no context
func CreatePVC(pvc *coreV1.PersistentVolumeClaim, nameSpace string) (*coreV1.PersistentVolumeClaim, error) {
	return gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Create(context.TODO(), pvc, metaV1.CreateOptions{})
}

// GetPVC Retrieve a PVC in default namespace, no options and no context
func GetPVC(volName string, nameSpace string) (*coreV1.PersistentVolumeClaim, error) {
	return gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
}

// ListPVCs retrieves pvc list from a given namespace.
func ListPVCs(nameSpace string) (*coreV1.PersistentVolumeClaimList, error) {
	pvcs, err := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list PVCs, error %v", err)
	}
	return pvcs, err
}

// DeletePVC Delete a PVC in default namespace, no options and no context
func DeletePVC(volName string, nameSpace string) error {
	return gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Delete(context.TODO(), volName, metaV1.DeleteOptions{})
}

// GetPV Retrieve a PV in default namespace, no options and no context
func GetPV(volName string) (*coreV1.PersistentVolume, error) {
	return gTestEnv.KubeInt.CoreV1().PersistentVolumes().Get(context.TODO(), volName, metaV1.GetOptions{})
}

func getMayastorScMap() (map[string]bool, error) {
	mayastorStorageClasses := make(map[string]bool)
	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	scs, err := ScApi().List(context.TODO(), metaV1.ListOptions{})
	if err == nil {
		for _, sc := range scs.Items {
			if sc.Provisioner == e2e_config.GetConfig().Product.CsiProvisioner {
				mayastorStorageClasses[sc.Name] = true
			}
		}
	}
	return mayastorStorageClasses, err
}

func CheckForPVCs() (bool, error) {
	logf.Log.Info("CheckForPVCs")
	foundResources := false

	mayastorStorageClasses, err := getMayastorScMap()
	if err != nil {
		return false, err
	}

	nameSpaces, err := gTestEnv.KubeInt.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err == nil {
		for _, ns := range nameSpaces.Items {
			if strings.HasPrefix(ns.Name, common.NSE2EPrefix) || ns.Name == common.NSDefault {
				pvcs, err := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(ns.Name).List(context.TODO(), metaV1.ListOptions{})
				if err == nil && pvcs != nil && len(pvcs.Items) != 0 {
					for _, pvc := range pvcs.Items {
						if !mayastorStorageClasses[*pvc.Spec.StorageClassName] {
							continue
						}
						logf.Log.Info("CheckForVolumeResources: found PersistentVolumeClaims",
							"PersistentVolumeClaim", pvc)
						foundResources = true
					}
				}
			}
		}
	}

	return foundResources, err
}

func CheckForPVs() (bool, error) {
	logf.Log.Info("CheckForPVs")
	foundResources := false

	mayastorStorageClasses, err := getMayastorScMap()
	if err != nil {
		return false, err
	}

	pvs, err := gTestEnv.KubeInt.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err == nil && pvs != nil && len(pvs.Items) != 0 {
		for _, pv := range pvs.Items {
			if !mayastorStorageClasses[pv.Spec.StorageClassName] {
				continue
			}
			logf.Log.Info("CheckForVolumeResources: found PersistentVolumes",
				"PersistentVolume", pv)
			foundResources = true
		}
	}
	return foundResources, err
}

func CreatePvc(createOpts *coreV1.PersistentVolumeClaim, errBuf *error, uuid *string, wg *sync.WaitGroup) {
	// Create the PVC.
	pvc, err := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(createOpts.ObjectMeta.Namespace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	*errBuf = err
	if pvc != nil {
		*uuid = string(pvc.UID)
	}
	wg.Done()
}

func DeletePvc(volName string, namespace string, errBuf *error, wg *sync.WaitGroup) {
	// Delete the PVC.
	err := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), volName, metaV1.DeleteOptions{})
	*errBuf = err
	wg.Done()
}

// RemovePVFinalizer remove the finalizer of a given PV
func RemovePVFinalizer(pvName string) error {
	var patch, prevData, newData []byte
	var pv *coreV1.PersistentVolume
	var err error

	logf.Log.Info("RemovePVFinalizer", "pv", pvName)
	pv, err = gTestEnv.KubeInt.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	pvClone := pv.DeepCopy()
	prevData, err = json.Marshal(pvClone)
	if err != nil {
		return err
	}

	pvClone.ObjectMeta.Finalizers = nil
	newData, err = json.Marshal(pvClone)
	if err != nil {
		return err
	}

	patch, err = strategicpatch.CreateTwoWayMergePatch(prevData, newData, pvClone)
	if err != nil {
		return err
	}

	_, err = gTestEnv.KubeInt.CoreV1().PersistentVolumes().Patch(context.TODO(), pvName, types.StrategicMergePatchType, patch, metaV1.PatchOptions{})
	logf.Log.Info("RemovePVFinalizer done", "pv", pvName)
	return err
}

// KillPV destroy PV with extreme prejudice
func KillPV(pvName string) error {
	var err error

	logf.Log.Info("KillPV", "pv", pvName)
	if err = gTestEnv.KubeInt.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metaV1.DeleteOptions{GracePeriodSeconds: &ZeroInt64}); err == nil {
		err = RemovePVFinalizer(pvName)
	}

	if k8serrors.IsNotFound(err) {
		err = nil
	}

	logf.Log.Info("KillPV done", "pv", pvName, "error", err)

	return err
}

func TryMkOversizedPVC(volSizeMb int, volName string, scName string, volType common.VolumeType, nameSpace string) (bool, error) {
	const timoSleepSecs = 5
	var foundOversizeError bool
	logf.Log.Info("Creating", "volume", volName, "storageClass", scName, "volume type", volType)
	volSizeMbStr := fmt.Sprintf("%dMi", volSizeMb)
	var err error

	// PVC create options
	createOpts := &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      volName,
			Namespace: nameSpace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes:      []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
			Resources: coreV1.ResourceRequirements{
				Requests: coreV1.ResourceList{
					coreV1.ResourceStorage: resource.MustParse(volSizeMbStr),
				},
			},
		},
	}

	switch volType {
	case common.VolFileSystem:
		var fileSystemVolumeMode = coreV1.PersistentVolumeFilesystem
		createOpts.Spec.VolumeMode = &fileSystemVolumeMode
	case common.VolRawBlock:
		var blockVolumeMode = coreV1.PersistentVolumeBlock
		createOpts.Spec.VolumeMode = &blockVolumeMode
	}

	// Create the PVC.
	PVCApi := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims
	_, createErr := PVCApi(nameSpace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	if createErr != nil {
		return foundOversizeError, fmt.Errorf("failed to create pvc: %s, error: %v", volName, createErr)
	}

	// Confirm the PVC has been created.
	pvc, getPvcErr := PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if getPvcErr != nil {
		return foundOversizeError, fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return foundOversizeError, fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}

	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	_, getScErr := ScApi().Get(context.TODO(), scName, metaV1.GetOptions{})
	if getScErr != nil {
		return foundOversizeError, fmt.Errorf("failed to get storageclass: %s, error: %v", scName, getScErr)
	}

	// Wait for the PVC to be bound.
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		var pvcPhase coreV1.PersistentVolumeClaimPhase
		pvcPhase, err = GetPvcStatusPhase(volName, nameSpace)
		if err == nil && pvcPhase == coreV1.ClaimPending {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return foundOversizeError, fmt.Errorf("failed to get pvc status, pvc: %s, namespace:  %s, error: %v", volName, nameSpace, err)
	}

	// Refresh the PVC contents, so that we can get the PV name.
	pvc, getPvcErr = PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if getPvcErr != nil {
		return foundOversizeError, fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return foundOversizeError, fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}

	// Wait for the PVC to have a failure event
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		options := metaV1.ListOptions{
			TypeMeta: metaV1.TypeMeta{Kind: "PersistentVolumeClaim"},
		}
		eventList, listerr := GetEvents(nameSpace, options)
		if listerr != nil {
			err = fmt.Errorf("failed to get namespace events, pvc: %s, namespace:  %s, error: %v", volName, nameSpace, listerr)
			_ = RmPVC(volName, scName, nameSpace)
			return foundOversizeError, err
		}
		for _, event := range eventList.Items {
			if event.Type == "Warning" && strings.Contains(event.Message, insufficientStorageEvent) {
				logf.Log.Info("Found event", "message", event.Message)
				foundOversizeError = true
			}
		}
		if foundOversizeError {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	return foundOversizeError, nil
}

type Pvc struct {
	object *coreV1.PersistentVolumeClaim
}

// PvcBuilder enables building an instance of PersistentVolumeClaim
type PvcBuilder struct {
	pvc  *Pvc
	errs []error
}

// NewPvcBuilder returns new instance of PvcBuilder
func NewPvcBuilder() *PvcBuilder {
	obj := PvcBuilder{pvc: &Pvc{object: &coreV1.PersistentVolumeClaim{}}}
	// set default ReadWriteOnce access mode
	pvcObject := obj.WithAccessModes(coreV1.ReadWriteOnce)
	return pvcObject
}

// WithAccessModes sets the accessModes field of pvc with provided argument.
func (b *PvcBuilder) WithAccessModes(mode coreV1.PersistentVolumeAccessMode) *PvcBuilder {
	if len(mode) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing pvc accessModes"))
		return b
	}
	b.pvc.object.Spec.AccessModes = append(b.pvc.object.Spec.AccessModes, mode)
	return b
}

// WithName sets the Name field of pvc with provided argument.
func (b *PvcBuilder) WithName(name string) *PvcBuilder {
	if len(name) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing pvc name"))
		return b
	}
	b.pvc.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of pvc with provided argument.
func (b *PvcBuilder) WithNamespace(namespace string) *PvcBuilder {
	if len(namespace) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing pvc namespace"))
		return b
	}
	b.pvc.object.Namespace = namespace
	return b
}

// WithStorageClass sets the storage class name field of pvc with provided argument.
func (b *PvcBuilder) WithStorageClass(scName string) *PvcBuilder {
	if len(scName) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing pvc sc name"))
		return b
	}
	b.pvc.object.Spec.StorageClassName = &scName
	return b
}

// WithPvcSize sets the size field of pvc with provided argument.
func (b *PvcBuilder) WithPvcSize(size string) *PvcBuilder {
	if len(size) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing pvc size"))
		return b
	}
	pvcSize := resource.MustParse(size)
	b.pvc.object.Spec.Resources.Requests = coreV1.ResourceList{
		coreV1.ResourceStorage: pvcSize,
	}
	return b
}

// WithDataSourceName sets the data source name field of pvc with provided argument.
func (b *PvcBuilder) WithDataSourceName(dataSourceName string) *PvcBuilder {
	if len(dataSourceName) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing data source name"))
		return b
	}
	if b.pvc.object.Spec.DataSource == nil {
		b.pvc.object.Spec.DataSource = &coreV1.TypedLocalObjectReference{}
	}
	b.pvc.object.Spec.DataSource.Name = dataSourceName
	return b
}

// WithDataSourceKind sets the data source kind field of pvc with provided argument.
func (b *PvcBuilder) WithDataSourceKind(dataSourceKind string) *PvcBuilder {
	if len(dataSourceKind) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing data source kind"))
		return b
	}
	if b.pvc.object.Spec.DataSource == nil {
		b.pvc.object.Spec.DataSource = &coreV1.TypedLocalObjectReference{}
	}
	b.pvc.object.Spec.DataSource.Kind = dataSourceKind
	return b
}

// WithDataSourceApiGroup sets the data source aoi group field of pvc with provided argument.
func (b *PvcBuilder) WithDataSourceApiGroup(dataSourceApiGroup string) *PvcBuilder {
	if len(dataSourceApiGroup) == 0 {
		b.errs = append(b.errs, errors.New("failed to build pvc: missing data source api group"))
		return b
	}
	if b.pvc.object.Spec.DataSource == nil {
		b.pvc.object.Spec.DataSource = &coreV1.TypedLocalObjectReference{}
	}
	b.pvc.object.Spec.DataSource.APIGroup = &dataSourceApiGroup
	return b
}

// WithStorageClass sets the storage class name field of pvc with provided argument.
func (b *PvcBuilder) WithVolumeMode(volType common.VolumeType) *PvcBuilder {
	var volMode coreV1.PersistentVolumeMode
	switch volType {
	case common.VolFileSystem:
		volMode = coreV1.PersistentVolumeFilesystem
	case common.VolRawBlock:
		volMode = coreV1.PersistentVolumeBlock
	default:
		volMode = coreV1.PersistentVolumeFilesystem
	}
	b.pvc.object.Spec.VolumeMode = &volMode
	return b
}

// Build returns the StorageClass API instance
func (b *PvcBuilder) Build() (*coreV1.PersistentVolumeClaim, error) {
	if len(b.errs) > 0 {
		return nil, errors.Errorf("%+v", b.errs)
	}

	return b.pvc.object, nil
}

// Build and create the StorageClass
func (b *PvcBuilder) BuildAndCreate() error {
	pvcObj, err := b.Build()
	if err == nil {
		_, err = CreatePVC(pvcObj, b.pvc.object.Namespace)
	}
	return err
}

// MkRestorePVC Create a PVC from snapshot source and verify that
//  1. The PVC status transitions to bound,
//  2. The associated PV is created and its status transitions bound
//  3. The associated MV is created and has a State "healthy"
func MkRestorePVC(pvcSizeMb int, pvcName string, scName string, nameSpace string, volType common.VolumeType, snapshotName string, skipVerification bool) (string, error) {
	snap, err := GetSnapshot(snapshotName, nameSpace)
	if err != nil {
		return "", fmt.Errorf("failed to get snapshot %s, error: %v", snapshotName, err)
	} else if snap == nil {
		return "", fmt.Errorf("snapshot %s not found, error: %v", snapshotName, err)
	}
	pvcSizeMbStr := fmt.Sprintf("%dMi", pvcSizeMb)
	logf.Log.Info("Creating", "pvc", pvcName, "storageClass", scName, "size", pvcSizeMbStr, "snapshot", snapshotName)
	// PVC create options
	snapshotKind := "VolumeSnapshot"
	snapshotGroup := snapshotv1.GroupName
	err = NewPvcBuilder().
		WithName(pvcName).
		WithNamespace(nameSpace).
		WithStorageClass(scName).
		WithPvcSize(pvcSizeMbStr).
		WithDataSourceApiGroup(snapshotGroup).
		WithDataSourceName(snapshotName).
		WithDataSourceKind(snapshotKind).
		WithVolumeMode(volType).
		BuildAndCreate()
	if err != nil {
		return "", err
	}

	// if skipVerification is set to True then don't verify pv amd mayastor volume creation
	// return PVC volume UUID and nil error
	if skipVerification {
		pvc, getPvcErr := GetPVC(pvcName, nameSpace)
		if getPvcErr != nil {
			return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, nameSpace, getPvcErr)
		} else if pvc == nil {
			return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, nameSpace)
		}
		return string(pvc.ObjectMeta.UID), nil
	}

	return VerifyMayastorVolumeProvision(pvcName, nameSpace)
}

// Wait for the PVC to bound
func WaitPvcToBound(pvcName string, nameSpace string) error {
	const timoSleepSecs = 1
	var pvcPhase coreV1.PersistentVolumeClaimPhase
	var err error
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		pvcPhase, err = GetPvcStatusPhase(pvcName, nameSpace)
		if err == nil && pvcPhase == coreV1.ClaimBound {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get pvc, pvc: %s, namespace:  %s, error: %v", pvcName, nameSpace, err)
	} else if pvcPhase != coreV1.ClaimBound {
		return fmt.Errorf("pvc not bound, pvc: %s, state:  %s", pvcName, pvcPhase)
	}
	return err
}

// Refresh the PVC contents, so that we can get the PV name.
func RefreshPvcToGetPvName(pvcName string, nameSpace string) (string, error) {
	const timoSleepSecs = 1
	var getPvcErr error
	var pvc *coreV1.PersistentVolumeClaim
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs && pvcName != ""; ix++ {
		pvc, getPvcErr = GetPVC(pvcName, nameSpace)
		if getPvcErr != nil {
			return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, nameSpace, getPvcErr)
		} else if pvc == nil {
			return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, nameSpace)
		}
		if pvc.Spec.VolumeName != "" {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if pvc == nil {
		return "", fmt.Errorf("PVC is nil, pvc: %s", pvcName)
	} else if pvc.Spec.VolumeName == "" {
		return "", fmt.Errorf("PVC spec.VolumeName is empty for %s", pvcName)
	}
	return pvc.Spec.VolumeName, getPvcErr
}

// Wait for the PV to be provisioned
func WaitForPvToProvision(pvName string) error {
	const timoSleepSecs = 1
	var pv *coreV1.PersistentVolume
	var err error
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		pv, err = GetPV(pvName)
		if err == nil && pv != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get pv, pvc.Spec.VolumeName: %s, error: %v", pvName, err)
	} else if pv == nil {
		return fmt.Errorf("PV is nil, pv: %s", pvName)
	}
	return err
}

// Wait for the PV to bound
func WaitPvToBound(pvName string) error {
	const timoSleepSecs = 1
	var pvPhase coreV1.PersistentVolumePhase
	var err error
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		pvPhase, err = GetPvStatusPhase(pvName)
		if err == nil && pvPhase == coreV1.VolumeBound {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get pv status, pv: %s, error: %v", pvName, err)
	} else if pvPhase != coreV1.VolumeBound {
		return fmt.Errorf("pv not bound, pv: %s, state:  %s", pvName, pvPhase)
	}
	return err
}

// Wait for the mayastor volume to be provisioned
func WaitForMayastorVolumeToProvision(msvName string) (*common.MayastorVolume, error) {
	var msv *common.MayastorVolume
	var err error
	const timoSleepSecs = 1
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		msv, err = GetMSV(msvName)
		if err == nil && msv != nil {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mayastor volume, uuid: %s, error: %v", msvName, err)
	} else if msv == nil {
		return nil, fmt.Errorf("mayastor volume not found, uuid: %s", msvName)
	}
	return msv, err
}

// Verify volume provisioning, return mayastor volume uuid and error
func VerifyMayastorVolumeProvision(pvcName string, namespace string) (string, error) {
	pvc, getPvcErr := GetPVC(pvcName, namespace)
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, namespace, getPvcErr)
	} else if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, namespace)
	}

	err := WaitPvcToBound(pvcName, namespace)
	if err != nil {
		return "", err
	}
	pvName, err := RefreshPvcToGetPvName(pvcName, namespace)
	if err != nil {
		return "", err
	}
	err = WaitForPvToProvision(pvName)
	if err != nil {
		return "", err
	}
	err = WaitPvToBound(pvName)
	if err != nil {
		return "", err
	}
	pvc, getPvcErr = GetPVC(pvcName, namespace)
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, namespace, getPvcErr)
	} else if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, namespace)
	}
	msv, err := WaitForMayastorVolumeToProvision(string(pvc.ObjectMeta.UID))
	if err != nil {
		return "", err
	}
	err = MsvConsistencyCheck(msv.Spec.Uuid)
	if err != nil {
		return "", fmt.Errorf("msv consistency check failed, msv uuid: %s, error: %v", msv.Spec.Uuid, err)
	}

	logf.Log.Info("Created", "volume", pvcName, "uuid", msv.Spec.Uuid, "storageClass", pvc.Spec.StorageClassName, "size", pvc.Size)
	return msv.Spec.Uuid, err
}

// VerifyVolumeProvision Verify volume provisioning, return volume uuid and error
func VerifyVolumeProvision(pvcName string, namespace string) (string, error) {
	pvc, getPvcErr := GetPVC(pvcName, namespace)
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, namespace, getPvcErr)
	} else if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, namespace)
	}

	err := WaitPvcToBound(pvcName, namespace)
	if err != nil {
		return "", err
	}
	pvName, err := RefreshPvcToGetPvName(pvcName, namespace)
	if err != nil {
		return "", err
	}
	err = WaitForPvToProvision(pvName)
	if err != nil {
		return "", err
	}
	err = WaitPvToBound(pvName)
	if err != nil {
		return "", err
	}
	pvc, getPvcErr = GetPVC(pvcName, namespace)
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", pvcName, namespace, getPvcErr)
	} else if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s", pvcName, namespace)
	}
	volUuid := fmt.Sprintf("%v", pvc.UID)
	logf.Log.Info("Created", "volume", pvcName, "uuid", volUuid, "storageClass", pvc.Spec.StorageClassName, "size", pvc.Size)
	return volUuid, err
}

func GetPvcEvents(pvcName string, namespace string) (*coreV1.EventList, error) {
	options := metaV1.ListOptions{
		TypeMeta:      metaV1.TypeMeta{Kind: "PersistentVolumeClaim"},
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", pvcName),
	}
	return GetEvents(namespace, options)
}

func CheckMsvIsDeleted(uuid string) error {
	const timoSleepSecs = 1
	isDeleted := false
	for i := 0; i < defTimeoutSecs/timoSleepSecs; i++ {
		isDeleted = IsMsvDeleted(uuid)
		if isDeleted {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if !isDeleted {
		return fmt.Errorf("mayastor volume not deleted, msv: %s", uuid)
	}
	return nil
}

// UpdatePvcSize updates the pvc resource request size
func UpdatePvcSize(pvcName string, namespace string, volSizeMb int) (*coreV1.PersistentVolumeClaim, error) {
	volSizeMbStr := fmt.Sprintf("%dMi", volSizeMb)
	var pvc *coreV1.PersistentVolumeClaim
	var err error

	// try again in case of conflict while updating PVC
	// maximum retry for five times
	retryCount := 5
	for i := 0; i < retryCount; i++ {
		pvc, err = GetPVC(pvcName, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get pvc %s, error: %v", pvcName, err)
		} else if *pvc.Status.Capacity.Storage() == resource.MustParse(volSizeMbStr) {
			logf.Log.Info("PVC size same as new requested size", "PVC", pvcName, "current size",
				pvc.Status.Capacity.Storage(),
				"requested size", volSizeMbStr)
			return pvc, nil
		}

		pvc.Spec.Resources.Requests = coreV1.ResourceList{
			coreV1.ResourceStorage: resource.MustParse(volSizeMbStr),
		}
		pvc, err = UpdatePVC(pvc, namespace)

		//https://github.com/kubernetes/kubernetes/blob/release-1.1/docs/devel/api-conventions.md#concurrency-control-and-consistency
		if k8serrors.IsConflict(err) {
			logf.Log.Info("PVC update conflict error", "error", err)
			logf.Log.Info("Sleep for 2 seconds before next retry")
			time.Sleep(2 * time.Second)
			continue
		} else {
			break
		}
	}

	return pvc, err
}

// UpdatePVC update a PVC in given namespace, no options and no context
func UpdatePVC(pvc *coreV1.PersistentVolumeClaim, nameSpace string) (*coreV1.PersistentVolumeClaim, error) {
	return gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims(nameSpace).Update(context.TODO(), pvc, metaV1.UpdateOptions{})
}

// GetPvcCapacity Retrieve status capacity of a Persistent Volume Claim
func GetPvcCapacity(volname string, namespace string) (*resource.Quantity, error) {
	var capacity *resource.Quantity
	pvc, err := GetPVC(volname, namespace)
	if err != nil {
		return capacity, fmt.Errorf("failed to get PVC %s, error: %v", volname, err)
	}

	if pvc == nil {
		return capacity, fmt.Errorf("PVC not found: %s", volname)
	}
	return pvc.Status.Capacity.Storage(), nil
}

// GetPvCapacity Retrieve status capacity of a Persistent Volume
func GetPvCapacity(pvName string) (*resource.Quantity, error) {
	var capacity *resource.Quantity
	pv, err := GetPV(pvName)
	if err != nil {
		return capacity, fmt.Errorf("failed to get PV %s, error: %v", pvName, err)
	}

	if pv == nil {
		return capacity, fmt.Errorf("PV not found: %s", pvName)
	}
	return pv.Spec.Capacity.Storage(), nil
}

// MakePVC Create a PVC and verify that
//  1. The PVC status transitions to bound,
//  2. The associated PV is created and its status transitions bound
//  3. The associated maaystor volume is created and has a State "healthy" if engine is mayastor
func MakePVC(volSizeMb int, volName string, scName string, volType common.VolumeType, nameSpace string, local bool, skipVolumeVerification bool) (string, error) {
	volSizeMbStr := fmt.Sprintf("%dMi", volSizeMb)
	logf.Log.Info("Creating", "volume", volName, "storageClass", scName, "volume type", volType, "size", volSizeMbStr)

	t0 := time.Now()
	// PVC create options
	createOpts := &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      volName,
			Namespace: nameSpace,
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes:      []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
			Resources: coreV1.ResourceRequirements{
				Requests: coreV1.ResourceList{
					coreV1.ResourceStorage: resource.MustParse(volSizeMbStr),
				},
			},
		},
	}

	switch volType {
	case common.VolFileSystem:
		var fileSystemVolumeMode = coreV1.PersistentVolumeFilesystem
		createOpts.Spec.VolumeMode = &fileSystemVolumeMode
	case common.VolRawBlock:
		var blockVolumeMode = coreV1.PersistentVolumeBlock
		createOpts.Spec.VolumeMode = &blockVolumeMode
	}

	// Create the PVC.
	PVCApi := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims
	_, createErr := PVCApi(nameSpace).Create(context.TODO(), createOpts, metaV1.CreateOptions{})
	if createErr != nil {
		return "", fmt.Errorf("failed to create pvc: %s, error: %v", volName, createErr)
	}

	// Confirm the PVC has been created.
	pvc, getPvcErr := PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if getPvcErr != nil {
		return "", fmt.Errorf("failed to get pvc: %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return "", fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}

	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	sc, getScErr := ScApi().Get(context.TODO(), scName, metaV1.GetOptions{})
	if getScErr != nil {
		return "", fmt.Errorf("failed to get storageclass: %s, error: %v", scName, getScErr)
	}
	if *sc.VolumeBindingMode == storagev1.VolumeBindingWaitForFirstConsumer {
		return string(pvc.ObjectMeta.UID), nil
	}

	if !skipVolumeVerification {
		// verify volume provision
		uuid, err := VerifyVolumeProvision(volName, nameSpace)
		if err != nil {
			return string(pvc.ObjectMeta.UID), err
		}
		if !local {
			uuid, err = VerifyMayastorVolumeProvision(volName, nameSpace)
			if err != nil {
				return string(pvc.ObjectMeta.UID), err
			}
		}
		logf.Log.Info("Created", "volume", volName, "uuid", pvc.ObjectMeta.UID, "storageClass", scName, "volume type", volType, "size", volSizeMbStr, "elapsed time", time.Since(t0))
		return uuid, nil
	}
	logf.Log.Info("Created", "volume", volName, "uuid", pvc.ObjectMeta.UID, "storageClass", scName, "volume type", volType, "size", volSizeMbStr, "elapsed time", time.Since(t0))
	return string(pvc.ObjectMeta.UID), nil
}

// RmLocalPVC Delete a PVC in the default namespace and verify that
//  1. The PVC is deleted
//  2. The associated PV is deleted
//  3. The associated mayastor volume is deleted if engine is mayastor
func RemovePVC(volName string, scName string, nameSpace string, local bool) error {
	const timoSleepSecs = 1
	logf.Log.Info("Removing volume", "volume", volName, "storageClass", scName)
	var isDeleted bool

	PVCApi := gTestEnv.KubeInt.CoreV1().PersistentVolumeClaims

	// Confirm the PVC has been deleted.
	pvc, getPvcErr := PVCApi(nameSpace).Get(context.TODO(), volName, metaV1.GetOptions{})
	if k8serrors.IsNotFound(getPvcErr) {
		return fmt.Errorf("PVC %s not found error, namespace: %s", volName, nameSpace)
	} else if getPvcErr != nil {
		return fmt.Errorf("failed to get pvc %s, namespace: %s, error: %v", volName, nameSpace, getPvcErr)
	} else if pvc == nil {
		return fmt.Errorf("PVC %s not found, namespace: %s", volName, nameSpace)
	}
	// Delete the PVC
	deleteErr := PVCApi(nameSpace).Delete(context.TODO(), volName, metaV1.DeleteOptions{})
	if deleteErr != nil {
		return fmt.Errorf("failed to delete PVC %s, namespace: %s, error: %v", volName, nameSpace, deleteErr)
	}

	logf.Log.Info("Waiting for PVC to be deleted", "volume", volName, "storageClass", scName)
	var err error
	// Wait for the PVC to be deleted.
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		isDeleted, err = IsPVCDeleted(volName, nameSpace)
		if isDeleted {
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return err
	} else if !isDeleted {
		return fmt.Errorf("pvc not deleted, pvc: %s, namespace: %s", volName, nameSpace)
	}

	// Wait for the PV to be deleted.
	logf.Log.Info("Waiting for PV to be deleted", "volume", volName, "storageClass", scName)
	for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
		// This check is required here because it will check for pv name
		// when pvc is in pending state at that time we will not
		// get pv name inside pvc spec i.e pvc.Spec.VolumeName
		if pvc.Spec.VolumeName != "" {
			isDeleted, err = IsPVDeleted(pvc.Spec.VolumeName)
			if isDeleted {
				break
			}
		} else {
			isDeleted = true
			break
		}
		time.Sleep(timoSleepSecs * time.Second)
	}
	if err != nil {
		return err
	} else if !isDeleted {
		return fmt.Errorf("PV not deleted, pv: %s", pvc.Spec.VolumeName)
	}

	// if it's replicated engine(mayastor), verify mayastor volume deletion
	if !local {
		// Wait for the mayastor to be deleted.
		for ix := 0; ix < defTimeoutSecs/timoSleepSecs; ix++ {
			isDeleted = IsMsvDeleted(string(pvc.ObjectMeta.UID))
			if isDeleted {
				break
			}
			time.Sleep(timoSleepSecs * time.Second)
		}
		if !isDeleted {
			return fmt.Errorf("mayastor volume not deleted, msv: %s", pvc.ObjectMeta.UID)
		}
	}
	return nil
}

func IsNormalPvcEventPresent(pvcName string, namespace string, errorSubstring string) (bool, error) {
	events, err := GetPvcEvents(pvcName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get pvc events in namespace %s, error: %v", namespace, err)
	}
	for _, event := range events.Items {
		logf.Log.Info("Found PVC event", "message", event.Message)
		if event.Type == "Normal" && strings.Contains(event.Message, errorSubstring) {
			return true, err
		}
	}
	return false, err
}

// Wait for the PVC normal event
// return true if pvc normal event found else return false , it return error in case of error
func WaitForPvcNormalEvent(pvcName string, namespace string, errorSubstring string) (bool, error) {
	const timeSleepSecs = 2
	var hasWarning bool
	var err error
	// Wait for the PVC event
	logf.Log.Info("Check Pvc event", "Event substring", errorSubstring)
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		hasWarning, err = IsNormalPvcEventPresent(pvcName, namespace, errorSubstring)
		if err == nil && hasWarning {
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return hasWarning, err
}
