package zfs_volume_provisioning

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/zfs"

	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var nodeConfig zfs.ZfsNodesDevicePoolConfig
var app k8stest.FioApplication
var err error

func TestZfsVolumeProvisioningTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "zfs_volume_provisioning", "zfs_volume_provisioning")
}

var _ = Describe("zfs_volume_provisioning", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {

		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err = app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		// Check resource leakage.
		err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())

	})

	// immediate binding
	It("zfs ext4 immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("ext4", common.Zfs, common.VolFileSystem, common.Ext4FsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs xfs immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("xfs", common.Zfs, common.VolFileSystem, common.XfsFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs btrfs immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("btrfs", common.Zfs, common.VolFileSystem, common.BtrfsFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs block immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("rb", common.Zfs, common.VolRawBlock, common.NoneFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs zfs immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("zfs", common.Zfs, common.VolFileSystem, common.ZfsFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})

	// late binding
	It("zfs ext4: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("ext4", common.Zfs, common.VolFileSystem, common.Ext4FsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs xfs: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("xfs", common.Zfs, common.VolFileSystem, common.XfsFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs btrfs: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("btrfs", common.Zfs, common.VolFileSystem, common.BtrfsFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs zfs: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("zfs", common.Zfs, common.VolFileSystem, common.ZfsFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("zfs block: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("rb", common.Zfs, common.VolRawBlock, common.NoneFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})

})

var _ = BeforeSuite(func() {
	err = e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with zfs pool
	nodeConfig, err = zfs.SetupZfsNodes("zfspv-pool", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup zfs pool")

})

var _ = AfterSuite(func() {
	// logf.Log.Info("remove node with device and zpool", "node config", nodeConfig)
	err = nodeConfig.RemoveConfiguredDeviceZfsPool()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
