package k8stest

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
	"github.com/openebs/openebs-e2e/common/mayastorclient"

	coreV1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs                         = 90 // in seconds
	DefVerifyReplicaWaitHealthyVolTimeSecs = 120
)

type dfaStatus struct {
	suffix     string
	fioTargets []string
	// volume was successfully created and volume is online
	createdVolume bool
	// pvc was created - volume may not be online
	createdPVC     bool
	importedVolume bool
	volUuid        string
	scName         string
	podName        string
	volName        string
	replicaCount   int
	msv            *common.MayastorVolume
	sessionId      string
	monitor        *common.E2eFioPodOutputMonitor
}

type FioApp struct {
	Decor      string
	DeployName string
	VolSizeMb  int
	VolType    common.VolumeType
	FsType     common.FileSystemType
	// FsPercent -> controls size of file allocated on FS
	// 0 -> default (lessby N blocks)
	// > 0 < 100 percentage of available blocks used
	FsPercent                           uint
	ReplicaCount                        int
	Runtime                             uint
	Loops                               int
	VerifyReplicasOnDelete              bool
	AddFioArgs                          []string
	StatusInterval                      int
	OutputFormat                        string
	AppNodeName                         string
	ZeroFill                            bool
	WipeReplicas                        bool
	ThinProvisioned                     bool
	VolWaitForFirstConsumer             bool
	Liveness                            bool
	BlockSize                           uint
	SnapshotName                        string
	CloneFsIdAsVolumeId                 common.CloneFsIdAsVolumeIdType
	FioDebug                            string
	SaveFioPodLog                       bool
	PostOpSleep                         uint
	MountOptions                        []string
	VerifyReplicaWaitHealthyVolTimeSecs uint
	SkipRestoreVolumeVerification       bool
	AllowVolumeExpansion                common.AllowVolumeExpansion
	NodeHasTopologyKey                  string
	NodeSpreadTopologyKey               string
	NodeAffinityTopologyLabel           map[string]string
	MaxSnapshots                        int
	PoolHasTopologyKey                  string
	PoolAffinityTopologyLabel           map[string]string
	status                              dfaStatus
}

func (dfa *FioApp) DeployApp() error {
	return dfa.DeployFio(common.DefaultFioArgs, "")
}

func (dfa *FioApp) DeployAppWithArgs(fioArgsSet common.FioAppArgsSet) error {
	return dfa.DeployFio(fioArgsSet, "")
}

func (dfa *FioApp) validate() error {
	if dfa.WipeReplicas && dfa.VolWaitForFirstConsumer {
		return fmt.Errorf("WipeReplicas is incompatible with VolWaitForConsumer")
	}
	return nil
}

func (dfa *FioApp) DeployFio(fioArgsSet common.FioAppArgsSet, podPrefix string) error {
	var err error
	err = dfa.validate()
	if err != nil {
		return err
	}
	if dfa.status.podName != "" {
		return fmt.Errorf("previous pod not deleted %s", dfa.status.podName)
	}

	// list disk pool in cluster
	if dfa.Runtime != 0 && dfa.Loops != 0 {
		return fmt.Errorf("cannot specify both Runtime and Loops")
	}

	if fioArgsSet == common.CustomFioArgs {
		return fmt.Errorf("custom fio args is not supported")
	}

	err = dfa.CreateVolume()
	if err != nil {
		return err
	}

	// dfa.status.suffix will have been set by dfa.CreateVolume
	decoration := strings.ToLower(dfa.Decor) + dfa.status.suffix
	if dfa.DeployName == "" {
		dfa.status.podName = podPrefix + decoration
	}
	

	efab := common.NewE2eFioArgsBuilder().WithArgumentSet(fioArgsSet).WithZeroFill(dfa.ZeroFill)
	if dfa.Liveness {
		efab = efab.WithLiveness()
	}
	if dfa.BlockSize != 0 {
		efab = efab.WithBlockSize(dfa.BlockSize)
	}

	if dfa.status.fioTargets == nil {
		if dfa.VolType == common.VolFileSystem {
			if dfa.FsPercent == 0 {
				efab = efab.WithDefaultFile()
			} else {
				if dfa.FsPercent > 100 {
					return fmt.Errorf("invalid FsPercent value, valid range is 1 - 100")
				}
				efab = efab.WithDefaultFileExt(common.FioFsAllocPercentage, dfa.FsPercent)
			}
		} else {
			efab = efab.WithDefaultRawBlock()
		}
	} else {
		efab = efab.WithTargets(dfa.status.fioTargets)
	}

	if dfa.StatusInterval != 0 {
		efab = efab.WithAdditionalArg(fmt.Sprintf("--status-interval=%d", dfa.StatusInterval))
	}
	if dfa.OutputFormat != "" {
		efab = efab.WithAdditionalArg(fmt.Sprintf("--output-format=%s", dfa.OutputFormat))
	}
	if dfa.Runtime != 0 {
		// time based loop "forever" timed SIGTERM with terminate
		efab = efab.WithRuntime(int(dfa.Runtime))
	}

	if dfa.Loops != 0 {
		efab = efab.WithAdditionalArg(fmt.Sprintf("--loops=%d", dfa.Loops))
	}
	efab = efab.WithAdditionalArgs(dfa.AddFioArgs)

	if dfa.FioDebug != "" {
		efab = efab.WithAdditionalArg(fmt.Sprintf("--debug=%s", dfa.FioDebug))
	}

	dfa.status.fioTargets = efab.GetTargets()

	podArgs, efabErr := efab.Build()
	if efabErr != nil {
		return fmt.Errorf("failed to compile fio commandline %v", efabErr)
	}
	dfa.status.sessionId = string(uuid.NewUUID())
	podArgs = append([]string{"sessionId", dfa.status.sessionId, ";"}, podArgs...)
	if dfa.PostOpSleep != 0 {
		podArgs = append([]string{"postopsleep", fmt.Sprintf("%d", dfa.PostOpSleep), ";"}, podArgs...)
	}
	logf.Log.Info("e2e-fio", "arguments", strings.Join(podArgs, " "))

	if dfa.DeployName == "" {

		// fio pod container
		container := MakeFioContainer(dfa.status.podName, podArgs)
		//	container.ImagePullPolicy = coreV1.PullAlways
		// volume claim details
		volume := coreV1.Volume{
			Name: "ms-volume",
			VolumeSource: coreV1.VolumeSource{
				PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
					ClaimName: dfa.status.volName,
				},
			},
		}
		// create the fio pod
		pod := NewPodBuilder("fio").
			WithName(dfa.status.podName).
			WithNamespace(common.NSDefault).
			WithRestartPolicy(coreV1.RestartPolicyNever).
			WithContainer(container).
			WithVolume(volume).
			WithVolumeDeviceOrMount(dfa.VolType)
		//		WithHostPath("tmp", "/tmp")

		if dfa.AppNodeName != "" {
			pod = pod.WithNodeName(dfa.AppNodeName)
		}
		podObj, err := pod.Build()
		if err != nil {
			return fmt.Errorf("generating fio pod definition %s, %v", dfa.status.podName, err)
		}
		if podObj == nil {
			return fmt.Errorf("failed to generate fio pod definition")
		}
		_, err = CreatePod(podObj, common.NSDefault)
		if err != nil {
			return fmt.Errorf("creating fio pod %s, %v", dfa.status.podName, err)
		}
		// wait for pod to transition to running or complete whichever is first
		var phase coreV1.PodPhase
		var podLogSynopsis *common.E2eFioPodLogSynopsis
		for secs := 0; secs < DefTimeoutSecs; secs++ {
			phase, podLogSynopsis, err = CheckFioPodCompleted(dfa.status.podName, common.NSDefault)
			if err != nil {
				return err
			}
			switch phase {
			case coreV1.PodSucceeded:
				return nil
			case coreV1.PodRunning:
				dfa.status.msv, _ = GetMSV(dfa.status.volUuid)
				logf.Log.Info("PodRunning", "msv", dfa.status.msv)
				return nil
			case coreV1.PodFailed:
				return fmt.Errorf("pod state is %v, %s", phase, podLogSynopsis)
			}
			time.Sleep(1 * time.Second)
		}

		return fmt.Errorf("pod state is %v, %s", phase, podLogSynopsis)
		
	} else {
		labelKey := "e2e-test"
		labelValue := "fio"
		labelselector := map[string]string{
			labelKey: labelValue,
		}
		deployment, err := NewDeploymentBuilder().
			WithName(dfa.DeployName).
            		WithNamespace(common.NSDefault).
            		WithLabelsNew(labelselector).
            		WithSelectorMatchLabelsNew(labelselector).
            		WithPodTemplateSpecBuilder(
                		NewPodtemplatespecBuilder().
                    		WithLabels(labelselector).
				WithRestartPolicy(coreV1.RestartPolicyNever).
                    		WithContainerBuildersNew(
                        		NewContainerBuilder().
                            		WithName(dfa.DeployName).
                            		WithImage(common.GetFioImage()).
                            		WithVolumeDeviceOrMount("ms-volume", dfa.VolType).
                            		WithImagePullPolicy(coreV1.PullAlways).
                            		WithArgumentsNew(podArgs)).
                    		WithVolumeBuilders(
                        		NewVolumeBuilder().
                            		WithName("ms-volume").
                            		WithPVCSource(dfa.status.volName),
                    		),
            		).Build()
		if err != nil {
            return fmt.Errorf("failed to create deployment %s definition object in %s namesppace", dfa.DeployName, common.NSDefault)
        	}
		
		if dfa.AppNodeName != "" {
			deployment.Spec.Template.Spec.NodeName = dfa.AppNodeName
		}

		// Create the deployment
		err = CreateDeployment(deployment)
		if err != nil {
			return fmt.Errorf("creating fio deployment %s, %v", dfa.DeployName, err)
		}
		var running bool
		// Wait for pods to be running
		for i := 0; i < DefTimeoutSecs; i++ {
			running, err = VerifyDeploymentReadyReplicaCount(dfa.DeployName, common.NSDefault, 1)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
		if running {
			dfa.status.msv, err = GetMSV(dfa.status.volUuid)
			if err != nil {
				return err
			}
			logf.Log.Info("DeploymentRunning", "msv", dfa.status.msv)
			pods, err := GetDeploymentPods(dfa.DeployName, common.NSDefault)
			if err != nil {
				return err
			}
			dfa.status.podName = pods.Items[0].Name
			// wait for pod to transition to running or complete whichever is first
			var phase coreV1.PodPhase
			var podLogSynopsis *common.E2eFioPodLogSynopsis
			for secs := 0; secs < DefTimeoutSecs; secs++ {
				phase, podLogSynopsis, err = CheckFioPodCompleted(dfa.status.podName, common.NSDefault)
				if err != nil {
					return err
				}
				switch phase {
				case coreV1.PodSucceeded:
					return nil
				case coreV1.PodRunning:
					dfa.status.msv, _ = GetMSV(dfa.status.volUuid)
					logf.Log.Info("PodRunning", "msv", dfa.status.msv)
					return nil
				case coreV1.PodFailed:
					return fmt.Errorf("pod state is %v, %s", phase, podLogSynopsis)
				}
				time.Sleep(1 * time.Second)
			}
		}
		return fmt.Errorf("deployment did not reach running state in time")
	}
}

func (dfa *FioApp) CreateVolume() error {
	var err error
	var poolsInCluster []common.MayastorPool
	const sleepTime = 3

	err = dfa.validate()
	if err != nil {
		return err
	}

	if dfa.status.createdPVC || dfa.status.importedVolume {
		return nil
	}

	for ix := 0; ix < (DefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
		poolsInCluster, err = ListMsPools()
		if err != nil {
			logf.Log.Info("ListMsPools", "Error", err)
			time.Sleep(sleepTime * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed to list disk pools, error %v", err)
	}

	if len(poolsInCluster) == 0 {
		return fmt.Errorf("no pools found in cluster")
	}
	if dfa.ReplicaCount == 0 {
		dfa.status.replicaCount = len(poolsInCluster)
	} else {
		dfa.status.replicaCount = dfa.ReplicaCount
	}

	if dfa.VolType.String() == "" {
		dfa.VolType = common.VolRawBlock
	}

	decoration := fmt.Sprintf("-%dr", dfa.status.replicaCount)
	if dfa.VolType == common.VolFileSystem {
		if dfa.FsType == common.NoneFsType {
			logf.Log.Info("default filesystem type to ext4")
			dfa.FsType = common.Ext4FsType
		}
		decoration += "-" + string(dfa.FsType)
	}
	if dfa.VolType == common.VolRawBlock {
		decoration += "-rb"
	}
	if dfa.ThinProvisioned {
		decoration += "-thin"
	} else {
		decoration += "-thick"
	}
	if dfa.VolWaitForFirstConsumer {
		decoration += "-bindlate"
	} else {
		decoration += "-bindimm"
	}
	dfa.status.suffix = decoration
	decoration = strings.ToLower(dfa.Decor) + decoration
	dfa.status.volName = decoration
	dfa.status.scName = decoration

	provisioning := common.ThickProvisioning
	volBindingMode := storageV1.VolumeBindingImmediate
	if dfa.ThinProvisioned {
		provisioning = common.ThinProvisioning
	}
	if dfa.VolWaitForFirstConsumer {
		volBindingMode = storageV1.VolumeBindingWaitForFirstConsumer
	}
	if dfa.FsType == common.BtrfsFsType {
		dfa.FsPercent = 95
		dfa.MountOptions = []string{"nodatacow"}
	}
	scBuilder := NewScBuilder().
		WithName(dfa.status.scName).
		WithReplicas(dfa.status.replicaCount).
		WithProtocol(common.ShareProtoNvmf).
		WithNamespace(common.NSDefault).
		WithVolumeBindingMode(volBindingMode).
		WithProvisioningType(provisioning).
		WithMountOptions(dfa.MountOptions).
		WithVolumeExpansion(dfa.AllowVolumeExpansion)

	if dfa.VolType == common.VolFileSystem {
		scBuilder = scBuilder.
			WithFileSystemType(dfa.FsType).
			WithCloneFsIdAsVolumeId(dfa.CloneFsIdAsVolumeId)
	}

	if dfa.NodeAffinityTopologyLabel != nil {
		scBuilder = scBuilder.WithNodeAffinityTopologyLabel(dfa.NodeAffinityTopologyLabel)
	}
	if dfa.NodeSpreadTopologyKey != "" {
		scBuilder = scBuilder.WithNodeSpreadTopologyKey(dfa.NodeSpreadTopologyKey)
	}
	if dfa.NodeHasTopologyKey != "" {
		scBuilder = scBuilder.WithNodeHasTopologyKey(dfa.NodeHasTopologyKey)
	}
	if dfa.MaxSnapshots != 0 {
		scBuilder = scBuilder.WithMaxSnapshots(dfa.MaxSnapshots)
	}
	if dfa.PoolHasTopologyKey != "" {
		scBuilder = scBuilder.WithPoolHasTopologyKey(dfa.PoolHasTopologyKey)
	}
	if dfa.PoolAffinityTopologyLabel != nil {
		scBuilder = scBuilder.WithPoolAffinityTopologyLabel(dfa.PoolAffinityTopologyLabel)
	}
	err = scBuilder.BuildAndCreate()
	if err != nil {
		return fmt.Errorf("failed to create storage class %s %v", dfa.status.scName, err)
	}
	// Create the volume
	if dfa.SnapshotName != "" {
		dfa.status.volUuid, err = MkRestorePVC(dfa.VolSizeMb, dfa.status.volName, dfa.status.scName, common.NSDefault, dfa.VolType, dfa.SnapshotName, dfa.SkipRestoreVolumeVerification)
		dfa.status.createdPVC = dfa.status.volUuid != ""
		if err != nil {
			return fmt.Errorf("failed to create pvc %s from snapshot source %s, %v", dfa.status.volName, dfa.SnapshotName, err)
		} else if dfa.SkipRestoreVolumeVerification {
			return err
		}
	} else {
		dfa.status.volUuid, err = MkPVC(dfa.VolSizeMb, dfa.status.volName, dfa.status.scName, dfa.VolType, common.NSDefault)
		dfa.status.createdPVC = dfa.status.volUuid != ""
		if err != nil {
			return fmt.Errorf("failed to create pvc %s, %v", dfa.status.volName, err)
		}
	}
	dfa.status.createdVolume = true
	logf.Log.Info("Volume", "uid", dfa.status.volUuid)
	if dfa.WipeReplicas {
		err = WipeVolumeReplicas(dfa.status.volUuid)
		if err != nil {
			return fmt.Errorf("failed to wipe volume replicas")
		}
	}
	return err
}

func (dfa *FioApp) Cleanup() error {
	var err error
	// delete pod and volume
	//	if dfa.SaveFioPodLog {
	// dfa.DumpPodLog()
	//	}

	podName := dfa.GetPodName()

	// If the pod name is "", we never deployed a pod.
	if podName != "" && dfa.DeployName == "" {
		// Dump fio pod logs only if pod phase is completed otherwise log collection til pod is completed
		podPhase, _, err := CheckFioPodCompleted(dfa.GetPodName(), common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to get fio pod %s phase : %v", dfa.GetPodName(), err)
		}
		if podPhase == coreV1.PodSucceeded || podPhase == coreV1.PodFailed {
			logf.Log.Info("Dump fio pod log", "pod", dfa.GetPodName(), "pod phase", podPhase)
			dfa.DumpPodLog()
		} else {
			logf.Log.Info("Skipping fio pod log collection", "pod", dfa.GetPodName(), "pod phase", podPhase)
		}

		err = dfa.DeletePod()
		if err != nil {
			return fmt.Errorf("failed to delete fio pod %s, err: %v", dfa.GetPodName(), err)
		}
		if dfa.VerifyReplicasOnDelete {
			err = dfa.VerifyReplicas()
			if err != nil {
				return fmt.Errorf("verify replicas before deleting PVC failed: %v", err)
			}
		}
	} else if dfa.DeployName != "" {
		err := DeleteDeployment(dfa.DeployName, common.NSDefault)
		if err != nil {
			return fmt.Errorf("failed to delete fio deployment %s, err: %v", dfa.DeployName, err)
		}
		const sleepTime = 3
		var pods *coreV1.PodList
		var podsDeleted bool
		podsDeleted = false
		for ix := 0; ix < (DefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
			pods, err = GetDeploymentPods(dfa.DeployName, common.NSDefault)
			logf.Log.Info("GetDeploymentPods", "error", err, "pods", pods)
			if err != nil && k8serrors.IsNotFound(err) {
				podsDeleted = true
				break
			}
			time.Sleep(sleepTime * time.Second)
		}
		if !podsDeleted {
			return fmt.Errorf("failed to delete fio deployment pods %s, err: %v", dfa.DeployName, err)
		}
	}
	// Only delete PVC and storage class if they were created by this instance
	if dfa.status.createdPVC {
		err = RmPVC(dfa.status.volName, dfa.status.scName, common.NSDefault)
		if err == nil {
			err = RmStorageClass(dfa.status.scName)
		}
	}

	return err
}

// ForcedCleanup - delete resources associated with the FioApp by name
// regardless of creation status
// This function should be used sparingly.
// FIXME: refactor so that is function can be replaced
// by simply calling Cleanup
func (dfa *FioApp) ForcedCleanup() {
	_ = DeletePod(dfa.status.podName, common.NSDefault)
	dfa.status.podName = ""
	_ = RmPVC(dfa.status.volName, dfa.status.scName, common.NSDefault)
	dfa.status.createdVolume = false
	dfa.status.createdPVC = false
	_ = RmStorageClass(dfa.status.scName)
	dfa.status.scName = ""
}

func (dfa *FioApp) WaitComplete(timeoutSecs int) error {
	msv, err := GetMSV(dfa.status.volUuid)
	if err == nil && msv != nil {
		dfa.status.msv = msv
	}
	logf.Log.Info("WaitComplete", "msv", dfa.status.msv)
	return WaitFioPodComplete(dfa.status.podName, 10, timeoutSecs)
}

func (dfa *FioApp) WaitRunning(timeoutSecs int) bool {
	logf.Log.Info("Wait for pod running", "pod", dfa.status.podName, "timeoutSecs", timeoutSecs)
	return WaitPodRunning(dfa.status.podName, common.NSDefault, timeoutSecs)
}

func (dfa *FioApp) GetPodStatus() (coreV1.PodPhase, error) {
	return GetPodStatus(dfa.status.podName, common.NSDefault)
}

func (dfa *FioApp) DeletePod() error {
	var err error
	if dfa.status.podName != "" {

		err = DeletePod(dfa.status.podName, common.NSDefault)
		if err == nil {
			dfa.status.podName = ""
		}
	}
	return err
}

func (dfa *FioApp) SetVolumeReplicaCount(replicaCount int) error {
	return SetMsvReplicaCount(dfa.status.volUuid, replicaCount)
}

// VolumeComparePreFlightChecks  volume replicas can only be compared meaningfully if the
// volume is not degraded - ie all replicas are online and rebuilds if any have been completed.
// rebuilds will not complete if the volume is unpublished so this function mounts the volume
// on a sleeper fio pio (no io) and waits for the volume to be healthy, within a time period.
func VolumeComparePreFlightChecks(volUuid string, volTyp common.VolumeType, volName string, waitSecs uint) error {
	var err error
	var msv *common.MayastorVolume

	nexus, _ := GetMsvNodes(volUuid)
	for endTime := time.Now().Add(60 * time.Second); nexus != "" && time.Now().Before(endTime); time.Sleep(10 * time.Second) {
		logf.Log.Info("VolumeComparePreFlightChecks: volume has: ", "nexus", nexus)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		return fmt.Errorf("attempting to verify replicas when nexus %v is present", nexus)
	}

	msv, err = GetMSV(volUuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve mayastor volume for k8s volume %s uuid=%s, %v ", volName, volUuid, err)
	}
	if msv.State.Status != controlplane.VolStateHealthy() {
		return fmt.Errorf("mayastor volume for k8s volume %s uuid=%s, Status.State=%s, is not Online", volName, volUuid, msv.State.Status)
	}

	// double check that volume is healthy
	sleeperPodName := fmt.Sprintf("sleeper-%s", string(uuid.NewUUID()))
	err = CreateSleepingFioPod(sleeperPodName, volName, volTyp)
	if err != nil {
		return fmt.Errorf("failed to create sleeping fio pod")
	}
	if !WaitPodRunning(sleeperPodName, common.NSDefault, 60) {
		_ = DeletePod(sleeperPodName, common.NSDefault)
		return fmt.Errorf("sleeping fio pod failed to start")
	}

	if waitSecs == 0 {
		waitSecs = uint(DefVerifyReplicaWaitHealthyVolTimeSecs)
	}

	for start := time.Now(); time.Since(start) < time.Duration(waitSecs)*time.Second; time.Sleep(10 * time.Second) {
		msv, err = GetMSV(volUuid)
		if err == nil {
			if msv.State.Status == controlplane.VolStateHealthy() {
				break
			} else {
				logf.Log.Info("VolumeComparePreFlightChecks:", "msv.State.Status", msv.State.Status)
			}
		} else {
			logf.Log.Info("VolumeComparePreFlightChecks: failed to retrieve msv", "error", err)
		}
	}
	logf.Log.Info("VolumeComparePreFlightChecks:", "msv.State.Status", msv.State.Status)
	if msv.State.Status != controlplane.VolStateHealthy() {
		return fmt.Errorf("volume is not healthy aborting replica comparison msv.State.Status is %s", msv.State.Status)
	}
	err = DeletePod(sleeperPodName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to delete sleeping fio pod")
	}
	nexus, _ = GetMsvNodes(volUuid)
	for start := time.Now(); nexus != "" && time.Since(start) < time.Second*120; time.Sleep(10 * time.Second) {
		logf.Log.Info("VolumeComparePreFlightChecks: volume has: ", "nexus", nexus)
		nexus, _ = GetMsvNodes(volUuid)
	}
	if nexus != "" {
		return fmt.Errorf("attempting to verify replicas when nexus %v is present", nexus)
	}

	// Wait for upto a minute for all replicas to be online
	var replicaCheckErr error
	for endTime := time.Now().Add(60 * time.Second); time.Now().Before(endTime); time.Sleep(10 * time.Second) {
		var replicaTopology common.ReplicaTopology
		replicaTopology, replicaCheckErr = GetMsvReplicaTopology(volUuid)
		if replicaCheckErr != nil {
			replicaCheckErr = fmt.Errorf("failed to retrieve replica topology for volume %s uuid=%s, %v ", volName, volUuid, replicaCheckErr)
		} else {
			for ix, repl := range replicaTopology {
				logf.Log.Info(fmt.Sprintf("VolumeComparePreFlightChecks: %s", ix), "replica", repl)
				if repl.State != controlplane.ReplicaStateOnline() {
					replicaCheckErr = fmt.Errorf("replica State not Online %v", repl)
				}
			}
		}
		if replicaCheckErr == nil {
			break
		}
	}
	if replicaCheckErr != nil {
		return replicaCheckErr
	}
	return nil
}

// VerifyReplicas verify replicas for volumes used by the FioApp instance
func (dfa *FioApp) VerifyReplicas() error {
	err := VolumeComparePreFlightChecks(dfa.status.volUuid, dfa.VolType, dfa.status.volName, dfa.VerifyReplicaWaitHealthyVolTimeSecs)
	if err != nil {
		return err
	}

	res := CompareVolumeReplicas(dfa.status.volName, common.NSDefault)
	switch res.Result {
	case common.CmpReplicasMismatch:
		_, cmp, _ := ByteCompareVolumeReplicas(dfa.status.volName, common.NSDefault)
		res.Description += "\n" + cmp
		logf.Log.Info("VerifyReplicas", "msv", dfa.status.msv)
		if dfa.status.msv != nil {
			logf.Log.Info("verifyReplicas", "nexusNode", dfa.status.msv.State.Target.Node, "replicas", dfa.status.msv.State.ReplicaTopology)
			res.Description += "\n" + fmt.Sprintf("last recorded MSV: nexusNode=%s, replicas=%v ",
				dfa.status.msv.State.Target.Node, dfa.status.msv.State.ReplicaTopology)
		}
		return fmt.Errorf("replica verification failed:\n%s", res.Description)
	case common.CmpReplicasFailed:
		return res.Err
	case common.CmpReplicasMatch:
		return nil
	default:
		return fmt.Errorf("unexpected replica comparison result value %d", res.Result)
	}
}

func (dfa *FioApp) DumpPodLog() {
	DumpPodLog(dfa.status.podName, common.NSDefault)
}

func (dfa *FioApp) DumpReplicas() {
	msv, err := GetMSV(dfa.status.volUuid)
	if err == nil {
		for ix, replica := range msv.State.ReplicaTopology {
			logf.Log.Info(fmt.Sprintf("Replica %s", ix), "node", replica.Node, "uri", ix, "state", replica.State)
		}
	} else {
		logf.Log.Info("failed to retrieve msv", "error", err)
	}
}

func wipeReplica(replica common.Replica, replicaUuid string, ch chan error, wg *sync.WaitGroup) {
	var err error
	var address *string

	address, err = GetNodeIPAddress(replica.Node)
	if err == nil {
		err = mayastorclient.WipeReplica(*address, replicaUuid, replica.Pool)
	}
	ch <- err

	wg.Done()
}

func (dfa *FioApp) SetAppNodeName(nodeName string) error {
	if nodeName == "" {
		return fmt.Errorf("nodeName cannot be empty")
	}
	dfa.AppNodeName = nodeName
	logf.Log.Info("AppNodename is ", "AppNodename", dfa.AppNodeName)
	return nil
}

func (dfa *FioApp) GetVolUuid() string {
	return dfa.status.volUuid
}

func (dfa *FioApp) GetAppUuid() string {
	return dfa.status.sessionId
}

func (dfa *FioApp) GetScName() string {
	return dfa.status.scName
}

func (dfa *FioApp) GetPodName() string {
	if dfa.DeployName == "" {
		return dfa.status.podName
	}
	pods, _ := GetDeploymentPods(dfa.DeployName, common.NSDefault)
	return pods.Items[0].Name
}

func (dfa *FioApp) GetVolName() string {
	return dfa.status.volName
}

func (dfa *FioApp) GetReplicaCount() int {
	return dfa.status.replicaCount
}

// PVC was created and PVC is accessbile
func (dfa *FioApp) IsVolumeCreated() bool {
	return dfa.status.createdVolume
}

// PVC was created and PVC may not be accessbile
func (dfa *FioApp) IsPVCCreated() bool {
	return dfa.status.createdPVC
}

func (dfa *FioApp) ScaleVolumeReplicas(v int) error {
	logf.Log.Info("Scaling volume replica count",
		"volname", dfa.status.volName,
		"vol uuid", dfa.status.volUuid,
		"from", dfa.status.replicaCount,
		"to", dfa.status.replicaCount+v,
	)
	err := SetMsvReplicaCount(dfa.status.volUuid, dfa.status.replicaCount+v)
	if err == nil {
		dfa.status.replicaCount += v
	}
	return err
}

func (dfa *FioApp) RefreshVolumeState() error {
	pvc, err := GetPVC(dfa.status.volName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to get pvc %s, error: %v", dfa.status.volName, err)
	} else if pvc == nil {
		return fmt.Errorf("pvc %s not found", dfa.status.volName)
	} else if *pvc.Spec.StorageClassName != dfa.status.scName {
		return fmt.Errorf("storage class %s not used to create pvc %s", dfa.status.scName, dfa.status.volName)
	} else if pvc.ObjectMeta.UID == "" {
		return fmt.Errorf("pvc %s does not have  pvc.ObjectMeta.UID non empty string", dfa.status.volName)
	}
	err = VerifyMayastorPvcIsUsable(pvc)
	if err == nil {
		// set dfa.status.volUuid
		dfa.status.volUuid = string(pvc.ObjectMeta.UID)
		//set dfa.status.createdPVC true
		dfa.status.createdPVC = true
		//set dfa.status.createdVolume true
		dfa.status.createdVolume = true
	} else {
		logf.Log.Info("RefreshVolumeState: Volume is not usable", "error", err)
	}
	return err
}

// MonitorPod - invokes MonitorE2EFioPod which launches a go thread
// to stream fio pod log output and scan that stream
// to populate fields in E2eFioPodOutputMonitor
func (dfa *FioApp) MonitorPod() (*common.E2eFioPodOutputMonitor, error) {
	var err error
	if dfa.status.monitor == nil {
		dfa.status.monitor, err = MonitorE2EFioPod(dfa.GetPodName(), common.NSDefault)
		if err != nil {
			dfa.status.monitor = nil
		}
	}
	return dfa.status.monitor, err
}

func (dfa *FioApp) WaitFioComplete(timeoutSecs int, pollTimeSecs int) (int, error) {
	mon, err := dfa.MonitorPod()
	if err != nil {
		return 0, err
	}

	timeout := time.Duration(int(time.Second) * timeoutSecs)
	sleeptime := time.Duration(int(time.Second) * pollTimeSecs)
	for endTime := time.Now().Add(timeout); time.Now().Before(endTime); time.Sleep(sleeptime) {
		switch len(mon.Synopsis.JsonRecords.ExitValues) {
		case 0:
			// no exit values found - fio is still running
			break
		case 1:
			// single exit value found - fio has completed
			logf.Log.Info("fio", "elapsed", *mon.Synopsis.JsonRecords.ExitValues[0].ElapsedSecs)
			return *mon.Synopsis.JsonRecords.ExitValues[0].ExitValue, nil
		default:
			// bug - in e2e-fio or monitoring code multiple exit values recorded for single fio run
			return *mon.Synopsis.JsonRecords.ExitValues[0].ExitValue, fmt.Errorf("multiple exit values found")
		}
	}

	return 0, fmt.Errorf("timed out waiting for fio completion")
}

func (dfa *FioApp) FioTargetSizes() (map[string]uint64, error) {
	mon, err := dfa.MonitorPod()
	if err != nil {
		return nil, err
	}
	sizes := make(map[string]uint64)
	for _, ftSize := range mon.Synopsis.JsonRecords.TargetSizes {
		sizes[*ftSize.Path] = *ftSize.Size
	}
	return sizes, err
}

func (dfa *FioApp) ImportVolume(volName string) error {
	pvc, err := GetPVC(volName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to get pvc %s, error: %v", volName, err)
	}
	if pvc == nil {
		return fmt.Errorf("pvc %s not found", volName)
	}

	dfa.status.volName = volName
	// set dfa.status.volUuid
	dfa.status.volUuid = string(pvc.ObjectMeta.UID)

	dfa.status.importedVolume = true
	dfa.status.suffix = "-imported-vol"

	return err
}

func (dfa *FioApp) ImportVolumeFromApp(srcDfa *FioApp) error {
	dfa.VolSizeMb = srcDfa.VolSizeMb
	dfa.VolType = srcDfa.VolType
	dfa.ReplicaCount = srcDfa.ReplicaCount
	return dfa.ImportVolume(srcDfa.GetVolName())
}
