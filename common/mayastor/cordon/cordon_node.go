package cordon

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	NodeCordonTimeout = 120
	VolSizeMb         = 500 // size in Mb
	CordonsCount      = 5
)

var DefTimeoutSecs = "120s"

func VerifyCordonNode(nodeName string, cordonLabel string) (bool, error) {
	logf.Log.Info("Verifying node cordon", "node", nodeName, "cordon label", cordonLabel)

	var isPresent bool
	cordonLabels, err := controlplane.GetCordonNodeLabels(nodeName)
	if err != nil {
		return isPresent, err
	}
	for _, label := range cordonLabels {
		if label == cordonLabel {
			isPresent = true
			break
		}
	}
	return isPresent, err
}

func VerifyUncordonNodesStatus() (bool, error) {
	logf.Log.Info("Verifying node uncordon status")
	nodes, err := k8stest.GetMayastorNodeNames()
	if err != nil {
		return false, fmt.Errorf("failed to list mayastor nodes, error %v", err)
	}
	for _, node := range nodes {
		cordonLabels, err := controlplane.GetCordonNodeLabels(node)
		if err != nil {
			return false, err
		}
		if len(cordonLabels) != 0 {
			logf.Log.Info("node cordon label present", "labels", cordonLabels)
			return false, err
		}
	}
	return true, err
}

func cancelCordons(nodeName string, labels []string, errIn error, labelType string) bool {
	if errIn != nil {
		logf.Log.Info("get labels failed", "node", nodeName, "type", labelType, "error", errIn)
		return false
	}
	var err error
	success := true
	for _, label := range labels {
		err = controlplane.UnCordonNode(nodeName, label)
		success = success && err == nil
		if err != nil {
			logf.Log.Info("uncordon failed", "node", nodeName, "error", err)
		}
	}
	return success
}

// CancelAllCordonsOnNode utility function to cancel all cordons on a node
func CancelAllCordonsOnNode(nodeName string) bool {
	var labels []string
	var err error
	success := true

	// Draining
	labels, _, err = controlplane.GetDrainNodeLabels(nodeName)
	success = success && cancelCordons(nodeName, labels, err, "draining")

	// Drained
	_, labels, err = controlplane.GetDrainNodeLabels(nodeName)
	success = success && cancelCordons(nodeName, labels, err, "drained")

	// Cordoned
	labels, err = controlplane.GetCordonNodeLabels(nodeName)
	success = success && cancelCordons(nodeName, labels, err, "cordoned")

	return success
}

func UncordonAllNodes() error {
	nodes, err := k8stest.GetIOEngineNodes()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if !CancelAllCordonsOnNode(node.NodeName) {
			err = fmt.Errorf("CancelAllCordonsOnNode(%s) failed; %v", node.NodeName, err)
		}
	}
	return err
}
