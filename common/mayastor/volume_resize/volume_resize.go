package volume_resize

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common/k8stest"
	"k8s.io/apimachinery/pkg/api/resource"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs            = 120 // in seconds
	DefSleepTime              = 30  // in seconds
	DefFioCompletionTime      = 420 // in seconds
	InsufficientStorageError  = "507 Insufficient Storage"
	PvcShrinkErrorSubString   = " Forbidden: field can not be less than previous value"
	SnapshotVolumeResizeError = "Volume can't be resized while it has snapshots, or it's a snapshot restore"
)

// VerifyVolumeResize verify:
// 1. pvc resize
// 2. pv resize
// 3. mayastor volume replicas resize
// 4. mayastor volume resize
// 5. verify volume state
func VerifyVolumeResize(pvcName string, volumeName string, namespace string, volSizeMb int) (bool, error) {
	logf.Log.Info("Verify volume resize status")
	pvc, err := k8stest.GetPVC(pvcName, namespace)
	if err != nil {
		return false, err
	}

	// verify msv replica capacity to new size
	logf.Log.Info("Verify replicas resize status")
	replicaResizeStatus, err := WaitForReplicaResize(volumeName, volSizeMb)
	if err != nil {
		return false, err
	} else if !replicaResizeStatus {
		logf.Log.Info("Mayastor volume replicas not resized", "MSV", volumeName)
		return false, err
	}

	// verify msv capacity to new size
	logf.Log.Info("Verify volume resize status")
	msvResizeStatus, err := WaitForVolumeResize(volumeName, volSizeMb)
	if err != nil {
		return false, err
	} else if !msvResizeStatus {
		logf.Log.Info("Mayastor volume not resized", "MSV", volumeName)
		return false, err
	}

	// verify pv capacity to new size
	logf.Log.Info("Verify pv resize status")
	pvResizeStatus, err := WaitForPvResize(pvc.Spec.VolumeName, volSizeMb)
	if err != nil {
		return false, err
	} else if !pvResizeStatus {
		logf.Log.Info("PV not resized", "PV", pvc.Spec.VolumeName)
		return false, err
	}

	// verify pvc capacity to new size
	logf.Log.Info("Verify pvc resize status")
	pvcResizeStatus, err := WaitForPvcResize(pvcName, namespace, volSizeMb)
	if err != nil {
		return false, err
	} else if !pvcResizeStatus {
		logf.Log.Info("PVC not resized", "PVC", pvcName)
		return false, err
	}

	// verify volume provisioning
	_, err = k8stest.VerifyVolumeProvision(pvcName, namespace)
	if err != nil {
		return false, err
	}

	return true, err
}

func WaitForReplicaResize(volName string, volSizeMb int) (bool, error) {
	const timeSleepSecs = 5
	var err error
	var isAllReplicaResized bool

	// Wait for the mayastor volume to be resized
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		replicas, err := k8stest.GetMsvReplicaTopology(volName)
		if err != nil {
			return false, fmt.Errorf("failed to get replica topology for volume %s, error: %v", volName, err)
		}
		isAllReplicaResized = true
		for _, replica := range replicas {
			if int64(k8stest.GetSizePerUnits(uint64(replica.Usage.Capacity), "MiB")) != int64(volSizeMb) {
				isAllReplicaResized = false
				break
			}
		}
		if isAllReplicaResized {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return isAllReplicaResized, err

}

func WaitForVolumeResize(volName string, volSizeMb int) (bool, error) {
	const timeSleepSecs = 5
	var capacity int64
	var err error

	// Wait for the mayastor volume to be resized
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		capacity, err = k8stest.GetMsvSize(volName)
		if err != nil {
			return false, fmt.Errorf("failed to get msv %s size, error: %v", volName, err)
		}
		if int64(k8stest.GetSizePerUnits(uint64(capacity), "MiB")) == int64(volSizeMb) {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return false, err
}

func WaitForPvResize(pvName string, volSizeMb int) (bool, error) {
	const timeSleepSecs = 5
	var capacity *resource.Quantity
	var err error

	// Wait for the PV to be resized
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		capacity, err = k8stest.GetPvCapacity(pvName)
		if err != nil {
			return false, fmt.Errorf("failed to get pvc %s capacity, error: %v", pvName, err)
		}
		capacityBytes, conversionStatus := capacity.AsInt64()
		if !conversionStatus {
			return false, fmt.Errorf("failed to convert PV capacity into int, pv size: %v", capacity)
		} else if int64(k8stest.GetSizePerUnits(uint64(capacityBytes), "MiB")) == int64(volSizeMb) {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return false, err
}

func WaitForPvcResize(pvcName string, namespace string, volSizeMb int) (bool, error) {
	const timeSleepSecs = 5
	var capacity *resource.Quantity
	var err error

	// Wait for the PV to be resized
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		capacity, err = k8stest.GetPvcCapacity(pvcName, namespace)
		if err != nil {
			return false, fmt.Errorf("failed to get pvc %s capacity, error: %v", pvcName, err)
		}

		capacityBytes, conversionStatus := capacity.AsInt64()
		if !conversionStatus {
			return false, fmt.Errorf("failed to convert PVC capacity into int, pvc size: %v", capacity)
		} else if int64(k8stest.GetSizePerUnits(uint64(capacityBytes), "MiB")) == int64(volSizeMb) {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return false, err
}

func WaitForVolumeUnPublish(volName string) (bool, error) {
	const timeSleepSecs = 2

	// Wait for the nexus to be removed
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		nexusUuid, err := k8stest.GetNexusUuid(volName)
		if err != nil {
			return false, fmt.Errorf("failed to get nexus uuid for volume %s, error: %v", volName, err)
		}
		if nexusUuid == "" {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return false, nil
}

func ByteSizeString(bytesize uint64) string {
	var gib uint64 = 1024 * 1024 * 1024
	var mib uint64 = gib / 1024
	var kib uint64 = mib / 1024

	if bytesize > gib {
		return fmt.Sprintf("%f GiB", float32(bytesize)/float32(gib))
	}

	if bytesize > mib {
		return fmt.Sprintf("%f MiB", float32(bytesize)/float32(mib))
	}

	if bytesize > kib {
		return fmt.Sprintf("%f KiB", float32(bytesize)/float32(kib))
	}

	return fmt.Sprintf("%d bytes", bytesize)
}
