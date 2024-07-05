package k8stest

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type patchStruct struct {
	Op    string            `json:"op"`
	Path  string            `json:"path"`
	Value map[string]string `json:"value"`
}

func patchNodeLabels(nodeName string, labels map[string]string) error {
	labelPatch := []patchStruct{
		{
			Op:    "replace",
			Path:  "/metadata/labels",
			Value: labels,
		},
	}
	labelPatchBytes, err := json.Marshal(labelPatch)
	if err == nil {
		_, err = gTestEnv.KubeInt.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType, labelPatchBytes, metaV1.PatchOptions{})
	}
	return err
}

func ContainsK8sControlPlaneLabels(labels map[string]string) bool {
	_, ok := labels["node-role.kubernetes.io/control-plane"]
	if !ok {
		_, ok = labels["node-role.kubernetes.io/master"]
	}
	return ok
}

// SetupControlNode set up a node as the control node,
// A node is selected if it does not have any labels in the set defined by excludeLabels.
// A control node which is not a k8s master/control-plane node will be tainted with NoSchedule.
// Workloads are scheduled on the control node by setting tolerations when deploying the workload.
func SetupControlNode(excludeLabels []string, testControlNodeLabel string) (bool, error) {
	controlNodeSet := false
	nodeList, err := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err == nil {
		// First check if a node with the label exists,
		// if it does assume that the control node was previously set up correctly
		for _, node := range nodeList.Items {
			labels := node.GetLabels()
			if _, ok := labels[testControlNodeLabel]; ok {
				log.Log.Info("Using node found with test control label", "node.Name", node.Name)
				return true, nil
			}
		}
		//
		for _, node := range nodeList.Items {
			exclude := false
			labels := node.GetLabels()
			for _, excludeLabel := range excludeLabels {
				if _, exclude = labels[excludeLabel]; exclude {
					break
				}
			}
			if !exclude {
				// found a node which matches selection criteria namely
				// it does not have a label in present in the exclusion list
				labels[testControlNodeLabel] = ""
				err := patchNodeLabels(node.Name, labels)
				if err == nil {
					isK8sControlPlane := ContainsK8sControlPlaneLabels(labels)
					if !isK8sControlPlane {
						err = AddControlNodeNoScheduleTaintOnNode(node.Name)
					}
					if err == nil {
						controlNodeSet = true
						log.Log.Info("Test control node is set", "node.Name", node.Name, "isK8sControlPlane", isK8sControlPlane)
					} else {
						log.Log.Info("Failed to set test control node", "node.Name", node.Name, "isK8sControlPlane", isK8sControlPlane)
					}
				}
				break
			}
		}
	}
	return controlNodeSet, err
}

// EnsureNodeLabelValues that existing labels on a node have a particular value
func EnsureNodeLabelValues(IOEngineLabelKey string, IOEngineLabelValue string, testControlNodeLabel string) error {
	// try setting up a control node which is not a k8s control plane node
	controlNodeSet, err := SetupControlNode(
		[]string{IOEngineLabelKey, "node-role.kubernetes.io/control-plane", "node-role.kubernetes.io/master"},
		testControlNodeLabel)
	if err != nil {
		return err
	}
	if !controlNodeSet {
		// no node which is not a k8s control plane node available for use ass the control node,
		// try setting up a k8s control plane node as the control node.
		controlNodeSet, err = SetupControlNode(
			[]string{IOEngineLabelKey},
			testControlNodeLabel)
		if err != nil {
			return err
		}
	}
	if !controlNodeSet {
		// there are no nodes available for use as a control node:
		// - no k8s control plane nodes
		// - all other nodes are labelled as IOEngine nodes
		return fmt.Errorf("unable to setup the test control node")
	}

	nodeList, fnErrV := gTestEnv.KubeInt.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if fnErrV != nil {
		return fmt.Errorf("failed to list nodes, error: %v", fnErrV)
	}

	for _, node := range nodeList.Items {
		var val string
		var ok bool
		labels := node.GetLabels()

		if _, ok = labels[testControlNodeLabel]; ok {
			continue
		}

		val, ok = labels[IOEngineLabelKey]
		if !ok || val == IOEngineLabelValue {
			// if the node does not have the label key, leave untouched
			// if the labelValue matches -> nothing to do
			continue
		}

		// just replacing/adding  the key value pair doesn't work
		// replacing the labels map does though.
		log.Log.Info("re-labelling ", "node.Name", node.Name, "", fmt.Sprintf("%s=%s", IOEngineLabelKey, IOEngineLabelValue))
		labels[IOEngineLabelKey] = IOEngineLabelValue
		err := patchNodeLabels(node.Name, labels)
		if err != nil {
			fnErrV = err
			log.Log.Info("EnsureNodeLabelValues: add label failed", "err", err)
		}
	}

	log.Log.Info("EnsureNodeLabelValues:", IOEngineLabelKey, IOEngineLabelValue, "err", fnErrV)
	return fnErrV
}

func AddControlNodeNoScheduleTaintOnNode(nodeName string) error {
	cmd := exec.Command("kubectl", "taint", "node", nodeName, "openebs-test-control=:NoSchedule")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add no schedule taint on control node %s, error: %v", nodeName, err)
	}
	return nil
}

func RemoveControlNodeNoScheduleTaintFromNode(nodeName string) error {
	cmd := exec.Command("kubectl", "taint", "node", nodeName, "openebs-test-control=:NoSchedule"+"-")
	cmd.Dir = ""
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove no schedule taint from control node %s, error: %v", nodeName, err)
	}
	return nil
}

// EnsureNodeLabels  add the label xxxxxxx/engine=mayastor to  all worker nodes so that K8s runs mayastor on them
// returns error is accessing the list of nodes fails.
func EnsureNodeLabels() error {
	var errors common.ErrorAccumulator
	nodes, err := getNodeLocs()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if !node.K8sControlPlane && !node.ControlNode {
			err = LabelNode(node.NodeName,
				e2e_config.GetConfig().Product.EngineLabel,
				e2e_config.GetConfig().Product.EngineLabelValue)
			if err != nil {
				errors.Accumulate(err)
			}
		}
	}
	return errors.GetError()
}
