package k8sinstall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/openebs/openebs-e2e/apps"
	"github.com/openebs/openebs-e2e/common"
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

func InstallProduct() error {
	err := installTheProduct()
	return err
}

func installTheProduct() error {

	var err error

	_ = CreateNamespace(common.NSOpenEBS())

	e2eCfg := e2e_config.GetConfig()
	logf.Log.Info("e2e_config.GetConfig()", "CONFIG", e2eCfg)

	_, err = k8stest.EnsureE2EAgent()
	if err != nil {
		return err
	}

	err = k8stest.GetNamespace(common.NSOpenEBS())
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"install",
		e2eCfg.Product.OpenEBSHelmReleaseName,
		"-n",
		common.NSOpenEBS(),
		e2eCfg.Product.OpenEBSHelmRepo,
	}

	// Remove the existing Helm repository before adding it back again.
	// This step ensures that any potential issues related to an outdated or corrupted repository configuration
	// are resolved. It also helps in cases where the repository URL or content has changed, ensuring that the
	// Helm repository is up-to-date with the latest charts. By removing and re-adding the repository, we make
	// sure that the subsequent Helm commands (like `helm install`) interact with the correct and current version
	// of the repository.
	err = apps.RemoveHelmRepository(e2e_config.GetConfig().Product.OpenEBSHelmRepoName, e2e_config.GetConfig().Product.OpenEBSHelmRepoUrl)
	if err != nil {
		logf.Log.Info("failed to remove helm repository")
	}

	err = apps.AddHelmRepository(e2e_config.GetConfig().Product.OpenEBSHelmRepoName, e2e_config.GetConfig().Product.OpenEBSHelmRepoUrl)
	if err != nil {
		logf.Log.Info("failed to add helm repository")
	}

	err = apps.UpdateHelmRepository(e2e_config.GetConfig().Product.OpenEBSHelmRepoName)
	if err != nil {
		logf.Log.Info("failed to update helm repository")
	}

	err = installOpenebs(common.NSOpenEBS(), cmdArgs)
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
		fmt.Sprintf("engines.replicated.mayastor.enabled=%v", e2eCfg.ReplicatedEngine),
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
