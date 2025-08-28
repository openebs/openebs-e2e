package v1

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// upgradeMetadata struct to hold plugin and registry information
type upgradeMetadata struct {
	KubectlPlugin    string
	PluginVersion    string
	CIRegistry       string
	IsMayastorPlugin bool
}

// upgradeOptions struct to hold various upgrade flags
type upgradeOptions struct {
	IsUpgradingToUnstableBranch   bool
	IsPartialRebuildDisableNeeded bool
	ExtraFlags                    []string
}

// Helper to fetch plugin info and CI_REGISTRY
func getUpgradeMetadata() (upgradeMetadata, error) {
	var meta upgradeMetadata
	var ok bool
	var err error
	meta.KubectlPlugin = GetPluginPath()         // This line is already present in the original code.
	meta.PluginVersion, err = GetPluginVersion() // This line is already present in the original code.
	if err != nil {
		return upgradeMetadata{}, fmt.Errorf("failed to get plugin version, err:%v", err)
	}
	meta.PluginVersion = strings.TrimSpace(meta.PluginVersion)
	meta.CIRegistry, ok = os.LookupEnv("CI_REGISTRY")
	if !ok {
		return upgradeMetadata{}, fmt.Errorf("environment variable CI_REGISTRY is not defined")
	}
	meta.IsMayastorPlugin = filepath.Base(meta.KubectlPlugin) == e2e_config.GetConfig().Product.MayastorPluginName

	return meta, nil
}

// Build the common upgrade command arguments
func buildUpgradeArgs(meta upgradeMetadata, opts upgradeOptions) []string {

	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade"}
	cmdArgs = append(cmdArgs, opts.ExtraFlags...)

	if opts.IsUpgradingToUnstableBranch {
		// for plugin version with rc tag i.e. release candidate
		// upgrade job images are not present in ci-registry, they are in
		// docker hub. so for those tags we dont need to add --registry flag.
		if !meta.IsMayastorPlugin {
			cmdArgs = append(cmdArgs, string(common.SkipUpgradePathValidationFlag))
		} else if strings.Contains(meta.PluginVersion, "-rc") {
			cmdArgs = append(cmdArgs, string(common.AllowUpgradeToUnstableBranchFlag))
		} else {
			cmdArgs = append(cmdArgs, "--registry", meta.CIRegistry, string(common.AllowUpgradeToUnstableBranchFlag))
		}
	}
	if opts.IsPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(common.DisablePartialRebuild))
	}
	return cmdArgs
}

// Generic upgrade command invoker
func (cp CPv1) runUpgrade(
	opts upgradeOptions) (string, error) {
	upgradeMeta, err := getUpgradeMetadata()
	if err != nil {
		return "", err
	}
	args := buildUpgradeArgs(
		upgradeMeta, upgradeOptions{
			IsUpgradingToUnstableBranch:   opts.IsUpgradingToUnstableBranch,
			IsPartialRebuildDisableNeeded: opts.IsPartialRebuildDisableNeeded,
			ExtraFlags:                    opts.ExtraFlags,
		},
	)
	cmd := exec.Command(upgradeMeta.KubectlPlugin, args...)
	logf.Log.Info("Executing", "command", strings.Join(cmd.Args, " "))
	var out []byte
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if out, err = cmd.Output(); err != nil {
		// If the command fails, return the error message.
		logf.Log.Info("Command failed", "command", strings.Join(cmd.Args, " "), "error", err.Error())
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr != "" {
			logf.Log.Info("Command stderr", "stderr", stderrStr)
			return stderrStr, fmt.Errorf("plugin failed to upgrade: %v", err)
		}
		logf.Log.Info("Command stdout", "output", string(out))
		return string(out), fmt.Errorf("plugin failed to upgrade: %v", err)
	}
	return string(out), nil
}

// Refactored upgrade commands -- easy to add new variants!
func (cp CPv1) Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error) {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
	}
	// No extra upgrade flags
	return cp.runUpgrade(opts)
}

func (cp CPv1) UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
		ExtraFlags:                    []string{string(common.SkipDataPlaneRestartFlag)},
	}
	_, err := cp.runUpgrade(opts)
	return err
}

func (cp CPv1) UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
		ExtraFlags:                    []string{string(common.SkipSingleReplicaValidationFlag)},
	}
	_, err := cp.runUpgrade(opts)
	return err
}

func (cp CPv1) UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
		ExtraFlags:                    []string{string(common.SkipReplicaRebuildFlag)},
	}
	_, err := cp.runUpgrade(opts)
	return err
}

func (cp CPv1) UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	opts := upgradeOptions{IsUpgradingToUnstableBranch: isUpgradingToUnstableBranch, IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded, ExtraFlags: []string{string(common.SkipCordonNodeValidationFlag)}}
	_, err := cp.runUpgrade(opts)
	return err
}

// Helper for command selection for upgrade status/delete
func getUpgradeStatusCmdArgs(cmdType string) ([]string, error) {
	kubectlPlugin := GetPluginPath()
	isMayastor := filepath.Base(kubectlPlugin) == e2e_config.GetConfig().Product.MayastorPluginName

	switch cmdType {
	case "status":
		if isMayastor {
			return []string{"-n", common.NSMayastor(), "get", "upgrade-status"}, nil
		}
		return []string{"-n", common.NSMayastor(), "upgrade", "status"}, nil

	case "delete":
		if isMayastor {
			return []string{"-n", common.NSMayastor(), "delete", "upgrade"}, nil
		}
		return []string{"upgrade", "-n", common.NSMayastor(), "delete"}, nil

	default:
		return nil, fmt.Errorf("invalid cmdType: %s", cmdType)
	}
}

// Unified upgrade status fetching
func (cp CPv1) GetUpgradeStatus() (string, error) {
	kubectlPlugin := GetPluginPath()
	args, err := getUpgradeStatusCmdArgs("status")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(kubectlPlugin, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("plugin failed to get upgrade status, error %v", err)
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "Upgrade Status") {
			parts := strings.Split(line, ":")
			logf.Log.Info("Output", "Upgrade Status", strings.TrimSpace(parts[len(parts)-1]))
			return strings.TrimSpace(parts[len(parts)-1]), nil
		}
	}
	return "", nil // or return an error if not found
}

func (cp CPv1) GetToUpgradeVersion() (string, error) {
	kubectlPlugin := GetPluginPath()
	args, err := getUpgradeStatusCmdArgs("status")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(kubectlPlugin, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("plugin failed to get `to upgrade` version, error %v", err)
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "Upgrade To") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return strings.TrimSpace(fields[len(fields)-1]), nil
			}
		}
	}
	return "", nil // or return error if not found
}

func (cp CPv1) DeleteUpgrade() error {
	kubectlPlugin := GetPluginPath()
	args, err := getUpgradeStatusCmdArgs("delete")
	if err != nil {
		return err
	}
	cmd := exec.Command(kubectlPlugin, args...)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("plugin failed to delete resources created by the upgrade process, error %v", err)
	}
	return nil
}
