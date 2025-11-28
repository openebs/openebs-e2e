package k8stest

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/storage/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// VolumeAttachmentWaitTimeout is the timeout to wait for volume attachment
	volumeAttachmentWaitTimeout = 5 // in seconds
)

// DeleteVolumeAttachments deletes volume attachments for a node
func DeleteVolumeAttachments(nodeName string) error {
	volumeAttachments, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list volume attachments, error: %v", err)
	}
	if len(volumeAttachments.Items) == 0 {
		return nil
	}
	for _, volumeAttachment := range volumeAttachments.Items {
		if volumeAttachment.Spec.NodeName != nodeName {
			continue
		}
		logf.Log.Info("DeleteVolumeAttachments: Deleting", "volumeAttachment", volumeAttachment.Name)
		delErr := gTestEnv.KubeInt.StorageV1().VolumeAttachments().Delete(context.TODO(), volumeAttachment.Name, metaV1.DeleteOptions{})
		if delErr != nil {
			logf.Log.Info("DeleteVolumeAttachments: failed to delete the volumeAttachment", "volumeAttachment", volumeAttachment.Name, "error", delErr)
			return delErr
		}
	}
	return nil
}

// ListVolumeAttachments list volume attachments
func ListVolumeAttachments() (*v1.VolumeAttachmentList, error) {
	volumeAttachments, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volume attachments, error: %v", err)
	}
	return volumeAttachments, nil
}

// ListVolumeAttachmentsOnNode list volume attachments on a node
func ListVolumeAttachmentsOnNode(nodeName string) (*v1.VolumeAttachmentList, error) {
	volumeAttachments, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volume attachments, error: %v", err)
	}
	if len(volumeAttachments.Items) == 0 {
		return nil, nil
	}
	var filteredVolumeAttachments v1.VolumeAttachmentList
	for _, volumeAttachment := range volumeAttachments.Items {
		if volumeAttachment.Spec.NodeName != nodeName {
			continue
		}
		filteredVolumeAttachments.Items = append(filteredVolumeAttachments.Items, volumeAttachment)
	}
	return &filteredVolumeAttachments, nil
}

// GetVolumeAttachment get volume attachment with source PV name
func GetVolumeAttachment(sourcePvName string) (v1.VolumeAttachment, error) {
	var volumeAttachment v1.VolumeAttachment
	volumeAttachments, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		if k8serror.IsNotFound(err) {
			return volumeAttachment, nil
		}
		return volumeAttachment, fmt.Errorf("failed to get volume attachment with source PV %s, error: %v", sourcePvName, err)
	}
	if len(volumeAttachments.Items) == 0 {
		return volumeAttachment, nil
	}
	for _, va := range volumeAttachments.Items {
		logf.Log.Info("GetVolumeAttachment: checking volume attachment", "volumeAttachment", va.Name, "sourcePV", va.Spec.Source.PersistentVolumeName)
		sourcePvNamePtr := va.Spec.Source.PersistentVolumeName
		if *sourcePvNamePtr == sourcePvName {
			return va, nil
		}
	}
	return volumeAttachment, nil
}

func GetVolumeAttachmentNameForAppPod(podName, namespace string) (string, error) {
	// get pvc name from pod
	pvcName, err := GetPvcNameFromPod(podName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get pvc name from pod %s/%s, error: %v", namespace, podName, err)
	} else if pvcName == "" {
		return "", fmt.Errorf("no pvc found in pod %s/%s", namespace, podName)
	}
	logf.Log.Info("found pvc for pod", "pod", podName, "namespace", namespace, "pvc", pvcName)
	// get pv name from pvc
	pvc, err := GetPVC(pvcName, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get pv name from pvc %s/%s, error: %v", namespace, pvcName, err)
	} else if pvc.Spec.VolumeName == "" {
		return "", fmt.Errorf("no pv found for pvc %s", pvcName)
	}

	logf.Log.Info("found pv for pvc", "pvc", pvcName, "namespace", namespace, "pv", pvc.Spec.VolumeName)
	// get volume attachment from pv name
	volumeAttachment, err := GetVolumeAttachment(pvc.Spec.VolumeName)
	if err != nil {
		return "", fmt.Errorf("failed to get volume attachment for pv %s, error: %v", pvc.Spec.VolumeName, err)
	}
	logf.Log.Info("found volume attachment for pv", "pv", pvc.Spec.VolumeName, "volumeAttachment", volumeAttachment.Name)
	return volumeAttachment.Name, nil
}

func WaitForVolumeAttachment(podName, namespace string, timeoutSecs int) (string, error) {

	for start := time.Now(); time.Since(start) < time.Duration(timeoutSecs)*time.Second; {
		volAttachName, err := GetVolumeAttachmentNameForAppPod(podName, namespace)
		if err != nil {
			logf.Log.Info("WaitForVolumeAttachment: volume attachment not found yet", "pod", podName, "namespace", namespace, "error", err)
			time.Sleep(volumeAttachmentWaitTimeout * time.Second)
			continue
		}
		if volAttachName != "" {
			return volAttachName, nil
		}
		time.Sleep(5 * time.Second)
	}
	return "", fmt.Errorf("volume attachment not found for pod %s/%s within timeout %d seconds", namespace, podName, timeoutSecs)
}

func GetVolumeAttachmentByName(volAttachName string) (v1.VolumeAttachment, error) {
	volumeAttachment, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().Get(context.TODO(), volAttachName, metaV1.GetOptions{})
	return *volumeAttachment, err
}

func CheckVolumeAttachmentDeleted(volAttachName string) (bool, error) {
	volAttachment, err := GetVolumeAttachmentByName(volAttachName)
	if err != nil {
		if k8serror.IsNotFound(err) {
			logf.Log.Info("CheckVolumeAttachmentDeleted: volume attachment not found", "volumeAttachment", volAttachName)
			return true, nil
		}
		return false, fmt.Errorf("failed to get volume attachment %s, error: %v", volAttachName, err)
	}
	logf.Log.Info("CheckVolumeAttachmentDeleted: volume attachment still exists", "volumeAttachment", volAttachName, "status", volAttachment.Status)
	return false, nil
}
