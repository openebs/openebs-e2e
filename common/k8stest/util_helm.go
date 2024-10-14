package k8stest

import (
	"fmt"
	"os/exec"
	"strings"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func UpgradeHelmChart(helmChart, namespace, releaseName, version string, values map[string]interface{}) error {
	var vals []string
	for k, v := range values {
		vals = append(vals, fmt.Sprintf("%s=%v", k, v))
	}
	setVals := strings.Join(vals, ",")
	logf.Log.Info("executing helm upgrade ", "releaseName: ", releaseName, ", chart: ", helmChart, "version", version, "namespace: ", namespace, ", values: ", setVals)
	// Define the Helm installation command.
	cmd := exec.Command("helm", "upgrade", releaseName, helmChart, "-n", namespace, "--reuse-values", "--set", setVals, "--version", version)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to upgrade with Helm: %v\n%s", err, output)
	}
	return nil
}
