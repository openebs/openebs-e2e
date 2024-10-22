package k8stest

import (
	// container "github.com/openebs/maya/pkg/kubernetes/container/v1alpha1"
	// volume "github.com/openebs/maya/pkg/kubernetes/volume/v1alpha1"

	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	errors "github.com/pkg/errors"
	coreV1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// k8sNodeLabelKeyHostname is the label key used by Kubernetes
	// to store the hostname on the node resource.
	K8sNodeLabelKeyHostname = "kubernetes.io/hostname"
	// timeout and sleep time in seconds
	timeout       = 300 // timeout in seconds
	timeSleepSecs = 10  // sleep time in seconds
)

type Pod struct {
	object *corev1.Pod
}

// PodBuilder is the builder object for Pod
type PodBuilder struct {
	pod            *Pod
	errs           []error
	readOnlyVolume bool
}

// NewPodBuilder returns new instance of Builder
func NewPodBuilder(appLabelValue string) *PodBuilder {
	return &PodBuilder{pod: &Pod{object: &corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Labels: map[string]string{"app": appLabelValue},
		},
	}}}
}

// registerErrorMessage helper function to register errors,
// the Build method will fail if any error message have been registered.
func (b *PodBuilder) registerErrorMessage(errMsg string) {
	b.errs = append(b.errs, errors.New(errMsg))
}

// WithTolerationsForTaints sets the Spec.Tolerations with provided taints.
func (b *PodBuilder) WithTolerationsForTaints(taints ...corev1.Taint) *PodBuilder {

	tolerations := []corev1.Toleration{}
	for i := range taints {
		var toleration corev1.Toleration
		toleration.Key = taints[i].Key
		toleration.Effect = taints[i].Effect
		if len(taints[i].Value) == 0 {
			toleration.Operator = corev1.TolerationOpExists
		} else {
			toleration.Value = taints[i].Value
			toleration.Operator = corev1.TolerationOpEqual
		}
		tolerations = append(tolerations, toleration)
	}

	b.pod.object.Spec.Tolerations = append(
		b.pod.object.Spec.Tolerations,
		tolerations...,
	)
	return b
}

// WithName sets the Name field of Pod with provided value.
func (b *PodBuilder) WithName(name string) *PodBuilder {
	if len(name) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing Pod name")
		return b
	}
	b.pod.object.Name = name
	return b
}

// WithNamespace sets the Namespace field of Pod with provided value.
func (b *PodBuilder) WithNamespace(namespace string) *PodBuilder {
	if len(namespace) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing namespace")
		return b
	}
	b.pod.object.Namespace = namespace
	return b
}

// WithNamespace sets the Namespace field of Pod with provided value.
func (b *PodBuilder) WithLabels(labels map[string]string) *PodBuilder {
	if len(labels) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing labels")
		return b
	}
	if b.pod.object.Labels != nil {
		for k, v := range labels {
			b.pod.object.Labels[k] = v
		}
	} else {
		b.pod.object.Labels = labels
	}
	return b
}

// WithRestartPolicy sets the RestartPolicy field in Pod with provided arguments
func (b *PodBuilder) WithRestartPolicy(
	restartPolicy corev1.RestartPolicy,
) *PodBuilder {
	b.pod.object.Spec.RestartPolicy = restartPolicy
	return b
}

// WithNodeName sets the NodeName field of Pod with provided value.
func (b *PodBuilder) WithNodeName(nodeName string) *PodBuilder {
	if len(nodeName) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing Pod node name")
		return b
	}
	b.pod.object.Spec.NodeName = nodeName
	return b
}

// WithNodeSelectorHostnameNew sets the Pod NodeSelector to the provided hostname value
// This function replaces (resets) the NodeSelector to use only hostname selector
func (b *PodBuilder) WithNodeSelectorHostnameNew(hostname string) *PodBuilder {
	if len(hostname) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing Pod hostname")
		return b
	}

	b.pod.object.Spec.NodeSelector = map[string]string{
		K8sNodeLabelKeyHostname: hostname,
	}

	return b
}

// WithContainers sets the Containers field in Pod with provided arguments
func (b *PodBuilder) WithContainers(containers []corev1.Container) *PodBuilder {
	if len(containers) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing containers")
		return b
	}
	b.pod.object.Spec.Containers = containers
	return b
}

// WithContainer sets the Containers field in Pod with provided arguments
func (b *PodBuilder) WithContainer(container corev1.Container) *PodBuilder {
	return b.WithContainers([]corev1.Container{container})
}

// WithVolumes sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolumes(volumes []corev1.Volume) *PodBuilder {
	if len(volumes) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing volumes")
		return b
	}
	if b.pod.object.Spec.Volumes == nil {
		b.pod.object.Spec.Volumes = volumes
	} else {
		b.pod.object.Spec.Volumes = append(b.pod.object.Spec.Volumes, volumes...)
	}
	return b
}

// WithVolume sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolume(volume corev1.Volume) *PodBuilder {
	return b.WithVolumes([]corev1.Volume{volume})
}

// WithMountReadOnly - mount filesystem volumes as read only
// fails at Build  if volume is raw block
func (b *PodBuilder) WithMountReadOnly(rdOnly bool) *PodBuilder {
	b.readOnlyVolume = rdOnly
	return b
}

// WithServiceAccountName sets the ServiceAccountName of Pod spec with
// the provided value
func (b *PodBuilder) WithServiceAccountName(serviceAccountName string) *PodBuilder {
	if len(serviceAccountName) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing Pod service account name")
		return b
	}
	b.pod.object.Spec.ServiceAccountName = serviceAccountName
	return b
}

// WithVolumeMounts sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolumeMounts(volMounts []corev1.VolumeMount) *PodBuilder {
	if len(volMounts) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing VolumeMount")
		return b
	}
	if b.pod.object.Spec.Containers[0].VolumeMounts == nil {
		b.pod.object.Spec.Containers[0].VolumeMounts = volMounts
	} else {
		b.pod.object.Spec.Containers[0].VolumeMounts = append(b.pod.object.Spec.Containers[0].VolumeMounts, volMounts...)
	}
	return b
}

// WithVolumeMount sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolumeMount(volMount corev1.VolumeMount) *PodBuilder {
	return b.WithVolumeMounts([]corev1.VolumeMount{volMount})
}

// WithVolumeDevices sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolumeDevices(volDevices []corev1.VolumeDevice) *PodBuilder {
	if len(volDevices) == 0 {
		b.registerErrorMessage("failed to build Pod object: missing VolumeDevices")
		return b
	}
	b.pod.object.Spec.Containers[0].VolumeDevices = volDevices
	return b
}

// WithVolumeDevice sets the Volumes field in Pod with provided arguments
func (b *PodBuilder) WithVolumeDevice(volDevice corev1.VolumeDevice) *PodBuilder {
	return b.WithVolumeDevices([]corev1.VolumeDevice{volDevice})
}

func (b *PodBuilder) WithVolumeDeviceOrMount(volType common.VolumeType) *PodBuilder {
	volMounts := coreV1.VolumeMount{
		Name:      "ms-volume",
		MountPath: common.FioFsMountPoint,
	}
	volDevices := coreV1.VolumeDevice{
		Name:       "ms-volume",
		DevicePath: common.FioBlockFilename,
	}
	if volType == common.VolRawBlock {
		b.WithVolumeDevice(volDevices)
	} else {
		b.WithVolumeMount(volMounts)
	}

	return b
}

func (b *PodBuilder) WithHostPath(name string, hostPath string) *PodBuilder {
	vHostPathDirectory := coreV1.HostPathDirectory
	b.WithVolume(coreV1.Volume{
		Name: name,
		VolumeSource: coreV1.VolumeSource{
			HostPath: &coreV1.HostPathVolumeSource{
				Path: hostPath,
				Type: &vHostPathDirectory,
			},
		},
	})

	b.WithVolumeMount(corev1.VolumeMount{
		Name:      name,
		MountPath: fmt.Sprintf("/mnt/host/%s", hostPath),
	})

	return b
}

// Build returns the Pod API instance
func (b *PodBuilder) Build() (*corev1.Pod, error) {
	if e2e_config.GetConfig().Platform.HostNetworkingRequired {
		b.pod.object.Spec.HostNetwork = true
	}

	// readonly volume can only be enforced for filesystem volumes
	if b.readOnlyVolume {
		if b.pod.object.Spec.Containers[0].VolumeDevices != nil || len(b.pod.object.Spec.Containers[0].VolumeDevices) != 0 {
			b.registerErrorMessage("read only volume is incompatible with raw block volumes")
		}
		// for now turn on readonly for all filesystem volumes
		for ix := range b.pod.object.Spec.Containers[0].VolumeMounts {
			b.pod.object.Spec.Containers[0].VolumeMounts[ix].ReadOnly = true
		}
	}

	if len(b.errs) > 0 {
		return nil, errors.Errorf("%+v", b.errs)
	}
	return b.pod.object, nil
}

// GetPod return requested pod by name in the given namespace
func GetPod(name, ns string) (*v1.Pod, error) {
	pod, err := gTestEnv.KubeInt.CoreV1().Pods(ns).Get(context.TODO(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, errors.New("failed to get pod")
	}
	return pod, nil
}

// ListPod return lis of pods in the given namespace
func ListPod(ns string) (*v1.PodList, error) {
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(ns).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list pods")
	}
	return pods, nil
}

// ListPodsWithLabel return list of pods with a given label in the given namespace
func ListPodsWithLabel(namespace string, labels map[string]string) (*v1.PodList, error) {
	pods, err := gTestEnv.KubeInt.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{
		LabelSelector: metaV1.FormatLabelSelector(&metaV1.LabelSelector{MatchLabels: labels}),
	})
	if err != nil {
		return nil, err
	}
	return pods, err
}

func VerifyPodsOnNode(podLabelsList []string, nodeName string, namespace string) error {
	for _, label := range podLabelsList {
		var err error
		var nodeList map[string]v1.PodPhase
		ok := false
		for ix := 0; ix < timeout/timeSleepSecs; ix++ {
			nodeList, err = GetNodeListForPods("app="+label, namespace)
			logf.Log.Info("VerifyPodsOnNode", "podLabel", label, "NodeList", nodeList, "error", err)
			if err == nil && len(nodeList) == 1 && nodeList[nodeName] == v1.PodRunning {
				ok = true
				break
			}
			time.Sleep(timeSleepSecs * time.Second)
		}
		if err != nil || !ok {
			return fmt.Errorf("failed to verify pod on node %s, podLabel: %s, error: %v", nodeName, label, err)
		}
	}
	return nil
}

// VerifyPodStatusWithAppLabel return true if all the pod with a given label are in running state
func VerifyPodStatusWithAppLabel(podLabel string, namespace string) (bool, error) {
	ok := false
	var err error
	// var nodeList map[string]v1.PodPhase
	podList, err := GetNodeListForPods("app="+podLabel, namespace)
	logf.Log.Info("VerifyPodsOnNode", "podLabel", podLabel, "error", err)
	if err != nil {
		return ok, fmt.Errorf("failed to list pod with label, podLabel: %s, error: %v", podLabel, err)
	}
	podRunningCount := 0
	for _, podPhase := range podList {
		if podPhase == v1.PodRunning {
			podRunningCount++
		}
	}
	if len(podList) == podRunningCount {
		ok = true
	}
	return ok, nil
}

// RestartPodByPrefix restart the pod by prefix name
func RestartPodByPrefix(prefix string) error {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) && pod.Status.Phase == v1.PodRunning {
			delErr := podApi(common.NSMayastor()).Delete(context.TODO(), pod.Name, metaV1.DeleteOptions{})
			if delErr != nil {
				logf.Log.Info("Failed to delete", "pod", pod.Name, "error", delErr)
				return delErr
			}
			logf.Log.Info("Restarted the ", "pod ", pod.Name)
		}
	}
	return nil
}

// CheckPodIsRunningByPrefix check pod is running by prefix name
func CheckPodIsRunningByPrefix(prefix string) bool {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return false
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) && pod.Status.Phase == v1.PodRunning {
			logf.Log.Info("pod is running ", "pod -> ", pod.Name)
			return true
		}
	}
	return false
}

// CheckMsPodOnNodeByPrefix check mayastor pod is running on a node with names matching the prefix
func CheckMsPodOnNodeByPrefix(nodeIP string, prefix string) bool {
	logf.Log.Info("CheckMsPodOnNodeByPrefix", "nodeIP", nodeIP, "prefix", prefix)
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(common.NSMayastor()).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return false
	}
	for _, pod := range pods.Items {
		if nodeIP == pod.Status.HostIP && strings.HasPrefix(pod.Name, prefix) {
			logf.Log.Info("Found the", "pod", pod.Name, "HostIP", pod.Status.HostIP)
			return pod.Status.Phase == v1.PodRunning
		}
	}
	return false
}

// CreateSleepingFioPod create fio pod definition and create pod  with sleep, it will not run fio
func CreateSleepingFioPod(fioPodName string, volName string, volType common.VolumeType) error {
	// fio pod labels
	label := map[string]string{
		"app": "fio",
	}
	// fio pod container
	firstPodContainer := coreV1.Container{
		Name:            fioPodName,
		Image:           common.GetFioImage(),
		ImagePullPolicy: coreV1.PullAlways,
		Args:            []string{"sleep", "1000000"},
	}

	// volume claim details
	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: volName,
			},
		},
	}

	podObj, err := NewPodBuilder("fio").
		WithName(fioPodName).
		WithNamespace(common.NSDefault).
		WithRestartPolicy(coreV1.RestartPolicyNever).
		WithContainer(firstPodContainer).
		WithVolume(volume).
		WithVolumeDeviceOrMount(volType).
		WithLabels(label).Build()

	if err != nil {
		return fmt.Errorf("failed to generate fio pod definition %s, error %v", fioPodName, err)
	} else if podObj == nil {
		return fmt.Errorf("failed to generate fio pod definition successfully %s", fioPodName)
	}
	// Create  fio pod
	_, err = CreatePod(podObj, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to create fio pod %s, error %v", fioPodName, err)
	}

	// Wait for the fio Pod to transition to running
	timeoutSecs := 90   // in seconds
	const sleepTime = 5 // in seconds
	for ix := 1; ix < timeoutSecs/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		if IsPodRunning(fioPodName, common.NSDefault) {
			break
		}
	}
	return nil
}

// CreateSleepingFioPodOnNode create fio pod definition and create pod on given node if node provided
// otherwise fio will be created on any suitable node with sleep, it will not run fio
func CreateSleepingFioPodOnNode(fioPodName string, volName string, volType common.VolumeType, nodename string) error {

	err := CreateSleepingFioPod(fioPodName, volName, volType)
	if err != nil {
		return err
	}

	pod, err := gTestEnv.KubeInt.CoreV1().Pods(common.NSDefault).Get(context.TODO(), fioPodName, metaV1.GetOptions{})
	if err != nil {
		return err
	}
	// add node selector to fio pod
	pod.Spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname": nodename,
	}

	// update fio pod
	_, err = gTestEnv.KubeInt.CoreV1().Pods(common.NSDefault).Update(context.TODO(), pod, metaV1.UpdateOptions{})
	if err == nil {
		return fmt.Errorf("failed to update fio pod %v with node selector with node %s, err %v", pod, nodename, err)
	}

	// Wait for the fio Pod to transition to running
	timeoutSecs := 90   // in seconds
	const sleepTime = 5 // in seconds
	for ix := 1; ix < timeoutSecs/sleepTime; ix++ {
		time.Sleep(sleepTime * time.Second)
		if IsPodRunning(fioPodName, common.NSDefault) {
			break
		}
	}
	return nil
}

// GetNodeForPodByPrefix return node where pod is running by pod prefix
func GetNodeForPodByPrefix(prefix string, namespace string) (string, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pod in %s namespace, error %v", common.NSMayastor(), err)
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			return pod.Spec.NodeName, nil
		}
	}
	return "", err
}

// GetPodStatusByPrefix return pod phase by pod prefix
func GetPodStatusByPrefix(prefix string, namespace string) (coreV1.PodPhase, error) {
	podApi := gTestEnv.KubeInt.CoreV1().Pods
	pods, err := podApi(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pod in %s namespace, error %v", namespace, err)
	}
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			return pod.Status.Phase, nil
		}
	}
	return "", err
}

func GetPodEvents(podName string, namespace string) (*coreV1.EventList, error) {
	options := metaV1.ListOptions{
		TypeMeta:      metaV1.TypeMeta{Kind: "Pod"},
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", podName),
	}
	return GetEvents(namespace, options)
}

// ListPodsByPrefix return list of pods in the given namespace with names that start with a prefix
// for example pods deployed on behalf of a daemonset
func ListPodsByPrefix(ns string, prefix string) ([]coreV1.Pod, error) {
	var pods []coreV1.Pod
	podList, err := gTestEnv.KubeInt.CoreV1().Pods(ns).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, errors.New("failed to list pods")
	}
	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

// ExecuteCommandInPod FIXME: Explore how to make the stream interactive and more flexible, etc.
func ExecuteCommandInPod(namespace, podName string, cmd string) (string, string, error) {
	logf.Log.Info("Executing command", "namespace", namespace, "podName", podName, "cmd", cmd)
	command := []string{"sh", "-c", cmd}
	// Setting up the request
	request := gTestEnv.KubeInt.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	// Command options
	option := &v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}
	request.VersionedParams(option, scheme.ParameterCodec)

	// Creating executor
	exec, err := remotecommand.NewSPDYExecutor(config.GetConfigOrDie(), "POST", request.URL())
	if err != nil {
		logf.Log.Error(err, "Failed to create SPDYExecutor")
		return "", "", err
	}

	// Setting up the output buffer
	var stdout, stderr bytes.Buffer
	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// Running the command
	err = exec.StreamWithContext(context.Background(), streamOptions)
	if err != nil {
		logf.Log.Error(err, "Error in streaming command")
		logf.Log.Info(strings.TrimRight(stderr.String(), "\n"), "source", "stderr")
		logf.Log.Info(strings.TrimRight(stdout.String(), "\n"), "source", "stdout")
		return "", "", err
	}

	logf.Log.Info("Command execution completed")
	logf.Log.Info(strings.TrimRight(stdout.String(), "\n"), "source", "stdout")
	logf.Log.Info(strings.TrimRight(stderr.String(), "\n"), "source", "stderr")
	return stdout.String(), stderr.String(), nil
}

// ExecuteCommandInContainer is needed for pods with more than one container
func ExecuteCommandInContainer(namespace string, podName string, containerName string, cmd string) (string, string, error) {
	logf.Log.Info("Executing command", "namespace", namespace, "podName", podName, "cmd", cmd)
	command := []string{"sh", "-c", cmd}
	// Setting up the request
	request := gTestEnv.KubeInt.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName)

	// Command options
	option := &v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}
	request.VersionedParams(option, scheme.ParameterCodec)

	// Creating executor
	exec, err := remotecommand.NewSPDYExecutor(config.GetConfigOrDie(), "POST", request.URL())
	if err != nil {
		logf.Log.Error(err, "Failed to create SPDYExecutor")
		return "", "", err
	}

	// Setting up the output buffer
	var stdout, stderr bytes.Buffer
	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// Running the command
	err = exec.StreamWithContext(context.Background(), streamOptions)
	if err != nil {
		logf.Log.Error(err, "Error in streaming command")
		logf.Log.Info(strings.TrimRight(stderr.String(), "\n"), "source", "stderr")
		logf.Log.Info(strings.TrimRight(stdout.String(), "\n"), "source", "stdout")
		return "", "", err
	}

	logf.Log.Info("Command execution completed")
	logf.Log.Info(strings.TrimRight(stdout.String(), "\n"), "source", "stdout")
	logf.Log.Info(strings.TrimRight(stderr.String(), "\n"), "source", "stderr")
	return stdout.String(), stderr.String(), nil
}

// XfsCheck returns output and exit status code of xfs_repair -n
func XfsCheck(nodeName string, deployName string, containerName string) (string, error) {
	// Get DaemonSet Name of csi-node
	execPod, err := GetMayastorPodNameonNodeByPrefix(deployName, nodeName)
	if err != nil {
		logf.Log.Info("GetMayastorPodNameonNodeByPrefix failed", "error", err)
		return "", err
	}

	// Get nvmename of the device
	cmd := "lsblk -fa|grep nvme|awk '{print $1}'"
	output, errOutput, err := ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running lsblk", "error", err)
		return "", err
	}
	logf.Log.Info("lsblk output", "Output", output)

	nvmeName := strings.TrimSpace(output)
	devicePath := "/dev/" + nvmeName

	// Generate random name for temporary directory to mount xfs filesystem
	cmd = "uuidgen"
	output, errOutput, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running uuidgen", "error", err)
		return "", err
	}
	outputString := strings.TrimSpace(output)

	// Create a temporary directory
	cmd = fmt.Sprintf("mkdir %s", outputString)
	_, errOutput, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running mkdir", "error", err)
		return "", err
	}

	// Mount xfs filesystem on temporary directory
	cmd = fmt.Sprintf("mount	%s %s -t xfs", devicePath, outputString)
	_, errOutput, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running mount", "error", err)
		return "", err
	}

	// Umount the directory
	cmd = fmt.Sprintf("umount %s", outputString)
	_, errOutput, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running unmount", "error", err)
		return "", err
	}

	// Delete the temporary directoy
	cmd = fmt.Sprintf("rm -rf %s", outputString)
	_, errOutput, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil || errOutput != "" {
		logf.Log.Info("ExecuteCommandInContainer failed while running rm", "error", err)
		return "", err
	}

	// Perform xfs_repair(fsck equivalent for xfs) on the nvme device with xfs filesystem
	cmd = fmt.Sprintf("xfs_repair -n	%s; echo $?", devicePath)
	output, _, err = ExecuteCommandInContainer(common.NSMayastor(), execPod, containerName, cmd)
	if err != nil {
		logf.Log.Info("ExecuteCommandInContainer failed while running xfs_repair", "error", err)
		return "", err
	}

	return output, nil
}
