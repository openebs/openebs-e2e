package apps

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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

func UpgradeHelmChart(helmChart, namespace, releaseName string, values map[string]interface{}) error {
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
