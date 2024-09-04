package lvm_shared_mount_volume

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/lvm"
	coreV1 "k8s.io/api/core/v1"
	imageutils "k8s.io/kubernetes/test/utils/image"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Background:
//     Given a k8s cluster is running with the product installed
// Scenario: shared mount volume creation for lvm volume
//     Given a volume has been successfully created with shared mount enabled via storage class
//     When two applications are deployed using the same volume
//     Then both applications should be in running state using the same volume

var nodeConfig lvm.LvmNodesDevicePvVgConfig

func setupLVM() {
	loopDevice := e2e_agent.LoopDevice{
		Size:   10737418240,
		ImgDir: "/tmp",
	}

	workerNodes, err := lvm.ListLvmNode(common.NSOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to list worker nodes")

	nodeConfig = lvm.LvmNodesDevicePvVgConfig{
		VgName:        "lvmvg",
		NodeDeviceMap: make(map[string]e2e_agent.LoopDevice),
	}
	for _, node := range workerNodes {
		nodeConfig.NodeDeviceMap[node] = loopDevice
	}

	logf.Log.Info("Setting up node with loop device, PV, and VG", "node config", nodeConfig)
	err = nodeConfig.ConfigureLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to setup node")
}

func volumeSharedMountTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool) {
	// Initialize FioApplication instance
	app := k8stest.FioApplication{
		Decor:                          decor,
		VolSizeMb:                      4096,
		OpenEbsEngine:                  engine,
		VolType:                        volType,
		FsType:                         fstype,
		Loops:                          2,
		VolWaitForFirstConsumer:        volBindModeWait,
		SkipPvcVerificationAfterCreate: true,
	}

	// Set up LVM configuration
	setupLVM()

	// Set up storage class parameters
	app.Lvm = k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
		Shared:        common.Yes,
	}

	// Create volume
	logf.Log.Info("create volume")
	err := app.CreateVolume()
	Expect(err).To(BeNil(), "failed to create volume")

	// Deploy first BusyBox pod and create file with MD5 checksum
	podName1 := "busybox"
	deployBusyBoxPod(podName1, app.GetPvcName(), app.VolType)

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

	out, err := executeCommandInPod(podName1, combinedCmd1)
	Expect(err).To(BeNil(), "error: %v, output: %s", err, out)

	// Deploy second BusyBox pod to verify data
	podName2 := "busybox-second"
	deployBusyBoxPod(podName2, app.GetPvcName(), app.VolType)

	combinedCmd2 := fmt.Sprintf(
		"md5sum %s > %s",
		filePath,
		md5FilePath2,
	)

	out, err = executeCommandInPod(podName2, combinedCmd2)
	Expect(err).To(BeNil(), "error: %v, output: %s", err, out)

	// Compare MD5 checksums from both pods
	md5sum1, err := executeCommandInPod(podName1, fmt.Sprintf("cat %s", md5FilePath1))
	Expect(err).To(BeNil(), "error %v", err)

	md5sum2, err := executeCommandInPod(podName2, fmt.Sprintf("cat %s", md5FilePath2))
	Expect(err).To(BeNil(), "error %v", err)

	Expect(md5sum1 == md5sum2).Should(BeTrue(), "MD5 verification failed. Data has been altered.")

	// Check the status of both pods
	for _, podName := range []string{podName1, podName2} {
		logf.Log.Info(fmt.Sprintf("Checking %s pod status", podName))
		phase, err := k8stest.GetPodStatusByPrefix(podName, common.NSDefault)
		Expect(err).To(BeNil(), "GetPodStatusByPrefix got error %v", err)
		Expect(phase == coreV1.PodRunning).Should(BeTrue(), fmt.Sprintf("%s pod is not in running state", podName))
	}

	// Clean up resources
	cleanUpResources([]string{podName1, podName2}, app.GetPvcName())
}

func deployBusyBoxPod(podName, pvcName string, volType common.VolumeType) *coreV1.Pod {
	args := []string{"sleep", "10000000"}
	podContainer := coreV1.Container{
		Name:            podName,
		Image:           imageutils.GetE2EImage(imageutils.BusyBox),
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

func executeCommandInPod(podName, cmd string) (string, error) {
	cmdArgs := []string{
		"exec",
		podName,
		"--",
		"sh",
		"-c",
		cmd,
	}

	execCmd := exec.Command("kubectl", cmdArgs...)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("Command execution failed, error: %v", err)
	}
	return string(output), nil
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
		).Should(Equal(false), "busybox pod 1 deletion failed")
	}

	err := k8stest.DeletePVC(pvcName, common.NSDefault)
	Expect(err).ToNot(HaveOccurred(), "failed to delete pvc %s err %v", pvcName, err)

	k8stest.RmStorageClass(pvcName)

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
		err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	It("lvm ext4: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.Ext4FsType, false)
	})
	It("lvm xfs: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.XfsFsType, false)
	})
	It("lvm btrfs: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.BtrfsFsType, false)
	})

	It("lvm ext4: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.Ext4FsType, true)
	})
	It("lvm xfs: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.XfsFsType, true)
	})
	It("lvm btrfs: should verify shared mount volume", func() {
		volumeSharedMountTest("lvm-volume-shared-mount", common.Lvm, common.VolFileSystem, common.BtrfsFsType, true)
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
