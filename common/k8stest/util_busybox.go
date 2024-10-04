package k8stest

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var defBusyboxTimeoutSecs = 120 // in seconds

func DeployBusyBoxPod(podName, pvcName string, volType common.VolumeType) error {
	args := []string{"sleep", "10000000"}
	podContainer := coreV1.Container{
		Name:            podName,
		Image:           "busybox",
		ImagePullPolicy: coreV1.PullAlways,
		Args:            args,
	}

	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
			},
		},
	}

	podObj, err := NewPodBuilder(podName).
		WithName(podName).
		WithNamespace(common.NSDefault).
		WithRestartPolicy(coreV1.RestartPolicyNever).
		WithContainer(podContainer).
		WithVolume(volume).
		WithVolumeDeviceOrMount(volType).Build()

	if err != nil {
		return fmt.Errorf("failed to generate busybox pod definition, error: %v", err)
	}
	if podObj == nil {
		return fmt.Errorf("busybox pod definition is nil")
	}

	_, err = CreatePod(podObj, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to create busybox pod, error: %v", err)
	}

	isPodRunning := WaitPodRunning(podName, common.NSDefault, defBusyboxTimeoutSecs)
	if !isPodRunning {
		return fmt.Errorf("failed to create busybox pod, error: %v", err)
	}
	logf.Log.Info(fmt.Sprintf("%s pod is running.", podName))
	return nil
}

func CleanUpBusyboxResources(pods []string, pvcName string) error {
	for _, pod := range pods {
		err := DeletePod(pod, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to delete pod %s err %v", pod, err)
		}

	}
	if pvcName != "" {
		pvc, err := GetPVC(pvcName, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to get pvc %s err %v", pvcName, err)
		}
		err = RemovePVC(pvcName, *pvc.Spec.StorageClassName, common.NSDefault, true)
		if err != nil {
			return fmt.Errorf("failed to delete pvc %s err %v", pvcName, err)
		}
		err = RmStorageClass(*pvc.Spec.StorageClassName)
		if err != nil {
			return fmt.Errorf("failed to delete sc from pvc %s err %v", pvcName, err)
		}
	}

	return nil
}
