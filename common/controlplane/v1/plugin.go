package v1

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

func GetPluginPath() string {
	if e2e_config.GetConfig().KubectlPluginDir == "" {
		panic("unspecified location of kubectl plugin")
	}
	pluginpath := fmt.Sprintf("%s/%s",
		e2e_config.GetConfig().KubectlPluginDir,
		e2e_config.GetConfig().Product.KubectlPluginName)
	return pluginpath
}

// GetPluginVersion returns the kubectl-maystor plugin version
// this function uses the kubectl-mayastor --version command
// and filter the output to get plugin version using `awk`
func GetPluginVersion() (string, error) {
	kubectlPlugin := GetPluginPath()

	// Define the kubectl-mayastor command with the "--version" flag
	cmd := exec.Command(kubectlPlugin, "--version")

	// Create a pipe to capture the command's output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running kubectl-mayastor --version, err:%v", err)
	}

	// Create an 'awk' command to extract the 6th column which is the plugin version
	// sample output for the command:-
	// # ./kubectl-mayastor --version
	// Kubectl Plugin (kubectl-mayastor) revision 571addcb58dd (2.5.0-0-main-unstable-2023-10-18-06-53-15-0+0)
	awkCmd := exec.Command("awk", "{print $6}")
	awkCmd.Stdin = strings.NewReader(string(output))

	// Capture the output of the 'awk' command
	awkOutput, err := awkCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running in awk command to filter plugin-version, err:%v", err)
	}

	if len(string(awkOutput)) == 0 {
		return "", fmt.Errorf("plugin version is empty in the output of kubectl-mayastor --version")
	}

	return string(awkOutput), nil
}

func CheckPluginError(jsonInput []byte, err error) error {
	// json error output trumps, error input
	if strings.Contains(string(jsonInput), ErrOutput) {
		return fmt.Errorf("%s", string(jsonInput))
	}
	return err
}
