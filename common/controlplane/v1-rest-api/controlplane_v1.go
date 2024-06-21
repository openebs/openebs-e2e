package v1_rest_api

import (
	"fmt"
	"os"
	"regexp"

	cpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/generated/openapi"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type CPv1RestApi struct {
	oa OAClientWrapper
}

func (cp CPv1RestApi) Version() string {
	return "1.0.0"
}

func (cp CPv1RestApi) MajorVersion() int {
	return 1
}

func (cp CPv1RestApi) VolStateHealthy() string {
	return string(openapi.VOLUMESTATUS_ONLINE)
}

func (cp CPv1RestApi) VolStateUnknown() string {
	return string(openapi.VOLUMESTATUS_UNKNOWN)
}

func (cp CPv1RestApi) VolStateDegraded() string {
	return string(openapi.VOLUMESTATUS_DEGRADED)
}

func (cp CPv1RestApi) VolStateFaulted() string {
	return string(openapi.VOLUMESTATUS_FAULTED)
}

func (cp CPv1RestApi) ChildStateOnline() string {
	return string(openapi.CHILDSTATE_ONLINE)
}

func (cp CPv1RestApi) ChildStateDegraded() string {
	return string(openapi.CHILDSTATE_DEGRADED)
}

func (cp CPv1RestApi) ChildStateUnknown() string {
	return string(openapi.CHILDSTATE_UNKNOWN)
}

func (cp CPv1RestApi) ChildStateFaulted() string {
	return string(openapi.CHILDSTATE_FAULTED)
}

func (cp CPv1RestApi) NexusStateUnknown() string {
	return string(openapi.CHILDSTATE_UNKNOWN)
}

func (cp CPv1RestApi) NexusStateOnline() string {
	return string(openapi.NEXUSSTATE_ONLINE)
}

func (cp CPv1RestApi) NexusStateDegraded() string {
	return string(openapi.NEXUSSTATE_DEGRADED)
}

func (cp CPv1RestApi) NexusStateFaulted() string {
	return string(openapi.NEXUSSTATE_FAULTED)
}

func (cp CPv1RestApi) MspStateOnline() string {
	return string(openapi.POOLSTATUS_ONLINE)
}

func (cp CPv1RestApi) NodeStateOffline() string {
	return string(openapi.NODESTATUS_OFFLINE)
}

func (cp CPv1RestApi) NodeStateOnline() string {
	return string(openapi.NODESTATUS_ONLINE)
}

func (cp CPv1RestApi) NodeStateUnknown() string {
	return string(openapi.NODESTATUS_UNKNOWN)
}

func (cp CPv1RestApi) ReplicaStateUnknown() string {
	return string(openapi.REPLICASTATE_UNKNOWN)
}

func (cp CPv1RestApi) ReplicaStateOnline() string {
	return string(openapi.REPLICASTATE_ONLINE)
}

func (cp CPv1RestApi) ReplicaStateDegraded() string {
	return string(openapi.REPLICASTATE_DEGRADED)
}

func (cp CPv1RestApi) ReplicaStateFaulted() string {
	return string(openapi.REPLICASTATE_FAULTED)
}

func (cp CPv1RestApi) NodeStateEmpty() string {
	return ""
}

// MakeCP make control plane object which uses the REST API
func MakeCP() (CPv1RestApi, error) {
	var addrs []string
	var err error

	addrs, err = cpV1.GetClusterRestAPINodeIPs()
	if err != nil {
		return CPv1RestApi{}, err
	}

	cp := CPv1RestApi{
		oa: MakeWrapper(addrs),
	}
	logf.Log.Info("Control Plane v1 - Rest API")
	return cp, err
}

var re = regexp.MustCompile(`(statusCode=408)`)

func (cp CPv1RestApi) IsTimeoutError(err error) bool {
	str := fmt.Sprintf("%v", err)
	frags := re.FindSubmatch([]byte(str))
	return len(frags) == 2 && string(frags[1]) == "statusCode=408"
}

func (cp CPv1RestApi) CreatePoolOnInstall() bool {
	noPoolInstall := os.Getenv("no_pool_install")
	return noPoolInstall != "true"
}
