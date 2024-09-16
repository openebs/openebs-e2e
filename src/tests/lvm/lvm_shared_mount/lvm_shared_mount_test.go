package lvm_shared_mount_volume

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/lvm"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Background:
//     Given a k8s cluster is running with the product installed
// Scenario: shared mount volume creation for lvm volume
//     Given a volume has been successfully created with shared mount enabled via storage class
//     When two applications are deployed using the same volume
//     Then both applications should be in running state using the same volume

var nodeConfig lvm.LvmNodesDevicePvVgConfig
var busyboxapp k8stest.FioApplication
var podNames []string

func fsVolumeSharedMountTest(decor string, engine common.OpenEbsEngine, fstype common.FileSystemType, volBindModeWait bool) {
	// FIXME: here we are using k8stest.FioApplication to use its functionality
	// and not deploying FIO, instead busybox application will be deployed.
	busyboxapp = k8stest.FioApplication{
		Decor:                          decor,
		VolSizeMb:                      4096,
		OpenEbsEngine:                  engine,
		VolType:                        common.VolFileSystem,
		FsType:                         fstype,
		Loops:                          2,
		VolWaitForFirstConsumer:        volBindModeWait,
		SkipPvcVerificationAfterCreate: true,
	}

	// Set up storage class parameters
	busyboxapp.Lvm = k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
		Shared:        common.Yes,
	}

	// Create volume
	logf.Log.Info("create volume")
	err := busyboxapp.CreateVolume()
	Expect(err).To(BeNil(), "failed to create volume")

	// Deploy first BusyBox pod and create file with MD5 checksum
	podName1 := "busybox"
	deployBusyBoxPod(podName1, busyboxapp.GetPvcName(), busyboxapp.VolType)
	podNames = append(podNames, podName1)

	filePath := "/volume/testfile.txt"
	fileContent := "This is some test data."
	md5FilePath1 := "/volume/md5sum1.txt"
	md5FilePath2 := "/volume/md5sum2.txt"
	combinedCmd1 := fmt.Sprintf(
		"echo '%s' > %s && md5sum %s > %s",
		fileContent,
		filePath,
		filePath,
		md5FilePath1,
	)

	out, _, err := k8stest.ExecuteCommandInPod(common.NSDefault, podName1, combinedCmd1)
	Expect(err).To(BeNil(), "error: %v, output: %s", err, out)

	// Deploy second BusyBox pod to verify data
	podName2 := "busybox-second"
	deployBusyBoxPod(podName2, busyboxapp.GetPvcName(), busyboxapp.VolType)
	podNames = append(podNames, podName2)

	combinedCmd2 := fmt.Sprintf(
		"md5sum %s > %s",
		filePath,
		md5FilePath2,
	)

	out, _, err = k8stest.ExecuteCommandInPod(common.NSDefault, podName2, combinedCmd2)
	Expect(err).To(BeNil(), "error: %v, output: %s", err, out)

	// Compare MD5 checksums from both pods
	md5sum1, _, err := k8stest.ExecuteCommandInPod(common.NSDefault, podName1, fmt.Sprintf("cat %s", md5FilePath1))
	Expect(err).To(BeNil(), "error %v", err)

	md5sum2, _, err := k8stest.ExecuteCommandInPod(common.NSDefault, podName2, fmt.Sprintf("cat %s", md5FilePath2))
	Expect(err).To(BeNil(), "error %v", err)

	Expect(md5sum1 == md5sum2).Should(BeTrue(), "MD5 verification failed. Data has been altered.")
	logf.Log.Info("MD5 sum verification passed")

	// Check the status of both pods
	for _, podName := range []string{podName1, podName2} {
		logf.Log.Info(fmt.Sprintf("Checking %s pod status", podName))
		phase, err := k8stest.GetPodStatusByPrefix(podName, common.NSDefault)
		Expect(err).To(BeNil(), "GetPodStatusByPrefix got error %v", err)
		Expect(phase == coreV1.PodRunning).Should(BeTrue(), fmt.Sprintf("%s pod is not in running state", podName))
	}
}

func deployBusyBoxPod(podName, pvcName string, volType common.VolumeType) *coreV1.Pod {
	args := []string{"sleep", "10000000"}
	podContainer := coreV1.Container{
		Name:            podName,
		Image:           "busybox",
		ImagePullPolicy: coreV1.PullAlways,
		Args:            args,
	}

	volume := coreV1.Volume{
		Name: "ms-volume",
		VolumeSource: coreV1.VolumeSource{
			PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
			},
		},
	}

	podObj, err := k8stest.NewPodBuilder(podName).
		WithName(podName).
		WithNamespace(common.NSDefault).
		WithRestartPolicy(coreV1.RestartPolicyNever).
		WithContainer(podContainer).
		WithVolume(volume).
		WithVolumeDeviceOrMount(volType).Build()
	Expect(err).ToNot(HaveOccurred(), "Generating pod definition, err: %v", err)
	Expect(podObj).ToNot(BeNil(), "failed to generate pod definition")

	_, err = k8stest.CreatePod(podObj, common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "Creating pod, err: %v", err)

	Eventually(func() bool {
		return k8stest.IsPodRunning(podName, common.NSDefault)
	}, k8stest.DefTimeoutSecs, "2s").Should(Equal(true))

	logf.Log.Info(fmt.Sprintf("%s pod is running.", podName))
	return podObj
}

func cleanUpResources(pods []string, pvcName string) {
	for _, pod := range pods {
		err := k8stest.DeletePod(pod, common.NSDefault)
		Expect(err).ToNot(HaveOccurred(), "failed to delete pod %s err %v", pod, err)

		// check if pod is deleted successfully
		Eventually(func() bool {
			return k8stest.IsPodRunning(pod, common.NSDefault)
		},
			k8stest.DefTimeoutSecs,
			"5s",
		).Should(Equal(false), "busybox pod deletion failed")
	}

	err := k8stest.DeletePVC(pvcName, common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to delete pvc %s err %v", pvcName, err)

	err = k8stest.RmStorageClass(pvcName)
	Expect(err).ToNot(HaveOccurred(), "Deleting storage class %s", pvcName)

	time.Sleep(10 * time.Second)
}

func TestLvmVolumeResizeTest(t *testing.T) {
	e2e_ginkgo.InitTesting(t, "lvm_shared_mount_volume", "lvm_shared_mount_volume")
}

var _ = Describe("lvm_shared_mount_volume", func() {
	BeforeEach(func() {
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Clean up after each test
		cleanUpResources(podNames, busyboxapp.GetPvcName())
		podNames = nil // Reset the pod list for the next test

		err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	It("lvm ext4: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.Ext4FsType, false)
	})
	It("lvm xfs: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.XfsFsType, false)
	})
	It("lvm btrfs: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.BtrfsFsType, false)
	})

	It("lvm ext4: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.Ext4FsType, true)
	})
	It("lvm xfs: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.XfsFsType, true)
	})
	It("lvm btrfs: should verify shared mount volume", func() {
		fsVolumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.BtrfsFsType, true)
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
	// NB This only tears down the local structures for talking to the cluster, not the kubernetes cluster itself.
	logf.Log.Info("remove node with device and vg", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
