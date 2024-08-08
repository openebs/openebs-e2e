package k8stest

import (
	"context"

	v1 "k8s.io/api/apps/v1"
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
