package apps

import (
	"encoding/json"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"os"
	"path/filepath"
	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	Standalone  Architecture = "standalone"
	Replicaset  Architecture = "replicaset"
	Replication Architecture = "replication"
)

type Architecture string

func (a Architecture) String() string {
	return string(a)
}

func RemoveHelmRepository(helmRepoName, helmRepoUrl string) error {
	cmd := exec.Command("helm", "repo", "remove", helmRepoName, helmRepoUrl)
	logf.Log.Info("executing helm remove repo ", "helm repo name: ", helmRepoName, ", helm repo url: ", helmRepoUrl)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove helm repo with url %s: %v\n%s", helmRepoUrl, err, output)
	}
	return nil
}

func AddHelmRepository(helmRepoName, helmRepoUrl string) error {
	cmd := exec.Command("helm", "repo", "add", helmRepoName, helmRepoUrl)
	logf.Log.Info("executing helm add repo ", "helm repo name: ", helmRepoName, ", helm repo url: ", helmRepoUrl)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add helm repo with url %s: %v\n%s", helmRepoUrl, err, output)
	}
	return nil
}

func UpdateHelmRepository(helmRepoName string) error {
	cmd := exec.Command("helm", "repo", "update", helmRepoName)
	logf.Log.Info("executing helm update repo ", "helm repo name: ", helmRepoName)
	// Execute the command.
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update helm repo")
	}
	return nil
}

func InstallHelmChart(helmChart, version, namespace, releaseName string, values map[string]interface{}) error {
	var vals []string
	for k, v := range values {
		vals = append(vals, fmt.Sprintf("%s=%v", k, v))
	}
	setVals := strings.Join(vals, ",")
	logf.Log.Info("executing helm install ", "releaseName: ", releaseName, ", chart: ", helmChart, ", version: ", version, ", namespace: ", namespace, ", values: ", setVals)
	// Define the Helm installation command.
	cmd := exec.Command("helm", "install", releaseName, helmChart, "--version", version, "-n", namespace, "--create-namespace", "--set", setVals)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install %s with Helm: %v\n%s", releaseName, err, output)
	}
	return nil
}

func GetInstalledProductChartVersionViaHelm(namespace string) (string, error) {

	// Execute the helm list command with YAML output
	cmd := exec.Command("helm", "list", "-n", namespace, "-o", "yaml")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to list release name in namespace: %s, err: %v", namespace, err)
	}

	// Convert output bytes to string
	output := string(out)

	// Split the output by newline character
	lines := strings.Split(output, "\n")

	// Initialize variable to store chart version
	var chartVersion string

	// Iterate through each line
	for _, line := range lines {
		// Trim leading and trailing whitespaces
		line = strings.TrimSpace(line)
		// Check if the line contains "chart:"
		if strings.Contains(line, "chart:") {
			// Extract the chart value
			chartValue := strings.TrimSpace(strings.TrimPrefix(line, "chart:"))
			// Split the chart value by "-"
			parts := strings.Split(chartValue, "-")
			if len(parts) >= 2 {
				// Extract the version part after "mayastor-"
				chartVersion = strings.Join(parts[1:], "-")
				break
			}
		}
	}

	// Print the app version
	logf.Log.Info("Installed Product", "Chart Version", chartVersion)

	if chartVersion == "" {
		return "", fmt.Errorf("chart version is empty")
	}

	return chartVersion, nil
}

func UpgradeHelmChartForValues(helmChart, namespace, releaseName string, values map[string]interface{}) error {
	var vals []string
	for k, v := range values {
		vals = append(vals, fmt.Sprintf("%s=%v", k, v))
	}
	setVals := strings.Join(vals, ",")
	logf.Log.Info("executing helm upgrade ", "releaseName: ", releaseName, ", chart: ", helmChart, ", namespace: ", namespace, ", values: ", setVals)
	// Define the Helm installation command.
	cmd := exec.Command("helm", "upgrade", releaseName, helmChart, "-n", namespace, "--reuse-values", "--set", setVals)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to upgrade with Helm: %v\n%s", err, output)
	}
	return nil
}

type Chart struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

func GetLatestHelmChartVersion(helmChart string) (Chart, error) {
	cmd := exec.Command("helm", "search", "repo", helmChart, "--versions", "-o", "json")
	var charts []Chart
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Chart{}, fmt.Errorf("failed to search repo %s: %v\n%s", helmChart, err, string(output))
	}
	err = json.Unmarshal(output, &charts)
	if err != nil {
		return Chart{}, err
	}
	if len(charts) == 0 {
		return Chart{}, fmt.Errorf("failed to find any chart version for repository %s", helmChart)
	}
	return charts[0], nil
}

func UninstallHelmRelease(releaseName, namespace string) error {
	// Define the Helm installation command.
	cmd := exec.Command("helm", "uninstall", releaseName, "-n", namespace)
	// Execute the command.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to uninstall release %s with Helm: %v\n%s", releaseName, err, output)
	}

	return nil

}

// Get values for a release from helm command
// This step will be needed before performing helm upgrade step
func HelmGetValues(releaseName string, namespace string)  ([]byte, error) {
	// Define the Helm installation command
	cmd := exec.Command("helm", "get", "values", releaseName, "-n", namespace)
	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output,fmt.Errorf("failed to helm get values for release %s : %v\n%s", releaseName, err, output)
	}
	return output, nil
}
// Save values for a release from helm command to a yaml file
// As an example "helm get values openebs -n openebs -o yaml > old-values.yaml"
func SaveHelmValues(releaseName string, namespace string, filePath string)  error {
	output, err := HelmGetValues(releaseName, namespace)
	if err != nil {
        	return err
    	}
	dir := filepath.Dir(filePath)
    	err = os.MkdirAll(dir, 0755) // Ensure intermediate directories are created
    	if err != nil {
        	return err
    	}
	// Write the output to a YAML file.
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output to file %s: %v", filePath, err)
	}
	return nil
}
// Upgrade a chart from a yaml file using helm command
// As an example "helm upgrade openebs openebs/openebs -n openebs -f old-values.yaml --version 4.1.1 \
//   --set openebs-crds.csi.volumeSnapshots.enabled=false"
func HelmUpgradeChartfromYaml(helmChart string, namespace string, releaseName string, setValues map[string]interface{}, filePath string, version string) (string, error) {
	var vals []string
	var setVals string
	// Check if setValues is empty
	if len(setValues) > 0 {
		for k, v := range setValues {
			vals = append(vals, fmt.Sprintf("%s=%v", k, v))
		}
		setVals = strings.Join(vals, ",")
	}
	
	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade", releaseName, helmChart, "-f", filePath, "--version", version}

	// Append arguments based on conditions
	if len(setValues) > 0 {
		cmdArgs = append(cmdArgs, "--set", setVals)
	}

	// Create the command
	cmd := exec.Command("helm", cmdArgs...)

	// Print the command that will be executed
	logf.Log.Info("Executing", "command", strings.Join(cmd.Args, " "))

	// Capture standard error output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		logf.Log.Info(stderr.String())
		return stderr.String(), fmt.Errorf("helm failed to upgrade, err:%v", err)
	}
	return stderr.String(), nil
}
