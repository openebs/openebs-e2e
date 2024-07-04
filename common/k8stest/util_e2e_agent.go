package k8stest

import (
	"context"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/locations"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func e2eReadyPodCount() int {
	daemonSet, err := gTestEnv.KubeInt.AppsV1().DaemonSets(common.NSE2EAgent).Get(
		context.TODO(),
		"e2e-rest-agent",
		metaV1.GetOptions{},
	)
	if err != nil {
		return -1
	}
	return int(daemonSet.Status.NumberAvailable)
}

// EnsureE2EAgent ensure that e2e agent daemonSet is running, if already deployed
// does nothing, otherwise creates the e2e agent namespace and deploys the daemonSet.
// asserts if creating the namespace fails. This function can be called repeatedly.
func EnsureE2EAgent() (bool, error) {
	const sleepTime = 5
	const duration = 60
	count := (duration + sleepTime - 1) / sleepTime
	ready := false
	err := EnsureNamespace(common.NSE2EAgent)
	if err != nil {
		return false, err
	}

	nodes, _ := GetIOEngineNodes()
	instances := len(nodes)

	if e2eReadyPodCount() == instances {
		return true, nil
	}

	err = KubeCtlApplyYaml("e2e-agent.yaml", locations.GetE2EAgentPath())

	if err != nil {
		return false, err
	}
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		ready = e2eReadyPodCount() == instances
	}
	return ready, nil
}

func UninstallE2eAgent() error {
	return KubeCtlDeleteYaml("e2e-agent.yaml", locations.GetE2EAgentPath())
}
