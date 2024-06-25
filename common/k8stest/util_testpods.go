package k8stest

// Utility functions for test pods.
import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// FIXME: this function runs fio with a bunch of parameters which are not configurable.
// sizeMb should be 0 for fio to use the entire block device
func RunFio(podName string, duration int, filename string, sizeMb int, args ...string) ([]byte, error) {
	argRuntime := fmt.Sprintf("--runtime=%d", duration)
	argFilename := fmt.Sprintf("--filename=%s", filename)

	logf.Log.Info("RunFio",
		"podName", podName,
		"duration", duration,
		"filename", filename,
		"args", args)

	cmdArgs := []string{
		"exec",
		"-it",
		podName,
		"--",
		"fio",
		"--name=benchtest",
		"--verify=crc32",
		"--verify_fatal=1",
		"--verify_async=2",
		argFilename,
		"--direct=1",
		"--rw=randrw",
		"--ioengine=libaio",
		"--bs=4k",
		"--iodepth=16",
		"--numjobs=1",
		"--time_based",
		argRuntime,
	}

	if sizeMb != 0 {
		sizeArgs := []string{fmt.Sprintf("--size=%dm", sizeMb)}
		cmdArgs = append(cmdArgs, sizeArgs...)
	}

	if args != nil {
		cmdArgs = append(cmdArgs, args...)
	}

	cmd := exec.Command(
		"kubectl",
		cmdArgs...,
	)
	cmd.Dir = ""
	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("Running fio failed", "error", err)
	}
	return output, err
}

// IsPodWithLabelsRunning expects that at any time only one application pod will be in running state
// if there are more then one pod in terminating state then it will return the last terminating pod.
func IsPodWithLabelsRunning(labels, namespace string) (string, bool, error) {
	var podName string
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labels})
	if err != nil {
		return "", false, err
	}
	if len(pods.Items) == 0 {
		return "", false, nil
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodRunning {
			return pod.Name, true, nil
		}
		podName = pod.Name
	}
	return podName, false, nil
}

// ForceDeleteTerminatingPods force deletes the pod this function is required because
// sometimes after powering off the node some pods stuck in terminating state.
func ForceDeleteTerminatingPods(labels, namespace string) error {
	deletionTime := int64(0)
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labels})
	if err != nil {
		return err
	}
	if len(pods.Items) == 0 {
		return nil
	}
	for _, pod := range pods.Items {
		if pod.DeletionTimestamp != nil {
			err = gTestEnv.KubeInt.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{GracePeriodSeconds: &deletionTime})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetNodeListForPods(labels, namespace string) (map[string]v1.PodPhase, error) {
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{LabelSelector: labels})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) == 0 {
		return nil, nil
	}
	nodeList := map[string]v1.PodPhase{}
	for _, pod := range pods.Items {
		nodeList[pod.Spec.NodeName] = pod.Status.Phase
	}
	return nodeList, nil
}

func GetNodeNameForScheduledPod(podName, nameSpace string) (string, error) {
	pod, err := gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("GetNodeNameForPod: %v", err)
	}
	return pod.Spec.NodeName, nil
}

func IsPodRunning(podName string, nameSpace string) bool {
	pod, err := gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return false
	}
	return pod.Status.Phase == v1.PodRunning
}

// WaitPodRunning wait for pod to transition to running with timeout,
// returns true of the pod is running, false otherwise.
func WaitPodRunning(podName string, nameSpace string, timeoutSecs int) bool {
	const sleepTime = 3
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime && !IsPodRunning(podName, nameSpace); ix++ {
		time.Sleep(sleepTime * time.Second)
	}
	return IsPodRunning(podName, nameSpace)
}

func GetPodScheduledStatus(podName string, nameSpace string) (coreV1.ConditionStatus, string, error) {
	pod, err := gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return coreV1.ConditionUnknown, "", fmt.Errorf("failed to get pod")
	}
	status := pod.Status
	for _, condition := range status.Conditions {
		if condition.Type == coreV1.PodScheduled {
			return condition.Status, condition.Reason, nil
		}
	}
	return coreV1.ConditionUnknown, "", fmt.Errorf("failed to find pod scheduled condition")
}

// CreatePod Create a Pod in the specified namespace, no options and no context
func CreatePod(podDef *coreV1.Pod, nameSpace string) (*coreV1.Pod, error) {
	logf.Log.Info("Creating", "pod", podDef.Name)
	return gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Create(context.TODO(), podDef, metaV1.CreateOptions{})
}

// DeletePod Delete a Pod in the specified namespace, no options and no context
func DeletePod(podName string, nameSpace string) error {
	logf.Log.Info("Deleting", "pod", podName)
	return gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Delete(context.TODO(), podName, metaV1.DeleteOptions{})
}

func DeletePodsByLabel(label string, nameSpace string) error {
	logf.Log.Info("Deleting", "pod with label", label)

	podList, err := gTestEnv.KubeInt.CoreV1().Pods(nameSpace).List(context.TODO(), metaV1.ListOptions{LabelSelector: label})
	if err != nil {
		return fmt.Errorf("failed to list pods with label: %s in namespace: %s", label, nameSpace)
	}
	if len(podList.Items) == 0 {
		logf.Log.Info("no pods found with label: %s in namespace: %s", label, nameSpace)
		return nil
	}

	for _, pod := range podList.Items {
		err = gTestEnv.KubeInt.CoreV1().Pods(nameSpace).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete pod: %s, err:%v", pod.Name, err)
		} else {
			logf.Log.Info("Deleted", "pod", pod.Name)
		}
	}

	return nil
}

// CreateFioPodDef  deprecated use MakeFioContainer and NewPodBuilder instead
// / Create a test fio pod in default namespace, no options and no context
// / for filesystem,  mayastor volume is mounted on /volume
// / for raw-block, mayastor volume is mounted on /dev/sdm
func CreateFioPodDef(podName string, volName string, volType common.VolumeType, nameSpace string) *coreV1.Pod {
	volMounts := []coreV1.VolumeMount{
		{
			Name:      "ms-volume",
			MountPath: common.FioFsMountPoint,
		},
	}
	volDevices := []coreV1.VolumeDevice{
		{
			Name:       "ms-volume",
			DevicePath: common.FioBlockFilename,
		},
	}

	podDef := coreV1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      podName,
			Namespace: nameSpace,
			Labels:    map[string]string{"app": "fio"},
		},
		Spec: coreV1.PodSpec{
			RestartPolicy: coreV1.RestartPolicyNever,
			Containers: []coreV1.Container{
				{
					Name:            podName,
					Image:           common.GetFioImage(),
					ImagePullPolicy: coreV1.PullAlways,
					Args:            []string{"sleep", "1000000"},
				},
			},
			Volumes: []coreV1.Volume{
				{
					Name: "ms-volume",
					VolumeSource: coreV1.VolumeSource{
						PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
							ClaimName: volName,
						},
					},
				},
			},
		},
	}
	if e2e_config.GetConfig().Platform.HostNetworkingRequired {
		podDef.Spec.HostNetwork = true
	}
	if volType == common.VolRawBlock {
		podDef.Spec.Containers[0].VolumeDevices = volDevices
	}
	if volType == common.VolFileSystem {
		podDef.Spec.Containers[0].VolumeMounts = volMounts
	}
	return &podDef
}

// CreateFioPod deprecated use MakeFioContainer, NewPodBuilder and CreatePod instead
// / Create a test fio pod in default namespace, no options and no context
// / mayastor volume is mounted on /volume
func CreateFioPod(podName string, volName string, volType common.VolumeType, nameSpace string) (*coreV1.Pod, error) {
	logf.Log.Info("Creating fio pod definition", "name", podName, "volume type", volType)
	podDef := CreateFioPodDef(podName, volName, volType, nameSpace)
	return CreatePod(podDef, common.NSDefault)
}

// CheckForTestPods Check if any test pods exist in the default and e2e related namespaces .
func CheckForTestPods() (bool, error) {
	logf.Log.Info("CheckForTestPods")
	foundPods := false

	nameSpaces, err := gTestEnv.KubeInt.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err == nil {
		for _, ns := range nameSpaces.Items {
			if strings.HasPrefix(ns.Name, common.NSE2EPrefix) || ns.Name == common.NSDefault {
				pods, err := gTestEnv.KubeInt.CoreV1().Pods(ns.Name).List(context.TODO(), metaV1.ListOptions{})
				if err == nil && pods != nil && len(pods.Items) != 0 {
					logf.Log.Info("CheckForTestPods",
						"Pods", pods.Items)
					foundPods = true
				}
			}
		}
	}

	return foundPods, err
}

// isPodHealthCheckCandidate is a filter function for health check on pod,
func isPodHealthCheckCandidate(podName string, namespace string) bool {
	if namespace == common.NSMayastor() {
		return !strings.HasPrefix(podName, "mayastor-etcd")
	}
	return true
}

// CheckTestPodsHealth Check test pods in a namespace for restarts and failed/unknown state
func CheckTestPodsHealth(namespace string) error {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	var errorStrings []string
	podList, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return errors.New("failed to list pods")
	}

	for _, pod := range podList.Items {
		if !isPodHealthCheckCandidate(pod.Name, namespace) {
			continue
		}
		containerStatuses := pod.Status.ContainerStatuses
		for _, containerStatus := range containerStatuses {
			if containerStatus.RestartCount != 0 {
				logf.Log.Info(pod.Name, "restarts", containerStatus.RestartCount)
				errorStrings = append(errorStrings, fmt.Sprintf("%s restarted %d times", pod.Name, containerStatus.RestartCount))
			}
			if pod.Status.Phase == coreV1.PodFailed || pod.Status.Phase == coreV1.PodUnknown {
				logf.Log.Info(pod.Name, "phase", pod.Status.Phase)
				errorStrings = append(errorStrings, fmt.Sprintf("%s phase is %v", pod.Name, pod.Status.Phase))
			}
		}
	}

	if len(errorStrings) != 0 {
		return errors.New(strings.Join(errorStrings[:], "; "))
	}
	return nil
}

func CheckPodCompleted(podName string, nameSpace string) (coreV1.PodPhase, error) {

	// Keeping this commented out code for the time being.
	//	podApi := gTestEnv.KubeInt.CoreV1().Pods
	//	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	//	if err != nil {
	//		return coreV1.PodUnknown, err
	//	}
	//	return pod.Status.Phase, err
	//
	return CheckPodContainerCompleted(podName, nameSpace)
}

var reCompileOnce sync.Once
var reFioLog *regexp.Regexp = nil
var reFioCritical *regexp.Regexp = nil

func ScanFioPodLogs(pod v1.Pod, synopsisIn *common.E2eFioPodLogSynopsis) *common.E2eFioPodLogSynopsis {
	var podLogSynopsis *common.E2eFioPodLogSynopsis
	if synopsisIn != nil {
		podLogSynopsis = synopsisIn
	} else {
		podLogSynopsis = &common.E2eFioPodLogSynopsis{
			CriticalFailure: false,
			Text:            []string{},
		}
	}
	reCompileOnce.Do(func() {
		var reErr error
		reFioLog, reErr = regexp.Compile("(verify: bad)|(error)")
		if reErr != nil {
			// probably not needed but safer than sorry
			reFioLog = nil
			logf.Log.Info("WARNING failed to compile regular expression for fio log search")
		}
		reFioCritical, reErr = regexp.Compile("verify: bad")
		if reErr != nil {
			// probably not needed but safer than sorry
			reFioCritical = nil
			logf.Log.Info("WARNING failed to compile regular expression for fio critical failure search")
		}
	})
	for _, container := range pod.Spec.Containers {
		opts := v1.PodLogOptions{}
		opts.Follow = true
		opts.Container = container.Name
		podLogs, err := gTestEnv.KubeInt.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &opts).Stream(context.TODO())
		if err != nil {
			podLogSynopsis.Err = err
			logf.Log.Info("Failed to stream logs for", "pod", pod, "err", err)
			return podLogSynopsis
		}
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			line := reader.Text()
			if reFioLog != nil && reFioLog.MatchString(line) {
				podLogSynopsis.Text = append(podLogSynopsis.Text, line)
			}
			if reFioCritical != nil && reFioCritical.MatchString(line) {
				podLogSynopsis.CriticalFailure = true
			}
			if strings.HasPrefix(line, "JSON") {
				jsondata := line[4:]
				fTSize := common.FioTargetSizeRecord{}
				fExit := common.FioExitRecord{}
				ju_err := json.Unmarshal([]byte(jsondata), &fTSize)
				if ju_err == nil && fTSize.Size != nil {
					podLogSynopsis.JsonRecords.TargetSizes = append(podLogSynopsis.JsonRecords.TargetSizes, fTSize)
				}
				ju_err = json.Unmarshal([]byte(jsondata), &fExit)
				if ju_err == nil && fExit.ExitValue != nil {
					podLogSynopsis.JsonRecords.ExitValues = append(podLogSynopsis.JsonRecords.ExitValues, fExit)
				}
			}
		}
		_ = podLogs.Close()
	}
	return podLogSynopsis
}

func ScanFioPodLogsByName(podName string, nameSpace string) (*common.E2eFioPodLogSynopsis, error) {
	var podLogSynopsis *common.E2eFioPodLogSynopsis
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return podLogSynopsis, err
	}
	return ScanFioPodLogs(*pod, nil), nil
}

// MonitorE2EFioPod launches a go thread to stream fio pod log output and scan that stream
// to populate fields in E2eFioPodOutputMonitor
func MonitorE2EFioPod(podName string, nameSpace string) (*common.E2eFioPodOutputMonitor, error) {
	var podOut common.E2eFioPodOutputMonitor
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	go func(synopsis *common.E2eFioPodLogSynopsis, pod v1.Pod) {
		ScanFioPodLogs(pod, synopsis)
		podOut.Completed = true
	}(&podOut.Synopsis, *pod)
	return &podOut, podOut.Synopsis.Err
}

func CheckFioPodCompleted(podName string, nameSpace string) (coreV1.PodPhase, *common.E2eFioPodLogSynopsis, error) {
	var podLogSynopsis *common.E2eFioPodLogSynopsis
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return coreV1.PodUnknown, podLogSynopsis, err
	}
	containerStatuses := pod.Status.ContainerStatuses
	for _, containerStatus := range containerStatuses {
		if containerStatus.Name == podName {
			if !containerStatus.Ready {
				if containerStatus.State.Terminated != nil &&
					containerStatus.State.Terminated.Reason == "Completed" {
					podLogSynopsis = ScanFioPodLogs(*pod, nil)
					if containerStatus.State.Terminated.ExitCode == 0 {
						return coreV1.PodSucceeded, podLogSynopsis, podLogSynopsis.Err
					} else {
						return coreV1.PodFailed, podLogSynopsis, podLogSynopsis.Err
					}
				}
			}
		}
	}
	if podLogSynopsis == nil || podLogSynopsis.Err != nil {
		if pod.Status.Phase != coreV1.PodRunning && pod.Status.Phase != coreV1.PodPending {
			podLogSynopsis = ScanFioPodLogs(*pod, nil)
		} else {
			podLogSynopsis = &common.E2eFioPodLogSynopsis{}
		}
	}
	return pod.Status.Phase, podLogSynopsis, podLogSynopsis.Err
}

func GetPodStatus(podName string, nameSpace string) (coreV1.PodPhase, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return coreV1.PodUnknown, err
	}
	containerStatuses := pod.Status.ContainerStatuses
	for _, containerStatus := range containerStatuses {
		if containerStatus.Name == podName {
			if !containerStatus.Ready {
				if containerStatus.State.Terminated != nil &&
					containerStatus.State.Terminated.Reason == "Completed" {
					if containerStatus.State.Terminated.ExitCode == 0 {
						return coreV1.PodSucceeded, err
					} else {
						return coreV1.PodFailed, err
					}
				}
			}
		}
	}
	return pod.Status.Phase, err
}

func DumpPodLog(podName string, nameSpace string) {
	var fLog *os.File
	var logfile string
	logsPath, err := common.GetTestCaseLogsPath()
	if err == nil {
		logsPath += "/runtime"
		_ = os.MkdirAll(logsPath, 0755)
		logfile = fmt.Sprintf("%s/%s.log", logsPath, podName)
		fLog, err = os.Create(logfile)
		if err != nil {
			logf.Log.Info("DumpPodLog: failed to open", "logfile", logfile, "error", err)
			fLog = nil
		} else {
			defer fLog.Close()
		}
	}
	if fLog == nil {
		return
	}
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return
	}
	for _, container := range pod.Spec.Containers {
		opts := v1.PodLogOptions{}
		opts.Follow = true
		opts.Container = container.Name
		podLogs, err := gTestEnv.KubeInt.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &opts).Stream(context.TODO())
		if err != nil {
			continue
		}
		ts := time.Now()
		chnl := make(chan bool)
		go func(chn chan bool) {
			reader := bufio.NewScanner(podLogs)
			for reader.Scan() {
				line := reader.Text()
				chn <- true
				_, _ = fLog.WriteString(fmt.Sprintln(line))
			}
			chn <- false
		}(chnl)
		wait := true
		for wait {
			select {
			case msg := <-chnl:
				wait = msg
				ts = time.Now()
			case <-time.After(10 * time.Second):
				if time.Since(ts) > 10 {
					wait = false
					logf.Log.Info("DumpPodLog: timeout: incomplete", "logfile", logfile, "error", err)
				}
			}
		}
		_ = podLogs.Close()
	}
}

// GetPodHostIp retrieve the IP address  of the node hosting a pod
func GetPodHostIp(podName string, nameSpace string) (string, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	}
	return pod.Status.HostIP, err
}

func CheckPodContainerCompleted(podName string, nameSpace string) (coreV1.PodPhase, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})
	if err != nil {
		return coreV1.PodUnknown, err
	}
	containerStatuses := pod.Status.ContainerStatuses
	for _, containerStatus := range containerStatuses {
		if containerStatus.Name == podName {
			if !containerStatus.Ready {
				if containerStatus.State.Terminated != nil &&
					containerStatus.State.Terminated.Reason == "Completed" {
					if containerStatus.State.Terminated.ExitCode == 0 {
						return coreV1.PodSucceeded, nil
					} else {
						return coreV1.PodFailed, nil
					}
				}
			}
		}
	}
	return pod.Status.Phase, err
}

func GetPodRestartCount(podName string, nameSpace string) (int32, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pod, err := podApi(nameSpace).Get(context.TODO(), podName, metaV1.GetOptions{})

	if err != nil {
		return -1, fmt.Errorf("failed to get pod: %v", err)
	}

	return pod.Status.ContainerStatuses[0].RestartCount, nil

}

var mayastorInitialPodCount int

func SetMayastorInitialPodCount(count int) {
	mayastorInitialPodCount = count
}

func GetMayastorInitialPodCount() int {
	return mayastorInitialPodCount
}

// List mayastor pod names, conditionally
//  1. No timestamp - all mayastor pods
//  2. With timestamp - all mayastor pods created after the timestamp which are Running.
func ListRunningMayastorPods(timestamp *time.Time) ([]string, error) {
	var podNames []string
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return podNames, err
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, "mayastor-etcd") {
			continue
		}

		// skip not running pods
		if pod.Status.Phase != v1.PodRunning {
			logf.Log.Info("not ready pod", "name", pod.Name, "phase", pod.Status.Phase)
			continue
		}
		// If timestamp != nil, then we return the list of pods which are
		// created after the timestamp
		if timestamp != nil {
			cs := pod.GetCreationTimestamp()
			if !cs.After(*timestamp) {
				continue
			}
		}
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}

func ListIOEnginePods() (*v1.PodList, error) {
	ioEngineLabel := e2e_config.GetConfig().Product.PodLabelKey + "=" + e2e_config.GetConfig().Product.IOEnginePodLabelValue
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{LabelSelector: ioEngineLabel})
	if err != nil {
		return nil, errors.New("failed to list pods")
	}

	return pods, nil
}

func ListControlAndDataPlanePods() (*coreV1.PodList, *coreV1.PodList, error) {

	podList, err := ListPod(common.NSMayastor())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list pods in namespace: %s, err: %v", common.NSMayastor(), err)
	}

	controlPlanePodList := &coreV1.PodList{}
	dataPlanePodList := &coreV1.PodList{}

	for _, pod := range podList.Items {
		if pod.Labels[e2e_config.GetConfig().Product.PodLabelKey] == e2e_config.GetConfig().Product.IOEnginePodLabelValue {
			dataPlanePodList.Items = append(dataPlanePodList.Items, pod)
		} else {
			controlPlanePodList.Items = append(controlPlanePodList.Items, pod)
		}
	}
	return controlPlanePodList, dataPlanePodList, nil
}

// RestartMayastorPods shortcut to reinstalling mayastor, especially useful on some
// platforms, for example calling this function after patching the installation to
// use different mayastor images, should allow us to have reasonable confidence that
// mayastor has been restarted with those images.
// Deletes all mayastor pods except for mayastor etcd pods,
// then waits upto specified time for new pods to be provisioned
// Simply deleting the pods and then waiting for daemonset ready checks do not work due to k8s latencies,
// for example it has been observed the mayastor-csi pods are deemed ready
// because they enter terminating state after we've checked for readiness
// Caller must perform readiness checks after calling this function.
func RestartMayastorPods(timeoutSecs int) error {
	var err error
	podApi := gTestEnv.KubeInt.CoreV1().Pods

	podNames, err := ListRunningMayastorPods(nil)
	if err != nil {
		return err
	}

	logf.Log.Info("Restarting", "pods", podNames)
	now := time.Now()
	time.Sleep(1 * time.Second)
	for _, podName := range podNames {
		delErr := podApi(common.NSMayastor()).Delete(context.TODO(), podName, metaV1.DeleteOptions{})
		if delErr != nil {
			logf.Log.Info("Failed to delete", "pod", podName, "error", delErr)
			err = delErr
		} else {
			logf.Log.Info("Deleted", "pod", podName)
		}
	}

	if err != nil {
		return err
	}
	var newPodNames []string
	const sleepTime = 10
	// Wait (with timeout) for all pods to have restarted
	logf.Log.Info("Waiting for all pods to restart", "timeoutSecs", timeoutSecs)
	for ix := 1; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		newPodNames, err = ListRunningMayastorPods(&now)
		if err == nil {
			logf.Log.Info("Restarted", "pods", newPodNames)
			if len(newPodNames) >= GetMayastorInitialPodCount() {
				logf.Log.Info("All pods have been restarted.")
				return nil
			}
		}
	}
	logf.Log.Info("Restart pods failed", "oldpods", podNames, "newpods", newPodNames)
	return fmt.Errorf("restart failed incomplete error=%v", err)
}

/*
func collectNatsPodNames() ([]string, error) {
	var podNames []string
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return podNames, err
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, e2e_config.GetConfig().Product.DataPlaneNats) {
			podNames = append(podNames, pod.Name)
		}
	}
	return podNames, nil
}

// RestartNatsPods restart the nats pods
func RestartNatsPods(timeoutSecs int) error {
	var err error
	podApi := gTestEnv.KubeInt.CoreV1().Pods

	podNames, err := collectNatsPodNames()
	if err != nil {
		return err
	}

	for _, podName := range podNames {
		delErr := podApi(common.NSMayastor()).Delete(context.TODO(), podName, metaV1.DeleteOptions{})
		if delErr != nil {
			logf.Log.Info("Failed to delete", "pod", podName, "error", delErr)
			err = delErr
		} else {
			logf.Log.Info("Deleted", "pod", podName)
		}
	}

	if err != nil {
		return err
	}
	const sleepTime = 5
	// Wait (with timeout) for all pods to have restarted
	// For this to work we rely on the fact that for daemonsets and deployments,
	// when a pod is deleted, k8s spins up a new pod with a different name.
	// So the check is comparison between
	//      1) the list of nats pods deleted
	//      2) a freshly generated list of nats pods
	// - the size of the fresh list >= size of the deleted list
	// - the names of the pods deleted do not occur in the fresh list
	for ix := 1; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		newPodNames, err := collectNatsPodNames()
		if err == nil {
			if len(podNames) <= len(newPodNames) {
				return nil
			}
			time.Sleep(sleepTime * time.Second)
		}
	}
	return fmt.Errorf("restart failed in some nebulous way! ")
}
*/

func restartMayastor(restartTOSecs int, readyTOSecs int, poolsTOSecs int) error {
	var err error
	ready := false

	// note: CleanUp removes replicas
	CleanUp()
	//	_, _ = DeleteAllPoolFinalizers()
	_ = DeleteAllPools()

	err = RestartMayastorPods(restartTOSecs)
	if err != nil {
		return fmt.Errorf("RestartMayastorPods failed %v", err)
	}

	ready, err = MayastorReady(5, readyTOSecs)
	if err != nil {
		return fmt.Errorf("failure waiting for mayastor to be ready %v", err)
	}
	if !ready {
		return fmt.Errorf("mayastor is not ready after deleting all pods")
	}

	// Pause to allow things to settle.
	time.Sleep(30 * time.Second)

	//	_, _ = DeleteAllPoolFinalizers()
	//_ = DeleteAllPools()
	err = CreateConfiguredPools()
	if err != nil {
		return err
	}

	const sleepTime = 10
	for ix := 0; ix < (poolsTOSecs+sleepTime-1)/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		err = custom_resources.CheckAllMsPoolsAreOnline()
		if err == nil {
			break
		}
	}

	err = custom_resources.CheckAllMsPoolsAreOnline()
	if err != nil {
		return fmt.Errorf("not all pools are online %v", err)
	}

	if mayastorclient.CanConnect() {
		err := RmReplicasInCluster()
		if err != nil {
			return fmt.Errorf("RmReplicasInCluster failed %v", err)
		}
	} else {
		logf.Log.Info("WARNING: gRPC calls to mayastor are not enabled, unable to clear orphan replicas")
	}
	return err
}

// RestartMayastor this function "restarts" mayastor by
//   - cleaning up all mayastor resource artefacts,
//   - deleting all mayastor pods
func RestartMayastor(restartTOSecs int, readyTOSecs int, poolsTOSecs int) error {
	var err error
	const restartRetries = 3
	// try to restart upto N times
	// chiefly this is a fudge to get restart to work on some platforms
	for ix := 1; ix <= restartRetries; ix++ {

		logf.Log.Info("Restarting mayastor", "try", fmt.Sprintf("%d/%d", ix, restartRetries))

		if err = restartMayastor(restartTOSecs, restartTOSecs, poolsTOSecs); err != nil {
			logf.Log.Info("Restarting failed, failed to restart pods", "error", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if err = CheckTestPodsHealth(common.NSMayastor()); err != nil {
			logf.Log.Info("Restarting failed, pods are not healthy", "error", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if err = ResourceCheck(true); err != nil {
			logf.Log.Info("Restarting failed, resource check failed", "error", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// @here all restarts and checks passed
		logf.Log.Info("Restarted mayastor successfully")
		break
	}

	return err
}

func GetCoreAgentNodeName() (string, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return "", err
	}
	controlPlaneAgent := e2e_config.GetConfig().Product.ControlPlaneAgent
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, controlPlaneAgent) && pod.Status.Phase == "Running" {
			return pod.Spec.NodeName, nil
		}
	}
	return "", nil
}

// MakeFioContainer returns a container object setup to use e2e-fio and run fio with appropriate permissions.
// Privileged: True, AllowPrivilegeEscalation: True, RunAsUser root,
// parameters:
//
//	name - name of the container (usually the pod name)
//	args - container arguments, if empty the defaults to "sleep", "1000000"
func MakeFioContainer(name string, args []string) coreV1.Container {
	containerArgs := args
	if len(containerArgs) == 0 {
		containerArgs = []string{"sleep", "1000000"}
	}
	var z64 int64 = 0
	var vTrue bool = true

	sc := coreV1.SecurityContext{
		Privileged:               &vTrue,
		RunAsUser:                &z64,
		AllowPrivilegeEscalation: &vTrue,
	}
	return coreV1.Container{
		Name:            name,
		Image:           common.GetFioImage(),
		ImagePullPolicy: coreV1.PullPolicy(e2e_config.GetConfig().ImagePullPolicy),
		Args:            containerArgs,
		SecurityContext: &sc,
	}
}

func fioPodBuilder(podName string, volName string, volType common.VolumeType, args []string) *PodBuilder {
	// fio pod container
	podContainer := MakeFioContainer(podName, args)

	// volume claim details
	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: volName,
			},
		},
	}

	return NewPodBuilder("fio").
		WithName(podName).
		WithNamespace(common.NSDefault).
		WithRestartPolicy(coreV1.RestartPolicyNever).
		WithContainer(podContainer).
		WithVolume(volume).
		WithVolumeDeviceOrMount(volType)
}

func CreateFioPodOnNode(podName string, volName string, nodeName string, args []string) error {
	podObj, err := fioPodBuilder(podName, volName, common.VolFileSystem, args).
		WithNodeName(nodeName).
		Build()

	if err != nil || podObj == nil {
		return err
	}
	// Create fio pod
	_, err = CreatePod(podObj, common.NSDefault)

	return err
}

func CreateFioPodWithNodeSelector(podName string, volName string, volType common.VolumeType, nodeName string, args []string) error {
	podObj, err := fioPodBuilder(podName, volName, volType, args).
		WithNodeSelectorHostnameNew(nodeName).
		Build()

	if err != nil || podObj == nil {
		return err
	}
	// Create fio pod
	_, err = CreatePod(podObj, common.NSDefault)

	return err
}

// DeleteMayastorPodOnNode deletes mayastor pods on a node with names matching the prefix
func DeleteMayastorPodOnNode(nodeIP string, prefix string) error {
	logf.Log.Info("DeleteMayastorPodOnNode", "nodeIP", nodeIP, "prefix", prefix)
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if nodeIP == pod.Status.HostIP && strings.HasPrefix(pod.Name, prefix) {
			logf.Log.Info("Deleting", "pod", pod.Name, "HostIP", pod.Status.HostIP)
			delErr := podApi(common.NSMayastor()).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{})
			if delErr != nil {
				logf.Log.Info("Failed to delete", "pod", pod.Name, "error", delErr)
				err = delErr
			} else {
				logf.Log.Info("Deleted", "pod", pod.Name)
			}
		}
	}
	return err
}

// DeleteFioFile deletes the file which was created by fio while writing data
func DeleteFioFile(podName string, filename string) ([]byte, error) {

	logf.Log.Info("Delete file",
		"podName", podName,
		"filename", filename)

	cmdArgs := []string{
		"exec",
		"-it",
		podName,
		"--",
		"rm",
		"-f",
		filename,
	}

	cmd := exec.Command(
		"kubectl",
		cmdArgs...,
	)
	cmd.Dir = ""
	output, err := cmd.CombinedOutput()
	if err != nil {
		logf.Log.Info("Deleting file failed", "error", err)
	}
	return output, err
}

// DeleteRestartedPods utility function to "clear" pods with restart counts
// use with care - this functions exists because the test framework systematically
// asserts that pods have not restarted to mark a test as passed.
// However there are "legitimate" reasons why mayastor pods may have restarted
// for example when a node has been rebooted
// this function should be used in such cases in conjunction with the MayastorReady
// function to ensure that the test cluster is back to "pristine" state.
func DeleteRestartedPods(namespace string) error {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	var errorStrings []string
	podList, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return errors.New("failed to list pods")
	}

	for _, pod := range podList.Items {
		if !isPodHealthCheckCandidate(pod.Name, namespace) {
			continue
		}
		containerStatuses := pod.Status.ContainerStatuses
		for _, containerStatus := range containerStatuses {
			if containerStatus.RestartCount != 0 {
				logf.Log.Info("DeleteRestartedPods: deleting", "pod", pod.Name)
				err := podApi(namespace).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{})
				if err != nil {
					errorStrings = append(errorStrings, fmt.Sprintf("%v", err))
				}
				break
			}
		}
	}

	if len(errorStrings) != 0 {
		return errors.New(strings.Join(errorStrings[:], "; "))
	}
	return nil
}

func restartMayastorPodsOnNode(timeoutSecs int, nodeName string) error {
	podApi := gTestEnv.KubeInt.CoreV1().Pods

	podNames, err := ListRunningMayastorPodsOnNode(nodeName)

	if err == nil {
		logf.Log.Info("Restarting", "pods", podNames)
		time.Sleep(1 * time.Second)
		for _, podName := range podNames {
			delErr := podApi(common.NSMayastor()).Delete(context.TODO(), podName, metaV1.DeleteOptions{})
			if delErr != nil {
				logf.Log.Info("Failed to delete", "pod", podName, "error", delErr)
				err = delErr
			} else {
				logf.Log.Info("Deleted", "pod", podName)
			}
		}
		// Allow time for the pods deletion to be effective and noticed
		time.Sleep(30 * time.Second)
	}
	return err
}

// List mayastor pod names scheduled on a given node
func ListRunningMayastorPodsOnNode(nodeName string) ([]string, error) {
	var podNames []string
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return podNames, fmt.Errorf("failed to list pod in %s namespace, error: %v", common.NSMayastor(), err)
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, "mayastor-etcd") {
			continue
		}
		if pod.Status.Phase != v1.PodRunning {
			logf.Log.Info("not ready pod", "name", pod.Name, "phase", pod.Status.Phase)
			continue
		}
		if pod.Spec.NodeName == nodeName {
			podNames = append(podNames, pod.Name)
		}
	}
	return podNames, nil
}

// RestartMayastorPodsOnNode restart all mayastor pods scheduled on a given node
func RestartMayastorPodsOnNode(restartTOSecs int, readyTOSecs int, poolsTOSecs int, nodeName string) error {
	ready := false

	err := restartMayastorPodsOnNode(restartTOSecs, nodeName)
	if err != nil {
		logf.Log.Info("Warning: RestartMayastorPodsOnNode failed", "error", err)
	}

	ready, err = MayastorReady(10, readyTOSecs)
	if err != nil {
		return fmt.Errorf("failure waiting for mayastor to be ready %v", err)
	}
	if !ready {
		return fmt.Errorf("mayastor is not ready after deleting all pods")
	}

	// Pause to allow things to settle.
	time.Sleep(30 * time.Second)

	const sleepTime = 10
	for ix := 0; ix < (poolsTOSecs+sleepTime-1)/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		err = custom_resources.CheckAllMsPoolsAreOnline()
		if err == nil {
			break
		}
	}

	return err
}

// GetMayastorPodNameonNodeByPrefix
func GetMayastorPodNameonNodeByPrefix(prefix string, nodeName string) (string, error) {
	//podNames, err := ListRunningMayastorPodsOnNode(nodeName)
	//Expect(err).To(BeNil())
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pod in %s namespace, error: %v", common.NSMayastor(), err)
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			if pod.Spec.NodeName == nodeName {
				return pod.Name, nil
			}
		}
	}
	return "", fmt.Errorf("failed to get mayastor pod with prefix %s on node %s", prefix, nodeName)
}
