package lvm_volume_provisioning

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/common/lvm"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var defFioCompletionTime = 240 // in seconds

func VolumeProvisioningTest(decor string,
	engine common.OpenEbsEngine,
	volType common.VolumeType,
	fstype common.FileSystemType,
	volBindModeWait bool,
	nodeConfig lvm.LvmNodesDevicePvVgConfig,
) (k8stest.FioApplication, error) {

	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               2048,
		OpenEbsEngine:           engine,
		VolType:                 volType,
		FsType:                  fstype,
		Loops:                   1,
		VolWaitForFirstConsumer: volBindModeWait,
	}

	// setup sc parameters
	app.Lvm = k8stest.LvmOptions{
		VolGroup:      nodeConfig.VgName,
		Storage:       "lvm",
		ThinProvision: common.No,
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
