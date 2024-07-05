package k8stest

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"

	v1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type FioStsApp struct {
	Decor     string
	VolSizeMb int
	ScName    string
	StsName   string
	VolName   string
	// StatefulSet allows you to relax its ordering guarantees while preserving
	// its uniqueness and identity guarantees via its .spec.podManagementPolicy field
	// values permitted: "Parallel" and "OrderedReady"
	StsPodManagementPolicy string
	StsAffinityGroup       string // "true" or "false" for anti affinity feature for volume replicas and target
	VolType                common.VolumeType
	StsReplicaCount        *int32
	VolReplicaCount        int
	VolUuid                string
	FioArgs                []string
	Runtime                uint
	Loops                  int
	VerifyReplicasOnDelete bool
	AppNodeName            string
}

var (
	StsDefTimeoutSecs = 90 // in seconds
)

func (dfa *FioStsApp) StsApp() error {

	var err error
	// list disk pool in cluster
	if dfa.Runtime != 0 && dfa.Loops != 0 {
		return fmt.Errorf("cannot specify both Runtime and Loops")
	}

	var poolsInCluster []common.MayastorPool
	const sleepTime = 3
	for ix := 0; ix < (StsDefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
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
	if dfa.VolReplicaCount == 0 {
		dfa.VolReplicaCount = len(poolsInCluster)
	}

	dfa.ScName = strings.ToLower(fmt.Sprintf("%s-%d-repl-sc", dfa.Decor, dfa.VolReplicaCount))
	dfa.StsName = strings.ToLower(fmt.Sprintf("%s-%d-repl-sts", dfa.Decor, *dfa.StsReplicaCount))
	dfa.VolName = strings.ToLower(fmt.Sprintf("%s-%d-repl-vol", dfa.Decor, dfa.VolReplicaCount))

	if dfa.VolType.String() == "" {
		dfa.VolType = common.VolRawBlock
	}

	err = NewScBuilder().
		WithName(dfa.ScName).
		WithReplicas(dfa.VolReplicaCount).
		WithProtocol(common.ShareProtoNvmf).
		WithStsAffinityGroup(common.StsAffinityGroup(dfa.StsAffinityGroup)).
		WithVolumeBindingMode(storageV1.VolumeBindingImmediate).
		BuildAndCreate()
	if err != nil {
		return fmt.Errorf("failed to create storage class %s %v", dfa.ScName, err)
	}

	if len(dfa.FioArgs) == 0 && dfa.VolType == common.VolRawBlock {
		dfa.FioArgs = append(dfa.FioArgs, "--")
		dfa.FioArgs = append(dfa.FioArgs, fmt.Sprintf("--filename=%s", common.FioBlockFilename))
		dfa.FioArgs = append(dfa.FioArgs, common.GetFioArgs()...)
	} else if len(dfa.FioArgs) == 0 && dfa.VolType == common.VolFileSystem {
		dfa.FioArgs = append(dfa.FioArgs, "--")
		dfa.FioArgs = append(dfa.FioArgs, fmt.Sprintf("--filename=%s", common.FioFsFilename))
		dfa.FioArgs = append(dfa.FioArgs, fmt.Sprintf("--size=%dm", dfa.VolSizeMb-200))
		dfa.FioArgs = append(dfa.FioArgs, common.GetFioArgs()...)
	}
	if dfa.Runtime != 0 {
		dfa.FioArgs = append(dfa.FioArgs, "--time_based")
		dfa.FioArgs = append(dfa.FioArgs, fmt.Sprintf("--runtime=%v", dfa.Runtime))
	}
	if dfa.Loops != 0 {
		dfa.FioArgs = append(dfa.FioArgs, fmt.Sprintf("--loops=%d", dfa.Loops))
	}

	logf.Log.Info("fio", "arguments", dfa.FioArgs)

	labelKey := "e2e-test"
	labelValue := "sts-anti-affinity"
	labelselector := map[string]string{
		labelKey: labelValue,
	}

	nodeKey := "kubernetes.io/hostname"
	nodeValue := dfa.AppNodeName
	nodeSelector := map[string]string{
		nodeKey: nodeValue,
	}

	volSizeMbStr := fmt.Sprintf("%dMi", dfa.VolSizeMb)

	// create the fio sts
	sts := NewStatefulsetBuilder().
		WithName(dfa.StsName).
		WithNamespace(common.NSDefault).
		WithLabels(labelselector).
		WithSelectorMatchLabels(labelselector).
		WithPodTemplateSpecBuilder(
			NewPodtemplatespecBuilder().
				WithLabels(labelselector).
				WithContainerBuildersNew(
					NewContainerBuilder().
						WithName(dfa.StsName).
						WithImage(common.GetFioImage()).
						WithImagePullPolicy(coreV1.PullAlways).
						WithArgumentsNew(dfa.FioArgs).
						WithVolumeDeviceOrMount(dfa.VolName, dfa.VolType))).
		WithReplicas(dfa.StsReplicaCount).
		WithPodManagementPolicy(v1.PodManagementPolicyType(dfa.StsPodManagementPolicy)).
		WithVolumeClaimTemplate(dfa.VolName, volSizeMbStr, dfa.ScName, dfa.VolType)

	if dfa.AppNodeName != "" {
		sts = sts.WithNodeSelector(nodeSelector)
	}

	stsObj, err := sts.Build()
	if err != nil {
		return fmt.Errorf("failed to create statefulset %s definition object in %s namespace", dfa.StsName, common.NSDefault)
	}

	logf.Log.Info("sts", "object", stsObj)
	err = CreateStatefulset(stsObj)
	if err != nil {
		return fmt.Errorf("failed to create statefulset %s, error %v", dfa.StsName, err)
	}

	logf.Log.Info("Verify if statefulset is ready", "statefulset", dfa.StsName, "namespace", common.NSDefault)

	var stsReady bool
	for ix := 0; ix < (StsDefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
		if StatefulSetReady(dfa.StsName, common.NSDefault) {
			stsReady = true
			break
		}
		time.Sleep(sleepTime * time.Second)
	}
	if !stsReady {
		return fmt.Errorf("fio statefulset %s not ready", dfa.StsName)
	}

	return nil
}

func (dfa *FioStsApp) Cleanup(replicaCount int) error {
	var err error
	// delete pod and volume
	err = DeleteStatefulset(dfa.StsName, common.NSDefault)
	if err == nil {
		for i := 0; i < replicaCount; i++ {
			// PVC name format when create with statefulset looks like this:
			// pvc name ==  claimName plus hyphen (-) and append with sts pod name
			// for e.g `testpvc-pod-0`
			// where testpvc is claim name and pod is statefulset name and pod-0 is pod name
			volName := dfa.VolName + "-" + dfa.StsName + "-" + strconv.Itoa(i)
			err = RmPVC(volName, dfa.ScName, common.NSDefault)
			if err != nil {
				return fmt.Errorf("failed to remove pvc %s", volName)
			}
		}
		err = RmStorageClass(dfa.ScName)
		if err != nil {
			return fmt.Errorf("failed to delete storage class %s", dfa.ScName)
		}
	}
	return err
}
