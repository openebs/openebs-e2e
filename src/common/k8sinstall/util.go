package k8sinstall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// func CreateNamespace(namespace string) error {
// 	cmd := exec.Command("kubectl", "create", "namespace", namespace)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		logf.Log.Info("Error", "output", out)
// 	}
// 	return err
// }

// func deleteNamespace(namespace string) error {
// 	cmd := exec.Command("kubectl", "delete", "namespace", namespace)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		logf.Log.Info("Error", "output", out)
// 	}
// 	return err
// }

func InstallProduct() error {
	err := installTheProduct()
	// if err != nil {
	// 	k8stest.GenerateInstallSupportBundle()
	// }
	return err
}

func installTheProduct() error {

	var err error
	e2eCfg := e2e_config.GetConfig()

	logf.Log.Info("e2e_config.GetConfig()", "CONFIG", e2eCfg)

	_, err = k8stest.EnsureE2EAgent()
	if err != nil {
		return err
	}

	err = k8stest.GetNamespace("openebs")
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"install",
		"openebs",
		"-n",
		"openebs",
	}

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

	ready, err := k8stest.OpenEBSReady(10, 540)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("openebs installation is not ready")
	}

	// OpenEBS has been installed and is now ready for use.
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

	cmd := exec.Command("helm", cmdArgs...)
	logf.Log.Info("installHelmChart: About to execute: helm", "arguments", cmdArgs)
	logf.Log.Info(strings.Join(cmdArgs, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install product using helm chart: namespace: %s  Output: %s : Error: %v", namespace, out, err)
	}
	return nil
}
