package v1

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common"
)

type upgradeFlags string

const (
	SkipDataPlaneRestartFlag         upgradeFlags = "--skip-data-plane-restart"
	SkipSingleReplicaValidationFlag  upgradeFlags = "--skip-single-replica-volume-validation"
	SkipReplicaRebuildFlag           upgradeFlags = "--skip-replica-rebuild"
	SkipCordonNodeValidationFlag     upgradeFlags = "--skip-cordoned-node-validation"
	AllowUpgradeToUnstableBranchFlag upgradeFlags = "--allow-unstable"
	DisablePartialRebuild            upgradeFlags = "agents.core.rebuild.partial.enabled=false"
)

// This function is to fire upgrade command with kubectl mayastor plugin
// Syntax is: `kubectl-mayastor upgrade`
// this function takes one boolean parameter, `isUpgradingToUnstableBranch`
// this parameters is passed as true if we want to test upgrade to unstable main branch
// and uses --allow-unstable flag with upgrade command.
func (cp CPv1) Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error) {
	kubectlPlugin := GetPluginPath()

	CIRegistry, ok := os.LookupEnv("CI_REGISTRY")
	if !ok {
		return "", fmt.Errorf("environment variable CI_REGISTRY is not defined")
	}

	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade"}

	// Append arguments based on conditions
	if isUpgradingToUnstableBranch {
		cmdArgs = append(cmdArgs, "--registry", CIRegistry, string(AllowUpgradeToUnstableBranchFlag))
	}
	if isPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
	}

	// Create the command
	cmd := exec.Command(kubectlPlugin, cmdArgs...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	// Capture standard error output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("plugin failed to upgrade, err:%v", err)
	}
	return stderr.String(), nil
}

// This function is for starting upgrade with a flag to skip data plane restart
// User can manually restart data plane later. During upgrade images will be upgraded
// but pod restart will be manual process when used this flag.
// Syntax is: `kubectl-mayastor upgrade --skip-data-plane-restart`
// this function takes one boolean parameter, `isUpgradingToUnstableBranch`
// this parameters is passed as true if we want to test upgrade to unstable main branch
// and uses --allow-unstable flag with upgrade command.
func (cp CPv1) UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {

	kubectlPlugin := GetPluginPath()

	CIRegistry, ok := os.LookupEnv("CI_REGISTRY")
	if !ok {
		return fmt.Errorf("environment varianble CI_REGISTRY is not defined")
	}

	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade", string(SkipDataPlaneRestartFlag)}

	// Append arguments based on conditions
	if isUpgradingToUnstableBranch {
		cmdArgs = append(cmdArgs, "--registry", CIRegistry, string(AllowUpgradeToUnstableBranchFlag))
	}
	if isPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
	}

	// Create the command
	cmd := exec.Command(kubectlPlugin, cmdArgs...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("plugin failed to upgrade when skip data plane restart flag is passsed , error %v", err)
	}
	return nil
}

// This function is to fire upgrade command with skip single replica volume validation flag
// Syntax is: `kubectl-mayastor upgrade --skip-single-replica-volume-validation`
// this function takes one boolean parameter, `isUpgradingToUnstableBranch`
// this parameters is passed as true if we want to test upgrade to unstable main branch
// and uses --allow-unstable flag with upgrade command.
func (cp CPv1) UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {

	kubectlPlugin := GetPluginPath()

	CIRegistry, ok := os.LookupEnv("CI_REGISTRY")
	if !ok {
		return fmt.Errorf("environment varianble CI_REGISTRY is not defined")
	}

	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade", string(SkipSingleReplicaValidationFlag)}

	// Append arguments based on conditions
	if isUpgradingToUnstableBranch {
		cmdArgs = append(cmdArgs, "--registry", CIRegistry, string(AllowUpgradeToUnstableBranchFlag))
	}
	if isPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
	}

	// Create the command
	cmd := exec.Command(kubectlPlugin, cmdArgs...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("plugin failed to upgrade when skip single replica volume flag is passsed , error %v", err)
	}
	return nil
}

// This function is to fire upgrade command with kubectl mayastor plugin with --skip-replica-rebuild flag
// Syntax is: `kubectl-mayastor upgrade --skip-replica-rebuild`
// this function takes one boolean parameter, `isUpgradingToUnstableBranch`
// this parameters is passed as true if we want to test upgrade to unstable main branch
// and uses --allow-unstable flag with upgrade command.
func (cp CPv1) UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	kubectlPlugin := GetPluginPath()

	CIRegistry, ok := os.LookupEnv("CI_REGISTRY")
	if !ok {
		return fmt.Errorf("environment varianble CI_REGISTRY is not defined")
	}

	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade", string(SkipReplicaRebuildFlag)}

	// Append arguments based on conditions
	if isUpgradingToUnstableBranch {
		cmdArgs = append(cmdArgs, "--registry", CIRegistry, string(AllowUpgradeToUnstableBranchFlag))
	}
	if isPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
	}

	// Create the command
	cmd := exec.Command(kubectlPlugin, cmdArgs...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("plugin failed to upgrade with --skip-rebuild-replica flag , error %v", err)
	}
	return nil
}

// this function takes one boolean parameter, `isUpgradingToUnstableBranch`
// this parameters is passed as true if we want to test upgrade to unstable main branch
// and uses --allow-unstable flag with upgrade command.
func (cp CPv1) UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	kubectlPlugin := GetPluginPath()

	CIRegistry, ok := os.LookupEnv("CI_REGISTRY")
	if !ok {
		return fmt.Errorf("environment varianble CI_REGISTRY is not defined")
	}

	// Construct the base command
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade", string(SkipCordonNodeValidationFlag)}

	// Append arguments based on conditions
	if isUpgradingToUnstableBranch {
		cmdArgs = append(cmdArgs, "--registry", CIRegistry, string(AllowUpgradeToUnstableBranchFlag))
		if isPartialRebuildDisableNeeded {
			cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
		}
	} else if isPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(DisablePartialRebuild))
	}

	// Create the command
	cmd := exec.Command(kubectlPlugin, cmdArgs...)

	// Print the command that will be executed
	fmt.Println("Executing command:", strings.Join(cmd.Args, " "))

	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("plugin failed to upgrade with skip cordon node validation flag, error %v", err)
	}
	return nil
}

// This function is for getting status of upgrade
// Syntax is: `kubectl-mayastor get upgrade-status`

func (cp CPv1) GetUpgradeStatus() (string, error) {
	kubectlPlugin := GetPluginPath()

	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "get", "upgrade-status")
	upgradeStatusInfo, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("plugin failed to get upgrade status, error %v", err)
	}

	out := strings.Split(string(upgradeStatusInfo), "\n")
	var upgradeStatus string
	for _, line := range out {
		if strings.Contains(line, "Upgrade Status") {
			parts := strings.Split(line, ":")
			upgradeStatus = parts[len(parts)-1]
		}
	}
	return upgradeStatus, nil
}

// This function is to get `to upgrade` version from upgrade-status infromation
// Syntax is: `kubectl-mayastor get upgrade-status`

func (cp CPv1) GetToUpgradeVersion() (string, error) {
	kubectlPlugin := GetPluginPath()

	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "get", "upgrade-status")
	toUpgradeVersionInfo, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("plugin failed to get `to upgrade` version, error %v", err)
	}

	out := strings.Split(string(toUpgradeVersionInfo), "\n")
	var toUpgradeVersion string
	for _, line := range out {
		if strings.Contains(line, "Upgrade To") {
			fields := strings.Fields(line)
			toUpgradeVersion = fields[len(fields)-1]
		}
	}
	return toUpgradeVersion, nil
}

func (cp CPv1) DeleteUpgrade() error {
	kubectlPlugin := GetPluginPath()

	cmd := exec.Command(kubectlPlugin, "-n", common.NSMayastor(), "delete", "upgrade")

	_, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("plugin failed to delete resources created by the upgrade process, error %v", err)
	}
	return nil
}
