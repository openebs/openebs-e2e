package v1

import (
	"fmt"
	"regexp"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type CPv1 struct {
	nodeIPAddresses []string
}

func (cp CPv1) Version() string {
	return "1.0.0"
}

func (cp CPv1) MajorVersion() int {
	return 1
}

var regExs = []*regexp.Regexp{
	regexp.MustCompile(`(Request Timeout)`),
	regexp.MustCompile(`(request timed out)`),
}

func (cp CPv1) IsTimeoutError(err error) bool {
	str := fmt.Sprintf("%v", err)
	for _, re := range regExs {
		frags := re.FindSubmatch([]byte(str))
		if len(frags) != 0 {
			return true
		}
	}
	return false
}

func (cp CPv1) VolStateHealthy() string {
	return "Online"
}

func (cp CPv1) VolStateDegraded() string {
	return "Degraded"
}

func (cp CPv1) VolStateFaulted() string {
	return "Faulted"
}

func (cp CPv1) VolStateUnknown() string {
	return "Unknown"
}

func (cp CPv1) ChildStateOnline() string {
	return "Online"
}

func (cp CPv1) ChildStateDegraded() string {
	return "Degraded"
}

func (cp CPv1) ChildStateUnknown() string {
	return "Unknown"
}

func (cp CPv1) ChildStateFaulted() string {
	return "Faulted"
}

func (cp CPv1) NexusStateUnknown() string {
	return "Unknown"
}

func (cp CPv1) NexusStateOnline() string {
	return "Online"
}

func (cp CPv1) NexusStateDegraded() string {
	return "Degraded"
}

func (cp CPv1) NexusStateFaulted() string {
	return "Faulted"
}

func (cp CPv1) ReplicaStateUnknown() string {
	return "Unknown"
}

func (cp CPv1) ReplicaStateOnline() string {
	return "Online"
}

func (cp CPv1) ReplicaStateDegraded() string {
	return "Degraded"
}

func (cp CPv1) ReplicaStateFaulted() string {
	return "Faulted"
}

func (cp CPv1) MspStateOnline() string {
	return "Online"
}

func MakeCP() (CPv1, error) {
	var addrs []string
	var err error

	addrs, err = GetClusterRestAPINodeIPs()
	if err != nil {
		return CPv1{}, err
	}

	logf.Log.Info("Control Plane v1 - Kubectl plugin")
	return CPv1{
		nodeIPAddresses: addrs,
	}, nil
}

func (cp CPv1) NodeStateOnline() string {
	return "Online"
}

func (cp CPv1) NodeStateOffline() string {
	return "Offline"
}

func (cp CPv1) NodeStateUnknown() string {
	return "Unknown"
}

func (cp CPv1) NodeStateEmpty() string {
	return ""
}
