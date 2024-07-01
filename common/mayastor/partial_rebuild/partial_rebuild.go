package partial_rebuild

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/apps"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	mayastorclient "github.com/openebs/openebs-e2e/common/mayastorclient"
	"github.com/openebs/openebs-e2e/common/platform"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs               = 180 // in seconds
	DefSleepTime                 = 30  //in seconds
	DefRebuildStatsTimeout       = 60  //in seconds
	DefRebuildTimeoutSecs        = 360 // in seconds
	IoEnginePodName              = e2e_config.GetConfig().Product.IOEnginePodName
	LokiStatefulsetOnControlNode = e2e_config.GetConfig().LokiStatefulsetOnControlNode
	LokiStatefulset              = e2e_config.GetConfig().Product.LokiStatefulset
	PartialRebuildCpTimeout      = e2e_config.GetConfig().Product.PartialRebuildCpTimeout
	FioRunTime                   = 1200 // in seconds
	PartialRebuildWaitPeriodArg  = "--faulted-child-wait-period"
	VolSizeMb                    = 8192 // in  MiB
	SleepTimeSecs                = 20   // in second
	ReconcileTimeSecs            = 30   // in seconds
	PodCompletionWaitTimeSecs    = 1200 // in seconds
)

var MsCpDeployment = []string{
	e2e_config.GetConfig().Product.ControlPlanePoolOperator,
	e2e_config.GetConfig().Product.ControlPlaneRestServer,
	e2e_config.GetConfig().Product.ControlPlaneCsiController,
	e2e_config.GetConfig().Product.ControlPlaneCoreAgent,
	e2e_config.GetConfig().Product.ControlPlaneLocalpvProvisioner,
	e2e_config.GetConfig().Product.ControlPlaneObsCallhome,
}

var MsAppLabels = []string{
	e2e_config.GetConfig().AppLabelControlPlanePoolOperator,
	e2e_config.GetConfig().AppLabelControlPlaneRestServer,
	e2e_config.GetConfig().AppLabelControlPlaneCsiController,
	e2e_config.GetConfig().AppLabelControlPlaneCoreAgent,
}

type TestApp struct {
	App                k8stest.FioApp
	NexusNode          string
	NexusNodeIp        *string
	NexusUuid          string
	FaultReplicaNode   string
	FaultReplicaNodeIp *string
	SrcFaultReplicaUri string
	DstReplicaUri      string
	InitialReplicas    map[string]string
}

type TestAppMongo struct {
	App                apps.MongoApp
	NexusNode          string
	NexusNodeIp        *string
	NexusUuid          string
	FaultReplicaNode   string
	FaultReplicaNodeIp *string
	SrcFaultReplicaUri string
	DstReplicaUri      string
	InitialReplicas    map[string]string
}

func WaitForRebuildInProgress(uuid string, timeoutSecs int) (bool, error) {
	var isRebuildInprogress bool
	var err error
	var msv *common.MayastorVolume
	const sleepTime = 3
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		msv, err = k8stest.GetMSV(uuid)
		if err == nil &&
			msv.State.Status == controlplane.VolStateDegraded() &&
			msv.State.Target.Rebuilds != 0 {
			isRebuildInprogress = true
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	return isRebuildInprogress, err
}

func WaitForReplicaRebuildInProgress(uuid string, replicaUri string, timeoutSecs int) (bool, error) {
	const sleepTime = 1
	timeout := time.Duration(timeoutSecs) * time.Second
	var err error

	startTime := time.Now()
	for {
		var nexusChildren []common.TargetChild
		nexusChildren, err = k8stest.GetMsvNexusChildren(uuid)
		if err != nil {
			logf.Log.Info("failed to get Nexus children", "error", err)
		} else {
			logf.Log.Info("Nexus", "Children", nexusChildren)
			for _, child := range nexusChildren {
				if child.Uri == replicaUri &&
					child.RebuildProgress != nil &&
					child.State == controlplane.ChildStateDegraded() {
					return true, nil
				}
			}
		}

		if time.Since(startTime) >= timeout {
			break
		}

		time.Sleep(sleepTime * time.Second)
	}

	return false, err
}
func WaitForRebuildComplete(uuid string, timeoutSecs int) (bool, error) {
	var isRebuildCompleted bool
	var err error
	var msv *common.MayastorVolume
	const sleepTime = 3
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		msv, err = k8stest.GetMSV(uuid)
		if err == nil &&
			msv.State.Status == controlplane.VolStateHealthy() &&
			msv.State.Target.Rebuilds == 0 &&
			/* The volume status may be reflected incorrectly as healthy,
			during transition periods when the number replicas is lower than that in the volume spec.
			Given enough time the status will transition to the expected value.*/
			len(msv.State.Target.Children) == msv.Spec.Num_replicas {
			isRebuildCompleted = true
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	return isRebuildCompleted, err
}

func GetRebuildHistoryRecord(nexusUuid string, nexusNodeIp string) ([]mayastorclient.RebuildHistoryRecord, error) {
	rebuildHistory, err := mayastorclient.GetRebuildHistory(nexusUuid, nexusNodeIp)
	if err == nil && rebuildHistory != nil {
		logf.Log.Info("GetRebuildHistoryRecord", "Records", rebuildHistory.GetRecords())
		return rebuildHistory.GetRecords(), err
	} else if err == nil && rebuildHistory == nil {
		return nil, fmt.Errorf("no rebuild history record found for volume with nexus uuid %s", nexusUuid)
	}
	return nil, err
}

func IsRebuildPartial(nexusUuid string, nexusNodeIp string, childUri string) (bool, error) {
	rebuildHistory, err := mayastorclient.GetRebuildHistory(nexusUuid, nexusNodeIp)
	if err != nil {
		return false, err
	}
	for _, record := range rebuildHistory.GetRecords() {
		if record.GetChildUri() == childUri {
			return record.IsPartial(), err
		}
	}
	return false, fmt.Errorf("failed to find rebuild history for child with uri %s, nexus: %s", childUri, nexusUuid)
}

func IsReplicaPresent(VolUuid string, replicaUri string) (bool, error) {
	var isPresent bool
	replicas, err := k8stest.GetMsvReplicas(VolUuid)
	if err != nil {
		return isPresent, fmt.Errorf("failed to get volume %s replicas, error: %v", VolUuid, err)
	}
	for _, replica := range replicas {
		if replicaUri == replica.Uri {
			isPresent = true
			break
		}
	}
	return isPresent, err
}

func ScheduleLokiOnNode(nodeName string) error {
	err := k8stest.ApplyNodeSelectorToStatefulset(LokiStatefulset, common.NSMayastor(), "kubernetes.io/hostname", nodeName)
	if err != nil {
		return fmt.Errorf("failed to reschedule loki %s on node %s , error: %v", LokiStatefulset, nodeName, err)
	}
	return nil
}

func ScheduleDeploymentOnNode(deployments []string, nodeName string) error {
	// reschedule control plane components to node
	var errs common.ErrorAccumulator
	for _, deploy := range deployments {
		isPresent, _ := k8stest.PodPresentOnNode(deploy, common.NSMayastor(), nodeName)
		if !isPresent {
			err := k8stest.ApplyNodeSelectorToDeployment(deploy, common.NSMayastor(), "kubernetes.io/hostname", nodeName)
			if err != nil {
				errs.Accumulate(err)
			}
		}
	}
	return errs.GetError()
}

func RescheduleCpAndShutdownNode(cpNode string, faultNode string) (string, error) {
	var shutdownNode string

	// reschedule control plane pods
	err := RescheduleControlPlaneToNode(cpNode)
	if err != nil {
		return shutdownNode, err
	}

	// shutdown non nexus replica node
	err = platform.Create().PowerOffNode(faultNode)
	if err != nil {
		return shutdownNode, err
	}
	shutdownNode = faultNode

	// verify node not ready
	nodeReady := VerifyNodeNotReady(faultNode, DefTimeoutSecs)
	if nodeReady {
		return shutdownNode, fmt.Errorf("node %s still online after %d seconds", faultNode, DefTimeoutSecs)
	}

	// verify mayastor ready
	ready, err := k8stest.MayastorReady(2, 360)
	if err != nil {
		return shutdownNode, err
	}
	if !ready {
		return shutdownNode, fmt.Errorf("mayastor installation not ready")
	}

	return shutdownNode, nil
}

func PowerOnNode(nodeName string) (string, error) {
	shutdownNode := nodeName
	// shutdown non nexus replica node
	err := platform.Create().PowerOnNode(nodeName)
	if err != nil {
		return shutdownNode, err
	}
	shutdownNode = ""

	// verify node not ready
	nodeReady := VerifyNodeReady(nodeName, DefTimeoutSecs)
	if !nodeReady {
		return shutdownNode, fmt.Errorf("node %s still not ready after %d seconds", nodeName, DefTimeoutSecs)
	}

	return shutdownNode, nil
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

func VerifyMayastorAndPoolReady() (bool, error) {
	return k8stest.VerifyMayastorAndPoolReady(DefTimeoutSecs)
}

func RebootNode(nodeName string) (string, error) {
	var rebootNode string
	// shutdown non nexus replica node
	err := platform.Create().PowerOffNode(nodeName)
	if err != nil {
		return rebootNode, err
	}
	rebootNode = nodeName

	// verify node not ready
	nodeReady := VerifyNodeNotReady(nodeName, DefTimeoutSecs)
	logf.Log.Info("node not ready", "status", nodeReady)
	if nodeReady {
		return rebootNode, fmt.Errorf("node %s still online after %d seconds", nodeName, DefTimeoutSecs)
	}

	// add sleep twice of control plane reconcile time period
	logf.Log.Info("Sleep", "double of reconcile time period", 2*ReconcileTimeSecs)
	time.Sleep(time.Duration(2*ReconcileTimeSecs) * time.Second)

	logf.Log.Info("node power on")

	// shutdown non nexus replica node
	err = platform.Create().PowerOnNode(nodeName)
	if err != nil {
		return rebootNode, err
	}
	rebootNode = ""

	// verify node not ready
	nodeNotReady := VerifyNodeReady(nodeName, DefTimeoutSecs)
	logf.Log.Info("nodeReady", "status", nodeNotReady)

	if !nodeNotReady {
		return rebootNode, fmt.Errorf("node %s still not ready after %d seconds", nodeName, DefTimeoutSecs)
	}

	return rebootNode, nil
}

func RescheduleControlPlaneToNode(cpNode string) error {
	//reschedule loki sts to nexus node
	if !LokiStatefulsetOnControlNode {
		err := ScheduleLokiOnNode(cpNode)
		if err != nil {
			return err
		}
	}

	//reschedule control plane components to nexus node
	err := ScheduleDeploymentOnNode(MsCpDeployment, cpNode)
	if err != nil {
		return err
	}

	//verify pods on node
	return k8stest.VerifyPodsOnNode(MsAppLabels, cpNode, common.NSMayastor())
}

func IsRebuildTypePartial(nexusUuid string, dstUri string, nexusNodeIp string) (bool, error) {
	logf.Log.Info("Rebuild stats", "nexus uuid", nexusUuid, "dstUri", dstUri)
	rebuildStats, err := mayastorclient.GetRebuildStats(nexusUuid, dstUri, nexusNodeIp)
	if err != nil {
		return false, err
	}
	return rebuildStats.IsRebuildPartial(), err
}

func SetPartialWaitTimeoutArg() error {
	config := e2e_config.GetConfig().Product
	_, retval := k8stest.SetContainerArgs(
		config.ControlPlaneCoreAgent,
		config.AgentCoreContainerName,
		common.NSMayastor(),
		PartialRebuildWaitPeriodArg,
		config.PartialRebuildCpTimeout,
	)
	return retval
}

// As SetPartialWaitTimeoutArg() but returns the original arguments.
func SetPartialWaitTimeoutArgGetOrig() ([]string, error) {
	config := e2e_config.GetConfig().Product
	return k8stest.SetContainerArgs(
		config.ControlPlaneCoreAgent,
		config.AgentCoreContainerName,
		common.NSMayastor(),
		PartialRebuildWaitPeriodArg,
		config.PartialRebuildCpTimeout,
	)
}

// Set all of the arguments for the core agent.
// Typically used for restoring arguments after modification.
func SetCoreAgentArgs(args []string) error {
	config := e2e_config.GetConfig().Product
	return k8stest.SetAllContainerArgs(
		config.ControlPlaneCoreAgent,
		config.AgentCoreContainerName,
		common.NSMayastor(),
		args,
	)
}

// GetReplicasUri returns map of replica uri and corresponding node name
func GetReplicasUri(uuid string) (map[string]string, error) {
	replicas, err := k8stest.GetMsvReplicas(uuid)

	replicasUri := make(map[string]string)
	if err != nil {
		return replicasUri, fmt.Errorf("failed to get mayastor volume %s replicas, error: %v", uuid, err)
	}
	for _, replica := range replicas {
		replicasUri[replica.Uri] = replica.Replica.Node
	}
	return replicasUri, err
}

// GetFullyRebuiltReplicasUri returns array of replica uri which are not present in original set
// of replicas for a volume
func GetFullyRebuiltReplicasUri(uuid string, originalReplicas map[string]string) ([]string, error) {
	logf.Log.Info("Original Replica", "List", originalReplicas)
	replicas, err := k8stest.GetMsvReplicas(uuid)
	fullRebuiltReplicasUri := make([]string, 0)
	if err != nil {
		return fullRebuiltReplicasUri, fmt.Errorf("failed to get mayastor volume %s replicas, error: %v", uuid, err)
	}
	logf.Log.Info("New Replica", "List", replicas)
	for _, replica := range replicas {
		if replica.Uri != "" {
			if _, ok := originalReplicas[replica.Uri]; !ok {
				fullRebuiltReplicasUri = append(fullRebuiltReplicasUri, replica.Uri)
			}
		} else {
			logf.Log.Info("Replica URI is empty", "replica node", replica.Replica.Node)
		}

	}
	logf.Log.Info("Fully Rebuilt Replica", "List", fullRebuiltReplicasUri)
	return fullRebuiltReplicasUri, err
}

// VerifyNexusExists return true if nexus node and uuid matches with initial nexus details
// otherwise it return false
func VerifyNexusExists(volumeUuid string, nexusUuid string, nexusNode string) (bool, error) {
	logf.Log.Info("Verify Nexus", "nexus uuid", nexusUuid, "nexus node", nexusNode)
	msv, err := k8stest.GetMSV(volumeUuid)
	if err != nil {
		return false, fmt.Errorf("failed to get mayastor volume %s, error: %v", volumeUuid, err)
	}
	nexNode := msv.State.Target.Node
	nexUuid := msv.State.Target.Uuid
	logf.Log.Info("Nexus", "nexus uuid", nexUuid, "nexus node", nexNode)
	if nexNode == "" || nexUuid == "" {
		logf.Log.Info("Nexus not found", "volume", volumeUuid)
		return false, nil
	} else if nexNode != nexusNode || nexUuid != nexusUuid {
		logf.Log.Info("New nexus not found", "volume", volumeUuid)
		return false, nil
	}
	return true, nil
}
