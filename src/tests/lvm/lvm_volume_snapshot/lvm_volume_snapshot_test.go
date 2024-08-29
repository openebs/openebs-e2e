package lvm_volume_snapshot

import (
	"fmt"
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"
	"github.com/openebs/openebs-e2e/common/mayastor/snapshot"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

/*
   Background:
        Given a k8s cluster is running with the product installed
        And snapshot CRDs exists in k8s cluster

   Scenario: Snapshot creation for lvm volume
        Given a volume has been successfully created
        When snapshot is created for the volume
        Then the snapshot should be successfully created
		And the snapshot  object should be ready
        And the snapshot content object associated with snapshot should be ready
*/

var nodeConfig lvm.LvmNodesDevicePvVgConfig

func volumeSnapshotTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {

	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               4096,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   5,
		VolWaitForFirstConsumer: volBindModeWait,
	}

	loopDevice := e2e_agent.LoopDevice{
		Size:   10737418240,
		ImgDir: "/tmp",
	}

	workerNodes, err := lvm.ListLvmNode(common.NSMayastor())
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

	logf.Log.Info("create sc, pvc, fio pod")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	time.Sleep(30 * time.Second)

	// snapshot steps
	snapshotClassName := fmt.Sprintf("snapshotclass-%s", app.GetPvcName())
	snapshotName := fmt.Sprintf("snapshot-%s", app.GetPvcName())
	snapshotNamespace := common.NSDefault
	logf.Log.Info("Create Snapshot", "Snapshot class", snapshotClassName, "Snapshot", snapshotName, "Namespace", snapshotNamespace)
	csiDriver := e2e_config.GetConfig().Product.LvmEngineProvisioner
	// create snapshot for volume
	snapshotObj, snapshotContentName, err := k8stest.CreateVolumeSnapshot(snapshotClassName, snapshotName, app.GetPvcName(), common.NSDefault, csiDriver)
	Expect(err).ToNot(HaveOccurred())
	logf.Log.Info("Snapshot Created ", "Snapshot", snapshotObj, "Snapshot Content Name", snapshotContentName)

	Expect(snapshotContentName).ToNot(BeEmpty(), "snapshot content name should not be empty for snapshot %s", snapshotName)
	Expect(snapshotObj).ShouldNot(BeNil())

	// verify Snapshot CR
	status, err := snapshot.VerifySuccessfulSnapshotCreation(snapshotName, snapshotContentName, snapshotNamespace, true)
	Expect(err).ToNot(HaveOccurred(), "error while verifying snapshot creation")
	Expect(status).Should(BeTrue(), "failed to verify successful snapshot %s creation", snapshotName)

	// Check fio pod status
	phase, podLogSysnopsis, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSysnopsis)

	// remove snapshot and snapshot class
	err = snapshot.DeleteVolumeSnapshot(snapshotClassName, snapshotName, common.NSDefault)
	Expect(err).ToNot(HaveOccurred())

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

}

func TestLvmVolumeSnapshotTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_volume_snapshot", "lvm_volume_snapshot")
}

var _ = Describe("lvm_volume_snapshot", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	// immediate binding
	It("lvm ext4 immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("lvm xfs immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("lvm-xfs", common.Lvm, common.VolFileSystem, common.XfsFsType, false)
	})
	It("lvm btrfs immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("lvm-btrfs", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false)
	})
	It("lvm block immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("lvm-rb", common.Lvm, common.VolRawBlock, common.NoneFsType, false)
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
