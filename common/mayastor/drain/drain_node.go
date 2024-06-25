package drain

import (
	"github.com/openebs/openebs-e2e/common/controlplane"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	VolSizeMb            = 1024
	DefTimeoutSecs       = 180 // in seconds
	DefDrainTimeoutSecs  = 180 // in seconds
	ZeroDrainTimeoutSecs = 0   // in seconds
)

// VerifyNodeDrainStatus return draining and drained status
func VerifyNodeDrainStatus(nodeName string, drainLabel string) (bool, bool, error) {
	var isDraining, isDrained bool
	logf.Log.Info("Verifying node drain status", "node", nodeName, "drain label", drainLabel)
	_drainingLabels, drainedLabels, err := controlplane.GetDrainNodeLabels(nodeName)
	if err != nil {
		return isDraining, isDrained, err
	}
	for _, label := range drainedLabels {
		if label == drainLabel {
			isDrained = true
		}
	}
	for _, label := range _drainingLabels {
		if label == drainLabel {
			isDraining = true
		}
	}
	return isDraining, isDrained, err
}
