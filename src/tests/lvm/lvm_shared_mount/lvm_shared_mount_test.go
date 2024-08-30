package lvm_shared_mount_volume

import (
	"testing"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"
	"github.com/openebs/openebs-e2e/common/mayastor/volume_resize"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Background:
//     Given a k8s cluster is running with the product installed
// Scenario: shared mount volume creation for lvm volume
//     Given a volume has been successfully created with shared mount enable via storage class
//     When two application are deployed using same volume
//     Then the both application should be in running state using same volume

var nodeConfig lvm.LvmNodesDevicePvVgConfig

func volumeSharedMountTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {

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

	workerNodes, err := lvm.ListLvmNode(common.NSOpenEBS())
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
		Shared:        common.Yes,
	}

	// create sc and pvc
	logf.Log.Info("create sc, pvc, fio pod with snapshot as source")
	err = app.CreateVolume()
	Expect(err).To(BeNil(), "failed to create pvc")

	//verify pvc and pv to be bound
	volUuid, err := k8stest.VerifyVolumeProvision(app.GetPvcName(), common.NSDefault)
	Expect(err).ToNot(HaveOccurred())
	Expect(volUuid).ToNot(BeEmpty())

	// deploy fio pod with created volume
	logf.Log.Info("deploy fio pod with created volume")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	app2 := app
	app2.Decor = app.Decor + "second-app"

	// deploy fio pod with created volume
	logf.Log.Info("deploy fio pod with created volume")
	err = app2.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	// Check fio pod 1 status
	logf.Log.Info("Check fio pod status")
	phase, podLogSynopsis, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSynopsis)

	// Check fio pod 2 status
	logf.Log.Info("Check fio pod status")
	phase, podLogSynopsis, err = k8stest.CheckFioPodCompleted(app2.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSynopsis)

	// wait for fio completion - monitoring log output
	exitValue, fErr := app.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)

	// wait for fio completion - monitoring log output
	exitValue, fErr = app2.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

	// remove app pod, pvc,sc
	err = app2.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")

}

func TestLvmVolumeResizeTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_shared_mount_volume", "lvm_shared_mount_volume")
}

var _ = Describe("lvm_shared_mount_volume", func() {

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

	It("lvm ext4: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.Ext4FsType, true)
	})
	It("lvm xfs: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.XfsFsType, true)
	})
	It("lvm btrfs: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.BtrfsFsType, true)
	})
	It("lvm block: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolRawBlock, common.NoneFsType, true)
	})

	// immediate binding
	It("lvm ext4: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("lvm xfs: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.XfsFsType, false)
	})
	It("lvm btrfs: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false)
	})
	It("lvm block: should verify volume resize", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolRawBlock, common.NoneFsType, false)
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