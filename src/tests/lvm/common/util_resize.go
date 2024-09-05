package common

import (
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/lvm"
	"github.com/openebs/openebs-e2e/common/mayastor/volume_resize"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var NodeConfig lvm.LvmNodesDevicePvVgConfig

func LvmVolumeResizeTest(decor string, engine common.OpenEbsEngine, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool, thinProvisioned common.YesNoVal) {

	var ftSize1, ftSize2 uint64
	app := k8stest.FioApplication{
		Decor:                          decor,
		VolSizeMb:                      4096,
		OpenEbsEngine:                  engine,
		VolType:                        volType,
		FsType:                         fstype,
		Loops:                          5,
		VolWaitForFirstConsumer:        volBindModeWait,
		SkipPvcVerificationAfterCreate: true,
		// after fio completes sleep of a long time
		PostOpSleep: 600000,
	}

	loopDevice := e2e_agent.LoopDevice{
		Size:   10737418240,
		ImgDir: "/tmp",
	}

	workerNodes, err := lvm.ListLvmNode(common.NSOpenEBS())
	Expect(err).ToNot(HaveOccurred(), "failed to list worker node")

	NodeConfig = lvm.LvmNodesDevicePvVgConfig{
		VgName:        "lvmvg",
		NodeDeviceMap: make(map[string]e2e_agent.LoopDevice), // Properly initialize the map
	}
	for _, node := range workerNodes {
		NodeConfig.NodeDeviceMap[node] = loopDevice
	}

	logf.Log.Info("setup node with loop device, pv and vg", "node config", NodeConfig)
	err = NodeConfig.ConfigureLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to setup node")

	// setup sc parameters
	app.Lvm = k8stest.LvmOptions{
		VolGroup:      NodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: thinProvisioned,
	}

	logf.Log.Info("create sc, pvc, fio pod")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	// sleep for 30 seconds before resizing volume
	logf.Log.Info("Sleep before resizing volume", "duration", volume_resize.DefSleepTime)
	time.Sleep(time.Duration(volume_resize.DefSleepTime) * time.Second)

	expandedVolumeSizeMb := app.VolSizeMb + 1024
	// expand volume by editing pvc size
	logf.Log.Info("Update volume size", "new size in MiB", expandedVolumeSizeMb)
	_, err = k8stest.UpdatePvcSize(app.GetPvcName(), common.NSDefault, expandedVolumeSizeMb)
	Expect(err).ToNot(HaveOccurred(), "failed to expand volume %s, error: %v", app.GetPvcName(), err)

	// verify pvc capacity to new size
	logf.Log.Info("Verify pvc resize status")
	pvcResizeStatus, err := volume_resize.WaitForPvcResize(app.GetPvcName(), common.NSDefault, expandedVolumeSizeMb)
	Expect(err).ToNot(HaveOccurred(), "failed to verify resized pvc %s, error: %v", app.GetPvcName(), err)
	Expect(pvcResizeStatus).To(BeTrue(), "failed to resized pvc %s, error: %v", app.GetPvcName(), err)

	// Check fio pod status
	logf.Log.Info("Check fio pod status")
	phase, _, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s", phase)

	// wait for fio completion - monitoring log output
	exitValue, fErr := app.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)
	Expect(exitValue == 0).Should(BeTrue(), "fio exit value is not 0")

	// print fio target sizes retrieved by monitoring log output
	ftSizes, ffErr := app.FioTargetSizes()
	Expect(ffErr).ToNot(HaveOccurred())
	for path, size := range ftSizes {
		logf.Log.Info("ftSize (poc_resize_1)", "path", path, "size", volume_resize.ByteSizeString(size), "bytes", size)
		ftSize1 = size
	}
	Expect(len(ftSizes)).To(BeNumerically("==", 1), "unexpected fio target sizes")

	// second instance of e2e-fio, volume parameters should be the same as the 1st app instance
	app2 := app
	app2.Decor = app.Decor + "second-app"

	// Before deploying the 2nd app instance - "import" the volume
	// from the first app
	err = app2.ImportVolumeFromApp(&app)
	Expect(err).ToNot(HaveOccurred(), "import volume failed")
	// then deploy
	err = app2.DeployApplication()
	Expect(err).ToNot(HaveOccurred(), "deploy 2nd app failed")

	exitValue, fErr = app2.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)
	Expect(exitValue == 0).Should(BeTrue(), "fio exit value is not 0")

	ftSizes, ffErr = app2.FioTargetSizes()
	Expect(ffErr).ToNot(HaveOccurred())
	for path, size := range ftSizes {
		logf.Log.Info("ftSize (poc_resize_2)", "path", path, "size", volume_resize.ByteSizeString(size), "bytes", size)
		ftSize2 = size
	}
	Expect(len(ftSizes)).To(BeNumerically("==", 1), "unexpected fio target sizes")

	logf.Log.Info("fio target sizes in bytes", "e2e fio 1", ftSize1, "e2e fio 2", ftSize2)

	Expect(ftSize2).To(BeNumerically(">", ftSize1))

	// second app should complete normally
	err = app2.WaitComplete(volume_resize.DefFioCompletionTime)
	Expect(err).ToNot(HaveOccurred(), "app2 did not complete")

	// cleanup the second instance of e2e-fio app
	err = app2.Cleanup()
	Expect(err).ToNot(HaveOccurred(), "app2 cleanup failed")

	// cleanup the first instance of e2e-fio app
	err = app.Cleanup()
	Expect(err).ToNot(HaveOccurred(), "app1 cleanup failed")

}
