package k8stest

import (
	"context"
	"fmt"
	"time"

	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
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
	// Define the deletion policy
	deletePolicy := metav1.DeletePropagationBackground
	err := gTestEnv.KubeInt.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
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

type JobBuilder struct {
	job *batchV1.Job
}

func NewJobBuilder(name string) *JobBuilder {
	return &JobBuilder{
		job: &batchV1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}
}

func (b *JobBuilder) WithNamespace(ns string) *JobBuilder {
	b.job.Namespace = ns
	return b
}

func (b *JobBuilder) WithParallelism(p int32) *JobBuilder {
	b.job.Spec.Parallelism = &p
	return b
}

func (b *JobBuilder) WithCompletions(c int32) *JobBuilder {
	b.job.Spec.Completions = &c
	return b
}

func (b *JobBuilder) WithIndexedMode() *JobBuilder {
	mode := batchV1.IndexedCompletion
	b.job.Spec.CompletionMode = &mode
	return b
}

func (b *JobBuilder) WithBackoffLimit(limit int32) *JobBuilder {
	b.job.Spec.BackoffLimit = &limit
	return b
}

func (b *JobBuilder) WithPodTemplate(pod *coreV1.Pod) *JobBuilder {
	// We take the Spec and Labels from the Pod created by your PodBuilder
	b.job.Spec.Template = coreV1.PodTemplateSpec{
		ObjectMeta: pod.ObjectMeta,
		Spec:       pod.Spec,
	}
	return b
}

func (b *JobBuilder) Build() (*batchV1.Job, error) {
	return b.job, nil
}

func CreateJob(job *batchV1.Job) error {
	_, err := gTestEnv.KubeInt.BatchV1().Jobs(job.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Job %s in namespace %s: %v", job.Name, job.Namespace, err)
	}
	return nil
}
