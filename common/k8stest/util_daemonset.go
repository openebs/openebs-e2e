package k8stest

import (
	"context"
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListDaemonSet(namespace string) (v1.DaemonSetList, error) {
	daemonSets, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).List(context.TODO(), metaV1.ListOptions{})
	return *daemonSets, dserr
}

func GetDaemonSet(name, namespace string) (v1.DaemonSet, error) {
	daemonSet, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metaV1.GetOptions{})
	return *daemonSet, dserr
}

func UpdateDaemonSet(ds appsV1.DaemonSet, namespace string) (v1.DaemonSet, error) {
	daemonSet, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Update(context.TODO(), &ds, metaV1.UpdateOptions{})
	return *daemonSet, dserr
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
		return fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
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
