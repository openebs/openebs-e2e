package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const ErrorResponse = "412 Precondition Failed"

type NodeCordonLabelsInfo struct {
	Id    string          `json:"id"`
	Spec  NodeCordonSpec  `json:"spec"`
	State NodeCordonState `json:"state"`
}
type CordonedState struct {
	CordonLabels []string `json:"cordonlabels"`
}
type DrainingState struct {
	CordonLabels []string `json:"cordonlabels"`
	DrainLabels  []string `json:"drainlabels"`
}
type DrainedState struct {
	CordonLabels []string `json:"cordonlabels"`
	DrainLabels  []string `json:"drainlabels"`
}
type CordonDrainState struct {
	CordonedState *CordonedState `json:"cordonedstate,omitempty"`
	DrainingState *DrainingState `json:"drainingstate,omitempty"`
	DrainedState  *DrainedState  `json:"drainedstate,omitempty"`
}
type NodeCordonSpec struct {
	GrpcEndpoint     string            `json:"grpcEndpoint"`
	Id               string            `json:"id"`
	CordonDrainState *CordonDrainState `json:"cordondrainstate,omitempty"`
}

type NodeCordonState struct {
	GrpcEndpoint string `json:"grpcEndpoint"`
	Id           string `json:"id"`
	Status       string `json:"status"`
}

func (cp CPv1) CordonNode(nodeName string, cordonLabel string) error {
	logf.Log.Info("Executing cordon node command", "node", nodeName, "cordon label", cordonLabel)
	kubectlPlugin := GetPluginPath()
	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "cordon", "node", nodeName, cordonLabel)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to cordon node %s with cordon label %s , error %v", nodeName, cordonLabel, err)
	} else if strings.Contains(out.String(), ErrorResponse) {
		return fmt.Errorf("REST api error,failed to cordon node %s with label %s, error %v", nodeName, cordonLabel, out.String())
	}
	return nil
}

func (cp CPv1) GetCordonNodeLabels(nodeName string) ([]string, error) {
	logf.Log.Info("Executing command to get cordon node labels", "node", nodeName)
	pluginPath := GetPluginPath()
	cmd := exec.Command(pluginPath, "-n", common.NSMayastor(), "get", "node", nodeName, "-ojson")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get cordon labels for node %s, error %v", nodeName, err)
	}
	outputString := out.String()
	var cordonLabelsInfo NodeCordonLabelsInfo
	err = json.Unmarshal([]byte(outputString), &cordonLabelsInfo)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal command output %s, error %v ", outputString, err)
	}
	if cordonLabelsInfo.Spec.CordonDrainState == nil {
		return []string{}, nil
	}
	if cordonLabelsInfo.Spec.CordonDrainState.CordonedState != nil {
		return cordonLabelsInfo.Spec.CordonDrainState.CordonedState.CordonLabels, err
	}
	// in case we ever want the cordon labels for a node which is draining/drained
	if cordonLabelsInfo.Spec.CordonDrainState.DrainingState != nil {
		return cordonLabelsInfo.Spec.CordonDrainState.DrainingState.CordonLabels, err
	}
	if cordonLabelsInfo.Spec.CordonDrainState.DrainedState != nil {
		return cordonLabelsInfo.Spec.CordonDrainState.DrainedState.CordonLabels, err
	}
	return nil, fmt.Errorf("Unexpected cordon spec for node %s", nodeName)
}

func (cp CPv1) UnCordonNode(nodeName string, cordonLabel string) error {
	logf.Log.Info("Executing uncordon node command", "node", nodeName, "cordon label", cordonLabel)
	kubectlPlugin := GetPluginPath()
	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "uncordon", "node", nodeName, cordonLabel)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to uncordon node %s with cordon label %s , error %v", nodeName, cordonLabel, err)
	} else if strings.Contains(out.String(), ErrorResponse) {
		return fmt.Errorf("REST api error,failed to uncordon node %s with label %s, error %v", nodeName, cordonLabel, out.String())
	}
	return nil
}
