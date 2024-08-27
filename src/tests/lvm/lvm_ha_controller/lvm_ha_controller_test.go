package lvm_ha_controller

import (
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"
	"github.com/openebs/openebs-e2e/common/mayastor/restore"

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
	Scenario:  Lvm controller in high availability mode
		Given Lvm controller deployed in non HA mode
		When Lvm controller is scaled down to zero replica
		And Lvm volume(pvc) is deployed
		Then Lvm volume(pvc) should remain in pending state
		When Lvm controller is scaled to two replica via helm upgrade
		Then lvm volume(pvc) should transition to bound state
		And Deploy fio application to use lvm volume
		And Taint all worker node with NoSchedule taint
		When Delete one of the lvm controller pod
		Then Lvm node lease should point to running lvm controller pod
*/

var nodeConfig lvm.LvmNodesDevicePvVgConfig
var defLeaseSwitchTime = 120 // in seconds
var nodesWithoutTaint []string
var lvmControllerOrgReplica int32
var volumeProvisionErrorMsg = "waiting for a volume to be created"

func controllerHaTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {

	e2e_config := e2e_config.GetConfig().Product
	app := k8stest.FioApplication{
		Decor:                          decor,
		VolSizeMb:                      4096,
		OpenEbsEngine:                  engine,
		VolType:                        volType,
		FsType:                         fstype,
		Loops:                          5,
		VolWaitForFirstConsumer:        volBindModeWait,
		SkipPvcVerificationAfterCreate: true,
	}

	loopDevice := e2e_agent.LoopDevice{
		Size:   10737418240,
		ImgDir: "/tmp",
	}

	workerNodes, err := lvm.ListLvmNode(common.NsOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to list worker node")

	nodeConfig = lvm.LvmNodesDevicePvVgConfig{
		VgName:        "lvmvg",
		NodeDeviceMap: make(map[string]e2e_agent.LoopDevice), // Properly initialize the map
	}
	for _, node := range workerNodes {
		nodeConfig.NodeDeviceMap[node] = loopDevice
	}

	logf.Log.Info("setup node with loop device, pv and vg", "node config", nodeConfig)
	err = nodeConfig.ConfigureLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to setup node")

	// setup sc parameters
	app.Lvm = k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
	}

	// list nodes without NoSchedule taint
	nodesWithoutTaint, err = k8stest.ListNodeWithoutNoScheduleTaint()
	Expect(err).ToNot(HaveOccurred(), "failed to list nodes without NoSchedule taint")
	logf.Log.Info("nodes without NoSchedule taint", "nodes", nodesWithoutTaint)

	// get the no of replicas in lvm-controller deployment
	// Scale down the lvm-controller deployment
	// Check that lvm-controller pods has been terminated successfully
	logf.Log.Info("Get the no of replicas in lvm-controller deployment")
	lvmControllerName := e2e_config.LvmEngineControllerDeploymentName
	lvmControllerOrgReplica, err = k8stest.ZeroDeploymentReplicas(lvmControllerName, common.NsOpenEBS(), 60)
	Expect(err).To(BeNil(), "failed to scale down deployment %s, error: %v", lvmControllerName, err)
	logf.Log.Info("Lvm controller deployment", "name", lvmControllerName, "original replica", lvmControllerOrgReplica)

	// create sc and pvc
	logf.Log.Info("create sc, pvc, fio pod with snapshot as source")
	err = app.CreateVolume()
	Expect(err).To(BeNil(), "failed to create pvc")

	// verify pvc to be in pending state
	pvcPhase, err := k8stest.GetPvcStatusPhase(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to get pvc phase")
	Expect(pvcPhase).Should(Equal(coreV1.ClaimPending), "pvc phase is not pending")

	// verify pvc pre condition failed events
	isEventPresent, err := restore.WaitForPvcWarningEvent(app.GetPvcName(), common.NSDefault, restore.PreConditionFailedErrorSubstring)
	Expect(err).To(BeNil())
	Expect(isEventPresent).To(BeTrue())

	// sleep for 15 seconds and again verify state and events
	logf.Log.Info("sleep for 15 seconds")
	time.Sleep(15 * time.Second)

	// verify pvc to be in pending state
	pvcPhase, err = k8stest.GetPvcStatusPhase(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to get pvc phase")
	Expect(pvcPhase).Should(Equal(coreV1.ClaimPending), "pvc phase is not pending")

	// verify pvc pre condition failed events
	isEventPresent, err = restore.WaitForPvcWarningEvent(app.GetPvcName(), common.NSDefault, restore.PreConditionFailedErrorSubstring)
	Expect(err).To(BeNil())
	Expect(isEventPresent).To(BeTrue())

	// Scale up the lvm-controller deployment replica to initial replica + 1
	err = k8stest.RestoreDeploymentReplicas(lvmControllerName, common.NsOpenEBS(), 120, lvmControllerOrgReplica+1)
	Expect(err).To(BeNil(), "failed to scale  deployment %s, error: %v", lvmControllerName, err)

	//verify pvc and pv to be bound
	volUuid, err := k8stest.VerifyVolumeProvision(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
	Expect(volUuid).ToNot(BeEmpty())

	// use created PVC which is deployed as part of restore app
	err = app.RefreshVolumeState()
	Expect(err).ToNot(HaveOccurred())

	// deploy fio pod with created volume
	logf.Log.Info("deploy fio pod with created volume")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	// Get the name of the controller pod replica which is active as master at present
	lease, err := k8stest.GetLease(e2e_config.LvmEngineLeaseName, "kube-system")
	Expect(err).To(BeNil(), "failed to get lease %s in kube-system namespace", e2e_config.LvmEngineLeaseName)
	Expect(lease).ToNot(BeNil(), "no lease found")

	// get lvm controller master pod i.e spec.holderIdentity from lease
	initialHolderIdentity := lease.Spec.HolderIdentity
	Expect(initialHolderIdentity).ToNot(BeNil(), "no lease HolderIdentity found")
	logf.Log.Info("Lvm controller", "Initial HolderIdentity", initialHolderIdentity)

	// taint all nodes so that after deleting one of the lvm controller pods, now pod should not scheduled
	for _, node := range nodesWithoutTaint {
		err = k8stest.AddNoScheduleTaintOnNode(node)
		Expect(err).To(BeNil(), "failed to taint node %s", node)
	}

	// delete lvm controller pod which is holding lease
	err = k8stest.DeletePod(*initialHolderIdentity, common.NsOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "fio to delete pod %s", *initialHolderIdentity)

	// wait for lease to switch to different pod
	Eventually(func() bool {
		lease, err := k8stest.GetLease(e2e_config.LvmEngineLeaseName, "kube-system")
		if err != nil {
			logf.Log.Info("failed to get lease in kube-system namespace", "lease name", e2e_config.LvmEngineLeaseName, "error", err)
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
		defFioCompletionTime,
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

}

func TestLvmVolumeProvisioningTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_ha_controller", "lvm_ha_controller")
}

var _ = Describe("lvm_ha_controller", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		err := e2e_ginkgo.AfterEachCheck()
		Expect(err).ToNot(HaveOccurred())
		if len(nodesWithoutTaint) != 0 {
			// remove taints form nodes
			for _, node := range nodesWithoutTaint {
				err = k8stest.RemoveNoScheduleTaintFromNode(node)
				Expect(err).To(BeNil(), "failed to taint node %s", node)
			}
		}

		// Scale up the lvm-controller deployment replica to initial replica + 1
		err = k8stest.RestoreDeploymentReplicas(e2e_config.GetConfig().Product.LvmEngineControllerDeploymentName, common.NsOpenEBS(), 120, lvmControllerOrgReplica+1)
		Expect(err).To(BeNil(), "failed to scale  deployment %s, error: %v", e2e_config.GetConfig().Product.LvmEngineControllerDeploymentName, err)
	})

	It("lvm ext4: should verify high availability mode", func() {
		controllerHaTest("lvm-ha", common.Lvm, common.VolFileSystem, common.Ext4FsType, true)
	})
	It("lvm block: should verify high availability mode", func() {
		controllerHaTest("lvm-ha", common.Lvm, common.VolRawBlock, common.NoneFsType, true)
	})

	// immediate binding
	It("lvm ext4 immediate binding: should verify high availability mode", func() {
		controllerHaTest("lvm-ha", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("lvm block immediate binding: should verify high availability mode", func() {
		controllerHaTest("lvm-ha", common.Lvm, common.VolRawBlock, common.NoneFsType, false)
	})
})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

})

var _ = AfterSuite(func() {
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	logf.Log.Info("remove node with device and vg", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})