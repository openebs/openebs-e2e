package lvm_volume_provisioning

import (
	"testing"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var nodeConfig lvm.LvmNodesDevicePvVgConfig
var defFioCompletionTime = 240 // in seconds

func volumeProvisioningTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {

	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               4096,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   1,
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

	logf.Log.Info("wait for fio pod to complete")
	err = app.WaitComplete(defFioCompletionTime)
	Expect(err).ToNot(HaveOccurred(), "fio app did not complete")

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

}

func TestLvmVolumeProvisioningTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "volume_provisioning", "volume_provisioning")
}

var _ = Describe("volume_provisioning", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		err := e2e_ginkgo.AfterEachCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	It("lvm ext4: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, true)
	})
	It("lvm xfs: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-xfs", common.Lvm, common.VolFileSystem, common.XfsFsType, true)
	})
	It("lvm btrfs: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-btrfs", common.Lvm, common.VolFileSystem, common.BtrfsFsType, true)
	})
	It("lvm block: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-rb", common.Lvm, common.VolRawBlock, common.NoneFsType, true)
	})

	// immediate binding
	It("lvm ext4 immediate binding: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("lvm xfs immediate binding: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-xfs", common.Lvm, common.VolFileSystem, common.XfsFsType, false)
	})
	It("lvm btrfs immediate binding: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-btrfs", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false)
	})
	It("lvm block immediate binding: should verify a volume can be created, used and deleted", func() {
		volumeProvisioningTest("lvm-rb", common.Lvm, common.VolRawBlock, common.NoneFsType, false)
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
