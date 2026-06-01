package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	mayastorUpgradeJobImageName = "mayastor-upgrade-job"
	openebsUpgradeJobImageName  = "openebs-upgrade-job"
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
	meta.KubectlPlugin = GetPluginPath()
	meta.PluginVersion, err = GetPluginVersion()
	if err != nil {
		return upgradeMetadata{}, fmt.Errorf("failed to get plugin version, err:%v", err)
	}
	meta.PluginVersion = strings.TrimSpace(meta.PluginVersion)
	meta.PluginVersion = strings.TrimPrefix(meta.PluginVersion, "(")
	if idx := strings.IndexByte(meta.PluginVersion, '+'); idx != -1 {
		meta.PluginVersion = meta.PluginVersion[:idx]
	}
	meta.CIRegistry, ok = os.LookupEnv("CI_REGISTRY")
	if !ok {
		return upgradeMetadata{}, fmt.Errorf("environment variable CI_REGISTRY is not defined")
	}
	meta.IsMayastorPlugin = filepath.Base(meta.KubectlPlugin) == e2e_config.GetConfig().Product.MayastorPluginName
	return meta, nil
}

// imageExistsOnDockerHub checks if org/image:tag exists on DockerHub
// using the Hub API, which works for public repos without authentication.
func imageExistsOnDockerHub(org, image, tag string) bool {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/%s/", org, image, tag)
	logf.Log.Info("Checking DockerHub for upgrade job image", "url", url)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return false
	}
	defer resp.Body.Close() //nolint:errcheck
	return resp.StatusCode == http.StatusOK
}

// imageExistsOnGhcr checks if namespace/image:tag is accessible on ghcr.io
// using an anonymous pull token and a manifest HEAD request.
func imageExistsOnGhcr(namespace, image, tag string) bool {
	tokenURL := fmt.Sprintf(
		"https://ghcr.io/token?scope=repository:%s/%s:pull&service=ghcr.io",
		namespace, image,
	)
	resp, err := http.Get(tokenURL) //nolint:noctx
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close() //nolint:errcheck

	var tokenResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil || tokenResp.Token == "" {
		return false
	}

	manifestURL := fmt.Sprintf("https://ghcr.io/v2/%s/%s/manifests/%s", namespace, image, tag)
	logf.Log.Info("Checking ghcr.io for upgrade job image", "url", manifestURL)
	req, err := http.NewRequest(http.MethodHead, manifestURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.Token)
	req.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp2.Body.Close() //nolint:errcheck
	return resp2.StatusCode == http.StatusOK
}

// detectUpgradeRegistry probes registries in priority order and returns the one
// that hosts the upgrade job image for the current plugin version.
// Returns "" for DockerHub (no --registry flag needed), common.GhcrRegistryName for ghcr.io,
// or meta.CIRegistry as the final fallback.
func detectUpgradeRegistry(meta upgradeMetadata) string {
	dockerOrg := e2e_config.GetConfig().Product.DockerOrganisation

	upgradeImage := mayastorUpgradeJobImageName
	ghcrNS := common.MayastorRegistryNameSpace
	if !meta.IsMayastorPlugin {
		upgradeImage = openebsUpgradeJobImageName
		ghcrNS = common.OpenEbsRegistryNameSpace
	}

	if imageExistsOnDockerHub(dockerOrg, upgradeImage, meta.PluginVersion) {
		logf.Log.Info("Upgrade job image found on DockerHub",
			"image", dockerOrg+"/"+upgradeImage, "tag", meta.PluginVersion)
		return ""
	}

	if imageExistsOnGhcr(ghcrNS, upgradeImage, meta.PluginVersion) {
		logf.Log.Info("Upgrade job image found on ghcr.io",
			"registry", common.GhcrRegistryName, "namespace", ghcrNS, "tag", meta.PluginVersion)
		return common.GhcrRegistryName
	}

	logf.Log.Info("Upgrade job image not found on DockerHub or ghcr.io, using CI registry",
		"registry", meta.CIRegistry, "tag", meta.PluginVersion)
	return meta.CIRegistry
}

// Build the common upgrade command arguments
func buildUpgradeArgs(meta upgradeMetadata, opts upgradeOptions) []string {
	cmdArgs := []string{"-n", common.NSMayastor(), "upgrade"}
	cmdArgs = append(cmdArgs, opts.ExtraFlags...)

	if opts.IsUpgradingToUnstableBranch {
		if !meta.IsMayastorPlugin {
			cmdArgs = append(cmdArgs, string(common.SkipUpgradePathValidationFlag))
		} else {
			cmdArgs = append(cmdArgs, "--registry", meta.CIRegistry, string(common.AllowUpgradeToUnstableBranchFlag))
		}
	} else {
		regNS := common.MayastorRegistryNameSpace
		setPrefix := "image"
		if !meta.IsMayastorPlugin {
			regNS = common.OpenEbsRegistryNameSpace
			setPrefix = "mayastor.image"
		}

		detectedReg := detectUpgradeRegistry(meta)
		if detectedReg != "" {
			cmdArgs = append(cmdArgs,
				"--registry", detectedReg,
				"--repo-namespace", regNS,
				"--set", fmt.Sprintf("%s.registry=%s,%s.repo=%s", setPrefix, detectedReg, setPrefix, regNS),
			)
		}
		// empty string → DockerHub: plugin uses it by default, no --registry flag needed
	}

	if opts.IsPartialRebuildDisableNeeded {
		cmdArgs = append(cmdArgs, "--set", string(common.DisablePartialRebuild))
	}
	return cmdArgs
}

// Generic upgrade command invoker
func (cp CPv1) runUpgrade(opts upgradeOptions) (string, error) {
	upgradeMeta, err := getUpgradeMetadata()
	if err != nil {
		return "", err
	}
	args := buildUpgradeArgs(upgradeMeta, opts)
	cmd := exec.Command(upgradeMeta.KubectlPlugin, args...)
	logf.Log.Info("Executing", "command", strings.Join(cmd.Args, " "))
	var out []byte
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if out, err = cmd.Output(); err != nil {
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

func (cp CPv1) Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error) {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
	}
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
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
		ExtraFlags:                    []string{string(common.SkipCordonNodeValidationFlag)},
	}
	_, err := cp.runUpgrade(opts)
	return err
}

// UpgradeWithExtraFlags issues an upgrade command with custom extra flags.
func (cp CPv1) UpgradeWithExtraFlags(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool, extraFlags []string) error {
	opts := upgradeOptions{
		IsUpgradingToUnstableBranch:   isUpgradingToUnstableBranch,
		IsPartialRebuildDisableNeeded: isPartialRebuildDisableNeeded,
		ExtraFlags:                    extraFlags,
	}
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
	return "", nil
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
	return "", nil
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
