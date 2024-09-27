package lvm_thin_volume_resize

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/lvm"
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

var nodeConfig lvm.LvmNodesDevicePvVgConfig

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
		after_err := e2e_ginkgo.AfterEachK8sCheck()
		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err := volumeResize.ResizeApp.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")
		err = volumeResize.ResizeApp2.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")
		if volumeResize.ThinPoolNode != "" {
			out, err := e2e_agent.LvmLvRemoveThinPool(volumeResize.ThinPoolNode, "lvmvg")
			Expect(err).To(BeNil(), "failed to remove lv thin pool on node %s with vg %s, output: %s", volumeResize.ThinPoolNode, "lvmvg", out)
		}

		Expect(after_err).ToNot(HaveOccurred())

	})

	It("lvm ext4: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.Ext4FsType, true, common.Yes)
	})
	It("lvm xfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.XfsFsType, true, common.Yes)
	})
	It("lvm btrfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.BtrfsFsType, true, common.Yes)
	})
	It("lvm block: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolRawBlock, common.NoneFsType, true, common.Yes)
	})

	// immediate binding
	It("lvm ext4: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.Ext4FsType, false, common.Yes)
	})
	It("lvm xfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.XfsFsType, false, common.Yes)
	})
	It("lvm btrfs: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolFileSystem, common.BtrfsFsType, false, common.Yes)
	})
	It("lvm block: should verify thin volume resize", func() {
		volumeResize.LvmVolumeResizeTest("lvm-thin-volume-resize", common.Lvm, nodeConfig.VgName, common.VolRawBlock, common.NoneFsType, false, common.Yes)
	})
})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with lvm pv and vg
	nodeConfig, err = lvm.SetupLvmNodes("lvmvg", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup lvm pv and vg")
	err = lvm.EnableLvmThinPoolAutoExpansion(75, 20)
	Expect(err).ToNot(HaveOccurred(), "failed to update thin_pool_autoextend_threshold and thin_pool_autoextend_percent in lvm.conf")

})

var _ = AfterSuite(func() {
	// logf.Log.Info("remove node with device and vg", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
