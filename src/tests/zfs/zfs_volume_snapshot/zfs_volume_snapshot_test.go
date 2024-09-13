package zfs_volume_snapshot

import (
	"fmt"
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/mayastor/snapshot"
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
        Given a k8s cluster is running with the product installed
        And snapshot CRDs exists in k8s cluster

   Scenario: Snapshot creation for zfs volume
        Given a volume has been successfully created
        When snapshot is created for the volume
        Then the snapshot should be successfully created
		And the snapshot  object should be ready
        And the snapshot content object associated with snapshot should be ready
*/

var nodeConfig zfs.ZfsNodesDevicePoolConfig
var app k8stest.FioApplication
var snapshotClassName, snapshotName, snapshotNamespace string

func volumeSnapshotTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {
	app = k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               2048,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   3,
		VolWaitForFirstConsumer: volBindModeWait,
	}

	// setup sc parameters
	app.Zfs = k8stest.ZfsOptions{
		PoolName:      nodeConfig.PoolName,
		ThinProvision: common.No,
		RecordSize:    "128k",
		Compression:   common.Off,
		DedUp:         common.Off,
	}

	if app.FsType == common.ZfsFsType || app.FsType == common.BtrfsFsType {
		app.FsPercent = 50
	}

	logf.Log.Info("create sc, pvc, fio pod")
	err := app.DeployApplication()
	Expect(err).ToNot(HaveOccurred())

	logf.Log.Info("Sleep for 30 seconds before taking snapshot")
	time.Sleep(30 * time.Second)

	// snapshot steps
	snapshotClassName = fmt.Sprintf("snapshotclass-%s", app.GetPvcName())
	snapshotName = fmt.Sprintf("snapshot-%s", app.GetPvcName())
	snapshotNamespace = common.NSDefault
	logf.Log.Info("Create Snapshot", "Snapshot class", snapshotClassName, "Snapshot", snapshotName, "Namespace", snapshotNamespace)
	csiDriver := e2e_config.GetConfig().Product.ZfsEngineProvisioner
	// create snapshot for volume
	snapshotObj, snapshotContentName, err := k8stest.CreateVolumeSnapshot(snapshotClassName, snapshotName, app.GetPvcName(), common.NSDefault, csiDriver)
	Expect(err).ToNot(HaveOccurred())
	logf.Log.Info("Snapshot Created ", "Snapshot", snapshotObj, "Snapshot Content Name", snapshotContentName)

	Expect(snapshotContentName).ToNot(BeEmpty(), "snapshot content name should not be empty for snapshot %s", snapshotName)
	Expect(snapshotObj).ShouldNot(BeNil())

	// verify successful snapshot CR creation
	status, err := snapshot.VerifySuccessfulSnapshotCreation(snapshotName, snapshotContentName, snapshotNamespace, true)
	Expect(err).ToNot(HaveOccurred(), "error while verifying snapshot creation")
	Expect(status).Should(BeTrue(), "failed to verify successful snapshot %s creation", snapshotName)

	// Check fio pod status
	phase, podLogSysnopsis, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSysnopsis)

	// remove snapshot and snapshot class
	err = snapshot.DeleteVolumeSnapshot(snapshotClassName, snapshotName, snapshotNamespace)
	Expect(err).ToNot(HaveOccurred())
	snapshotName = ""

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

}

func TestZfsVolumeSnapshotTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "zfs_volume_snapshot", "zfs_volume_snapshot")
}

var _ = Describe("zfs_volume_snapshot", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		afterEacherr := e2e_ginkgo.AfterEachK8sCheck()
		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		// remove snapshot and snapshot class
		if snapshotName != "" {
			err := snapshot.DeleteVolumeSnapshot(snapshotClassName, snapshotName, snapshotNamespace)
			Expect(err).ToNot(HaveOccurred())
		}

		err := app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")
		Expect(afterEacherr).ToNot(HaveOccurred())
	})

	It("zfs: ext4 immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("ext4", common.Zfs, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("zfs: xfs immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("xfs", common.Zfs, common.VolFileSystem, common.XfsFsType, false)
	})
	It("zfs: btrfs immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("btrfs", common.Zfs, common.VolFileSystem, common.BtrfsFsType, false)
	})
	It("zfs: zfs immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("btrfs", common.Zfs, common.VolFileSystem, common.ZfsFsType, false)
	})
	It("zfs: block immediate binding: should verify a volume snapshot", func() {
		volumeSnapshotTest("rb", common.Zfs, common.VolRawBlock, common.NoneFsType, false)
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
