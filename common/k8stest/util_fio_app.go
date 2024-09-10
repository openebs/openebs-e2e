package k8stest

import (
	"fmt"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type FioApplicationStatus struct {
	suffix         string
	fioTargets     []string
	pvcName        string
	scName         string
	fioPodName     string
	sessionId      string
	createdPVC     bool
	monitor        *common.E2eFioPodOutputMonitor
	importedVolume bool
}

type FioApplication struct {
	Decor         string
	VolSizeMb     int
	VolType       common.VolumeType
	FsType        common.FileSystemType
	OpenEbsEngine common.OpenEbsEngine
	// FsPercent -> controls size of file allocated on FS
	// 0 -> default (lessby N blocks)
	// > 0 < 100 percentage of available blocks used
	FsPercent                      uint
	Runtime                        uint
	Loops                          int
	AddFioArgs                     []string
	StatusInterval                 int
	OutputFormat                   string
	AppNodeName                    string
	VolWaitForFirstConsumer        bool
	ScMountOptions                 []string
	ScReclaimPolicy                v1.PersistentVolumeReclaimPolicy
	Liveness                       bool
	BlockSize                      uint
	FioDebug                       string
	SaveFioPodLog                  bool
	PostOpSleep                    uint
	AllowVolumeExpansion           common.AllowVolumeExpansion
	Lvm                            LvmOptions
	Zfs                            ZfsOptions
	HostPath                       HostPathOptions
	status                         FioApplicationStatus
	SkipPvcVerificationAfterCreate bool
}

type LvmOptions struct {
	AllowedTopologies []v1.TopologySelectorTerm
	Shared            common.YesNoVal
	VgPattern         string
	Storage           string
	VolGroup          string
	ThinProvision     common.YesNoVal
}

type ZfsOptions struct {
	AllowedTopologies []v1.TopologySelectorTerm
	RecordSize        string
	Compression       common.OnOffVal
	DedUp             common.OnOffVal
	PoolName          string
	ThinProvision     common.YesNoVal
	VolBlockSize      string
	Shared            common.YesNoVal
}

type HostPathOptions struct {
	Annotations map[string]string
}

func (dfa *FioApplication) DeployApplication() error {
	return dfa.DeployFio(common.DefaultFioArgs, "")
}

func (dfa *FioApplication) DeployAppWithArgs(fioArgsSet common.FioAppArgsSet) error {
	return dfa.DeployFio(fioArgsSet, "")
}

func (dfa *FioApplication) DeployFio(fioArgsSet common.FioAppArgsSet, podPrefix string) error {
	var err error
	if dfa.status.fioPodName != "" {
		return fmt.Errorf("previous pod not deleted %s", dfa.status.fioPodName)
	}

	// list disk pool in cluster
	if dfa.Runtime != 0 && dfa.Loops != 0 {
		return fmt.Errorf("cannot specify both Runtime and Loops")
	}

	if fioArgsSet == common.CustomFioArgs {
		return fmt.Errorf("custom fio args is not supported")
	}

	if dfa.OpenEbsEngine.String() == "" {
		return fmt.Errorf("openebs engine not specified")
	}

	err = dfa.CreateVolume()
	if err != nil {
		return err
	}

	// dfa.status.suffix will have been set by dfa.CreateVolume
	decoration := strings.ToLower(dfa.Decor) + dfa.status.suffix
	dfa.status.fioPodName = podPrefix + decoration

	efab := common.NewE2eFioArgsBuilder().WithArgumentSet(fioArgsSet)
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

	// fio pod container
	container := MakeFioContainer(dfa.status.fioPodName, podArgs)
	//	container.ImagePullPolicy = coreV1.PullAlways
	// volume claim details
	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: dfa.status.pvcName,
			},
		},
	}
	// create the fio pod
	pod := NewPodBuilder("fio").
		WithName(dfa.status.fioPodName).
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
		return fmt.Errorf("generating fio pod definition %s, %v", dfa.status.fioPodName, err)
	}
	if podObj == nil {
		return fmt.Errorf("failed to generate fio pod definition")
	}
	_, err = CreatePod(podObj, common.NSDefault)
	if err != nil {
		return fmt.Errorf("creating fio pod %s, %v", dfa.status.fioPodName, err)
	}
	// wait for pod to transition to running or complete whichever is first
	var phase coreV1.PodPhase
	var podLogSynopsis *common.E2eFioPodLogSynopsis
	for secs := 0; secs < DefTimeoutSecs; secs++ {
		phase, podLogSynopsis, err = CheckFioPodCompleted(dfa.status.fioPodName, common.NSDefault)
		if err != nil {
			return err
		}
		switch phase {
		case coreV1.PodSucceeded:
			return nil
		case coreV1.PodRunning:
			return nil
		case coreV1.PodFailed:
			return fmt.Errorf("pod state is %v, %s", phase, podLogSynopsis)
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("pod state is %v, %s", phase, podLogSynopsis)
}

func (dfa *FioApplication) CreateVolume() error {
	var err error

	decoration := dfa.OpenEbsEngine.String()

	if dfa.VolType.String() == "" {
		dfa.VolType = common.VolRawBlock
	}

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

	if dfa.VolWaitForFirstConsumer {
		decoration += "-bindlate"
	} else {
		decoration += "-bindimm"
	}
	dfa.status.suffix = decoration + dfa.OpenEbsEngine.String()
	decoration = strings.ToLower(dfa.Decor) + decoration
	dfa.status.pvcName = decoration
	dfa.status.scName = decoration

	err = dfa.CreateSc()
	if err != nil {
		return fmt.Errorf("failed to create sc %s, %v", dfa.status.scName, err)
	}
	var localEngine bool
	if dfa.OpenEbsEngine != common.Mayastor {
		localEngine = true
	}

	// Create the volume
	_, err = MakePVC(dfa.VolSizeMb, dfa.status.pvcName, dfa.status.scName, dfa.VolType, common.NSDefault, localEngine, dfa.SkipPvcVerificationAfterCreate)

	if err != nil {
		return fmt.Errorf("failed to create pvc %s, %v", dfa.status.pvcName, err)
	}
	dfa.status.createdPVC = true
	return err
}

func (dfa *FioApplication) CreateSc() error {
	volBindingMode := storageV1.VolumeBindingImmediate

	if dfa.VolWaitForFirstConsumer {
		volBindingMode = storageV1.VolumeBindingWaitForFirstConsumer
	}
	var provisioner string
	productConfig := e2e_config.GetConfig().Product
	switch dfa.OpenEbsEngine {
	case common.Zfs:
		provisioner = productConfig.ZfsEngineProvisioner
	case common.Lvm:
		provisioner = productConfig.LvmEngineProvisioner
	case common.Hostpath:
		provisioner = productConfig.HostPathEngineProvisioner
	}

	if provisioner == "" {
		return fmt.Errorf("storage class provisioner not found for engine %s", dfa.OpenEbsEngine)
	}

	scBuilder := NewStorageClassBuilder().
		WithName(dfa.status.scName).
		WithNamespace(common.NSDefault).
		WithVolumeBindingMode(volBindingMode).
		WithVolumeExpansion(dfa.AllowVolumeExpansion).
		WithProvisioner(provisioner)

	if dfa.VolType == common.VolFileSystem {
		scBuilder = scBuilder.
			WithFileSystemType(dfa.FsType)
	}
	if len(dfa.ScMountOptions) != 0 {
		scBuilder = scBuilder.WithMountOptions(dfa.ScMountOptions)
	}
	if dfa.ScReclaimPolicy != "" {
		scBuilder = scBuilder.WithReclaimPolicy(dfa.ScReclaimPolicy)
	}

	if dfa.OpenEbsEngine == common.Lvm {
		if dfa.Lvm.Shared.String() != "" {
			scBuilder = scBuilder.WithLvmShared(dfa.Lvm.Shared.String())
		}
		if dfa.Lvm.Storage != "" {
			scBuilder = scBuilder.WithLvmStorage(dfa.Lvm.Storage)
		}
		if dfa.Lvm.ThinProvision.String() != "" {
			scBuilder = scBuilder.WithLvmThinVol(dfa.Lvm.ThinProvision.String())
		}
		if dfa.Lvm.VgPattern != "" {
			scBuilder = scBuilder.WithLvmVgPattern(dfa.Lvm.VgPattern)
		}
		if dfa.Lvm.VolGroup != "" {
			scBuilder = scBuilder.WithLvmVolGroup(dfa.Lvm.VolGroup)
		}
		if dfa.Lvm.AllowedTopologies != nil {
			scBuilder = scBuilder.WithAllowedTopologies(dfa.Lvm.AllowedTopologies)
		}
	} else if dfa.OpenEbsEngine == common.Zfs {

		if dfa.Zfs.RecordSize != "" {
			scBuilder = scBuilder.WithZfsRecordSize(dfa.Zfs.RecordSize)
		}
		if dfa.Zfs.AllowedTopologies != nil {
			scBuilder = scBuilder.WithAllowedTopologies(dfa.Zfs.AllowedTopologies)
		}
		if dfa.Zfs.PoolName != "" {
			scBuilder = scBuilder.WithZfsPoolName(dfa.Zfs.PoolName)
		} else {
			return fmt.Errorf("zfs pool name not specified")
		}
		if dfa.Zfs.Shared.String() != "" {
			scBuilder = scBuilder.WithZfsShared(dfa.Zfs.Shared.String())
		}
		if dfa.Zfs.ThinProvision.String() != "" {
			scBuilder = scBuilder.WithZfsThinVol(dfa.Zfs.ThinProvision.String())
		}
		if dfa.Zfs.Compression.String() != "" {
			scBuilder = scBuilder.WithZfsCompression(dfa.Zfs.Compression.String())
		}
		if dfa.Zfs.DedUp.String() != "" {
			scBuilder = scBuilder.WithZfsDeDUp(dfa.Zfs.DedUp.String())
		}

	} else if dfa.OpenEbsEngine == common.Hostpath {
		if dfa.HostPath.Annotations != nil {
			scBuilder = scBuilder.WithAnnotations(dfa.HostPath.Annotations)
		}
	}

	logf.Log.Info("SC", "param", scBuilder)

	err := scBuilder.BuildAndCreate()
	if err != nil {
		return fmt.Errorf("failed to create storage class %s %v", dfa.status.scName, err)
	}
	return nil
}

func (dfa *FioApplication) DumpPodLog() {
	DumpPodLog(dfa.status.fioPodName, common.NSDefault)
}

func (dfa *FioApplication) Cleanup() error {
	var err error
	// delete pod and volume
	//	if dfa.SaveFioPodLog {
	// dfa.DumpPodLog()
	//	}

	podName := dfa.GetPodName()

	// If the pod name is "", we never deployed a pod.
	if podName != "" {
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
		dfa.status.fioPodName = ""

	}
	var localEngine bool
	if dfa.OpenEbsEngine != common.Mayastor {
		localEngine = true
	}
	// Only delete PVC and storage class if they were created by this instance
	if dfa.status.createdPVC {
		err = RemovePVC(dfa.status.pvcName, dfa.status.scName, common.NSDefault, localEngine)
		if err == nil {
			dfa.status.createdPVC = false
			err = RmStorageClass(dfa.status.scName)
		}
	}

	return err
}

// ForcedCleanup - delete resources associated with the FioApplicationlication by name
// regardless of creation status
// This function should be used sparingly.
// FIXME: refactor so that is function can be replaced
// by simply calling Cleanup
func (dfa *FioApplication) ForcedCleanup() {
	_ = DeletePod(dfa.status.fioPodName, common.NSDefault)
	dfa.status.fioPodName = ""
	_ = RmPVC(dfa.status.pvcName, dfa.status.scName, common.NSDefault)
	dfa.status.createdPVC = false
	_ = RmStorageClass(dfa.status.scName)
	dfa.status.scName = ""
}

func (dfa *FioApplication) WaitComplete(timeoutSecs int) error {
	return WaitFioPodComplete(dfa.status.fioPodName, 10, timeoutSecs)
}

func (dfa *FioApplication) WaitRunning(timeoutSecs int) bool {
	logf.Log.Info("Wait for pod running", "pod", dfa.status.fioPodName, "timeoutSecs", timeoutSecs)
	return WaitPodRunning(dfa.status.fioPodName, common.NSDefault, timeoutSecs)
}

func (dfa *FioApplication) GetPodStatus() (coreV1.PodPhase, error) {
	return GetPodStatus(dfa.status.fioPodName, common.NSDefault)
}

func (dfa *FioApplication) DeletePod() error {
	var err error
	if dfa.status.fioPodName != "" {

		err = DeletePod(dfa.status.fioPodName, common.NSDefault)
		if err == nil {
			dfa.status.fioPodName = ""
		}
	}
	return err
}

func (dfa *FioApplication) SetAppNodeName(nodeName string) error {
	if nodeName == "" {
		return fmt.Errorf("nodeName cannot be empty")
	}
	dfa.AppNodeName = nodeName
	logf.Log.Info("AppNodename is ", "AppNodename", dfa.AppNodeName)
	return nil
}

func (dfa *FioApplication) GetAppUuid() string {
	return dfa.status.sessionId
}

func (dfa *FioApplication) GetScName() string {
	return dfa.status.scName
}

func (dfa *FioApplication) GetPodName() string {
	return dfa.status.fioPodName
}

func (dfa *FioApplication) GetPvcName() string {
	return dfa.status.pvcName
}

// PVC was created and PVC may not be accessbile
func (dfa *FioApplication) IsPVCCreated() bool {
	return dfa.status.createdPVC
}

func (dfa *FioApplication) RefreshVolumeState() error {
	pvc, err := GetPVC(dfa.status.pvcName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to get pvc %s, error: %v", dfa.status.pvcName, err)
	} else if pvc == nil {
		return fmt.Errorf("pvc %s not found", dfa.status.pvcName)
	} else if *pvc.Spec.StorageClassName != dfa.status.scName {
		return fmt.Errorf("storage class %s not used to create pvc %s", dfa.status.scName, dfa.status.pvcName)
	} else if pvc.ObjectMeta.UID == "" {
		return fmt.Errorf("pvc %s does not have  pvc.ObjectMeta.UID non empty string", dfa.status.pvcName)
	}
	return nil
}

// MonitorPod - invokes MonitorE2EFioPod which launches a go thread
// to stream fio pod log output and scan that stream
// to populate fields in E2eFioPodOutputMonitor
func (dfa *FioApplication) MonitorPod() (*common.E2eFioPodOutputMonitor, error) {
	var err error
	if dfa.status.monitor == nil {
		dfa.status.monitor, err = MonitorE2EFioPod(dfa.GetPodName(), common.NSDefault)
		if err != nil {
			dfa.status.monitor = nil
		}
	}
	return dfa.status.monitor, err
}

func (dfa *FioApplication) WaitFioComplete(timeoutSecs int, pollTimeSecs int) (int, error) {
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

func (dfa *FioApplication) FioTargetSizes() (map[string]uint64, error) {
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

func (dfa *FioApplication) ImportVolume(volName string) error {
	pvc, err := GetPVC(volName, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to get pvc %s, error: %v", volName, err)
	}
	if pvc == nil {
		return fmt.Errorf("pvc %s not found", volName)
	}
	dfa.status.pvcName = volName

	dfa.status.importedVolume = true
	dfa.status.suffix = "-imported-vol"

	return err
}

func (dfa *FioApplication) ImportVolumeFromApp(srcDfa *FioApplication) error {
	dfa.VolSizeMb = srcDfa.VolSizeMb
	dfa.VolType = srcDfa.VolType

	return dfa.ImportVolume(srcDfa.GetPvcName())
}
