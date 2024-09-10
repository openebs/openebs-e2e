package zfs_volume_provisioning

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/zfs"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var defFioCompletionTime = 240 // in seconds

func VolumeProvisioningTest(decor string,
	engine common.OpenEbsEngine,
	volType common.VolumeType,
	fstype common.FileSystemType,
	volBindModeWait bool,
	nodeConfig zfs.ZfsNodesDevicePoolConfig,
) (k8stest.FioApplication, error) {

	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               1024,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   1,
		VolWaitForFirstConsumer: volBindModeWait,
	}

	// setup sc parameters
	app.Zfs = k8stest.ZfsOptions{
		PoolName:      nodeConfig.PoolName,
		ThinProvision: common.No,
		RecordSize:    "128k",
		Compression:   common.Off,
		DedUp:         common.Off,
	}

	if app.FsType == common.ZfsFsType || app.FsType == common.BtrfsFsType {
		app.FsPercent = 60
	}

	logf.Log.Info("create sc, pvc, fio pod")
	err := app.DeployApplication()
	if err != nil {
		return app, err
	}

	logf.Log.Info("wait for fio pod to complete")
	err = app.WaitComplete(defFioCompletionTime)
	if err != nil {
		return app, err
	}
	// remove app pod, pvc,sc
	err = app.Cleanup()
	if err != nil {
		return app, err
	}

	return app, nil
}
