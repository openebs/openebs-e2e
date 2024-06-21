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

// DrainNode drain the given node with label
func (cp CPv1) DrainNode(nodeName string, drainLabel string, drainTimeOut int) error {
	logf.Log.Info("Executing drain node command", "node", nodeName, "drain label", drainLabel, "timeout", drainTimeOut)
	kubectlPlugin := GetPluginPath()
	// #FIXME remove drain timeout drain command is fuctional
	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "drain", "node", nodeName, drainLabel, "--drain-timeout", fmt.Sprintf("%ds", drainTimeOut))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s plugin failed to drain node %s with drain label %s , error %v", kubectlPlugin, nodeName, drainLabel, err)
	} else if strings.Contains(out.String(), ErrorResponse) {
		return fmt.Errorf("REST api error, failed to drain node %s with label %s, error %v", nodeName, drainLabel, out.String())
	}
	return nil
}

// GetDrainNodeLabels returns draining, drained labels and error
func (cp CPv1) GetDrainNodeLabels(nodeName string) ([]string, []string, error) {
	logf.Log.Info("Executing command to get drain node labels", "node", nodeName)
	pluginPath := GetPluginPath()
	cmd := exec.Command(pluginPath, "-n", common.NSMayastor(), "get", "node", nodeName, "-ojson")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("%s failed to get cordon labels for node %s, error %v", pluginPath, nodeName, err)
	}
	outputString := out.String()
	var cordonLabelsInfo NodeCordonLabelsInfo
	err = json.Unmarshal([]byte(outputString), &cordonLabelsInfo)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal command output %s, error %v ", outputString, err)
	}
	if cordonLabelsInfo.Spec.CordonDrainState == nil {
		return []string{}, []string{}, nil
	}

	if cordonLabelsInfo.Spec.CordonDrainState.DrainingState != nil {
		return cordonLabelsInfo.Spec.CordonDrainState.DrainingState.DrainLabels, []string{}, nil
	}

	if cordonLabelsInfo.Spec.CordonDrainState.DrainedState != nil {
		return []string{}, cordonLabelsInfo.Spec.CordonDrainState.DrainedState.DrainLabels, nil
	}

	if cordonLabelsInfo.Spec.CordonDrainState.CordonedState != nil {
		return []string{}, []string{}, nil // we are not interested in cordon labels
	}

	return []string{}, []string{}, fmt.Errorf("internal error, unexpected state")
}
