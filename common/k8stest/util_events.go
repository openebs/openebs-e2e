package k8stest

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetEvents retrieves events for a specific namespace
func GetEvents(nameSpace string, listOptions metaV1.ListOptions) (*v1.EventList, error) {
	return gTestEnv.KubeInt.CoreV1().Events(nameSpace).List(context.TODO(), listOptions)
}

// GetEvent retrieves event for a specific resource in namespace
func GetEvent(nameSpace string, name string, getOptions metaV1.GetOptions) (*v1.Event, error) {
	return gTestEnv.KubeInt.CoreV1().Events(nameSpace).Get(context.TODO(), name, getOptions)
}
