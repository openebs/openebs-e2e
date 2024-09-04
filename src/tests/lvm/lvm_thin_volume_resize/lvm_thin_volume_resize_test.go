package lvm_thin_volume_resize

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"
	volumeResize "github.com/openebs/openebs-e2e/src/tests/lvm/common"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Background:
// Given a k8s cluster is running with the product installed

// Scenario: thin volume resize for lvm volume
//      Given a thin volume has been successfully created
//      And application is using that volume
//      When volume size is updated in increased size in pvc spec
//      Then volume should be resized to desired capacity
//      And pvc and pv objects should verify that capacity
//      And application should be able to use that resized space

func TestLvmThinVolumeResizeTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_thin_volume_resize", "lvm_thin_volume_resize")
}

var _ = Describe("lvm_thin_volume_resize", func() {

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

	It("lvm ext4: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.Ext4FsType, true, common.Yes)
	})
	It("lvm xfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.XfsFsType, true, common.Yes)
	})
	It("lvm btrfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.BtrfsFsType, true, common.Yes)
	})
	It("lvm block: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolRawBlock, common.NoneFsType, true, common.Yes)
	})

	// immediate binding
	It("lvm ext4: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.Ext4FsType, false, common.Yes)
	})
	It("lvm xfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.XfsFsType, false, common.Yes)
	})
	It("lvm btrfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false, common.Yes)
	})
	It("lvm block: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, common.VolRawBlock, common.NoneFsType, false, common.Yes)
	})
})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

})

var _ = AfterSuite(func() {
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	logf.Log.Info("remove node with device and vg", "node config", volumeResize.NodeConfig)
	err := volumeResize.NodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
