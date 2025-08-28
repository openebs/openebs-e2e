package k8stest

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type EventType int

const (
	Critical EventType = iota
	Warning  EventType = iota
	Normal   EventType = iota
)

func (volType EventType) String() string {
	switch volType {
	case Critical:
		return "Critical"
	case Warning:
		return "Warning"
	case Normal:
		return "Normal"
	default:
		return "Unknown"
	}
}

// GetEvents retrieves events for a specific namespace
func GetEvents(nameSpace string, listOptions metaV1.ListOptions) (*v1.EventList, error) {
	return gTestEnv.KubeInt.CoreV1().Events(nameSpace).List(context.TODO(), listOptions)
}

// GetEvent retrieves event for a specific resource in namespace
func GetEvent(nameSpace string, name string, getOptions metaV1.GetOptions) (*v1.Event, error) {
	return gTestEnv.KubeInt.CoreV1().Events(nameSpace).Get(context.TODO(), name, getOptions)
}

func IsEventExist(events *v1.EventList, eventType EventType, eventMessage string) bool {
	for _, event := range events.Items {
		if event.Type == eventType.String() && strings.Contains(event.Message, eventMessage) {
			logf.Log.Info("Found event", "message", event.Message, "type", event.Type)
			return true
		}
	}
	return false
}
