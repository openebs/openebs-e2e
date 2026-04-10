package k8stest

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var defEventTimeoutSecs = 120

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

func IsWarningPvcEventPresent(pvcName string, namespace string, errorSubstring string) (bool, error) {
	events, err := GetPvcEvents(pvcName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get pvc events in namespace %s, error: %v", namespace, err)
	}
	for _, event := range events.Items {
		logf.Log.Info("Found PVC event", "message", event.Message)
		if event.Type == Warning.String() && strings.Contains(event.Message, errorSubstring) {
			return true, err
		}
	}
	return false, err
}

func IsWarningPodEventPresent(podName string, namespace string, errorSubstring string) (bool, error) {
	events, err := GetPodEvents(podName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get pod events in namespace %s, error: %v", namespace, err)
	}
	for _, event := range events.Items {
		logf.Log.Info("Found Pod event", "message", event.Message)
		if event.Type == "Warning" && strings.Contains(event.Message, errorSubstring) {
			return true, err
		}
	}
	return false, err
}

// Wait for the PVC warning event
// return true if pvc warning event found else return false , it return error in case of error
func WaitForPvcWarningEvent(pvcName string, namespace string, errorSubstring string) (bool, error) {
	const timeSleepSecs = 1
	var hasWarning bool
	var err error
	// Wait for the PVC event
	logf.Log.Info("Check Pvc event", "Event substring", errorSubstring)
	for ix := 0; ix < defEventTimeoutSecs/timeSleepSecs; ix++ {
		hasWarning, err = IsWarningPvcEventPresent(pvcName, namespace, errorSubstring)
		if err == nil && hasWarning {
			break
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return hasWarning, err
}
