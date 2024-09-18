package zfs_ha_controller

import (
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8sinstall"
	"github.com/openebs/openebs-e2e/common/zfs"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

/*
	Background:
		Given that the product is installed on a kubernetes cluster
	Scenario:  Zfs controller in high availability mode
		Given Zfs controller deployed in non HA mode
		When Zfs controller is scaled down to zero replica
		And Zfs volume(pvc) is deployed
		Then Zfs volume(pvc) should remain in pending state
		When Zfs controller is scaled to two replica via helm upgrade
		Then Zfs volume(pvc) should transition to bound state
		And Deploy fio application to use zfs volume
		And Taint all worker node with NoSchedule taint
		When Delete one of the zfs controller pod
		Then Zfs node lease should point to running zfs controller pod
*/

var nodeConfig zfs.ZfsNodesDevicePoolConfig
var defLeaseSwitchTime = 120 // in seconds
var nodesWithoutTaint []string
var zfsControllerOrgReplica int32
var app k8stest.FioApplication
var volumeProvisionErrorMsg = "waiting for a volume to be created, either by external provisioner"

func controllerHaTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {
	var err error
	e2e_config := e2e_config.GetConfig().Product

	app = k8stest.FioApplication{
		Decor:                          decor,
		VolSizeMb:                      1024,
		OpenEbsEngine:                  engine,
		VolType:                        volType,
		FsType:                         fstype,
		Loops:                          5,
		VolWaitForFirstConsumer:        volBindModeWait,
		SkipPvcVerificationAfterCreate: true,
	}

	// setup sc parameters
	app.Zfs = k8stest.ZfsOptions{
		PoolName:      nodeConfig.PoolName,
		ThinProvision: common.No,
		RecordSize:    "128k",
		Compression:   common.Off,
		DedUp:         common.Off,
	}

	// list nodes without NoSchedule taint
	nodesWithoutTaint, err = k8stest.ListNodesWithoutNoScheduleTaint()
	Expect(err).ToNot(HaveOccurred(), "failed to list nodes without NoSchedule taint")
	logf.Log.Info("nodes without NoSchedule taint", "nodes", nodesWithoutTaint)

	// get the no of replicas in zfs-controller deployment
	// Scale down the zfs-controller deployment
	// Check that zfs-controller pods has been terminated successfully
	logf.Log.Info("Scale down and get the no of replicas in zfs-controller deployment")
	zfsControllerName := e2e_config.ZfsEngineControllerDeploymentName
	// ScaleZfsControllerViaHelm return original replica count of zfs controller deployment
	// before scale operation which will be used once test finishes to reset to original replica count
	zfsControllerOrgReplica, err = k8sinstall.ScaleZfsControllerViaHelm(0)
	Expect(err).To(BeNil(), "failed to scale down deployment %s, error: %v", zfsControllerName, err)
	Expect(zfsControllerOrgReplica).ShouldNot(BeZero(), "zfs controller replica count should not be zero")
	logf.Log.Info("zfs controller deployment", "name", zfsControllerName, "original replica", zfsControllerOrgReplica)

	// create sc and pvc
	logf.Log.Info("create sc, pvc, fio pod with snapshot as source")
	err = app.CreateVolume()
	Expect(err).To(BeNil(), "failed to create pvc")

	// sleep for 30 seconds
	logf.Log.Info("Sleep for 30 seconds")
	time.Sleep(30 * time.Second)

	// verify pvc to be in pending state
	pvcPhase, err := k8stest.GetPvcStatusPhase(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to get pvc phase")
	Expect(pvcPhase).Should(Equal(coreV1.ClaimPending), "pvc phase is not pending")

	// verify pvc pre condition failed events
	isEventPresent, err := k8stest.WaitForPvcNormalEvent(app.GetPvcName(), common.NSDefault, volumeProvisionErrorMsg)
	Expect(err).To(BeNil())
	Expect(isEventPresent).To(BeTrue())

	// sleep for 15 seconds and again verify state and events
	logf.Log.Info("sleep for 30 seconds")
	time.Sleep(30 * time.Second)

	// verify pvc to be in pending state
	pvcPhase, err = k8stest.GetPvcStatusPhase(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to get pvc phase")
	Expect(pvcPhase).Should(Equal(coreV1.ClaimPending), "pvc phase is not pending")

	// verify pvc pre condition failed events
	isEventPresent, err = k8stest.WaitForPvcNormalEvent(app.GetPvcName(), common.NSDefault, volumeProvisionErrorMsg)
	Expect(err).To(BeNil())
	Expect(isEventPresent).To(BeTrue())

	// Scale up the zfs-controller deployment replica to initial replica + 1
	logf.Log.Info("Scale up zfs-controller deployment")
	_, err = k8sinstall.ScaleZfsControllerViaHelm(zfsControllerOrgReplica + 1)
	Expect(err).To(BeNil(), "failed to scale deployment %s, error: %v", zfsControllerName, err)

	//verify pvc and pv to be bound
	volUuid, err := k8stest.VerifyVolumeProvision(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
	Expect(volUuid).ToNot(BeEmpty())

	// verify volume state
	err = app.RefreshVolumeState()
	Expect(err).ToNot(HaveOccurred())

	// deploy fio pod with created volume
	logf.Log.Info("deploy fio pod with created volume")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	// Get the name of the controller pod replica which is active as master at present
	lease, err := k8stest.GetLease(e2e_config.ZfsEngineLeaseName, common.NSOpenEBS())
	Expect(err).To(BeNil(), "failed to get lease %s in %s namespace", e2e_config.ZfsEngineLeaseName, common.NSOpenEBS())
	Expect(lease).ToNot(BeNil(), "no lease found")

	// get zfs controller master pod i.e spec.holderIdentity from lease
	initialHolderIdentity := lease.Spec.HolderIdentity
	Expect(initialHolderIdentity).ToNot(BeNil(), "no lease HolderIdentity found")
	logf.Log.Info("Zfs controller", "Initial HolderIdentity", initialHolderIdentity)

	// taint all nodes so that after deleting zfs controller pod which holds lease, new pod should not get scheduled
	for _, node := range nodesWithoutTaint {
		err = k8stest.AddNoScheduleTaintOnNode(node)
		Expect(err).To(BeNil(), "failed to taint node %s", node)
	}

	// delete zfs controller pod which is holding lease
	err = k8stest.DeletePod(*initialHolderIdentity, common.NSOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to delete pod %s", *initialHolderIdentity)

	// wait for lease to switch to different pod
	Eventually(func() bool {
		lease, err := k8stest.GetLease(e2e_config.ZfsEngineLeaseName, common.NSOpenEBS())
		if err != nil {
			logf.Log.Info("failed to get lease", "lease name", e2e_config.ZfsEngineLeaseName, "namespace", common.NSOpenEBS(), "error", err)
			return false
		}
		if *lease.Spec.HolderIdentity == "" {
			logf.Log.Info("lease HolderIdentity can not be empty string")
			return false
		} else if *lease.Spec.HolderIdentity != *initialHolderIdentity {
			return true
		}

		return false
	},
		defLeaseSwitchTime,
		"5s",
	).Should(BeTrue())
	Expect(err).To(BeNil(), "failed to get lease")

	// Check fio pod status
	phase, podLogSynopsis, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSynopsis)

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

	// remove taints form nodes
	for _, node := range nodesWithoutTaint {
		err = k8stest.RemoveNoScheduleTaintFromNode(node)
		Expect(err).To(BeNil(), "failed to taint node %s", node)
	}
	nodesWithoutTaint = []string{}

	ready, err := k8stest.OpenEBSReady(10, 540)
	Expect(err).To(BeNil(), "failed to verify openebs pods running state")
	Expect(ready).To(BeTrue(), "some of the openebs pods are not running")
}

func TestZfsControllerHaTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "zfs_ha_controller", "zfs_ha_controller")
}

var _ = Describe("zfs_ha_controller", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		after_err := e2e_ginkgo.AfterEachK8sCheck()
		if len(nodesWithoutTaint) != 0 {
			// remove taints form nodes
			for _, node := range nodesWithoutTaint {
				err := k8stest.RemoveNoScheduleTaintFromNode(node)
				Expect(err).To(BeNil(), "failed to taint node %s", node)
			}
			ready, err := k8stest.OpenEBSReady(10, 540)
			Expect(err).To(BeNil(), "failed to verify openebs pods running state")
			Expect(ready).To(BeTrue(), "some of the openebs pods are not running")
		}

		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err := app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		// Scale up the zfs-controller deployment replica to initial replica
		logf.Log.Info("Scale up zfs-controller deployment")
		_, err = k8sinstall.ScaleZfsControllerViaHelm(zfsControllerOrgReplica)
		Expect(err).To(BeNil(), "failed to scale  deployment %s, error: %v", e2e_config.GetConfig().Product.ZfsEngineControllerDeploymentName, err)

		Expect(after_err).ToNot(HaveOccurred())
	})

	It("zfs ext4 immediate binding: should verify high availability mode", func() {
		controllerHaTest("zfs-ha", common.Zfs, common.VolFileSystem, common.Ext4FsType, false)
	})

})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with zfs pool
	nodeConfig, err = zfs.SetupZfsNodes("zfspv-pool", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup zfs pool")

})

var _ = AfterSuite(func() {

	logf.Log.Info("remove node with device and zpool", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredDeviceZfsPool()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
