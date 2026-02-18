package k8stest

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func ListDaemonSet(namespace string) (v1.DaemonSetList, error) {
	daemonSets, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).List(context.TODO(), metaV1.ListOptions{})
	return *daemonSets, dserr
}

func GetDaemonSet(name, namespace string) (v1.DaemonSet, error) {
	daemonSet, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	return *daemonSet, dserr
}

func UpdateDaemonSet(ds v1.DaemonSet, namespace string) (v1.DaemonSet, error) {
	daemonSet, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Update(context.TODO(), &ds, metaV1.UpdateOptions{})
	return *daemonSet, dserr
}

func DeleteDaemonSet(name, namespace string) error {
	err := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
	return err
}

// update daemonset container env
func UpdateDemonsetContainerEnv(daemonsetName string, containerName string, namespace string, envName string, envValue string) ([]coreV1.EnvVar, error) {
	var old_env []coreV1.EnvVar
	var err error
	daemonset, err := GetDaemonSet(daemonsetName, namespace)
	if err != nil {
		return old_env, fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
			daemonsetName,
			namespace,
			err)
	}

	containers := daemonset.Spec.Template.Spec.Containers
	for i, container := range containers {
		if container.Name == containerName {
			old_env = containers[i].Env
			containers[i].Env = replaceOrAppendEnv(container.Env, envName, envValue)
			break
		}
	}

	daemonset.Spec.Template.Spec.Containers = containers

	_, err = UpdateDaemonSet(daemonset, namespace)

	if err != nil {
		return old_env, fmt.Errorf("failed to set container env to daemonset, name: %s, container %s, namespace: %s, env name: %s, error: %v",
			daemonsetName,
			containerName,
			namespace,
			envName,
			err)
	}
	return old_env, nil
}

func replaceOrAppendEnv(envList []coreV1.EnvVar, envName string, envValue string) []coreV1.EnvVar {
	for ix, env := range envList {
		if env.Name == envName {
			env.Value = envValue
			envList[ix] = env
			return envList
		}
	}
	env := coreV1.EnvVar{
		Name:  envName,
		Value: envValue,
	}
	return append(envList, env)
}

// update daemonset container all env
func UpdateDemonsetContainerAllEnv(daemonsetName string, containerName string, namespace string, env []coreV1.EnvVar) error {

	var err error
	daemonset, err := GetDaemonSet(daemonsetName, namespace)
	if err != nil {
		return fmt.Errorf("failed to get daemonset, name: %s, namespace: %s, error: %v",
			daemonsetName,
			namespace,
			err)
	}

	containers := daemonset.Spec.Template.Spec.Containers
	for i, container := range containers {
		if container.Name == containerName {
			containers[i].Env = env
			break
		}
	}

	daemonset.Spec.Template.Spec.Containers = containers

	_, err = UpdateDaemonSet(daemonset, namespace)

	if err != nil {
		return fmt.Errorf("failed to set container env to daemonset, name: %s, container %s, namespace: %s, error: %v",
			daemonsetName,
			containerName,
			namespace,
			err)
	}
	return nil
}

func WaitForDaemonsetReady(dsName string, namespace string, sleepTime int, duration int) bool {
	ready := false
	count := (duration + sleepTime - 1) / sleepTime

	logf.Log.Info("DaemonsetReadyCheck", "Daemonet", dsName, "namespace", namespace)
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)

		ready = DaemonSetReady(dsName, namespace)
		logf.Log.Info("DaemonSetReady: ", "daemonset", dsName, "ready", ready)
	}

	if !ready {
		logf.Log.Info("Daemonset not ready", "Daemonset", dsName, "namespace", namespace)
		return false
	}

	dsPodList, err := ListPodsByPrefix(namespace, dsName)
	if err != nil {
		logf.Log.Info("Failed to list pods with daemonset prefix", "daemonset", dsName, "error", err)
		return false
	}
	for _, pod := range dsPodList {
		// verify pod running
		err := WaitForPodRunning(pod.Name, namespace, duration)
		if err != nil {
			logf.Log.Info("Pod not ready", "Pod name", pod.Name, "namespace", namespace)
			return false
		}

	}
	return ready
}

func WaitForDaemonsetLatestPodReady(dsName string, namespace string, sleepTime int, duration int) bool {
	ready := false
	count := (duration + sleepTime - 1) / sleepTime

	logf.Log.Info("DaemonSetReadyCheck", "DaemonSet", dsName, "namespace", namespace)
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)

		ready = DaemonSetReady(dsName, namespace)
		logf.Log.Info("DaemonSetReady: ", "DaemonSet", dsName, "ready", ready)
	}

	if !ready {
		logf.Log.Info("DaemonSet not ready", "DaemonSet", dsName, "namespace", namespace)
		return false
	}

	// sleep for 60 seconds before checking the pods state because daemonset restart and new pod will be created
	// and previous pod will be deleted
	logf.Log.Info("DaemonSet pod restarting, sleep for 60 seconds")
	time.Sleep(60 * time.Second)
	// get the daemonset pods
	dsPodList, err := ListPodsByPrefixSortedByCreationTimeStamp(namespace, dsName)
	if err != nil {
		logf.Log.Info("Failed to list pods with daemonset prefix", "DaemonSet", dsName, "error", err)
		return false
	}
	logf.Log.Info("DaemonSet", "pod count", len(dsPodList))

	if len(dsPodList) == 0 {
		logf.Log.Info("No pods found with daemonset prefix", "DaemonSet", dsName)
		return false
	}
	// get the latest pod
	latestPod := dsPodList[len(dsPodList)-1]
	fmt.Println("Latest pod:", latestPod.Name)

	// verify pod running
	err = WaitForPodRunning(latestPod.Name, namespace, duration)
	if err != nil {
		logf.Log.Info("Pod not ready", "Pod name", latestPod.Name, "namespace", namespace)
		return false
	}

	return ready
}

func DeleteDaemonsetAndWaitPodDeletion(dsName string, namespace string, sleepTime int, duration time.Duration) bool {

	ds, err := GetDaemonSet(dsName, namespace)
	if err != nil {
		logf.Log.Info("Failed to get daemonset", "daemonset", dsName, "namespace", namespace, "error", err)
		return false
	}

	dsPodList, err := ListPodsByPrefix(ds.Namespace, ds.Name)
	if err != nil {
		logf.Log.Info("Failed to list pods with daemonset prefix", "daemonset", dsName, "error", err)
		return false
	}

	err = DeleteDaemonSet(dsName, namespace)
	if err != nil {
		logf.Log.Info("Failed to delete daemonset", "daemonset", dsName, "namespace", namespace, "error", err)
		return false
	}

	for _, pod := range dsPodList {
		// verify pod running
		isPodDeleted, err := WaitForPodDeletion(pod.Name, namespace, duration)
		if err != nil {
			logf.Log.Info("failed to verify pod not deletion", "Pod name", pod.Name, "namespace", namespace, "error", err)
			return false
		}
		if !isPodDeleted {
			logf.Log.Info("Pod not deleted", "Pod name", pod.Name, "namespace", namespace)
			return isPodDeleted
		}

	}
	return true
}

func DaemonSetReady(daemonName string, namespace string) bool {
	daemon, err := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Get(
		context.TODO(),
		daemonName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get daemonset", "error", err)
		return false
	}

	status := daemon.Status
	logf.Log.Info("DaemonSet "+daemonName, "status", status)
	return status.DesiredNumberScheduled == status.CurrentNumberScheduled &&
		status.DesiredNumberScheduled == status.NumberReady &&
		status.DesiredNumberScheduled == status.NumberAvailable
}
