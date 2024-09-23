package common

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/mayastor/volume_resize"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/gomega"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var ResizeApp k8stest.FioApplication
var ResizeApp2 k8stest.FioApplication
var defFioCompletionTime = 240 // in seconds
var ThinPoolNode string

func LvmVolumeResizeTest(decor string, engine common.OpenEbsEngine, vgName string, volType common.VolumeType, fstype common.FileSystemType, volBindModeWait bool, thinProvisioned common.YesNoVal) {

	var ftSize1, ftSize2 uint64
	// setup sc parameters
	lvmScOptions := k8stest.LvmOptions{
		VolGroup:      vgName,
		Storage:       "lvm",
		ThinProvision: thinProvisioned,
	}

	ResizeApp = k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               1024,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   3,
		VolWaitForFirstConsumer: volBindModeWait,
		Lvm:                     lvmScOptions,
		// after fio completes sleep of a long time
		PostOpSleep:          600000,
		AllowVolumeExpansion: common.AllowVolumeExpansionEnable,
	}

	if ResizeApp.FsType == common.BtrfsFsType {
		ResizeApp.FsPercent = 60
	}

	logf.Log.Info("create sc, pvc, fio pod")
	err := ResizeApp.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	var node string
	if thinProvisioned == common.Yes {
		node, err = k8stest.GetNodeNameForScheduledPod(ResizeApp.GetPodName(), common.NSDefault)
		Expect(err).To(BeNil(), "failed to get node name for %s app", ResizeApp.GetPodName())
		nodeIp, err := k8stest.GetNodeIPAddress(node)
		Expect(err).To(BeNil(), "failed to get node %s  ip", node)
		ThinPoolNode = *nodeIp
		logf.Log.Info("App node", "name", node, "IP", ThinPoolNode)
		out, err := e2e_agent.LvmLvChangeMonitor(ThinPoolNode, vgName)
		Expect(err).To(BeNil(), "failed to set up lv change monitor on node %s with vg %s, output: %s", node, vgName, out)
	}

	// sleep for 30 seconds before resizing volume
	logf.Log.Info("Sleep before resizing volume", "duration", volume_resize.DefSleepTime)
	time.Sleep(time.Duration(volume_resize.DefSleepTime) * time.Second)

	expandedVolumeSizeMb := ResizeApp.VolSizeMb + 1024
	// expand volume by editing pvc size
	logf.Log.Info("Update volume size", "new size in MiB", expandedVolumeSizeMb)
	_, err = k8stest.UpdatePvcSize(ResizeApp.GetPvcName(), common.NSDefault, expandedVolumeSizeMb)
	Expect(err).ToNot(HaveOccurred(), "failed to expand volume %s, error: %v", ResizeApp.GetPvcName(), err)

	// verify pvc capacity to new size
	logf.Log.Info("Verify pvc resize status")
	pvcResizeStatus, err := volume_resize.WaitForPvcResize(ResizeApp.GetPvcName(), common.NSDefault, expandedVolumeSizeMb)
	Expect(err).ToNot(HaveOccurred(), "failed to verify resized pvc %s, error: %v", ResizeApp.GetPvcName(), err)
	Expect(pvcResizeStatus).To(BeTrue(), "failed to resized pvc %s, error: %v", ResizeApp.GetPvcName(), err)

	// Check fio pod status
	logf.Log.Info("Check fio pod status")
	phase, _, err := k8stest.CheckFioPodCompleted(ResizeApp.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s", phase)

	// wait for fio completion - monitoring log output
	exitValue, fErr := ResizeApp.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)
	Expect(exitValue == 0).Should(BeTrue(), "fio exit value is not 0")

	// print fio target sizes retrieved by monitoring log output
	ftSizes, ffErr := ResizeApp.FioTargetSizes()
	Expect(ffErr).ToNot(HaveOccurred())
	for path, size := range ftSizes {
		logf.Log.Info("ftSize (poc_resize_1)", "path", path, "size", volume_resize.ByteSizeString(size), "bytes", size)
		ftSize1 = size
	}
	Expect(len(ftSizes)).To(BeNumerically("==", 1), "unexpected fio target sizes")

	// second instance of e2e-fio, volume parameters should be the same as the 1st app instance
	ResizeApp2 = k8stest.FioApplication{
		Decor:     fmt.Sprintf("%s-2", ResizeApp.Decor),
		Loops:     2,
		FsPercent: ResizeApp.FsPercent,
	}

	// Before deploying the 2nd app instance - "import" the volume from the first app
	err = ResizeApp2.ImportVolumeFromApp(&ResizeApp)
	Expect(err).ToNot(HaveOccurred(), "import volume failed")

	// then deploy second fio app to use resized volume
	err = ResizeApp2.DeployApplication()
	Expect(err).ToNot(HaveOccurred(), "deploy 2nd app failed")

	exitValue, fErr = ResizeApp2.WaitFioComplete(volume_resize.DefFioCompletionTime, 5)
	Expect(fErr).ToNot(HaveOccurred())
	logf.Log.Info("fio complete", "exit value", exitValue)
	Expect(exitValue == 0).Should(BeTrue(), "fio exit value is not 0")

	ftSizes, ffErr = ResizeApp2.FioTargetSizes()
	Expect(ffErr).ToNot(HaveOccurred())
	for path, size := range ftSizes {
		logf.Log.Info("ftSize (poc_resize_2)", "path", path, "size", volume_resize.ByteSizeString(size), "bytes", size)
		ftSize2 = size
	}
	Expect(len(ftSizes)).To(BeNumerically("==", 1), "unexpected fio target sizes")

	logf.Log.Info("fio target sizes in bytes", "e2e fio 1", ftSize1, "e2e fio 2", ftSize2)

	Expect(ftSize2).To(BeNumerically(">", ftSize1))

	// second app should complete normally
	err = ResizeApp2.WaitComplete(defFioCompletionTime)
	Expect(err).ToNot(HaveOccurred(), "ResizeApp2 did not complete")

	// cleanup the second instance of e2e-fio app
	err = ResizeApp2.Cleanup()
	Expect(err).ToNot(HaveOccurred(), "app2 cleanup failed")

	// cleanup the first instance of e2e-fio app
	err = ResizeApp.Cleanup()
	Expect(err).ToNot(HaveOccurred(), "app1 cleanup failed")

	if thinProvisioned == common.Yes {
		out, err := e2e_agent.LvmLvRemoveThinPool(ThinPoolNode, vgName)
		Expect(err).To(BeNil(), "failed to remove lv thin pool on node %s with vg %s, output: %s", node, vgName, out)
		ThinPoolNode = ""
	}
}
