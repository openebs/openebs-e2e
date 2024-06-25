package upgrade

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common/mayastor/partial_rebuild"

	mcpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type userPromptMessages string

const (
	RebuildWarning                 userPromptMessages = "The cluster is rebuilding replica of some volumes"
	SkipSingleReplicaVolumeWarning userPromptMessages = "These single replica volumes may not be accessible during upgrade"
	CordonedNodeWarning            userPromptMessages = "One or more nodes in this cluster are in a Mayastor cordoned state"
)

const (
	VolSizeMb                        = 8192 // in Mb
	DefTimeoutSecs                   = 180  // in seconds
	WaitForRebuildTriggerTimeoutSecs = 60   // in seconds
	UpgradeJobCompeletionTimeOutSecs = 1800 // in seconds
	DefRebuildTimeoutSecs            = 600  // in seconds
	sleepTime                        = 3    // in seconds
	FioRunTime                       = 1800 // in seconds
	toLocalpvProvisionerImage        = "3.5.0"
)

// DisablePartialRebuildUpgradeVersions contains list of product
// version which requires disabling partial rebuild while upgrade
var DisablePartialRebuildUpgradeVersions = []string{
	"2.2.0", "2.3.0", "2.4.0", "2.5.0",
}

// here for no of loops we set default value as 10 and
// volsize as 4096. Now accroding to VolSizeMb in tests loops can be
// configured. For e.g. in upgrade tests we have 8192Mb size of volume
// then no of loops will become 5. Same way for lesser volume, no of loops will increase.
var NoOfFioRunLoops = CalculateNoOfFioRunLoops(VolSizeMb)

func CalculateNoOfFioRunLoops(VolSizeMb int) int {
	var loops int
	if VolSizeMb > 10*4096 {
		loops = 1
	} else {
		loops = (10 * 4096) / VolSizeMb
	}
	return loops
}

var (
	UpgradingControlPlane string = " Upgrading " + cases.Title(language.Und).String(e2e_config.GetConfig().Product.ProductName) + " control-plane"
	UpgradingDataPlane    string = " Upgrading " + cases.Title(language.Und).String(e2e_config.GetConfig().Product.ProductName) + " data-plane"
	UpgradeCompleted      string = " Successfully upgraded " + cases.Title(language.Und).String(e2e_config.GetConfig().Product.ProductName)
)

var MSDeployment = []string{
	e2e_config.GetConfig().Product.ControlPlanePoolOperator,
	e2e_config.GetConfig().Product.ControlPlaneRestServer,
	e2e_config.GetConfig().Product.ControlPlaneCsiController,
	e2e_config.GetConfig().Product.ControlPlaneCoreAgent,
	e2e_config.GetConfig().Product.ControlPlaneLocalpvProvisioner,
	e2e_config.GetConfig().Product.ControlPlaneObsCallhome,
}

var MSAppLabels = []string{
	e2e_config.GetConfig().AppLabelControlPlanePoolOperator,
	e2e_config.GetConfig().AppLabelControlPlaneRestServer,
	e2e_config.GetConfig().AppLabelControlPlaneCsiController,
	e2e_config.GetConfig().AppLabelControlPlaneCoreAgent,
}

var PodPrefixOfUpgradeJob = e2e_config.GetConfig().Product.ProductName + "-upgrade"
var LokiStatefulset = e2e_config.GetConfig().Product.LokiStatefulset
var LokiStatefulsetOnControlNode = e2e_config.GetConfig().LokiStatefulsetOnControlNode

type TestApp struct {
	App k8stest.FioApp
}

func AreContainerImagesUpgraded(podList *coreV1.PodList, toUpgradeImageTag, dockerImageOrgName string) (bool, error) {
	for _, pod := range podList.Items {
		for _, container := range pod.Spec.Containers {
			if strings.Contains(container.Image, dockerImageOrgName) {
				imageTag := strings.Split(container.Image, ":")
				if len(imageTag) != 2 {
					return false, fmt.Errorf("image didn't split successfully in name and tag parts")
				}
				logf.Log.Info("Container images are", "container name: ", container.Name, " container image: ", container.Image)

				if container.Name != "mayastor-localpv-provisioner" {
					if !strings.Contains(imageTag[len(imageTag)-1], toUpgradeImageTag) {
						return false, nil
					}
				} else if !strings.Contains(imageTag[len(imageTag)-1], toLocalpvProvisionerImage) {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func AreControlAndDataPlaneUpgraded(toUpgradeImageTag, dockerImageOrgName string) (bool, error) {
	podList, _ := k8stest.ListPod(common.NSMayastor())
	return AreContainerImagesUpgraded(podList, toUpgradeImageTag, dockerImageOrgName)
}

func IsControlPlaneUpgraded(toUpgradeImageTag, dockerImageOrgName string) (bool, error) {
	controlPlanePodList, _, err := k8stest.ListControlAndDataPlanePods()
	if err != nil {
		return false, fmt.Errorf("failed to list control plane pods in namespace: %s, err: %v", common.NSMayastor(), err)
	}
	return AreContainerImagesUpgraded(controlPlanePodList, toUpgradeImageTag, dockerImageOrgName)
}

func IsDataPlaneUpgraded(toUpgradeImageTag, dockerImageOrgName string) (bool, error) {
	_, dataPlanePodList, err := k8stest.ListControlAndDataPlanePods()
	if err != nil {
		return false, fmt.Errorf("failed to list data plane pods in namespace: %s, err: %v", common.NSMayastor(), err)

	}
	return AreContainerImagesUpgraded(dataPlanePodList, toUpgradeImageTag, dockerImageOrgName)
}

func VerifyMayastorAndPoolReady() (bool, error) {
	// mayastor ready check
	return k8stest.VerifyMayastorAndPoolReady(DefTimeoutSecs)
}

func RestartDataPlane(volUuid []string) error {
	ioEngineNodeList, err := k8stest.ListIOEngineNodes()
	if err != nil {
		return fmt.Errorf("failed to list nodes with io-engine label present: %v", err)
	}

	ioEnginePodList, err := k8stest.ListIOEnginePods()
	if err != nil {
		return fmt.Errorf("failed to list io-engine pods: %v", err)
	}

	if len(ioEnginePodList.Items) != len(ioEngineNodeList.Items) {
		return fmt.Errorf("no of io-engine pods didn't matched with desired count based on no. of nodes with io-engine label")
	}

	for _, ioEnginePod := range ioEnginePodList.Items {

		logf.Log.Info("restarting data plane pods: ", "io-engine pod: ", ioEnginePod.Name)
		err := k8stest.DeletePod(ioEnginePod.Name, common.NSMayastor())
		if err != nil {
			return fmt.Errorf("failed to delete pod %s, err: %v", ioEnginePod.Name, err)
		}

		isPodDeleted, err := checkIfIOEnginePodDeletedSuccessfully(ioEnginePod.Name)
		if err != nil {
			return fmt.Errorf("failed to check if io-engine pod %s deleted: %v", ioEnginePod.Name, err)
		}

		if isPodDeleted {
			logf.Log.Info("io-engine pod deleted successfully,", "podName: ", ioEnginePod.Name)
		} else {
			return fmt.Errorf("failed to delete io-engine pod %s deleted: %v", ioEnginePod.Name, err)
		}

		// Checking for the desired no of io-engine pods are present and should be in running state
		// before restarting the next io-engine pod.
		areIOPodsRunning, err := AreDesiredNoOfIOEnginePodsRunning(len(ioEngineNodeList.Items))
		if err != nil {
			return fmt.Errorf("failed to check if io-engines pods are running: %v", err)
		}

		if areIOPodsRunning {
			logf.Log.Info("all io-engine pods are running")
		} else {
			return fmt.Errorf("all io-engines pods are not in running phase: %v", err)
		}

		// wait for rebuild to complete if any replica is in rebuilding phase.
		logf.Log.Info("wait for rebuild to complete")
		for _, vol := range volUuid {
			isRebuildCompleted, err := partial_rebuild.WaitForRebuildComplete(vol, DefRebuildTimeoutSecs)
			if err != nil {
				return fmt.Errorf("failed to check rebuild completion, got error: %v", err)
			}

			if isRebuildCompleted {
				logf.Log.Info("rebuild got completed", "volume uuid: ", vol)
			} else {
				return fmt.Errorf("rebuild couldn't get completed for volume uuid: %q", vol)
			}
		}
	}
	return nil
}

func checkIfIOEnginePodDeletedSuccessfully(ioEnginePod string) (bool, error) {
	startTime := time.Now()
	for time.Since(startTime) < time.Duration(DefTimeoutSecs)*time.Second {
		newIOEnginePodList, err := k8stest.ListIOEnginePods()
		if err != nil {
			return false, fmt.Errorf("failed to list io-engine pods: %v", err)
		}

		isPodDeleted := true

		for _, pod := range newIOEnginePodList.Items {
			if pod.Name == ioEnginePod {
				isPodDeleted = false
				break
			}
		}

		if isPodDeleted {
			return true, nil
		}

		time.Sleep(sleepTime * time.Second)
	}
	return false, nil
}

func AreDesiredNoOfIOEnginePodsRunning(podCount int) (bool, error) {
	logf.Log.Info("Check status for the io-engine pods")
	startTime := time.Now()
	for time.Since(startTime) < time.Duration(DefTimeoutSecs)*time.Second {

		ioEnginePodList, err := k8stest.ListIOEnginePods()
		if err != nil {
			return false, fmt.Errorf("get list of io-engine pods was not successful: %v", err)
		}

		noOfIOEnginePods := len(ioEnginePodList.Items)
		allIOEnginePodsRunning := true

		for _, ioEnginePod := range ioEnginePodList.Items {
			if ioEnginePod.Status.Phase != coreV1.PodRunning {
				allIOEnginePodsRunning = false
				break
			}
		}

		if allIOEnginePodsRunning && noOfIOEnginePods == podCount {
			return true, nil
		}

		time.Sleep(sleepTime * time.Second)
	}
	return false, nil
}

func DeleteUpgradeResources() error {
	// delete upgrade resources created by the upgrade process
	logf.Log.Info("Deleting upgrade resources created by the upgrade process")
	err := controlplane.DeleteUpgrade()
	if err != nil {
		return fmt.Errorf("failed to delete upgrade resources, err:%v", err)
	}

	// delete the upgrade pod, created by the upgrade job
	logf.Log.Info("Delete the upgrade pod which was created by the upgrade job")
	upgradePodLabel := "app=" + e2e_config.GetConfig().Product.UpgradePodLabelValue
	err = k8stest.DeletePodsByLabel(upgradePodLabel, common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to delete upgrade pod, err:%v", err)
	}

	return nil
}

func BeforeEach() error {
	err := e2e_ginkgo.BeforeEachCheck()
	custom_resources.ClearDiskPoolCRDSelection()
	return err
}

func PostUpgradeWaitForRebuildCompletion(volUuid string) error {
	// verify rebuild is triggered
	// here we wait for 60 seconds to check if any rebuild is triggerd.
	// upgrade job also has the same timeout seconds to check for rebuilds.
	logf.Log.Info("wait for rebuild to get triggered")
	rebuildInProgress, err := partial_rebuild.WaitForRebuildInProgress(volUuid, WaitForRebuildTriggerTimeoutSecs)
	if err != nil {
		return fmt.Errorf("failed to check rebuild, err:%v", err)
	}
	// wait for rebuild to complete
	if rebuildInProgress {
		logf.Log.Info("wait for rebuild to complete")
		isRebuildCompleted, err := partial_rebuild.WaitForRebuildComplete(volUuid, DefRebuildTimeoutSecs)
		if err != nil {
			return fmt.Errorf("failed to check completion of rebuild")
		}
		if !isRebuildCompleted {
			return fmt.Errorf("rebuild not completed in given time of %d seconds", DefRebuildTimeoutSecs)
		}
	}
	return nil
}

// CheckIfUpgradingToUstableBranch checks the plugin version
// and then verifies if plugin version matches with the regex for
// unstable tag. if it matches then plugin will start upgrade to
// unstable tag with --allow-unstable flag. otherwise it will proceed
// with normal upgrade without any flag.
func CheckIfUpgradingToUnstableBranch() (string, bool, error) {

	pluginVersion, err := mcpV1.GetPluginVersion()
	if err != nil {
		return "", false, fmt.Errorf("failed to get plugin version, err:%v", err)
	}

	pluginVersion = strings.TrimSpace(pluginVersion)
	logf.Log.Info("kubectl mayastor plugin", "version", pluginVersion)

	// here starting "v?" part means that "v" is optional
	// regex will work for both tags, with and without starting v
	tagRegexForDevelopBranch := `v?[0-9]+\.[0-9]+\.[0-9]+-0-main-unstable(-[0-9]+){6}-0`
	tagRegexForPreReleaseTesting := `v?[0-9]+\.[0-9]+\.[0-9]+-0-release-unstable(-[0-9]+){6}-0`
	tagRegexForReleaseBranch := `v?[0-9]+\.[0-9]+\.[0-9]`

	// Construct a new regular expression by surrounding tagRegex with parentheses and appending "+0)"
	// this is needed because in the output format of kubectl-mayastor --version
	// tag is not present in exactly the same form of tagRegex
	pluginVersionOutputFormatRegexForDevelopBranch := fmt.Sprintf(`\(%s\+0\)`, tagRegexForDevelopBranch)
	pluginVersionOutputFormatRegexForPreReleaseTesting := fmt.Sprintf(`\(%s\+0\)`, tagRegexForPreReleaseTesting)
	pluginVersionOutputFormatRegexForReleaseBranch := fmt.Sprintf(`\(%s\+0\)`, tagRegexForReleaseBranch)

	// Create a regular expression object for the plugin version format regex for develop branch
	pluginRegexDevelopBranch, err := regexp.Compile(pluginVersionOutputFormatRegexForDevelopBranch)
	if err != nil {
		return pluginVersion, false, fmt.Errorf("failed to create valid regex, err:%v", err)
	}

	// Create a regular expression object for the plugin version format regex for pre release testing
	pluginRegexPreReleaseTesting, err := regexp.Compile(pluginVersionOutputFormatRegexForPreReleaseTesting)
	if err != nil {
		return pluginVersion, false, fmt.Errorf("failed to create valid regex, err:%v", err)
	}

	// Create a regular expression object for the plugin version format regex for release branch
	pluginRegexReleaseBranch, err := regexp.Compile(pluginVersionOutputFormatRegexForReleaseBranch)
	if err != nil {
		return pluginVersion, false, fmt.Errorf("failed to create valid regex, err:%v", err)
	}

	// Match plugin version with regular expressions
	// for develop branch and pre-release testing, version should be
	// considered as unstable and will need --allow-unstable flag
	if pluginRegexDevelopBranch.MatchString(pluginVersion) || pluginRegexPreReleaseTesting.MatchString(pluginVersion) {
		return pluginVersion, true, nil
	}

	// for released versions it should run without any flags
	if pluginRegexReleaseBranch.MatchString(pluginVersion) {
		return pluginVersion, false, nil
	}

	return pluginVersion, false, fmt.Errorf("plugin version does not match with expected format")

}

// VerifyUpgradedComponentsImages verifies that after upgrade
// component images are upgraded as expected. we match the image tags
// with the `Upgrade To` filed of kubectl-mayastor get upgrade-status command.
func VerifyUpgradedComponentsImages() (bool, error) {
	// Get the ToUpgrade version
	logf.Log.Info("Fetching `to upgrade` version from upgrade-status information")
	toUpgradeImageTag, err := controlplane.GetToUpgradeVersion()
	if err != nil {
		return false, fmt.Errorf("to upgrade version couldn't fetch successfully")
	}
	logf.Log.Info("`to upgrade` version is : ", "toUpgradeImageTag", toUpgradeImageTag)

	// check upgraded image tag from mayastor pods
	logf.Log.Info("Checking mayastor pods container images")
	AreControlAndDataPlaneUpgraded, err := AreControlAndDataPlaneUpgraded(toUpgradeImageTag, e2e_config.GetConfig().Product.DockerOrganisation)
	if err != nil {
		return false, fmt.Errorf("failed to verify upgrades images of components")
	}

	if !AreControlAndDataPlaneUpgraded {
		return false, nil
	}
	return true, nil
}

// CheckIfDiskPoolCRDVersionUpgraded verifies that after upgrade
// diskpool crd is in expected apiVersion
func CheckIfDiskPoolCRDVersionUpgraded(pluginVersion string, isUpgradingToUnstable bool) (bool, error) {

	logf.Log.Info("kubectl plugin", "version", pluginVersion)

	// pluginVersion argument passsed in this function is in format of `(v2.5.0+0)`
	// use strings trim function to get the exact value for plugin version
	var version string
	version = strings.TrimPrefix(pluginVersion, "(v")
	// in case if version doesn't have starting "v" then only trim "("
	version = strings.TrimPrefix(version, "(")
	version = strings.TrimSuffix(version, "+0)")

	var expectedAPIVersion string
	if !isUpgradingToUnstable {
		expectedAPIVersion = e2e_config.GetConfig().Product.DiskPoolAPIVersionMap[version]
	} else {
		expectedAPIVersion = e2e_config.GetConfig().Product.DiskPoolAPIVersionMap["develop"]
	}

	logf.Log.Info("diskpool crd version expected", "apiVersion", expectedAPIVersion)

	// Retrieve DiskPool CRD versions
	crdVersions, err := custom_resources.GetDiskpoolCrdVersions()
	if err != nil {
		return false, err
	}
	logf.Log.Info("diskpool crd version", "after upgrade", crdVersions)

	// Check if the DiskPool CRD is upgraded to the expected API version
	// api versin is stored in DiskPoolAPIVersionMap in this format --
	// `<CRDGroupName>/<version>`
	// here version can be -- v1beta2, v1beta1, v1alpha1 etc.
	crdGroupName := e2e_config.GetConfig().Product.CrdGroupName
	isDiskPoolCRDUpgraded := false
	for _, ver := range crdVersions {
		ver = crdGroupName + "/" + ver
		if ver == expectedAPIVersion {
			isDiskPoolCRDUpgraded = true
			break
		}
	}

	if !isDiskPoolCRDUpgraded {
		return false, nil
	}

	return true, nil
}

func IsPartialRebuildDisableNeeded(appVersion string) bool {
	for _, ver := range DisablePartialRebuildUpgradeVersions {
		if ver == appVersion {
			return true
		}
	}
	return false
}
