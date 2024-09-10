package lvm_volume_provisioning

import (
	"testing"

	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var err error
var nodeConfig lvm.LvmNodesDevicePvVgConfig
var app k8stest.FioApplication

func TestLvmVolumeProvisioningTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "volume_provisioning", "volume_provisioning")
}

var _ = Describe("volume_provisioning", func() {

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

	It("lvm ext4: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm xfs: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-xfs", common.Lvm, common.VolFileSystem, common.XfsFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm btrfs: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-btrfs", common.Lvm, common.VolFileSystem, common.BtrfsFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm block: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-rb", common.Lvm, common.VolRawBlock, common.NoneFsType, true, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})

	// immediate binding
	It("lvm ext4 immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm xfs immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-xfs", common.Lvm, common.VolFileSystem, common.XfsFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm btrfs immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-btrfs", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
	It("lvm block immediate binding: should verify a volume can be created, used and deleted", func() {
		app, err = VolumeProvisioningTest("lvm-rb", common.Lvm, common.VolRawBlock, common.NoneFsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})
})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with lvm pv and vg
	nodeConfig, err = lvm.SetupLvmNodes("lvmvg", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup lvm pv and vg")
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
