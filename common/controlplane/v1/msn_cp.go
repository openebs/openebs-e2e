package v1

// Utility functions for Mayastor control plane volume
import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type MayastorCpNode struct {
	Spec  msnSpec  `json:"spec"`
	State msnState `json:"state"`
}

type msnSpec struct {
	GrpcEndpoint string `json:"grpcEndpoint"`
	ID           string `json:"id"`
	Node_nqn     string `json:"node_nqn"`
}

type msnState struct {
	GrpcEndpoint string `json:"grpcEndpoint"`
	ID           string `json:"id"`
	Status       string `json:"status"`
	Node_nqn     string `json:"node_nqn"`
}

func GetMayastorCpNode(nodeName string) (*MayastorCpNode, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "node", nodeName)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response MayastorCpNode
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		msg := string(jsonInput)
		if !HasNotFoundRestJsonError(msg) {
			logf.Log.Info("Failed to unmarshal (get node)", "string", msg, "node", nodeName)
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return &response, nil
}

func ListMayastorCpNodes() ([]MayastorCpNode, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "nodes")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []MayastorCpNode
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		logf.Log.Info("Failed to unmarshal (get nodes)", "string", string(jsonInput))
		return []MayastorCpNode{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func GetMayastorNodeStatus(nodeName string) (string, error) {
	msn, err := GetMayastorCpNode(nodeName)
	if err == nil {
		return msn.State.Status, nil
	}
	return "", err
}

func cpNodeToMsn(cpMsn *MayastorCpNode) common.MayastorNode {
	return common.MayastorNode{
		Name: cpMsn.Spec.ID,
		Spec: common.MayastorNodeSpec{
			ID:           cpMsn.Spec.ID,
			GrpcEndpoint: cpMsn.Spec.GrpcEndpoint,
			Node_nqn:     cpMsn.Spec.Node_nqn,
		},
		State: common.MayastorNodeState{
			ID:           cpMsn.State.ID,
			Status:       cpMsn.State.Status,
			GrpcEndpoint: cpMsn.State.GrpcEndpoint,
			Node_nqn:     cpMsn.State.Node_nqn,
		},
	}
}

// GetMSN Get pointer to a mayastor control plane volume
// returns nil and no error if the msn is in pending state.
func (cp CPv1) GetMSN(nodeName string) (*common.MayastorNode, error) {
	cpMsn, err := GetMayastorCpNode(nodeName)
	if err != nil {
		return nil, fmt.Errorf("GetMSN: %v", err)
	}

	if cpMsn == nil {
		logf.Log.Info("Msn not found", "node", nodeName)
		return nil, nil
	}

	msn := cpNodeToMsn(cpMsn)
	return &msn, nil
}

func (cp CPv1) ListMsns() ([]common.MayastorNode, error) {
	var msns []common.MayastorNode
	list, err := ListMayastorCpNodes()
	if err == nil {
		for _, item := range list {
			msns = append(msns, cpNodeToMsn(&item))
		}
	}
	return msns, err
}

func (cp CPv1) GetMsNodeStatus(nodeName string) (string, error) {
	cpMsn, err := GetMayastorCpNode(nodeName)
	if err != nil {
		return "", fmt.Errorf("GetMsNodeStatus: %v", err)
	}
	return cpMsn.State.Status, nil
}

// UpdateNodeLabel adds or remove labels from nodes
func (cp CPv1) UpdateNodeLabel(nodeName string, labelKey, labelValue string) error {
	pluginPath := GetPluginPath()
	args := []string{"label", "node", nodeName}

	// Check if a label value is provided
	if labelValue != "" {
		// Add the label key-value pair to the arguments
		args = append(args, fmt.Sprintf("%s=%s", labelKey, labelValue))
	} else {
		// Add the label key to be removed to the arguments
		args = append(args, fmt.Sprintf("%s-", labelKey))
	}

	cmd := exec.Command(pluginPath, args...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		// Print the error message if the command fails
		return fmt.Errorf("plugin failed to update node label %s: %v", nodeName, err)
	}
	return nil
}
