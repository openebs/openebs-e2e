package k8sinstall

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/locations"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func CreateNamespace(namespace string) error {
	cmd := exec.Command("kubectl", "create", "namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("Error", "output", out)
	}
	return err
}

func deleteNamespace(namespace string) error {
	cmd := exec.Command("kubectl", "delete", "namespace", namespace)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("Error", "output", out)
	}
	return err
}

func InstallProduct() error {
	err := installTheProduct()
	if err != nil {
		k8stest.GenerateInstallSupportBundle()
	}
	return err
}

func installTheProduct() error {

	var err error
	e2eCfg := e2e_config.GetConfig()

	logf.Log.Info("e2e_config.GetConfig()", "CONFIG", e2eCfg)

	cpVersion := controlplane.Version()
	logf.Log.Info("Control Plane", "version", cpVersion)

	_, err = k8stest.EnsureE2EAgent()
	if err != nil {
		return err
	}

	err = k8stest.KubeCtlApplyYaml("e2e-proxy.yaml", proxyPath.GetE2EProxyPath())
	if err != nil {
		return err
	}

	err = k8stest.GetNamespace(common.NSOpenebs())
	if err != nil {
		return err
	}

	productName := e2e_config.GetConfig().Product.ProductName

	// if productName == "openebspro" {
	// 	// verify kubernetes secrets
	// 	err = k8stest.GetKubernetesSecret("login", common.NSMayastor())
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// chartDir, err := locations.GetHelmChartsDir()
	// if err != nil {
	// 	return err
	// }

	// // Get version.json unmarshalled data
	// version, err := locations.ParseVersionFile()
	// if err != nil {
	// 	return err
	// }
	//bundleVersion := version["version"]
	// installTypeMap, exists := version["install_type"]
	// var installType string
	// if exists {
	// 	installType = string(installTypeMap)
	// } else {
	// 	installType = Contained
	// }

	helmRegistry := string(version["helm_registry_url"])
	chartVersion := e2e_config.GetConfig().Product.ChartVersion
	if _, haveChartVersion := version["chart_version"]; haveChartVersion {
		chartVersion = string(version["chart_version"])
	}
	err = updateHelmChartDependency(chartDir)
	if err != nil {
		return err
	}

	outputDir := locations.GetGeneratedHelmYamlsDir()

	// generate install yamls from helm chart
	// generated yamls directory: /artifacts/sessions/{session-id}/charts/generated-yamls/bolt
	err = generateHelmInstallYamls(chartDir, common.NSMayastor(), outputDir)
	if err != nil {
		return err
	}

	// copy values.yaml file to session directory
	// values.yaml path: /artifacts/sessions/{session-id}/charts/generated-yamls
	err = copyHelmValuesYaml(chartDir, outputDir)
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"install",
		"openebs",
		"-n",
		"openebs",
	}

	// Install using helm registry if bundleVersion is 1.0.2

	logf.Log.Info("About to execute: helm repo remove")
	cmd := exec.Command("helm", "repo", "remove", "openebs")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("failed to remove helm repository", "output", out)
	}

	logf.Log.Info("About to execute: helm repo add", "openebs", "https://openebs.github.io/openebs")
	cmd = exec.Command("helm", "repo", "add", "openebs", "https://openebs.github.io/openebs")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add helm repository. Output: %s ", out)
	}

	err = installOpenebs("openebs", cmdArgs)
	if err != nil {
		return err
	}

	ready, err := k8stest.MayastorReady(10, 540)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation is not ready")
	}

	ready = k8stest.ControlPlaneReady(10, 180)
	if !ready {
		return fmt.Errorf("mayastor control plane installation is not ready")
	}

	// wait for mayastor node to be ready
	nodeReady, err := k8stest.MayastorNodesReady(5, 180)
	if err != nil {
		return err
	}
	if !nodeReady {
		return fmt.Errorf("all mayastor node are not ready")
	}

	logf.Log.Info("Checking whether product needs license", "product", productName)
	if productName == "openebspro" {
		if err = license.LicenseInstallIfRequired(); err != nil {
			//			return err
			logf.Log.Info("Ignoring license installation failure", "err", err)
		}
	}

	err = diskPoolConfigMap()
	if err != nil {
		return fmt.Errorf("failed to create diskpool config map %v", err)
	}
	if controlplane.CreatePoolOnInstall() {
		if !WaitForPoolCrd() {
			return fmt.Errorf("mayastor pool CRD is undefined")
		}

		// Now create configured pools on all nodes.
		err = k8stest.CreateConfiguredPools()
		if err != nil {
			return err
		}

		// Wait for pools to be online
		const timoSecs = 240
		const timoSleepSecs = 10
		for ix := 0; ix < timoSecs/timoSleepSecs; ix++ {
			time.Sleep(timoSleepSecs * time.Second)
			err = custom_resources.CheckAllMsPoolsAreOnline()
			if err == nil {
				break
			}
		}
		if err != nil {
			return fmt.Errorf("one or more pools are offline %v", err)
		}
	}
	if e2e_config.GetConfig().SetNexusRebuildVerifyOnInstall {
		err = k8stest.SetNexusRebuildVerify(true)
		if err != nil {
			return fmt.Errorf("set nexus rebuild verify failed %v", err)
		}
	}

	// Delete any storage classes created by the helm charts
	k8stest.ClearStorageClasses()

	// Mayastor/Bolt has been installed and is now ready for use.
	return nil
}

func installOpenebs(namespace string, cmdArgs []string) error {
	e2eCfg := e2e_config.GetConfig()
	cmdArgs = append(cmdArgs,
		"--set",
		"engines.replicated.mayastor.enabled=false",
	)

	if e2eCfg.ImagePullPolicy != "" {
		cmdArgs = append(cmdArgs,
			"--set",
			fmt.Sprintf("image.pullPolicy=%s", e2eCfg.ImagePullPolicy),
		)
	}

	// if e2eCfg.ImageTag != "" {
	// 	cmdArgs = append(cmdArgs,
	// 		"--set",
	// 		fmt.Sprintf("image.tag=%s", e2eCfg.ImageTag),
	// 	)
	// }

	cmd := exec.Command("helm", cmdArgs...)
	logf.Log.Info("installHelmChart: About to execute: helm", "arguments", cmdArgs)
	logf.Log.Info(strings.Join(cmdArgs, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install product using helm chart: namespace: %s  Output: %s : Error: %v", namespace, out, err)
	}
	return nil
}
