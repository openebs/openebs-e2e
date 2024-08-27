package k8stest

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	mcpV1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	crtypes "github.com/openebs/openebs-e2e/common/custom_resources/types"
	agent "github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/locations"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	//"github.com/go-openapi/strfmt"

	errors "github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Helper for passing yaml from the specified directory to kubectl
func KubeCtlApplyYaml(filename string, dir string) error {
	cmd := exec.Command("kubectl", "apply", "-f", filename)
	cmd.Dir = dir
	logf.Log.Info("kubectl apply ", "yaml file", filename, "path", cmd.Dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply yaml file %s : Output: %s : Error: %v", filename, out, err)
	}
	return nil
}

// Helper for passing yaml from the specified directory to kubectl
func KubeCtlDeleteYaml(filename string, dir string) error {
	cmd := exec.Command("kubectl", "delete", "-f", filename)
	cmd.Dir = dir
	logf.Log.Info("kubectl delete ", "yaml file", filename, "path", cmd.Dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply yaml file %s : Output: %s : Error: %v", filename, out, err)
	}
	return nil
}

// create a storage class with default volume binding mode i.e. not specified
func MkStorageClass(scName string, scReplicas int, protocol common.ShareProto, nameSpace string) error {
	return NewScBuilder().
		WithName(scName).
		WithReplicas(scReplicas).
		WithProtocol(protocol).
		WithNamespace(nameSpace).
		BuildAndCreate()
}

// remove a storage class
func RmStorageClass(scName string) error {
	logf.Log.Info("Deleting storage class", "name", scName)
	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	deleteErr := ScApi().Delete(context.TODO(), scName, metaV1.DeleteOptions{})
	if k8serrors.IsNotFound(deleteErr) {
		return nil
	}
	return deleteErr
}

func CheckForStorageClasses() (bool, error) {
	found := false
	ScApi := gTestEnv.KubeInt.StorageV1().StorageClasses
	scs, err := ScApi().List(context.TODO(), metaV1.ListOptions{})
	for _, sc := range scs.Items {
		if sc.Provisioner == e2e_config.GetConfig().Product.CsiProvisioner {
			found = true
		}
	}
	return found, err
}

func MkNamespace(nameSpace string) error {
	logf.Log.Info("Creating", "namespace", nameSpace)
	nsSpec := coreV1.Namespace{ObjectMeta: metaV1.ObjectMeta{Name: nameSpace}}
	_, err := gTestEnv.KubeInt.CoreV1().Namespaces().Create(context.TODO(), &nsSpec, metaV1.CreateOptions{})
	return err
}

// EnsureNamespace ensure that a namespace exists, creates namespace if not found.
func EnsureNamespace(nameSpace string) error {
	_, err := gTestEnv.KubeInt.CoreV1().Namespaces().Get(context.TODO(), nameSpace, metaV1.GetOptions{})
	if err == nil {
		return nil
	}
	return MkNamespace(nameSpace)
}

func RmNamespace(nameSpace string) error {
	logf.Log.Info("Deleting", "namespace", nameSpace)
	err := gTestEnv.KubeInt.CoreV1().Namespaces().Delete(context.TODO(), nameSpace, metaV1.DeleteOptions{})
	return err
}

// Add a node selector to the given pod definition
func ApplyNodeSelectorToPodObject(pod *coreV1.Pod, label string, value string) {
	if pod.Spec.NodeSelector == nil {
		pod.Spec.NodeSelector = make(map[string]string)
	}
	pod.Spec.NodeSelector[label] = value
}

// Add a node selector to the deployment spec and apply
func ApplyNodeSelectorToDeployment(deploymentName string, namespace string, label string, value string) error {
	depApi := gTestEnv.KubeInt.AppsV1().Deployments
	deployment, err := depApi(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s : ns: %s : Error: %v", deploymentName, namespace, err)
	}
	if deployment.Spec.Template.Spec.NodeSelector == nil {
		deployment.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	deployment.Spec.Template.Spec.NodeSelector[label] = value
	_, err = depApi(namespace).Update(context.TODO(), deployment, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply node selector to deployment %s : ns: %s : Error: %v", deploymentName, namespace, err)
	}
	return nil
}

// Remove all node selectors from the deployment spec and apply
func RemoveAllNodeSelectorsFromDeployment(deploymentName string, namespace string) error {
	depApi := gTestEnv.KubeInt.AppsV1().Deployments
	deployment, err := depApi(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment %s : ns: %s : Error: %v", deploymentName, namespace, err)
	}
	if deployment.Spec.Template.Spec.NodeSelector != nil {
		deployment.Spec.Template.Spec.NodeSelector = nil
		_, err = depApi(namespace).Update(context.TODO(), deployment, metaV1.UpdateOptions{})
	}
	if err != nil {
		return fmt.Errorf("failed to remove node selector from deployment %s : ns: %s : Error: %v", deploymentName, namespace, err)
	}
	return nil
}

func SetReplication(appLabel string, namespace string, replicas *int32) error {
	depAPI := gTestEnv.KubeInt.AppsV1().Deployments
	stsAPI := gTestEnv.KubeInt.AppsV1().StatefulSets

	labels := "app=" + appLabel
	deployments, err := depAPI(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labels})
	if err != nil {
		return fmt.Errorf("failed to list deployment, namespace: %s, error: %v", namespace, err)
	}
	sts, err := stsAPI(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labels})
	if err != nil {
		return fmt.Errorf("failed to list statefulset, namespace: %s, error: %v", namespace, err)
	}
	if len(deployments.Items) == 1 {
		err = SetDeploymentReplication(deployments.Items[0].Name, namespace, replicas)
		if err != nil {
			return err
		}
	} else if len(sts.Items) == 1 {
		err = SetStatefulsetReplication(sts.Items[0].Name, namespace, replicas)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("app %s is not deployed as a deployment or sts", appLabel)
	}
	return nil
}

// Wait until all instances of the specified pod are absent from the given node
func WaitForPodAbsentFromNode(podNameRegexp string, namespace string, nodeName string, timeoutSeconds int) error {
	var validID = regexp.MustCompile(podNameRegexp)
	var podAbsent bool = false

	podApi := gTestEnv.KubeInt.CoreV1().Pods

	for i := 0; i < timeoutSeconds && !podAbsent; i++ {
		podAbsent = true
		time.Sleep(time.Second)
		podList, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			return errors.New("failed to list pods")
		}
		for _, pod := range podList.Items {
			if pod.Spec.NodeName == nodeName {
				if validID.MatchString(pod.Name) {
					podAbsent = false
					break
				}
			}
		}
	}
	if !podAbsent {
		return errors.New("timed out waiting for pod")
	}
	return nil
}

// Get the execution status of the given pod on node, or nil if it does not exist
func getPodOnNodeStatus(podNameRegexp string, namespace string, nodeName string) (*coreV1.PodPhase, error) {
	var validID = regexp.MustCompile(podNameRegexp)
	podAPI := gTestEnv.KubeInt.CoreV1().Pods
	podList, err := podAPI(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods , namespace: %s, error: %v", namespace, err)
	}
	for _, pod := range podList.Items {
		if pod.Spec.NodeName == nodeName && validID.MatchString(pod.Name) {
			return &pod.Status.Phase, nil
		}
	}
	return nil, nil // pod not found
}

// Get the execution status of the given pod, or nil if it does not exist
func getPodStatus(podNameRegexp string, namespace string) (*coreV1.PodPhase, error) {
	var validID = regexp.MustCompile(podNameRegexp)
	podAPI := gTestEnv.KubeInt.CoreV1().Pods
	podList, err := podAPI(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods , namespace: %s, error: %v", namespace, err)
	}
	for _, pod := range podList.Items {
		if validID.MatchString(pod.Name) {
			return &pod.Status.Phase, nil
		}
	}
	return nil, nil // pod not found
}

func GetPodAddress(podName string, nameSpace string) (string, error) {
	pod, err := gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Status.PodIP, err
}

// Wait until the instance of the specified pod is present and in the running
// state on the given node
func WaitForPodRunningOnNode(podNameRegexp string, namespace string, nodeName string, timeoutSeconds int) error {
	for i := 0; i < timeoutSeconds; i++ {
		stat, err := getPodOnNodeStatus(podNameRegexp, namespace, nodeName)
		if err != nil {
			return fmt.Errorf("failed to get pod status, podRegexp: %s, namespace: %s, nodename: %s, error: %v",
				podNameRegexp,
				namespace,
				nodeName,
				err,
			)
		}
		if stat != nil && *stat == coreV1.PodRunning {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timed out waiting for pod to be running")
}

// Wait until the instance of the specified pod is absent or not in the running
// state on the given node
func WaitForPodNotRunningOnNode(podNameRegexp string, namespace string, nodeName string, timeoutSeconds int) error {
	for i := 0; i < timeoutSeconds; i++ {
		stat, err := getPodOnNodeStatus(podNameRegexp, namespace, nodeName)
		if err != nil {
			return fmt.Errorf("failed to get pod status, podRegexp: %s, namespace: %s, nodename: %s, error: %v",
				podNameRegexp,
				namespace,
				nodeName,
				err,
			)
		}
		if stat == nil || *stat != coreV1.PodRunning {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timed out waiting for pod to stop running")
}

// Wait until the instance of the specified pod with prefix is present and in the running state
func WaitForPodRunning(podNameRegexp string, namespace string, timeoutSeconds int) error {
	logf.Log.Info("WaitForPodRunning", "podNameRegexp", podNameRegexp, "namespace", namespace)
	for i := 0; i < timeoutSeconds; i++ {
		stat, err := getPodStatus(podNameRegexp, namespace)
		if err != nil {
			return fmt.Errorf("failed to get pod status, podRegexp: %s, namespace: %s, error: %v",
				podNameRegexp,
				namespace,
				err,
			)
		}
		logf.Log.Info("WaitForPodRunning", "podNameRegexp", podNameRegexp, "namespace", namespace, "status", stat)
		if stat != nil && *stat == coreV1.PodRunning {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timed out waiting for pod to be running")
}

// PodPresentOnNode returns true if the pod is present on the given node
func PodPresentOnNode(prefix string, namespace string, nodeName string) (bool, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	podList, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list pod, podRegexp: %s, namespace: %s, nodename: %s, error: %v",
			prefix,
			namespace,
			nodeName,
			err,
		)
	}

	for _, pod := range podList.Items {
		if pod.Spec.NodeName == nodeName {
			if strings.HasPrefix(pod.Name, prefix) {
				logf.Log.Info("checking", "pod ->", pod.Name, "is on node->", pod.Spec.NodeName)
				return true, nil
			}
		}
	}
	return false, nil
}

func DeploymentReadyCount(deploymentName string, namespace string) (int, error) {
	deployment, err := gTestEnv.KubeInt.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get deployment", "name", deploymentName)
		return 0, fmt.Errorf("failed to get deployment %s", deploymentName)
	}
	logf.Log.Info("deployment", "name", deploymentName, "available instances", deployment.Status.ReadyReplicas)
	return int(deployment.Status.ReadyReplicas), nil
}

func DeploymentReady(deploymentName, namespace string) bool {
	deployment, err := gTestEnv.KubeInt.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get deployment", "error", err)
		return false
	}

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsV1.DeploymentAvailable {
			if condition.Status == coreV1.ConditionTrue {
				return true
			}
		}
	}
	return false
}

func DaemonSetReady(daemonName string, namespace string) bool {
	daemon, err := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).Get(
		context.TODO(),
		daemonName,
		metaV1.GetOptions{},
	)
	if err != nil {
		logf.Log.Info("Failed to get daemonset", "error", err)
		return false
	}

	status := daemon.Status
	logf.Log.Info("DaemonSet "+daemonName, "status", status)
	return status.DesiredNumberScheduled == status.CurrentNumberScheduled &&
		status.DesiredNumberScheduled == status.NumberReady &&
		status.DesiredNumberScheduled == status.NumberAvailable
}

func ControlPlaneReady(sleepTime int, duration int) bool {
	ready := false
	count := (duration + sleepTime - 1) / sleepTime

	if controlplane.MajorVersion() < 1 {
		logf.Log.Info("unsupported control plane", "version", controlplane.MajorVersion())
		return ready
	}
	nonControlPlaneComponents := []string{
		e2e_config.GetConfig().Product.DaemonsetName,
		e2e_config.GetConfig().Product.CsiDaemonsetName,
	}

	logf.Log.Info("ControlPlaneReady: ", "count", count, "sleepTime", sleepTime)
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		deployments, err := gTestEnv.KubeInt.AppsV1().Deployments(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			continue
		}
		daemonsets, err := gTestEnv.KubeInt.AppsV1().DaemonSets(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			continue
		}
		statefulsets, err := gTestEnv.KubeInt.AppsV1().StatefulSets(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
		if err != nil {
			continue
		}
		ready = true
		for _, deployment := range deployments.Items {
			if contains(nonControlPlaneComponents, deployment.Name) {
				logf.Log.Info("ControlPlaneReady: skipping data plane", "deployment", deployment.Name)
				continue
			}
			tmp := DeploymentReady(deployment.Name, common.NSMayastor())
			logf.Log.Info("ControlPlaneReady: mayastor control plane", "deployment", deployment.Name, "ready", tmp)
			ready = ready && tmp
		}
		for _, daemon := range daemonsets.Items {
			if contains(nonControlPlaneComponents, daemon.Name) {
				logf.Log.Info("ControlPlaneReady: skipping data plane", "daemonset", daemon.Name)
				continue
			}
			tmp := DaemonSetReady(daemon.Name, common.NSMayastor())
			logf.Log.Info("ControlPlaneReady: mayastor control plane", "daemonset", daemon.Name, "ready", tmp)
			ready = ready && tmp
		}
		for _, statefulSet := range statefulsets.Items {
			if contains(nonControlPlaneComponents, statefulSet.Name) {
				logf.Log.Info("ControlPlaneReady: skipping data plane", "statefulset", statefulSet.Name)
				continue
			}
			tmp := StatefulSetReady(statefulSet.Name, common.NSMayastor())
			logf.Log.Info("ControlPlaneReady: mayastor control plane", "statefulset", statefulSet.Name, "ready", tmp)
			ready = ready && tmp
		}
		logf.Log.Info("ControlPlaneReady: mayastor control plane", "ready", ready)
	}
	return ready
}

func contains(list []string, str string) bool {
	for _, elem := range list {
		if elem == str {
			return true
		}
	}
	return false
}

// MayastorReady checks if the product installation is ready
func MayastorReady(sleepTime int, duration int) (bool, error) {
	logf.Log.Info("MayastorReady", "sleepTime", sleepTime, "duration", duration)
	var ready bool
	var err error
	count := (duration + sleepTime - 1) / sleepTime
	ready, err = readyCheck(common.NSMayastor(), true)
	// variable to throttle verbosity
	verboseTick := int(math.Max(1, 10/float64(sleepTime)))
	for ix := 1; ix <= count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		ready, err = readyCheck(common.NSMayastor(), ix%verboseTick == 0)
	}
	if !ready {
		_, _ = readyCheck(common.NSMayastor(), true)
	}
	logf.Log.Info("MayastorReady", "ready", ready, "error", err)
	if !ready {
		return ready, err
	}
	var podReady bool
	// check for all pods to be running
	for ix := 1; ix <= 30 && !podReady; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		podReady, err = PodReadyCheck(common.NSMayastor())
	}
	logf.Log.Info("MayastorPodReady", "podReady", podReady, "error", err)
	return ready, err
}

// OpenEBSReady checks if the product installation is ready
func OpenEBSReady(sleepTime int, duration int) (bool, error) {
	logf.Log.Info("OpenEBSReady", "sleepTime", sleepTime, "duration", duration)
	var ready bool
	var err error
	count := (duration + sleepTime - 1) / sleepTime
	ready, err = readyCheck(common.NSOpenEBS(), true)
	// variable to throttle verbosity
	verboseTick := int(math.Max(1, 10/float64(sleepTime)))
	for ix := 1; ix <= count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		ready, err = readyCheck(common.NSOpenEBS(), ix%verboseTick == 0)
	}
	if !ready {
		_, _ = readyCheck(common.NSOpenEBS(), true)
	}
	logf.Log.Info("OpenEBSReady", "ready", ready, "error", err)
	if !ready {
		return ready, err
	}
	var podReady bool
	// check for all pods to be running
	for ix := 1; ix <= 30 && !podReady; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		podReady, err = PodReadyCheck(common.NSOpenEBS())
	}
	logf.Log.Info("OpenEBSPodsReady", "podReady", podReady, "error", err)
	return ready, err
}

func getClusterMayastorNodeIPAddrs() ([]string, error) {
	var nodeAddrs []string
	nodes, err := GetIOEngineNodes()
	if err != nil {
		return nodeAddrs, err
	}

	for _, node := range nodes {
		nodeAddrs = append(nodeAddrs, node.IPAddress)
	}
	return nodeAddrs, err
}

// ListPoolsInCluster use mayastorclient to enumerate the set of mayastor pools present in the cluster
func ListPoolsInCluster() ([]mayastorclient.MayastorPool, error) {
	nodeAddrs, err := getClusterMayastorNodeIPAddrs()
	if err == nil {
		return mayastorclient.ListPools(nodeAddrs)
	}
	return []mayastorclient.MayastorPool{}, err
}

// ListNexusesInCluster use mayastorclient to enumerate the set of mayastor nexuses present in the cluster
func ListNexusesInCluster() ([]mayastorclient.MayastorNexus, error) {
	nodeAddrs, err := getClusterMayastorNodeIPAddrs()
	if err == nil {
		return mayastorclient.ListNexuses(nodeAddrs)
	}
	return []mayastorclient.MayastorNexus{}, err
}

// ListReplicasInCluster use mayastorclient to enumerate the set of mayastor replicas present in the cluster
func ListReplicasInCluster() ([]mayastorclient.MayastorReplica, error) {
	nodeAddrs, err := getClusterMayastorNodeIPAddrs()
	if err == nil {
		return mayastorclient.ListReplicas(nodeAddrs)
	}
	return []mayastorclient.MayastorReplica{}, err
}

// RmReplicasInCluster use mayastorclient to remove mayastor replicas present in the cluster
func RmReplicasInCluster() error {
	nodeAddrs, err := getClusterMayastorNodeIPAddrs()
	if err == nil {
		return mayastorclient.RmNodeReplicas(nodeAddrs)
	}
	return err
}

// ListNvmeControllersInCluster use mayastorclient to enumerate the set of mayastor nvme controllers present in the cluster
func ListNvmeControllersInCluster() ([]mayastorclient.NvmeController, error) {
	nodeAddrs, err := getClusterMayastorNodeIPAddrs()
	if err == nil {
		return mayastorclient.ListNvmeControllers(nodeAddrs)
	}
	return []mayastorclient.NvmeController{}, err
}

// GetPoolUsageInCluster use mayastorclient to enumerate the set of pools and sum up the pool usage in the cluster
func GetPoolUsageInCluster() (uint64, error) {
	var poolUsage uint64
	pools, err := ListPoolsInCluster()
	if err == nil {
		for _, pool := range pools {
			poolUsage += pool.GetUsed()
		}
	}
	return poolUsage, err
}

func CreateConfiguredPoolsAndWait(timeoutsecs int) error {
	logf.Log.Info("CreateConfiguredPoolsAndWait")
	err := CreateConfiguredPools()
	if err != nil {
		return err
	}
	const sleepTime = 5
	for ix := 1; ix < timeoutsecs/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		err := custom_resources.CheckAllMsPoolsAreOnline()
		if err == nil {
			break
		}
	}
	return custom_resources.CheckAllMsPoolsAreOnline()
}

func DeleteConfiguredPools() error {
	logf.Log.Info("DeleteConfiguredPools")
	var err error
	/*
		_, err = DeleteAllPoolFinalizers()
		if err != nil {
			return fmt.Errorf("failed to delete pool finalizers, error; %v", err)
		}
	*/

	deletedAllPools := DeleteAllPools()
	if !deletedAllPools {
		replicas, err := mayastorclient.ListReplicas(GetMayastorNodeIPAddresses())
		if err != nil {
			logf.Log.Error(err, err.Error())
		}
		if len(replicas) != 0 {
			for _, replica := range replicas {
				logf.Log.Info("DeleteAllPools: found replica:", "uuid", replica.GetUuid(), "pool", replica.GetPool(), "uri", replica.GetUri())
			}
		}
		return fmt.Errorf("failed to delete all pools")
	}

	const sleepTime = 5
	pools := []mayastorclient.MayastorPool{}
	for ix := 1; ix < 120/sleepTime; ix++ {
		pools, err = mayastorclient.ListPools(GetMayastorNodeIPAddresses())
		if err != nil {
			logf.Log.Info("ListPools", "error", err)
		}
		if len(pools) == 0 && err == nil {
			return nil
		}
		time.Sleep(sleepTime * time.Second)
		logf.Log.Info("DeleteConfiguredPools", "pools", pools)
	}
	if err != nil {
		return fmt.Errorf("failed to list pools, error; %v", err)
	}

	if len(pools) != 0 {
		return fmt.Errorf("failed to destroy all pools, existing pool: %v,error; %v", pools, err)
	}
	return err
}

// RestoreConfiguredPools (re)create pools as defined by the configuration.
// As part of the tests we may modify the pools, in such test cases
// the test should delete all pools and recreate the configured set of pools.
const defaultTimeoutSecs = 120

func RestoreConfiguredPools() error {
	logf.Log.Info("RestoreConfiguredPools")
	err := DeleteConfiguredPools()
	if err != nil {
		return err
	}
	err = CreateConfiguredPoolsAndWait(defaultTimeoutSecs)
	return err
}

func WaitForPoolsToBeOnline(timeoutSeconds int) error {
	const sleepTime = 5
	for ix := 1; ix < (timeoutSeconds+sleepTime)/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		err := custom_resources.CheckAllMsPoolsAreOnline()
		if err == nil {
			return nil
		}
	}
	return custom_resources.CheckAllMsPoolsAreOnline()
}

// WaitPodComplete waits until pod is in completed state
func WaitPodComplete(podName string, sleepTimeSecs, timeoutSecs int) error {
	var podPhase coreV1.PodPhase
	var err error

	logf.Log.Info("Waiting for pod to complete", "name", podName, "timeout secs", timeoutSecs)
	for elapsedTime := 0; elapsedTime <= timeoutSecs; elapsedTime += sleepTimeSecs {
		time.Sleep(time.Duration(sleepTimeSecs) * time.Second)
		podPhase, err = CheckPodCompleted(podName, common.NSDefault)
		logf.Log.Info("WaitPodComplete",
			"elapsed", fmt.Sprintf("%d/%d", elapsedTime+sleepTimeSecs, timeoutSecs),
			"podPhase", podPhase, "error", err)
		if err != nil {
			return fmt.Errorf("failed to access pod status %s %v", podName, err)
		}
		if podPhase == coreV1.PodSucceeded {
			return nil
		} else if podPhase == coreV1.PodFailed {
			break
		}
	}
	return errors.Errorf("pod did not complete, phase %v", podPhase)
}

// WaitFioPodComplete waits until pod is in completed state
func WaitFioPodComplete(podName string, sleepTimeSecs, timeoutSecs int) error {
	var podPhase coreV1.PodPhase
	var err error
	var podSynopsis *common.E2eFioPodLogSynopsis

	logf.Log.Info("Waiting for pod to complete", "name", podName, "timeout secs", timeoutSecs)
	for elapsedTime := 0; elapsedTime <= timeoutSecs; elapsedTime += sleepTimeSecs {
		time.Sleep(time.Duration(sleepTimeSecs) * time.Second)
		podPhase, podSynopsis, err = CheckFioPodCompleted(podName, common.NSDefault)
		logf.Log.Info("WaitFioPodComplete",
			"elapsed", fmt.Sprintf("%d/%d", elapsedTime+sleepTimeSecs, timeoutSecs),
			"podPhase", podPhase, "error", err)
		if err != nil {
			return fmt.Errorf("failed to access pod status %s %v", podName, err)
		}
		if podPhase == coreV1.PodSucceeded {
			return nil
		} else if podPhase == coreV1.PodFailed {
			break
		}
	}
	return errors.Errorf("pod did not complete, phase %v. %s", podPhase, podSynopsis)
}

// DeleteVolumeAttachmets deletes volume attachments for a node
func DeleteVolumeAttachments(nodeName string) error {
	volumeAttachments, err := gTestEnv.KubeInt.StorageV1().VolumeAttachments().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list volume attachments, error: %v", err)
	}
	if len(volumeAttachments.Items) == 0 {
		return nil
	}
	for _, volumeAttachment := range volumeAttachments.Items {
		if volumeAttachment.Spec.NodeName != nodeName {
			continue
		}
		logf.Log.Info("DeleteVolumeAttachments: Deleting", "volumeAttachment", volumeAttachment.Name)
		delErr := gTestEnv.KubeInt.StorageV1().VolumeAttachments().Delete(context.TODO(), volumeAttachment.Name, metaV1.DeleteOptions{})
		if delErr != nil {
			logf.Log.Info("DeleteVolumeAttachments: failed to delete the volumeAttachment", "volumeAttachment", volumeAttachment.Name, "error", delErr)
			return err
		}
	}
	return nil
}

// CheckAndSetControlPlane checks which deployments exists and sets config control plane setting
func CheckAndSetControlPlane() error {
	var err error
	var foundCoreAgents = false
	var version string

	// Check for core-agents either as deployment or statefulset to correctly handle older builds of control plane
	// which use core-agents deployment and newer builds which use core-agents statefulset
	_, err = gTestEnv.KubeInt.AppsV1().Deployments(common.NSMayastor()).Get(
		context.TODO(),
		e2e_config.GetConfig().Product.ControlPlaneAgent,
		metaV1.GetOptions{},
	)
	if err == nil {
		foundCoreAgents = true
	}

	_, err = gTestEnv.KubeInt.AppsV1().StatefulSets(common.NSMayastor()).Get(
		context.TODO(),
		e2e_config.GetConfig().Product.ControlPlaneAgent,
		metaV1.GetOptions{},
	)
	if err == nil {
		foundCoreAgents = true
	}

	if !foundCoreAgents {
		return fmt.Errorf("restful Control plane components are absent, name: %s, namespace: %s", e2e_config.GetConfig().Product.ControlPlaneAgent, common.NSMayastor())
	}

	version = "1.0.0"

	logf.Log.Info("CheckAndSetControlPlane", "version", version)
	if !e2e_config.SetControlPlane(version) {
		return fmt.Errorf("failed to setup config control plane to %s", version)
	}
	return nil
}

// MayastorNodesReady Checks if the requisite number of mayastor nodes are online.
func MayastorNodesReady(sleepTime int, duration int) (bool, error) {
	ready := false
	count := (duration + sleepTime - 1) / sleepTime
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		readyCount := 0
		// list mayastor node
		nodeList, err := ListMsns()
		if err != nil {
			logf.Log.Info("MayastorNodesReady: failed to list mayastor node", "error", err)
			return ready, err
		}
		for _, node := range nodeList {
			if node.State.Status == controlplane.NodeStateOnline() {
				readyCount++
			} else {
				logf.Log.Info("Not ready node", "node", node.Name, "status", node.State.Status)
			}
		}
		ready = DaemonSetReady(e2e_config.GetConfig().Product.DaemonsetName, common.NSMayastor())
		logf.Log.Info("mayastor node status",
			"MayastorNodes", len(nodeList),
			"MaystorNodeReadyCount", readyCount,
			"Daemonset ready", ready,
		)
	}
	return ready, nil
}

func GetNamespace(namespace string) error {
	_, err := gTestEnv.KubeInt.CoreV1().Namespaces().Get(context.TODO(), namespace, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace: %s, error: %v",
			namespace,
			err,
		)
	}
	return nil
}

func GetSystemNamespaceUuid() (string, error) {
	ns, err := gTestEnv.KubeInt.CoreV1().Namespaces().Get(context.TODO(), "kube-system", metaV1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get namespace: %s, error: %v",
			"kube-system",
			err,
		)
	}
	return string(ns.ObjectMeta.UID), nil
}

func GetKubernetesSecret(secret string, namespace string) error {
	_, err := gTestEnv.KubeInt.CoreV1().Secrets(namespace).Get(context.TODO(), secret, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get kuberenetes secret: %s, namespace, namespace: %s, error: %v",
			secret,
			namespace,
			err,
		)
	}
	return nil
}

func DeleteKubernetesSecret(secret string, namespace string) error {
	err := gTestEnv.KubeInt.CoreV1().Secrets(namespace).Delete(context.TODO(), secret, metaV1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete kuberenetes secret: %s, namespace: %s, error: %v",
			secret,
			namespace,
			err,
		)
	}
	return nil
}

func CreatePoolOnNodeAndWait(nodeName string, poolName string, timeoutsecs int) (bool, error) {
	err := CreatePoolOnNode(nodeName, poolName)
	state := false
	if err != nil {
		return state, err
	}
	const sleepTime = 5
	for ix := 1; ix < timeoutsecs/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		state, err := custom_resources.CheckDiskPoolsIsOnline(poolName)
		if state && err == nil {
			break
		}
	}
	return custom_resources.CheckDiskPoolsIsOnline(poolName)
}

// CreatePoolOnNode create pool on a specific node.
// No check is made on the status of pools
func CreatePoolOnNode(nodeName string, poolName string) error {
	disks, err := GetConfiguredNodePoolDevices(nodeName)
	if err == nil {
		// NO check is made on the status of pools
		var pool crtypes.DiskPool
		pool, err = custom_resources.CreateMsPool(poolName, nodeName, disks)
		if err == nil {
			logf.Log.Info("Created", "pool", pool)
		} else {
			err = fmt.Errorf("failed to create pool on %v , disks: %s, error: %v", nodeName, disks, err)
		}
	} else {
		logf.Log.Info("failed to retrieve configured diskpool devices", "node", nodeName)
	}
	return err
}

// Create meta.json file with release platform and install bundle details
func GenerateMetaJson() error {
	var err error
	testLogDir, ok := os.LookupEnv("e2etestlogdir")
	if !ok {
		testLogDir = "/tmp/e2e/logs"
	}

	imageName, imageOk := os.LookupEnv("e2e_image_tag")
	if !imageOk {
		logf.Log.Info("imageName not present", "error", imageOk)
	}
	platformConfig, platformOk := os.LookupEnv("e2e_platform_config_file")
	if !platformOk {
		logf.Log.Info("platformConfig not present", "error", platformOk)
	}
	regexPattern := `v(\d+\.\d+\.\d+)`
	re := regexp.MustCompile(regexPattern)
	match := re.FindStringSubmatch(imageName)
	var release string
	if strings.Contains(imageName, "develop") || strings.Contains(imageName, "unstable") {
		release = "develop"
	} else if len(match) == 0 {
		err = fmt.Errorf("failed to extract release information from install bundle name")
		return err
	} else {
		release = match[1]
		if match[1] == "0.0.0" {
			release = "develop"
		}
	}
	regexPattern = `\b(\w+)\.`
	re = regexp.MustCompile(regexPattern)
	match = re.FindStringSubmatch(platformConfig)

	// run e2e-meta-json.sh
	bashCmd := fmt.Sprintf("%s/e2e-meta-json.sh  --destdir '%s' --release '%s' --platform '%s' --bundle '%s'",
		locations.GetE2EScriptsPath(),
		testLogDir,
		release,
		match[1],
		imageName,
	)
	logf.Log.Info("About to execute", "command", bashCmd)
	cmd := exec.Command("bash", "-c", bashCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
		logf.Log.Info(out.String())
	}
	return err
}

// Generate Application Pod logs
func GenerateAppPodLogs(token string) {
	var err error
	logDir, ok := os.LookupEnv("e2etestlogdir")
	if !ok {
		logDir = "/tmp/e2e/logs"
	}
	t0 := time.Now().UTC()
	ts := fmt.Sprintf("%v%d%v%v%v", t0.Year(), t0.Month(), t0.Hour(), t0.Minute(), t0.Second())
	// generate app pod logs in a directory specific for the current test case
	testLogDir := fmt.Sprintf("%s/%s/%s", logDir, strings.Map(common.SanitizePathname, token), ts)
	// run e2e-pod-logs.sh first - it will make the test log directory if it does not exist
	bashCmd := fmt.Sprintf("%s/e2e-pod-logs.sh  --destdir '%s'",
		locations.GetE2EScriptsPath(),
		testLogDir,
	)
	logf.Log.Info("About to execute", "command", bashCmd)
	cmd := exec.Command("bash", "-c", bashCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
		logf.Log.Info(out.String())
	}
}

// GenerateSupportBundle generate a support bundle for the cluster
func GenerateSupportBundle(testLogDir string) {
	var err error
	// run e2e-cluster-dump.sh first - it will make the test log directory if it does not exist
	bashCmd := fmt.Sprintf("%s/e2e-cluster-dump.sh  --destdir '%s' --plugin '%s'",
		locations.GetE2EScriptsPath(),
		testLogDir,
		mcpV1.GetPluginPath(),
	)
	logf.Log.Info("About to execute", "command", bashCmd)
	cmd := exec.Command("bash", "-c", bashCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
		logf.Log.Info(out.String())
	}
	bashCmd = fmt.Sprintf("%s -n %s get volume-replica-topologies -o json > %s/%s",
		mcpV1.GetPluginPath(),
		common.NSMayastor(),
		testLogDir,
		"replica-topologies.json")
	logf.Log.Info("About to execute", "command", bashCmd)
	cmd = exec.Command("bash", "-c", bashCmd)
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
		logf.Log.Info(out.String())
	}
	cmd = exec.Command(mcpV1.GetPluginPath(), "dump", "system", "-n", common.NSMayastor(), "-d", testLogDir)
	logf.Log.Info("About to execute", "command", cmd)
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
	}
}

// GenerateInstallSupportBundle generate a support bundle for the cluster
func GenerateInstallSupportBundle() {
	var err error
	logDir, ok := os.LookupEnv("e2etestlogdir")
	if !ok {
		logDir = "/tmp/e2e/logs"
	}
	t0 := time.Now().UTC()
	ts := fmt.Sprintf("%v%d%v%v%v", t0.Year(), t0.Month(), t0.Hour(), t0.Minute(), t0.Second())
	// create support bundle in a directory specific for the current test case
	testLogDir := fmt.Sprintf("%s/%s/%s", logDir, strings.Map(common.SanitizePathname, "install"), ts)
	// run e2e-cluster-dump.sh first - it will make the test log directory if it does not exist
	bashCmd := fmt.Sprintf("%s/e2e-cluster-dump.sh  --clusteronly --destdir '%s'",
		locations.GetE2EScriptsPath(),
		testLogDir,
	)
	logf.Log.Info("About to execute", "command", bashCmd)
	cmd := exec.Command("bash", "-c", bashCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		logf.Log.Info("command failed", "error", err)
		logf.Log.Info(out.String())
	}
}

// DiscoverProduct return the name of the product installed on the test cluster
// this cannot be used by the product install test
// This function may be called very early during initialisation of a test ot
// tool run and should have no dependencies of state
func DiscoverProduct() (string, error) {
	productMap := e2e_config.GetProductsSpecsMap()
	var err error

	restConfig := config.GetConfigOrDie()
	if restConfig == nil {
		return "", fmt.Errorf("failed to create *rest.Config for talking to a Kubernetes apiserver : GetConfigOrDie")
	}
	kubeInt := kubernetes.NewForConfigOrDie(restConfig)
	if kubeInt == nil {
		return "", fmt.Errorf("failed to create new Clientset for the given config : NewForConfigOrDie")
	}

	for name, product := range productMap {

		// Check for core-agents either as deployment or statefulset to correctly handle older builds of control plane
		// which use core-agents deployment and newer builds which use core-agents statefulset
		fmt.Fprintln(os.Stderr, product.ProductNamespace)
		_, err = kubeInt.AppsV1().Deployments(product.ProductNamespace).Get(
			context.TODO(),
			product.ControlPlaneAgent,
			metaV1.GetOptions{},
		)
		if err == nil {
			return name, nil
		}

		_, err = kubeInt.AppsV1().StatefulSets(product.ProductNamespace).Get(
			context.TODO(),
			product.ControlPlaneAgent,
			metaV1.GetOptions{},
		)
		if err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("Product not found on cluster")
}

func CreateDiskPartitionsMib(addr string, count int, partitionSizeInMiB int, diskPath string) error {
	start := 1
	partedArgs := "--script " + diskPath + " mklabel gpt"
	logf.Log.Info("Labelling disk before partitioning", "addr", addr, "arguments", partedArgs)
	err := agent.DiskPartition(addr, partedArgs)
	if err != nil {
		return fmt.Errorf("failed to label disk on node %s , error:  %v", addr, err)
	}
	for i := 0; i < count; i++ {
		partedArgs := "--script " + diskPath + " mkpart primary ext4 " + strconv.Itoa(start) + "MiB" + " " + strconv.Itoa(start+partitionSizeInMiB) + "MiB"
		start += partitionSizeInMiB
		logf.Log.Info("Creating partition", "addr", addr, "arguments", partedArgs)
		err := agent.DiskPartition(addr, partedArgs)
		if err != nil {
			return fmt.Errorf("disk Partitioning failed for node %s , error:  %v", addr, err)
		}
	}
	return nil
}

func CreateDiskPartitions(addr string, count int, partitionSizeInGiB int, diskPath string) error {
	return CreateDiskPartitionsMib(addr, count, partitionSizeInGiB*1024, diskPath)
}

func DeleteDiskPartitions(addr string, count int, diskPath string) error {
	for i := 0; i < count; i++ {
		partedArgs := "--script " + diskPath + " rm " + strconv.Itoa(i+1)
		logf.Log.Info("Deleting partition", "addr", addr, "arguments", partedArgs)
		err := agent.DiskPartition(addr, partedArgs)
		if err != nil {
			return fmt.Errorf("failed to delete disk Partition on node %s, error:  %v", addr, err)
		}
	}
	return nil
}

// ApplyNodeSelectorToDaemonset add node selector to the daemonset spec and apply
func ApplyNodeSelectorToDaemonset(daemonsetName string, namespace string, label string, value string) error {
	dsApi := gTestEnv.KubeInt.AppsV1().DaemonSets
	daemonset, err := dsApi(namespace).Get(context.TODO(), daemonsetName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get daemonset %s : ns: %s : Error: %v", daemonsetName, namespace, err)
	}
	if daemonset.Spec.Template.Spec.NodeSelector == nil {
		daemonset.Spec.Template.Spec.NodeSelector = make(map[string]string)
	}
	daemonset.Spec.Template.Spec.NodeSelector[label] = value
	_, err = dsApi(namespace).Update(context.TODO(), daemonset, metaV1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to apply node selector to daemonset %s : ns: %s : Error: %v", daemonsetName, namespace, err)
	}
	return nil
}

// RemoveAllNodeSelectorsFromDaemonset removes all node selectors from the daemonset spec and apply
func RemoveAllNodeSelectorsFromDaemonset(daemonsetName string, namespace string) error {
	dsApi := gTestEnv.KubeInt.AppsV1().DaemonSets
	daemonset, err := dsApi(namespace).Get(context.TODO(), daemonsetName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get daemonset %s : ns: %s : Error: %v", daemonsetName, namespace, err)
	}
	if daemonset.Spec.Template.Spec.NodeSelector != nil {
		daemonset.Spec.Template.Spec.NodeSelector = nil
		_, err = dsApi(namespace).Update(context.TODO(), daemonset, metaV1.UpdateOptions{})
	}
	if err != nil {
		return fmt.Errorf("failed to remove node selector from deployment %s : ns: %s : Error: %v", daemonsetName, namespace, err)
	}
	return nil
}

// ScFioVolumeName holds the fio pod, storage class, volume names and uuid
type ScFioVolumeName struct {
	VolName    string
	Uuid       string
	ScName     string
	FioPodName string
	NexusNode  string
}

// CreateScVolumeAndFio creates sc, volume and fio pod and return struct of volume name, volume uuid, sc name,fio pod name
func CreateScVolumeAndFio(prefix string, replicaCount int, volSizeMb int, nodeName string, volType common.VolumeType) (ScFioVolumeName, error) {
	var scFioVolumeName ScFioVolumeName
	scFioVolumeName.ScName = fmt.Sprintf("%s-repl-%d", prefix, replicaCount)
	err := NewScBuilder().
		WithName(scFioVolumeName.ScName).
		WithNamespace(common.NSDefault).
		WithProtocol(common.ShareProtoNvmf).
		WithReplicas(replicaCount).
		BuildAndCreate()
	if err != nil {
		return scFioVolumeName, fmt.Errorf("failed to create storage class %s, error: %v", scFioVolumeName.ScName, err)
	}
	scFioVolumeName.VolName = fmt.Sprintf("vol-%s", scFioVolumeName.ScName)
	uid, err := MkPVC(volSizeMb, scFioVolumeName.VolName, scFioVolumeName.ScName, volType, common.NSDefault)
	if err != nil {
		return scFioVolumeName, fmt.Errorf("failed to create pvc %s, errorr: %v", scFioVolumeName.VolName, err)
	}
	scFioVolumeName.Uuid = uid
	scFioVolumeName.FioPodName = fmt.Sprintf("fio-%s", scFioVolumeName.VolName)
	if nodeName != "" {
		err = CreateSleepingFioPodOnNode(scFioVolumeName.FioPodName, scFioVolumeName.VolName, volType, nodeName)
		if err != nil {
			return scFioVolumeName, fmt.Errorf("failed to create sleeping fio %s, error: %v", scFioVolumeName.FioPodName, err)
		}
	} else {
		err = CreateSleepingFioPod(scFioVolumeName.FioPodName, scFioVolumeName.VolName, volType)
		if err != nil {
			return scFioVolumeName, fmt.Errorf("failed to create sleeping fio %s, error: %v", scFioVolumeName.FioPodName, err)
		}
	}

	// get nexus node
	nexusNode, err := GetNexusNode(uid)
	if err != nil {
		return scFioVolumeName, fmt.Errorf("failed to get nexus node for volume %s, error %v", uid, err)
	}
	scFioVolumeName.NexusNode = nexusNode
	return scFioVolumeName, nil
}

func DeleteScVolumeAndFio(scFioVolumeName ScFioVolumeName) error {
	// Delete the fio pod
	err := DeletePod(scFioVolumeName.FioPodName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to delete fio pod %s, error: %v", scFioVolumeName.FioPodName, err)
	}
	// Delete the volume
	err = RmPVC(scFioVolumeName.VolName, scFioVolumeName.ScName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to delete volume %s, error: %v", scFioVolumeName.VolName, err)
	}
	err = RmStorageClass(scFioVolumeName.ScName)
	if err != nil {
		return fmt.Errorf("failed to delete sc %s, error: %v", scFioVolumeName.ScName, err)
	}
	return nil
}

func readyCheck(namespace string, verbose bool) (bool, error) {
	var err error
	allReady := true
	{
		daemonsets, dserr := gTestEnv.KubeInt.AppsV1().DaemonSets(namespace).List(context.TODO(), metaV1.ListOptions{})
		if dserr == nil {
			for _, ds := range daemonsets.Items {
				ready := ds.Status.DesiredNumberScheduled != 0 &&
					ds.Status.DesiredNumberScheduled == ds.Status.CurrentNumberScheduled &&
					ds.Status.DesiredNumberScheduled == ds.Status.NumberReady &&
					ds.Status.DesiredNumberScheduled == ds.Status.NumberAvailable
				if verbose {
					logf.Log.Info("DaemonSet",
						"ready", ready,
						"name", ds.Name,
						"DesiredNumberScheduled", ds.Status.DesiredNumberScheduled,
						"CurrentNumberScheduled", ds.Status.CurrentNumberScheduled,
						"NumberReady", ds.Status.NumberReady,
						"NumberAvailable", ds.Status.NumberAvailable,
					)
				}
				allReady = allReady && ready
			}
		} else {
			err = dserr
		}
	}

	{
		statefulsets, stserr := gTestEnv.KubeInt.AppsV1().StatefulSets(namespace).List(context.TODO(), metaV1.ListOptions{})
		if stserr == nil {
			for _, sts := range statefulsets.Items {
				ready := sts.Status.Replicas == sts.Status.ReadyReplicas && sts.Status.ReadyReplicas == sts.Status.CurrentReplicas && sts.Status.ReadyReplicas != 0
				if verbose {
					logf.Log.Info("StatefulSet",
						"ready", ready,
						"name", sts.Name,
						"Replicas", sts.Status.Replicas,
						"ReadyReplicas", sts.Status.ReadyReplicas,
						"CurrentReplicas", sts.Status.CurrentReplicas,
					)
				}
				allReady = allReady && ready
			}
		} else {
			err = stserr
		}
	}

	{
		deployments, dperr := gTestEnv.KubeInt.AppsV1().Deployments(namespace).List(context.TODO(), metaV1.ListOptions{})
		if dperr == nil {
			for _, deployment := range deployments.Items {
				ready := false
				for _, condition := range deployment.Status.Conditions {
					if condition.Type == appsV1.DeploymentAvailable {
						if condition.Status == coreV1.ConditionTrue {
							ready = true
						}
					}
				}
				if verbose {
					logf.Log.Info("Deployment",
						"ready", ready,
						"name", deployment.Name,
					)
				}
				allReady = allReady && ready
			}
		} else {
			err = dperr
		}
	}

	if verbose {
		logf.Log.Info("ReadyCheck", "namespace", namespace, "allReady", allReady, "error", err)
	}
	return allReady, err
}

func ReadyCheck(namespace string) (bool, error) {
	logf.Log.Info("ReadyCheck", "namespace", namespace)
	return readyCheck(namespace, true)
}

// list all pods in given namespace and return true if all are in running state else return false
func PodReadyCheck(namespace string) (bool, error) {
	logf.Log.Info("PodReadyCheck", "namespace", namespace)
	return podReadyCheck(namespace)
}

// list all pods in
func podReadyCheck(namespace string) (bool, error) {
	pods, err := ListPod(namespace)
	if err != nil {
		logf.Log.Error(err, "can't get pods", "namespace", namespace)
		return false, err
	}
	var notRunningpods []string
	for _, p := range pods.Items {
		if p.Status.Phase != coreV1.PodRunning {
			logf.Log.Info("pod not in running state", "pod", p.Name, "namespace", p.Namespace)
			notRunningpods = append(notRunningpods, p.Name)
		}
	}
	if len(notRunningpods) == 0 {
		logf.Log.Info("pods not in running state", "pods", notRunningpods, "namespace", namespace)
		return true, nil
	}
	return false, nil
}

// GetDeploymentReplicaCount return deployment replica count
func GetDeploymentReplicaCount(deploymentName string, namespace string) (int32, error) {
	depAPI := gTestEnv.KubeInt.AppsV1().Deployments
	deployment, err := depAPI(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return *deployment.Spec.Replicas, fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
			deploymentName,
			namespace,
			err)
	}
	return *deployment.Spec.Replicas, err
}

func SetPromtailTolerations(tolerations []coreV1.Toleration, promtailDsName string, namespace string) error {
	dsAPI := gTestEnv.KubeInt.AppsV1().DaemonSets
	var err error

	// this is to cater for a race condition, occasionally seen,
	// when the deployment is changed between Get and Update
	for attempts := 0; attempts < 10; attempts++ {
		ds, err := dsAPI(namespace).Get(context.TODO(), promtailDsName, metaV1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get daemonset, name: %s, namespace: %s, error: %v",
				promtailDsName,
				namespace,
				err)
		}
		ds.Spec.Template.Spec.Tolerations = tolerations

		_, err = dsAPI(namespace).Update(context.TODO(), ds, metaV1.UpdateOptions{})
		if err == nil {
			break
		}
		logf.Log.Info("Re-trying update attempt due to error", "error", err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to add node toleration to daemonset, name: %s, namespace: %s, error: %v",
			promtailDsName,
			namespace,
			err)
	}
	return nil
}

// VerifyIoEnginePodDeletionFromNode verify io engine pod removal from node
func VerifyIoEnginePodDeletionFromNode(nodeIP string, podPrefix string, timeoutSecs int) bool {
	const sleepTime = 3
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		if CheckMsPodOnNodeByPrefix(nodeIP, podPrefix) {
			time.Sleep(sleepTime * time.Second)
			continue
		}
		return true
	}
	return false
}

// VerifyIoEnginePodCreationOnNode verify io engine pod creation on node
func VerifyIoEnginePodCreationOnNode(nodeIP string, podPrefix string, timeoutSecs int) bool {
	const sleepTime = 3
	time.Sleep(sleepTime * time.Second)
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		if !CheckMsPodOnNodeByPrefix(nodeIP, podPrefix) {
			time.Sleep(sleepTime * time.Second)
			continue
		}
		return true
	}
	return false
}

// UnscheduleIoEngineFromNode remove io engine label and verify pod removal from node
func UnscheduleIoEngineFromNode(nodeIP string, nodeName string, podPrefix string, timeoutSecs int) error {
	// remove io engine label from node
	engineLabel := e2e_config.GetConfig().Product.EngineLabel

	logf.Log.Info("Unlabel node", "nodeName", nodeName, "EngineLabel", engineLabel)
	err := UnlabelNode(nodeName, engineLabel)
	if err != nil {
		return fmt.Errorf("failed to remove label %s from node %s , error %v", engineLabel, nodeName, err)
	}

	// wait for io engine pod to be removed from node
	isDeleted := VerifyIoEnginePodDeletionFromNode(nodeIP, podPrefix, timeoutSecs)
	if !isDeleted {
		return fmt.Errorf("pod %s present in cluster after deletion", podPrefix)
	}

	return nil
}

// ApplyIoEngineLabelToNode add io engine label and verify pod creation on node
func ApplyIoEngineLabelToNode(nodeIP string, nodeName string, podPrefix string, timeoutSecs int) error {
	// remove io engine label from node
	engineLabel := e2e_config.GetConfig().Product.EngineLabel
	engineLabelValue := e2e_config.GetConfig().Product.EngineLabelValue
	logf.Log.Info("Label node", "nodeName", nodeName, "EngineLabel", engineLabel, "EngineLabelValue", engineLabelValue)
	err := LabelNode(nodeName, engineLabel, engineLabelValue)
	if err != nil {
		return fmt.Errorf("failed to add label %s to node %s , error %v", engineLabel, nodeName, err)
	}

	// wait for io engine pod to be scheduled on node
	isPresent := VerifyIoEnginePodCreationOnNode(nodeIP, podPrefix, timeoutSecs)
	if !isPresent {
		return fmt.Errorf("failed to schedule pod %s  on node %s", podPrefix, nodeName)
	}

	return nil
}

func VerifyMayastorAndPoolReady(timeoutSecs int) (bool, error) {
	// mayastor ready check
	mayastorReady, err := MayastorReady(2, 360)
	if err != nil {
		return mayastorReady, fmt.Errorf("failed to verify mayastor installation , error: %v", err)
	}
	if !mayastorReady {
		return mayastorReady, fmt.Errorf("mayastor installation not ready")
	}
	// disk pools online check
	err = WaitForPoolsToBeOnline(timeoutSecs)
	if err == nil && mayastorReady {
		return mayastorReady, err
	}
	return false, fmt.Errorf("mayastor pods are ready but all mayastor disk pools are not online, error: %v", err)
}

// Add argument to container in deployment
func SetContainerArgs(deploymentName string, containerName string, namespace string, argKey string, argValue string) ([]string, error) {
	var old_args []string
	depAPI := gTestEnv.KubeInt.AppsV1().Deployments
	var err error
	deployment, err := depAPI(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return old_args, fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
			deploymentName,
			namespace,
			err)
	}

	containers := deployment.Spec.Template.Spec.Containers
	for i, container := range containers {
		if container.Name == containerName {
			old_args = containers[i].Args
			containers[i].Args = replaceOrAppend(container.Args, argKey, argValue)
			break
		}
	}

	deployment.Spec.Template.Spec.Containers = containers

	_, err = depAPI(namespace).Update(context.TODO(), deployment, metaV1.UpdateOptions{})

	if err != nil {
		return old_args, fmt.Errorf("failed to set container argument to deployment, name: %s, container %s, namespace: %s, argument key: %s, error: %v",
			deploymentName,
			containerName,
			namespace,
			argKey,
			err)
	}
	return old_args, nil
}

// Set all of the arguments for a container. Used for restoring the original set.
func SetAllContainerArgs(deploymentName string, containerName string, namespace string, args []string) error {
	depAPI := gTestEnv.KubeInt.AppsV1().Deployments
	var err error
	deployment, err := depAPI(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment, name: %s, namespace: %s, error: %v",
			deploymentName,
			namespace,
			err)
	}

	containers := deployment.Spec.Template.Spec.Containers
	found := false
	for i, container := range containers {
		if container.Name == containerName {
			containers[i].Args = args
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("failed to find container, name: %s", containerName)
	}

	deployment.Spec.Template.Spec.Containers = containers

	_, err = depAPI(namespace).Update(context.TODO(), deployment, metaV1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to set container argument to deployment, name: %s, container %s, namespace: %s, error: %v",
			deploymentName,
			containerName,
			namespace,
			err)
	}
	return nil
}

func replaceOrAppend(args []string, argKey string, argValue string) []string {
	newArg := fmt.Sprintf("%s=%s", argKey, argValue)
	for i, str := range args {
		if strings.Contains(str, argKey) {
			args[i] = newArg
			return args
		}
	}
	return append(args, newArg)
}

// CheckIfAllResourcesAreDifferent function can be used to check if all elements passed
// as arguments in the lists are different. So for e.g. it will be useful in tests where we
// want to verify that all pools are on different nodes or all
// nexuses are on different nodes or all replicas have different pools.
func CheckIfAllResourcesAreDifferent(resourceList []string) bool {
	nodes := make(map[string]bool)
	allDifferent := true
	for _, node := range resourceList {
		nodeName := node
		if nodes[nodeName] {
			allDifferent = false
			break
		}
		nodes[nodeName] = true
	}
	return allDifferent
}

func RestrictMayastorToSingleNode() (string, []string, error) {
	var ioEngineNode string
	var ioEngineNodes []string
	msn, err := GetMayastorNodeNames()
	if err != nil {
		return ioEngineNode, ioEngineNodes, fmt.Errorf("failed to get mayastor node names, error: %v", err)
	} else if len(msn) == 0 {
		return ioEngineNode, ioEngineNodes, fmt.Errorf("mayastor nodes not found")
	}
	for _, node := range msn {
		if ioEngineNode == "" {
			ioEngineNode = node
		} else {
			ioEngineNodes = append(ioEngineNodes, node)
			logf.Log.Info("remove io-engine label", "node", node)
			nodeIp, err := GetNodeIPAddress(node)
			if err != nil {
				return ioEngineNode, ioEngineNodes, fmt.Errorf("failed to get node %s IP, error: %v", node, err)
			}
			err = UnscheduleIoEngineFromNode(*nodeIp,
				node,
				e2e_config.GetConfig().Product.IOEnginePodName,
				DefTimeoutSecs,
			)
			if err != nil {
				return ioEngineNode, ioEngineNodes, fmt.Errorf("failed remove io engine label from node %s, error: %v", node, err)
			}

		}
	}
	return ioEngineNode, ioEngineNodes, err
}

// returns the percentage of the disk pool in the selected units (1.00 = 100%)
func GetPoolSizeFraction(pool common.MayastorPool, percentCapacity float64, unit string) int {
	capacityInUnits := GetSizePerUnits(pool.Status.Capacity, unit)
	return int(capacityInUnits * percentCapacity)
}

func GetSizePerUnits(b uint64, unit string) float64 {
	m := map[string]float64{"": 0, "KiB": 1, "MiB": 2, "GiB": 3, "TiB": 4, "PiB": 5}
	bf := float64(b)
	bf /= math.Pow(1024, m[unit])
	return math.Round(bf)
}

func DsSetContainerEnv(dsName string, nameSpace string, containerName string, envName string, envValue string) (bool, error) {
	dsApi := gTestEnv.KubeInt.AppsV1().DaemonSets(nameSpace)
	var ok, updateRequired bool
	var err error
	var ds *appsV1.DaemonSet
	envSpec := coreV1.EnvVar{
		Name:      envName,
		Value:     envValue,
		ValueFrom: nil,
	}
	logf.Log.Info("Daemonset Set Container Env: ", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "env", envSpec)

	for !ok {
		ds, err = dsApi.Get(context.TODO(), dsName, metaV1.GetOptions{})
		if err != nil {
			return updateRequired, err
		}
		updateRequired = false
		foundContainer := false
		for cix, container := range ds.Spec.Template.Spec.Containers {
			if container.Name == containerName {
				foundContainer = true
				foundEnv := false
				envs := container.Env
				for eix, env := range envs {
					if env.Name == envName {
						foundEnv = true
						if env != envSpec {
							logf.Log.Info("Daemonset Set Container Env: Modify ", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "from", env, "to", envSpec)
							envs[eix] = envSpec
							ds.Spec.Template.Spec.Containers[cix].Env = envs
							updateRequired = true
						} else {
							logf.Log.Info("DaemonSet Set Container Env: No Change ", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "env", env)
							ok = true
						}
					}
				}
				if !foundEnv {
					envs = append(envs, envSpec)
					ds.Spec.Template.Spec.Containers[cix].Env = envs
					logf.Log.Info("Dameonset Set Container Env: Add", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "env", envSpec)
					updateRequired = true
				}
			}
		}
		if !foundContainer {
			return updateRequired, fmt.Errorf("container %s is not part of pod %s(%s)", containerName, dsName, nameSpace)
		}
		if updateRequired {
			_, err = dsApi.Update(context.TODO(), ds, metaV1.UpdateOptions{})
			logf.Log.Info("Dameonset Set Container Env: Updated daemonset", "name", dsName, "namespace", nameSpace, "error", err)
			ok = err == nil
		}
	}
	return updateRequired, nil
}

func DsUnsetContainerEnv(dsName string, nameSpace string, containerName string, envName string) (bool, error) {
	dsApi := gTestEnv.KubeInt.AppsV1().DaemonSets(nameSpace)
	var ok, updateRequired bool
	var err error
	var ds *appsV1.DaemonSet
	logf.Log.Info("DaemonSet Unset Container Env", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "envName", envName)
	for !ok {
		updateRequired = false
		ds, err = dsApi.Get(context.TODO(), dsName, metaV1.GetOptions{})
		if err != nil {
			return updateRequired, err
		}
		foundContainer := false
		for cix, container := range ds.Spec.Template.Spec.Containers {
			if container.Name == containerName {
				foundContainer = true
				foundEnv := false
				var envs []coreV1.EnvVar
				for _, env := range container.Env {
					if env.Name == envName {
						foundEnv = true
						updateRequired = true
						logf.Log.Info("DaemonSet Unset Container Env: Remove", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "env", env)
					} else {
						envs = append(envs, env)
					}
				}
				if !foundEnv {
					logf.Log.Info("DaemonSet Unset Container Env: No Change", "DaemonSet", dsName, "namespace", nameSpace, "container", containerName, "envName", envName)
					ok = true
				} else {
					ds.Spec.Template.Spec.Containers[cix].Env = envs
				}
			}
		}
		if !foundContainer {
			return updateRequired, fmt.Errorf("container %s is not part of pod %s(%s)", containerName, dsName, nameSpace)
		}
		if updateRequired {
			_, err = dsApi.Update(context.TODO(), ds, metaV1.UpdateOptions{})
			logf.Log.Info("DaemonSet Unset Container Env: Updated daemonset", "name", dsName, "namespace", nameSpace, "error", err)
			ok = err == nil
		}
	}
	return updateRequired, nil
}

// DSUpdateStrategy set a daemonset update strategy to rolling update on on delete
func DsUpdateStrategy(dsName string, nameSpace string, rollingUpdate bool) error {
	dsApi := gTestEnv.KubeInt.AppsV1().DaemonSets(nameSpace)
	var ok bool
	var err error
	var ds *appsV1.DaemonSet
	logf.Log.Info("DaemonSet Set Update Strategy", "DaemonSet", dsName, "namespace", nameSpace, "rollingUpdate", rollingUpdate)
	for !ok {
		ds, err = dsApi.Get(context.TODO(), dsName, metaV1.GetOptions{})
		if err != nil {
			return err
		}
		if rollingUpdate {
			val0 := intstr.FromInt(0)
			val1 := intstr.FromInt(1)
			rtu := appsV1.RollingUpdateDaemonSet{
				MaxUnavailable: &val1,
				MaxSurge:       &val0,
			}
			ds.Spec.UpdateStrategy = appsV1.DaemonSetUpdateStrategy{
				Type:          "RollingUpdate",
				RollingUpdate: &rtu,
			}
		} else {
			ds.Spec.UpdateStrategy = appsV1.DaemonSetUpdateStrategy{
				Type:          "OnDelete",
				RollingUpdate: nil,
			}
		}
		_, err = dsApi.Update(context.TODO(), ds, metaV1.UpdateOptions{})
		ok = err == nil
	}
	return nil
}

// SetNexusRebuildVerify sets or removes environment variable NEXUS_REBUILD_VERIFY for IO engine daemonset,
// then deletes all the IO engine pods and waits for Mayastor to be ready and pools to be online
func SetNexusRebuildVerify(on bool) error {
	var err error
	var updated bool
	if on {
		updated, err = DsSetContainerEnv("mayastor-io-engine", common.NSMayastor(), "io-engine", "NEXUS_REBUILD_VERIFY", "panic")
	} else {
		updated, err = DsUnsetContainerEnv("mayastor-io-engine", common.NSMayastor(), "io-engine", "NEXUS_REBUILD_VERIFY")
	}
	if err != nil {
		return err
	}
	if updated {
		logf.Log.Info("SetNexusRebuildVerify: mayastor-io-engine daemonset was updated, deleting all IO-engine pods")
		podApi := gTestEnv.KubeInt.CoreV1().Pods
		ioEnginePods, err := ListPodsByPrefix(common.NSMayastor(), e2e_config.GetConfig().Product.IOEnginePodName)
		if err != nil {
			return err
		}
		for _, pod := range ioEnginePods {
			delErr := podApi(common.NSMayastor()).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{})
			if delErr != nil {
				if err != nil {
					err = fmt.Errorf("%v ; %v", err, delErr)
				} else {
					err = delErr
				}
			}
		}
		if err != nil {
			return err
		}
		time.Sleep(20 * time.Second)

		var ready bool
		ready, err = MayastorReady(10, 120)
		if err != nil {
			return err
		}
		if !ready {
			return fmt.Errorf("mayastor is not ready after modifying io-engine pod environment variable")
		}

		err = WaitForPoolsToBeOnline(120)
		if err != nil {
			return err
		}
	}
	return err
}

func WaitForDeploymentReady(deploymentName string, namespace string, sleepTime int, duration int) bool {
	ready := false
	count := (duration + sleepTime - 1) / sleepTime

	logf.Log.Info("DeploymentReadyCheck", "Deployment", deploymentName, "namespace", namespace)
	for ix := 0; ix < count && !ready; ix++ {
		time.Sleep(time.Duration(sleepTime) * time.Second)

		ready = DeploymentReady(deploymentName, namespace)
		logf.Log.Info("DeploymentReady: ", "deployment", deploymentName, "ready", ready)
	}

	if !ready {
		logf.Log.Info("Deployment not ready", "Deployment", deploymentName, "namespace", namespace)
		return false
	}

	// verify pod running
	err := WaitForPodRunning(deploymentName, namespace, duration)
	if err != nil {
		logf.Log.Info("Pod not ready", "Pod prefix", deploymentName, "namespace", namespace)
		return false
	}
	return ready
}

func UpdateDeploymentReplica(deployName string, namespace string, replica int32) error {
	logf.Log.Info("Scale deployment", "deployement", deployName, "namespace", namespace, "replica", replica)
	err := SetDeploymentReplication(deployName, namespace, &replica)
	if err != nil {
		return fmt.Errorf("failed to scale up %s deployment with replica %s,error: %v", deployName, namespace, err)
	}

	// verify mayastor-csi-controller deployment is scaled up
	ready, err := WaitForDeploymentToMatchExpectedState(deployName, namespace, int(replica))
	if err != nil {
		return fmt.Errorf("failed to verify scaling up %s deployment with replica %s,error: %v", deployName, namespace, err)
	} else if !ready {
		return fmt.Errorf("%s deployment not ready", deployName)
	}
	return nil

}

func WaitForDeploymentToMatchExpectedState(deployName string, namespace string, expectedReplica int) (bool, error) {
	const timeSleepSecs = 5

	var err error
	// Wait for the mayastor volume to be resized
	for ix := 0; ix < DefTimeoutSecs/timeSleepSecs; ix++ {
		replica, err := DeploymentReadyCount(deployName, namespace)
		if err != nil {
			return false, fmt.Errorf("failed to get deplyment %s ready count, error: %v", deployName, err)
		}
		if replica == expectedReplica {
			return true, nil
		}
		time.Sleep(timeSleepSecs * time.Second)
	}
	return false, err
}
