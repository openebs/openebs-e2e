package disk_pool_operator

import (
	"fmt"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/custom_resources"
	"github.com/openebs/openebs-e2e/common/custom_resources/types"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"

	coreV1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	DefTimeoutSecs  = 90 // in seconds
	DefSleepTime    = 30 //in seconds
	IoEnginePodName = e2e_config.GetConfig().Product.IOEnginePodName
)

const (
	CrStateCreating           = "Creating"
	CrStateCreated            = "Created"
	CrStateTerminated         = "Terminating"
	CpStateOnline             = "Online"
	PoolNotFoundErrorResponse = "404 Not Found"
)

type DspOperator struct {
	PoolName string
	NodeName string
	NodeIp   *string
	NodeDisk []string
}

// CreatePoolOnNonIoEngineNode configures a node so that the IO engine is not running and
// then creates a pool targeted for that node
func (dsp *DspOperator) CreatePoolOnNonIoEngineNode() error {
	// list disk pool in cluster
	poolsInCluster, err := k8stest.ListMsPools()
	if err != nil {
		return fmt.Errorf("failed to list disk pools, error %v", err)
	}
	if len(poolsInCluster) == 0 {
		return fmt.Errorf("no pools found in cluster")
	}

	// Get pool name which will be deleted on a given node
	for _, pool := range poolsInCluster {
		dsp.PoolName = pool.Name
		dsp.NodeName = pool.Spec.Node
		dsp.NodeDisk = pool.Spec.Disks
		break
	}
	logf.Log.Info("pool", "count", len(poolsInCluster))
	logf.Log.Info("pool", "to delete", dsp.PoolName)
	logf.Log.Info("pool", "dsp", dsp)

	// delete pool on node
	err = custom_resources.DeleteMsPool(dsp.PoolName)
	if err != nil {
		return fmt.Errorf("failed to delete pool %s, error %v", dsp.PoolName, err)
	}

	// check for deleted pool
	isDeleted := VerifyDeletedPool(dsp.PoolName, DefTimeoutSecs)

	if !isDeleted {
		return fmt.Errorf("pool %s present in cluster after deletion", dsp.PoolName)
	}

	dsp.NodeIp, err = k8stest.GetNodeIPAddress(dsp.NodeName)
	if err != nil {
		return fmt.Errorf("failed to get node %s ip address %v", dsp.NodeName, err)
	}
	if len(*dsp.NodeIp) == 0 {
		return fmt.Errorf("node %s IP address is empty", dsp.NodeName)
	}

	err = k8stest.UnscheduleIoEngineFromNode(*dsp.NodeIp, dsp.NodeName, IoEnginePodName, DefTimeoutSecs)
	if err != nil {
		return err
	}

	// create disk pool on node where io engine label was removed
	_, err = custom_resources.CreateMsPool(dsp.PoolName, dsp.NodeName, dsp.NodeDisk)
	if err != nil {
		return fmt.Errorf("failed to create pool %s on node %s with disks %s, error: %v",
			dsp.PoolName,
			dsp.NodeName,
			dsp.NodeDisk,
			err,
		)
	}

	return nil
}

// VerifyDeletedPool verify pool removal from cluster
func VerifyDeletedPool(poolName string, timeoutSecs int) bool {
	const sleepTime = 3
	time.Sleep(sleepTime * time.Second)
	var err error
	var dsp types.DiskPool
	for ix := 0; ix < (timeoutSecs+sleepTime-1)/sleepTime; ix++ {
		dsp, err = custom_resources.GetMsPool(poolName)
		logf.Log.Info("dsp", "It", ix, "dsp", dsp.GetName())
		if err != nil && !k8serrors.IsNotFound(err) {
			logf.Log.Info("pool exist in cluster", "pool details", dsp.GetName())
			time.Sleep(sleepTime * time.Second)
			continue
		}
		if dsp.GetName() == "" {
			return true
		}
	}
	return false
}

type TestCase struct {
	Decor          string
	VolSizeMb      int
	ScName         string
	PodName        string
	VolName        string
	VolType        common.VolumeType
	ReplicaCount   int
	PoolsInCluster []common.MayastorPool
	Uid            string
}

func (tc *TestCase) DeployApp() error {
	var err error
	// list disk pool in cluster
	tc.PoolsInCluster, err = k8stest.ListMsPools()
	if err != nil {
		return fmt.Errorf("failed to list disk pools, error %v", err)
	}

	if len(tc.PoolsInCluster) == 0 {
		return fmt.Errorf("no pools found in cluster")
	}
	if tc.ReplicaCount == 0 {
		tc.ReplicaCount = len(tc.PoolsInCluster)
	}

	tc.ScName = strings.ToLower(fmt.Sprintf("%s-%drepl", tc.Decor, tc.ReplicaCount))
	tc.PodName = strings.ToLower(fmt.Sprintf("%s-%drepl", tc.Decor, tc.ReplicaCount))
	tc.VolName = strings.ToLower(fmt.Sprintf("%s-%drepl", tc.Decor, tc.ReplicaCount))
	tc.VolType = common.VolRawBlock

	err = k8stest.NewScBuilder().
		WithName(tc.ScName).
		WithReplicas(tc.ReplicaCount).
		WithProtocol(common.ShareProtoNvmf).
		WithNamespace(common.NSDefault).
		WithVolumeBindingMode(storageV1.VolumeBindingImmediate).
		BuildAndCreate()
	if err != nil {
		return fmt.Errorf("failed to create storage class %s %v", tc.ScName, err)
	}
	// Create the volume
	tc.Uid, err = k8stest.MkPVC(tc.VolSizeMb, tc.VolName, tc.ScName, tc.VolType, common.NSDefault)
	if err != nil {
		return fmt.Errorf("failed to create pvc %s, %v", tc.VolName, err)
	}
	logf.Log.Info("Volume", "uid", tc.Uid)
	args, efabErr := common.NewE2eFioArgsBuilder().WithDefaultRawBlock().WithDefaultArgs().Build()
	if efabErr != nil {
		return fmt.Errorf("failed to compile fio commandline %v", efabErr)
	}
	logf.Log.Info("fio", "arguments", args)
	// fio pod container
	container := k8stest.MakeFioContainer(tc.PodName, args)
	// volume claim details
	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: tc.VolName,
			},
		},
	}
	// create the fio pod
	podObj, err := k8stest.NewPodBuilder("fio").
		WithName(tc.PodName).
		WithNamespace(common.NSDefault).
		WithRestartPolicy(coreV1.RestartPolicyNever).
		WithContainer(container).
		WithVolume(volume).
		WithVolumeDeviceOrMount(tc.VolType).Build()
	if err != nil {
		return fmt.Errorf("generating fio pod definition %s, %v", tc.PodName, err)
	}
	if podObj == nil {
		return fmt.Errorf("failed to generate fio pod definition")
	}
	_, err = k8stest.CreatePod(podObj, common.NSDefault)
	if err != nil {
		return fmt.Errorf("creating fio pod %s, %v", tc.PodName, err)
	}
	// wait for pod to transition to running
	running := k8stest.WaitPodRunning(tc.PodName, common.NSDefault, DefTimeoutSecs)
	if !running {
		return fmt.Errorf("pod is not running")
	}
	return nil
}

func (tc *TestCase) Cleanup() error {
	var err error
	// delete pod and volume
	err = k8stest.DeletePod(tc.PodName, common.NSDefault)
	if err == nil {
		err = k8stest.RmPVC(tc.VolName, tc.ScName, common.NSDefault)
		if err == nil {
			err = k8stest.RmStorageClass(tc.ScName)
		}
	}
	return err
}
