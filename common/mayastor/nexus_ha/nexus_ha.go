package nexus_ha

import (
	"fmt"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/apps"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/mayastorclient"
	"github.com/openebs/openebs-e2e/common/platform"
	"github.com/openebs/openebs-e2e/common/platform/types"

	coreV1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	TimeoutSecs             = 60  // in seconds
	DefTimeoutSecs          = 120 // in seconds
	ReplicaCheckTimeoutSecs = 180 // in seconds
	RefreshTimeout          = 15  // in seconds
	DefaultvolSizeMb        = 4096
	DefaultFsVolSizeMb      = 512
	DefaultFioRunTime       = 600   // in seconds
	DefaultFioSleepTime     = 10000 // in seconds
	HaNodeTestingLabelKey   = "mayastor.testing/e2e-ha-nexus"
	HaNodeTestingLabelValue = "ha-testing"
	DefRebuildTimeoutSecs   = 240 // in seconds
)

var (
	EngineLabel                  = e2e_config.GetConfig().Product.EngineLabel
	EngineLabelValue             = e2e_config.GetConfig().Product.EngineLabelValue
	lokiStatefulset              = e2e_config.GetConfig().Product.LokiStatefulset
	lokiStatefulsetOnControlNode = e2e_config.GetConfig().LokiStatefulsetOnControlNode
	IoEnginePodName              = e2e_config.GetConfig().Product.IOEnginePodName
)

var msDeployment = []string{
	e2e_config.GetConfig().Product.ControlPlanePoolOperator,
	e2e_config.GetConfig().Product.ControlPlaneRestServer,
	e2e_config.GetConfig().Product.ControlPlaneCsiController,
	e2e_config.GetConfig().Product.ControlPlaneCoreAgent,
	e2e_config.GetConfig().Product.ControlPlaneLocalpvProvisioner,
	e2e_config.GetConfig().Product.ControlPlaneObsCallhome,
}

var msAppLabels = []string{
	e2e_config.GetConfig().AppLabelControlPlanePoolOperator,
	e2e_config.GetConfig().AppLabelControlPlaneRestServer,
	e2e_config.GetConfig().AppLabelControlPlaneCsiController,
	e2e_config.GetConfig().AppLabelControlPlaneCoreAgent,
}

type NexusHa struct {
	TestPodName      string
	TestNode         string
	ScName           string
	TestDeployName   string
	TestPodAppLabels string
	VolName          string
	NexusNode        string
	Uuid             string
	Platform         types.Platform
	NonNexusNode     string
}

func (c *NexusHa) CreateFioDeploy(fioArgs []string) error {
	labelKey := "e2e-test"
	labelValue := "maxsnap"
	labelselector := map[string]string{
		labelKey: labelValue,
	}
	mount := coreV1.VolumeMount{
		Name:      "ms-volume",
		MountPath: common.FioFsMountPoint,
	}
	var volMounts []coreV1.VolumeMount
	volMounts = append(volMounts, mount)

	logf.Log.Info("fio", "arguments", fioArgs)
	deployObj, err := k8stest.NewDeploymentBuilder().
		WithName(c.TestDeployName).
		WithNamespace(common.NSDefault).
		WithLabelsNew(labelselector).
		WithSelectorMatchLabelsNew(labelselector).
		WithPodTemplateSpecBuilder(
			k8stest.NewPodtemplatespecBuilder().
				WithLabels(labelselector).
				WithContainerBuildersNew(
					k8stest.NewContainerBuilder().
						WithName(c.TestDeployName).
						WithImage(common.GetFioImage()).
						WithVolumeMountsNew(volMounts).
						WithImagePullPolicy(coreV1.PullAlways).
						WithArgumentsNew(fioArgs)).
				WithVolumeBuilders(
					k8stest.NewVolumeBuilder().
						WithName("ms-volume").
						WithPVCSource(c.VolName),
				),
		).
		Build()

	if err != nil {
		return fmt.Errorf("failed to create deployment %s definittion object in %s namesppace", c.TestDeployName, common.NSDefault)
	}

	err = k8stest.CreateDeployment(deployObj)
	if err != nil {
		return fmt.Errorf("failed to create deployment %s, error %v", c.TestDeployName, err)
	}
	labels := fmt.Sprintf("%s=%s", labelKey, labelValue)
	c.TestPodAppLabels = labels
	return verifyApplicationPodRunning(c.TestDeployName, common.NSDefault, labels, true)
}

// CreatePodAndNexusOnDifferentNode A function to create pod and nexus on different node
func (c *NexusHa) CreatePodAndNexusOnDifferentNode(replicas int, prefix string, fioArgs []string, volSizeMb int) error {
	var errs common.ErrorAccumulator
	nodeName, err := UnlabelIoEngineFromNode("")
	if err != nil {
		return fmt.Errorf("failed to unlabel the pod %v", err)
	}
	logf.Log.Info("io-engine is removed from", "node", nodeName)

	err = c.CreateScAndVolume(prefix, replicas, volSizeMb)
	if err != nil {
		return fmt.Errorf("failed to create storage %v", err)
	}
	logf.Log.Info("storage", "ScName", c.ScName)
	fioPodName := "fio-" + c.VolName
	c.TestPodName = fioPodName
	// create fio pod on other node
	err = k8stest.CreateFioPodOnNode(fioPodName, c.VolName, nodeName, fioArgs)
	if err != nil {
		return fmt.Errorf("failed to create pod on node %v", err)
	}

	if !k8stest.WaitPodRunning(fioPodName, common.NSDefault, TimeoutSecs) {
		return fmt.Errorf("failed to start pod on node")
	}

	msv, err := k8stest.GetMSV(c.Uuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve msv %v", err)
	}
	logf.Log.Info("Current Nexus location", "->", msv)
	c.TestNode = nodeName
	c.TestPodName = fioPodName
	c.NexusNode = msv.State.Target.Node

	// make io engine pod online
	err = LabelIoEngineOnNode(nodeName)
	if err != nil {
		return fmt.Errorf("failed to label node %s, error %v", nodeName, err)
	}
	logf.Log.Info("Delay for 5 seconds")
	time.Sleep(5 * time.Second)

	//Get non nexus and non application node
	lokiNode, err := GetNonNexusNonTestPodNode(c.NexusNode, c.TestNode)
	if err != nil {
		return fmt.Errorf("failed to get non nexus non test node %v", err)
	}
	c.NonNexusNode = lokiNode

	if !lokiStatefulsetOnControlNode {
		// reschedule loki stateful set to different node
		err = k8stest.ApplyNodeSelectorToStatefulset(lokiStatefulset, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
		if err != nil {
			return fmt.Errorf("failed to reschedule loki %v", err)
		}
	}

	// reschedule control plane components to different node
	for _, deploy := range msDeployment {
		isPresent, _ := k8stest.PodPresentOnNode(deploy, common.NSMayastor(), lokiNode)
		if !isPresent {
			err = k8stest.ApplyNodeSelectorToDeployment(deploy, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
			if err != nil {
				errs.Accumulate(err)
			}
		}
	}

	logf.Log.Info("loki node scheduled on", "node->", lokiNode, "err", errs.GetError())

	err = k8stest.VerifyPodsOnNode(msAppLabels, lokiNode, common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to verify node reschedule %v", err)
	}

	logf.Log.Info("Mayastor ready check")
	// verify mayastor ready check
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation not ready")
	}

	logf.Log.Info(fmt.Sprintf("Checking mayastor volume %s for control plane readiness within %d seconds", c.Uuid, DefTimeoutSecs))
	t1 := time.Now()
	logf.Log.Info("start", "time", t1)
	for time.Since(t1).Seconds() < float64(DefTimeoutSecs) {
		_, err = k8stest.GetMSV(c.Uuid)
		if err == nil {
			break
		}
		logf.Log.Info("failed to get volume replicas", "uuid", c.Uuid, "Error", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get volume %s, error: %v", c.Uuid, err)
	}
	c.Platform = platform.Create()
	if err != nil {
		return fmt.Errorf("failed to create platform %v", err)
	}
	logf.Log.Info("Nexus and test pod on different node", "nexus node->", c.NexusNode, "app node->", nodeName)
	return nil
}

// CreateMongoPodAndNexusOnDifferentNode A function to create pod and nexus on different node
func (c *NexusHa) CreateMongoPodAndNexusOnDifferentNode(replicas int, result *string, volSizeMb int, wg *sync.WaitGroup) (apps.MongoApp, error) {
	var errs common.ErrorAccumulator
	nodeName, err := UnlabelIoEngineFromNode("")
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to unlabel the pod %v", err)
	}
	logf.Log.Info("io-engine is removed from", "node", nodeName)

	err = c.CreateScAndVolume("mongo", replicas, volSizeMb)
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to create storage %v", err)
	}
	logf.Log.Info("storage", "ScName", c.ScName)
	app, err := apps.NewMongoBuilder().
		WithYcsb().
		WithOwnStorageClass(c.ScName).
		WithPvc(c.VolName).
		WithNodeSelector(nodeName).
		Build()
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to deploy app %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Ycsb.BenchmarkParams.RecordCount = 1_000_000
		app.Ycsb.BenchmarkParams.ThreadCount = 1
		err = app.Ycsb.LoadYcsbApp()
		err = app.Ycsb.RunYcsbApp(result)
	}()

	msv, err := k8stest.GetMSV(c.Uuid)
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to retrieve msv %v", err)
	}
	logf.Log.Info("Current Nexus location", "->", msv)
	c.TestNode = nodeName
	c.TestPodName = app.Mongo.Pod.Name
	c.NexusNode = msv.State.Target.Node

	// check if nexus is not empty
	n, err := k8stest.GetNexusNode(c.Uuid)
	if err != nil {
		return apps.MongoApp{}, err
	}
	if n == "" {
		return apps.MongoApp{}, fmt.Errorf("nexus is empty for volUuid %s", c.Uuid)
	}

	//restore io-engine pod on the node once the nexus is created on other node
	err = LabelIoEngineOnNode(nodeName)
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to label node %s, error %v", nodeName, err)
	}
	logf.Log.Info("Delay for 5 seconds")
	time.Sleep(5 * time.Second)

	//Get non nexus and non application node
	lokiNode, err := GetNonNexusNonTestPodNode(c.NexusNode, c.TestNode)
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to get non nexus non test node %v", err)
	}
	c.NonNexusNode = lokiNode

	if !lokiStatefulsetOnControlNode {
		// reschedule loki stateful set to different node
		err = k8stest.ApplyNodeSelectorToStatefulset(lokiStatefulset, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
		if err != nil {
			return apps.MongoApp{}, fmt.Errorf("failed to reschedule loki %v", err)
		}
	}

	// reschedule control plane components to different node
	for _, deploy := range msDeployment {
		isPresent, _ := k8stest.PodPresentOnNode(deploy, common.NSMayastor(), lokiNode)
		if !isPresent {
			err = k8stest.ApplyNodeSelectorToDeployment(deploy, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
			if err != nil {
				errs.Accumulate(err)
			}
		}
	}

	logf.Log.Info("loki node scheduled on", "node->", lokiNode, "err", errs.GetError())

	err = k8stest.VerifyPodsOnNode(msAppLabels, lokiNode, common.NSMayastor())
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to verify node reschedule %v", err)
	}

	logf.Log.Info("Mayastor ready check")
	// verify mayastor ready check
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return apps.MongoApp{}, err
	}
	if !ready {
		return apps.MongoApp{}, fmt.Errorf("mayastor installation not ready")
	}

	logf.Log.Info(fmt.Sprintf("Checking mayastor volume %s for control plane readiness within %d seconds", c.Uuid, DefTimeoutSecs))
	t1 := time.Now()
	logf.Log.Info("start", "time", t1)
	for time.Since(t1).Seconds() < float64(DefTimeoutSecs) {
		_, err = k8stest.GetMSV(c.Uuid)
		if err == nil {
			break
		}
		logf.Log.Info("failed to get volume replicas", "uuid", c.Uuid, "Error", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to get volume %s, error: %v", c.Uuid, err)
	}
	c.Platform = platform.Create()
	if err != nil {
		return apps.MongoApp{}, fmt.Errorf("failed to create platform %v", err)
	}
	logf.Log.Info("Nexus and test pod on different node", "nexus node->", c.NexusNode, "app node->", nodeName)
	return app, nil
}

// CreateAnotherPodAndNexus A function to create pod and nexus on different node
func (c *NexusHa) CreateAnotherPodAndNexus(replicas int, prefix string, fioArgs []string) error {
	nodeName, err := GetNonNexusNonTestPodNode(c.NexusNode, c.TestNode)
	if err != nil {
		return fmt.Errorf("failed to get non nexus non test pod node %v", err)
	}

	nodeName, err = UnlabelIoEngineFromNode(nodeName)
	if err != nil {
		return fmt.Errorf("failed to ulabel the pod  %v", err)
	}
	logf.Log.Info("io-engine is removed from", "node", nodeName)

	err = c.CreateScAndVolume(prefix+"-sec", replicas, DefaultvolSizeMb)
	if err != nil {
		return fmt.Errorf("failed to create storage %v", err)
	}
	logf.Log.Info("storage", "ScName", c.ScName)
	fioPodName := "fio-" + c.VolName
	c.TestPodName = fioPodName

	// create fio pod on other node
	err = k8stest.CreateFioPodOnNode(fioPodName, c.VolName, nodeName, fioArgs)
	if err != nil {
		return fmt.Errorf("failed to create pod on node %v", err)
	}

	if !k8stest.WaitPodRunning(fioPodName, common.NSDefault, TimeoutSecs) {
		return fmt.Errorf("failed to start pod on node")
	}

	msv, err := k8stest.GetMSV(c.Uuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve msv %v", err)
	}
	logf.Log.Info("Current Nexus location", "->", msv)
	c.TestNode = nodeName
	c.NexusNode = msv.State.Target.Node

	//restore io-engine pod on the node once the nexus is created on other node
	err = LabelIoEngineOnNode(nodeName)
	if err != nil {
		return fmt.Errorf("failed to label the pod %v", err)
	}
	logf.Log.Info("Delay for 5 seconds")
	time.Sleep(5 * time.Second)
	// verify mayastor ready check
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation not ready")
	}
	c.Platform = platform.Create()

	logf.Log.Info("Sec Nexus,app pod on different node", "nexus node->", c.NexusNode, "app node->", nodeName)
	return nil
}

func (c *NexusHa) CreateScAndVolume(prefix string, replicas int, volSizeMb int) error {
	ScName := fmt.Sprintf("%s-repl-%d", prefix, replicas)
	err := k8stest.NewScBuilder().
		WithName(ScName).
		WithNamespace(common.NSDefault).
		WithProtocol(common.ShareProtoNvmf).
		WithReplicas(replicas).
		WithVolumeBindingMode(storageV1.VolumeBindingImmediate).
		WithReclaimPolicy(coreV1.PersistentVolumeReclaimDelete).
		BuildAndCreate()
	if err != nil {
		return fmt.Errorf("failed to create storage class %s", ScName)
	}
	VolName := fmt.Sprintf("vol-%s", ScName)
	uid, err := k8stest.MkPVC(volSizeMb, VolName, ScName, common.VolFileSystem, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to create pvc %s", VolName)
	}
	c.ScName = ScName
	c.VolName = VolName
	c.Uuid = uid
	return nil
}

func (c *NexusHa) DeleteScAndVolume() error {
	pvc, err := k8stest.GetPVC(c.VolName, common.NSDefault)
	if err != nil {
		logf.Log.Error(err, "failed to get pvc", "pvc", c.VolName)
		return err
	}
	logf.Log.Info("pvc", "pvName ->", pvc.Spec.VolumeName, "volume->", c.VolName)

	//Delete the volume
	err = k8stest.RmPVC(c.VolName, c.ScName, common.NSDefault)
	if err != nil {
		logf.Log.Error(err, "failed to delete", "pvc", c.VolName)
		return err
	}

	err = k8stest.RmStorageClass(c.ScName)
	if err != nil {
		logf.Log.Error(err, "failed to delete ", "storage class", c.ScName)
		return err
	}
	return nil
}

func (c *NexusHa) CleanUp() error {
	// Delete the fio pod
	if c.TestDeployName != "" {
		err := deleteDeployment(c.TestDeployName, common.NSDefault, c.TestPodAppLabels)
		if err != nil {
			return fmt.Errorf("failed to delete pod %s", c.TestPodName)
		}
	}
	if c.TestPodName != "" {
		err := k8stest.DeletePod(c.TestPodName, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to delete pod %s", c.TestPodName)
		}
	}
	return c.DeleteScAndVolume()
}

func (c *NexusHa) VerifyReplicasAdnCleanUp() error {
	// Delete the fio pod
	if c.TestDeployName != "" {
		err := deleteDeployment(c.TestDeployName, common.NSDefault, c.TestPodAppLabels)
		if err != nil {
			return fmt.Errorf("failed to delete pod %s", c.TestPodName)
		}
	}
	if c.TestPodName != "" {
		err := k8stest.DeletePod(c.TestPodName, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to delete pod %s", c.TestPodName)
		}
	}
	res := k8stest.CompareVolumeReplicas(c.VolName, common.NSDefault)
	switch res.Result {
	case common.CmpReplicasMismatch:
		return fmt.Errorf("replica contents do not match %s", res.Description)
	case common.CmpReplicasFailed:
		return fmt.Errorf("replica comparison failed %s ; %v", res.Description, res.Err)
	case common.CmpReplicasMatch:
		return c.DeleteScAndVolume()
	default:
		return fmt.Errorf("unexpected result value %d", res.Result)
	}
}

// GetNonNexusNode return node where nexus is absent
func GetNonNexusNode() (string, string) {
	nodeList, err := k8stest.GetIOEngineNodeLocsMap()
	if err != nil {
		logf.Log.Error(err, "failed to get node location %v", err)
		return "", ""
	}
	var nodeName string
	var address string
	numOfRecords := len(nodeList)
	skip := false
	if numOfRecords > 2 {
		skip = true
	}

	for _, node := range nodeList {
		logf.Log.Info("node", "->", node)
		nexuses, err := mayastorclient.ListNexuses([]string{node.IPAddress})
		if len(nexuses) > 0 || err != nil || skip {
			skip = false
			continue
		}
		nodeName = node.NodeName
		address = node.IPAddress
		break
	}
	return nodeName, address
}

func UnlabelIoEngineFromNode(nodeName string) (string, error) {
	if len(nodeName) == 0 {
		nodeName, _ = GetNonNexusNode()
		if nodeName == "" {
			return "", fmt.Errorf("failed to retrieve node or IP %v", nodeName)
		}
	}

	logf.Log.Info("unlabel node->", "nodeName->", nodeName, "EngineLabel->", EngineLabel)
	// remove the io-engine pod
	err := k8stest.UnlabelNode(nodeName, EngineLabel)
	if err != nil {
		return "", err
	}
	logf.Log.Info(fmt.Sprintf("Pause %d seconds for control plane refresh", RefreshTimeout))
	time.Sleep(RefreshTimeout * time.Second)
	return nodeName, nil
}

func LabelIoEngineOnNode(nodeName string) error {
	logf.Log.Info("Restore io engine")
	logf.Log.Info("label node->", "nodeName->", nodeName, "EngineLabel->", EngineLabel)
	// add the io-engine pod
	err := k8stest.LabelNode(nodeName, EngineLabel, EngineLabelValue)
	if err != nil {
		return err
	}
	logf.Log.Info(fmt.Sprintf("Pause %d seconds for control plane refresh", RefreshTimeout))
	time.Sleep(RefreshTimeout * time.Second)
	return nil
}

// VerifyNodeReady verify node ready status
func VerifyNodeReady(nodeName string, timeoutSecs int) bool {
	const sleepTime = 3
	time.Sleep(sleepTime * time.Second)
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		readyStatus, err := k8stest.IsNodeReady(nodeName, nil)
		logf.Log.Info("node ready", "status ->", readyStatus)
		if err != nil || !readyStatus {
			time.Sleep(sleepTime * time.Second)
			continue
		}
		return readyStatus
	}
	return false
}

// VerifyNodeNotReady verify node ready status
func VerifyNodeNotReady(nodeName string, timeoutSecs int) bool {
	const sleepTime = 3
	time.Sleep(sleepTime * time.Second)
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		readyStatus, err := k8stest.IsNodeReady(nodeName, nil)
		logf.Log.Info("node ready", "status ->", readyStatus)
		if err != nil || readyStatus {
			time.Sleep(sleepTime * time.Second)
			continue
		}
		return readyStatus
	}
	return true
}

// VerifyNewNexusCreated check new nexus is created on other node
func VerifyNewNexusCreated(uuid string, oldNexus string) bool {
	newNexus, err := k8stest.GetMSV(uuid)
	logf.Log.Info("New", "nexus ->", newNexus)
	if newNexus == nil || err != nil {
		logf.Log.Info("failed to get nexus", "err", err)
		return false
	}

	if newNexus != nil && (newNexus.State.Target.Node == oldNexus || len(newNexus.State.Target.Node) == 0) {
		logf.Log.Info("New Nexus", "node->", newNexus.State.Target.Node)
		logf.Log.Info("Old Nexus", "node->", oldNexus)
		return false
	}
	return true
}

// VerifyOldNexusReCreated check old nexus is re created on same node
func VerifyOldNexusReCreated(uuid string, oldNexus string) bool {
	newNexus, err := k8stest.GetMSV(uuid)
	logf.Log.Info("New", "nexus ->", newNexus)
	if newNexus == nil || err != nil {
		logf.Log.Info("failed to get nexus", "err", err)
		return false
	}

	if newNexus != nil && (newNexus.State.Target.Node != oldNexus || len(newNexus.State.Target.Node) == 0) {
		logf.Log.Info("New Nexus", "node->", newNexus.State.Target.Node)
		logf.Log.Info("Old Nexus", "node->", oldNexus)
		return false
	}
	return true
}

// VerifyNexusRemoved verify old nexus has been removed
func VerifyNexusRemoved(uuid string, oldNexus string) bool {
	newNexus, _ := k8stest.GetMSV(uuid)
	if newNexus == nil {
		logf.Log.Info("failed to get volume with uuid %s", uuid)
		return false
	}
	// Ensure that no nexus present in the node
	nodeIP, err := k8stest.GetNodeIPAddress(oldNexus)
	if err != nil {
		return false
	}
	ips := []string{*nodeIP}
	logf.Log.Info("Node IP", "address->", ips)
	// verify : old node has been removed the nexus
	nexuses, err := mayastorclient.ListNexuses(ips)
	if err != nil {
		return false
	}
	if len(nexuses) > 0 {
		logf.Log.Info("nexus", "list", nexuses)
		return false
	}
	return true
}

// GetNonNexusNonTestPodNode Get non nexus , non test pod node
func GetNonNexusNonTestPodNode(nexusNode string, testPodNode string) (string, error) {
	nodeList, err := k8stest.GetIOEngineNodeLocsMap()
	if err != nil {
		logf.Log.Error(err, "failed to get node location %v", err)
		return "", err
	}

	for _, node := range nodeList {
		if node.NodeName != nexusNode && node.NodeName != testPodNode {
			logf.Log.Info("Non nexus and non test", "node->", node.NodeName)
			return node.NodeName, nil
		}
	}
	return "", fmt.Errorf("failed to find non nexus and nod test pod node")
}

func VerifyVolumeHealth(uuid string) bool {
	msv, err := k8stest.GetMSV(uuid)
	if err == nil {
		if msv.State.Status == controlplane.VolStateHealthy() {
			logf.Log.Info("Volume health", "state ->", msv.State.Status)
			return true
		}
	}
	return false
}

func AddHaLabelToNodes() error {
	nodes, err := k8stest.GetMayastorNodeNames()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err = k8stest.LabelNode(node, HaNodeTestingLabelKey, HaNodeTestingLabelValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveHaLabelFromNodes() error {
	nodes, err := k8stest.GetMayastorNodeNames()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err = k8stest.UnlabelNode(node, HaNodeTestingLabelKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFioArguments(fioRunTime int) ([]string, error) {
	return common.NewE2eFioArgsBuilder().
		WithDefaultFile().
		WithRuntime(fioRunTime).
		WithDefaultArgs().
		Build()
}

// CreateDeploymentAndNexusOnDifferentNode A function to create fio deployment and nexus on different node
func (c *NexusHa) CreateDeploymentAndNexusOnDifferentNode(replicas int, prefix string, fioArgs []string, volSizeMb int) error {
	var errs common.ErrorAccumulator
	nodeName, err := UnlabelIoEngineFromNode("")
	if err != nil {
		return fmt.Errorf("failed to unlabel the pod %v", err)
	}
	logf.Log.Info("io-engine is removed from", "node", nodeName)

	err = c.CreateScAndVolume(prefix, replicas, volSizeMb)
	if err != nil {
		return fmt.Errorf("failed to create storage %v", err)
	}
	logf.Log.Info("storage", "ScName", c.ScName)
	fioDeployName := "fio-" + c.VolName
	c.TestDeployName = fioDeployName
	// create fio deployment on other node
	err = c.createFioDeployment(fioArgs, nodeName)
	if err != nil {
		return fmt.Errorf("failed to create fio deployment %s on node %v", fioDeployName, err)
	}

	msv, err := k8stest.GetMSV(c.Uuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve msv %v", err)
	}
	logf.Log.Info("Current Nexus location", "->", msv)
	c.TestNode = nodeName
	c.NexusNode = msv.State.Target.Node

	// make io engine pod online
	err = LabelIoEngineOnNode(nodeName)
	if err != nil {
		return fmt.Errorf("failed to label node %s, error %v", nodeName, err)
	}
	logf.Log.Info("Delay for 5 seconds")
	time.Sleep(5 * time.Second)

	//Get non nexus and non application node
	lokiNode, err := GetNonNexusNonTestPodNode(c.NexusNode, c.TestNode)
	if err != nil {
		return fmt.Errorf("failed to get non nexus non test node %v", err)
	}
	c.NonNexusNode = lokiNode

	if !lokiStatefulsetOnControlNode {
		// reschedule loki stateful set to different node
		err = k8stest.ApplyNodeSelectorToStatefulset(lokiStatefulset, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
		if err != nil {
			return fmt.Errorf("failed to reschedule loki %v", err)
		}
	}

	// reschedule control plane components to different node
	for _, deploy := range msDeployment {
		isPresent, _ := k8stest.PodPresentOnNode(deploy, common.NSMayastor(), lokiNode)
		if !isPresent {
			err = k8stest.ApplyNodeSelectorToDeployment(deploy, common.NSMayastor(), "kubernetes.io/hostname", lokiNode)
			if err != nil {
				errs.Accumulate(err)
			}
		}
	}

	logf.Log.Info("loki node scheduled on", "node->", lokiNode, "err", errs.GetError())

	err = k8stest.VerifyPodsOnNode(msAppLabels, lokiNode, common.NSMayastor())
	if err != nil {
		return fmt.Errorf("failed to verify node reschedule %v", err)
	}

	logf.Log.Info("Mayastor ready check")
	// verify mayastor ready check
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return err
	}
	if !ready {
		return fmt.Errorf("mayastor installation not ready")
	}

	logf.Log.Info(fmt.Sprintf("Checking mayastor volume %s for control plane readiness within %d seconds", c.Uuid, DefTimeoutSecs))
	t1 := time.Now()
	logf.Log.Info("start", "time", t1)
	for time.Since(t1).Seconds() < float64(DefTimeoutSecs) {
		_, err = k8stest.GetMSV(c.Uuid)
		if err == nil {
			break
		}
		logf.Log.Info("failed to get volume replicas", "uuid", c.Uuid, "Error", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to get volume %s, error: %v", c.Uuid, err)
	}

	c.Platform = platform.Create()
	if err != nil {
		return fmt.Errorf("failed to create platform %v", err)
	}
	logf.Log.Info("Nexus and test pod on different node", "nexus node->", c.NexusNode, "app node->", nodeName)
	return nil
}

func (c *NexusHa) createFioDeployment(fioArgs []string, nodeName string) error {
	labelKey := "e2e-test"
	labelValue := "ha"
	labelselector := map[string]string{
		labelKey: labelValue,
	}
	nodeSelector := map[string]string{
		"kubernetes.io/hostname": nodeName,
	}
	mount := coreV1.VolumeMount{
		Name:      "ms-volume",
		MountPath: common.FioFsMountPoint,
	}
	var volMounts []coreV1.VolumeMount
	volMounts = append(volMounts, mount)

	logf.Log.Info("fio", "arguments", fioArgs)
	deployObj, err := k8stest.NewDeploymentBuilder().
		WithName(c.TestDeployName).
		WithNamespace(common.NSDefault).
		WithLabelsNew(labelselector).
		WithSelectorMatchLabelsNew(labelselector).
		WithPodTemplateSpecBuilder(
			k8stest.NewPodtemplatespecBuilder().
				WithNodeSelector(nodeSelector).
				WithLabels(labelselector).
				WithContainerBuildersNew(
					k8stest.NewContainerBuilder().
						WithName(c.TestDeployName).
						WithImage(common.GetFioImage()).
						WithVolumeMountsNew(volMounts).
						WithImagePullPolicy(coreV1.PullAlways).
						WithArgumentsNew(fioArgs)).
				WithVolumeBuilders(
					k8stest.NewVolumeBuilder().
						WithName("ms-volume").
						WithPVCSource(c.VolName),
				),
		).
		Build()

	if err != nil {
		return fmt.Errorf("failed to create deployment %s definittion object in %s namesppace", c.TestDeployName, common.NSDefault)
	}

	err = k8stest.CreateDeployment(deployObj)
	if err != nil {
		return fmt.Errorf("failed to create deployment %s, error %v", c.TestDeployName, err)
	}
	labels := fmt.Sprintf("%s=%s", labelKey, labelValue)
	c.TestPodAppLabels = labels
	return verifyApplicationPodRunning(c.TestDeployName, common.NSDefault, labels, true)
}

func verifyApplicationPodRunning(deployName string, namespace string, labels string, state bool) error {
	logf.Log.Info("Verify application deployment ready", "deployment", deployName, "namespace", namespace)

	const sleepTime = 3
	var deployReady bool
	for ix := 0; ix < (DefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
		if state == k8stest.DeploymentReady(deployName, namespace) {
			deployReady = true
		}
		if deployReady {
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	if !deployReady {
		return fmt.Errorf("fio deployment %s not ready", deployName)
	}

	logf.Log.Info("Verify application pod running", "labels", labels)

	var podReady bool
	for ix := 0; ix < (DefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
		_, ready, err := k8stest.IsPodWithLabelsRunning(labels, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to check pod with labels %s running status, error: %v", labels, err)
		}
		if ready == state {
			podReady = true
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	if !podReady {
		return fmt.Errorf("fio app pod with label %v not ready", deployName)
	}
	return nil
}

func deleteDeployment(deployName string, namespace string, labels string) error {
	err := k8stest.DeleteDeployment(deployName, namespace)
	if err != nil {
		return fmt.Errorf("failed to delete deployment %s in namespcae %s, error: %v", deployName, namespace, err)
	}
	return verifyApplicationPodRunning(deployName, namespace, labels, false)
}

func WaitForVolumeToBeHealthy(uuid string) (bool, error) {
	var msv *common.MayastorVolume
	var isHealthy bool
	var err error
	t1 := time.Now()
	logf.Log.Info("start", "time", t1)
	logf.Log.Info(fmt.Sprintf("Checking mayastor volume state to be Healthy within %v seconds", DefRebuildTimeoutSecs))
	for time.Since(t1).Seconds() < float64(DefRebuildTimeoutSecs) {
		msv, err = k8stest.GetMSV(uuid)
		if err == nil && msv.State.Status == controlplane.VolStateHealthy() {
			isHealthy = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if msv != nil {
		logf.Log.Info("mayastor volume", "Status.State", msv.State.Status, "after", time.Since(t1))
	}
	return isHealthy, err
}

func GetOnlineReplicasCount(uuid string) (int, error) {
	var err error
	var onlineReplicaCount int
	var replicas []common.MsvReplica
	logf.Log.Info(fmt.Sprintf("Checking mayastor volume %s online replicas within %d seconds", uuid, DefTimeoutSecs))
	replicas, err = k8stest.GetMsvReplicas(uuid)
	if err != nil {
		return onlineReplicaCount, fmt.Errorf("failed to get online replica count for volume %s, error: %v", uuid, err)
	}
	for _, replica := range replicas {
		if replica.Replica.State == controlplane.ReplicaStateOnline() {
			onlineReplicaCount++
		} else {
			logf.Log.Info("Replica", "node", replica.Replica.Node, "state", replica.Replica.State, "uri", replica.Uuid)
		}
	}
	return onlineReplicaCount, err
}

func VerifyMayastorAndPoolReady() (bool, error) {
	return k8stest.VerifyMayastorAndPoolReady(DefTimeoutSecs)
}
