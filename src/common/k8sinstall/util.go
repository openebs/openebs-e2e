package k8sinstall

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

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

	// err = k8stest.KubeCtlApplyYaml("e2e-proxy.yaml", proxyPath.GetE2EProxyPath())
	// if err != nil {
	// 	return err
	// }

	err = k8stest.GetNamespace("openebs")
	if err != nil {
		return err
	}

	// productName := e2e_config.GetConfig().Product.ProductName

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

	// helmRegistry := string(version["helm_registry_url"])
	// // chartVersion := e2e_config.GetConfig().Product.ChartVersion
	// // if _, haveChartVersion := version["chart_version"]; haveChartVersion {
	// // 	chartVersion = string(version["chart_version"])
	// // }
	// err = updateHelmChartDependency(chartDir)
	// if err != nil {
	// 	return err
	// }

	// outputDir := locations.GetGeneratedHelmYamlsDir()

	// // generate install yamls from helm chart
	// // generated yamls directory: /artifacts/sessions/{session-id}/charts/generated-yamls/bolt
	// err = generateHelmInstallYamls(chartDir, common.NSMayastor(), outputDir)
	// if err != nil {
	// 	return err
	// }

	// // copy values.yaml file to session directory
	// // values.yaml path: /artifacts/sessions/{session-id}/charts/generated-yamls
	// err = copyHelmValuesYaml(chartDir, outputDir)
	// if err != nil {
	// 	return err
	// }

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

	logf.Log.Info("About to execute: helm repo update", "openebs")
	cmd = exec.Command("helm", "repo", "update", "openebs")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update helm repository. Output: %s ", out)
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

func updateHelmChartDependency(dir string) error {
	var out []byte
	var err error
	for tryCount := 0; tryCount < 3; tryCount++ {
		time.Sleep(time.Second * time.Duration(tryCount*30))
		cmd := exec.Command("helm", "dependency", "update")
		cmd.Dir = dir
		logf.Log.Info("About to execute: helm dependency update ,", "path", cmd.Dir)
		out, err = cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		logf.Log.Info(fmt.Sprintf("failed to update helm dependency: directory: %s Output: %s : Error: %v", cmd.Dir, out, err))
	}
	return fmt.Errorf("failed to update helm dependency: directory: %s Output: %s : Error: %v", dir, out, err)
}
