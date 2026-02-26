package k8stest

import (
	"context"
	"fmt"
	"time"

	batchV1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const sleepInterval = 1 // seconds

// get kubernetes Job
func GetJob(name, namespace string) (*batchV1.Job, error) {
	job, err := gTestEnv.KubeInt.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Job %s in namespace %s: %v", name, namespace, err)
	}
	return job, nil
}

// delete kubernetes Job
func DeleteJob(name, namespace string) error {
	err := gTestEnv.KubeInt.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Job %s in namespace %s: %v", name, namespace, err)
	}
	return nil
}

// IsJobPresent checks if a Kubernetes Job with the given name and namespace exists.
// It returns true if the Job is found, false if it is not found, and an error if there is an issue during retrieval.
func IsJobPresent(name, namespace string) (bool, error) {
	_, err := gTestEnv.KubeInt.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsJobCompleted checks if a Kubernetes Job with the given name and namespace has completed successfully.
func IsJobCompleted(name, namespace string) (bool, error) {
	job, err := GetJob(name, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get Job %s in namespace %s: %v", name, namespace, err)
	}
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchV1.JobComplete && condition.Status == "True" {
			return true, nil
		}
	}
	return false, nil
}

// wait for Job completion
func WaitForJobCompletion(name, namespace string, timeoutSeconds int) (bool, error) {
	for i := 0; i < timeoutSeconds; i++ {
		completed, err := IsJobCompleted(name, namespace)
		if err != nil {
			return false, fmt.Errorf("error checking job completion: %v", err)
		}
		if completed {
			logf.Log.Info("Job completed successfully", "name", name, "namespace", namespace)
			return true, nil
		}
		time.Sleep(sleepInterval * time.Second)
	}
	return false, fmt.Errorf("timeout waiting for Job %s in namespace %s to complete", name, namespace)
}

// wait for job deletion
func WaitForJobDeletion(name, namespace string, timeoutSeconds int) (bool, error) {
	for i := 0; i < timeoutSeconds; i++ {
		present, err := IsJobPresent(name, namespace)
		if err != nil {
			return false, fmt.Errorf("error checking job presence: %v", err)
		}
		if !present {
			logf.Log.Info("Job deleted successfully", "name", name, "namespace", namespace)
			return true, nil
		}
		time.Sleep(sleepInterval * time.Second)
	}
	return false, fmt.Errorf("timeout waiting for Job %s in namespace %s to be deleted", name, namespace)
}
