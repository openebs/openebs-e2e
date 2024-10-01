package hostpath_volume_provisioning

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	"github.com/openebs/openebs-e2e/common/hostpath"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func VolumeProvisioningTest(decor string,
	fstype common.FileSystemType,
) (string, string, error) {

	// FIXME: here we are using k8stest.FioApplication to use its functionality
	// and not deploying FIO, instead busybox application will be deployed.
	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               1024,
		OpenEbsEngine:           common.Hostpath,
		VolType:                 common.VolFileSystem,
		FsType:                  fstype,
		VolWaitForFirstConsumer: true,
	}

	hostpathConfig, err := hostpath.GetHostPathAnnotationConfig()
	if err != nil {
		return "", "", err
	}

	// setup sc parameters
	app.HostPath = k8stest.HostPathOptions{
		Annotations: hostpathConfig,
	}

	// Create volume
	logf.Log.Info("create volume")
	err = app.CreateVolume()
	if err != nil {
		return "", "", err
	}

	// Deploy BusyBox pod and create file with MD5 checksum
	podName := "busybox"
	err = k8stest.DeployBusyBoxPod(podName, app.GetPvcName(), app.VolType)
	if err != nil {
		return "", app.GetPvcName(), err
	}

	filePath := "/volume/testfile.txt"
	fileContent := "This is some test data."
	md5FilePath1 := "/volume/md5sum1.txt"
	combinedCmd1 := fmt.Sprintf(
		"echo '%s' > %s && md5sum %s > %s",
		fileContent,
		filePath,
		filePath,
		md5FilePath1,
	)

	out, _, err := k8stest.ExecuteCommandInPod(common.NSDefault, podName, combinedCmd1)
	if err != nil {
		logf.Log.Info("Failed to execute command", "Error", err, "Output", out)
		return podName, app.GetPvcName(), err
	}
	return podName, app.GetPvcName(), nil
}
