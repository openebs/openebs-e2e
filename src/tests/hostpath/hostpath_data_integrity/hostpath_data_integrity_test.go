package hostpath_data_integrity

import (
	"testing"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/hostpath"

	"github.com/openebs/openebs-e2e/common/k8stest"
	coreV1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHostpathDataIntegrityTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "hostpath_data_integrity", "hostpath_data_integrity")
}
func dataIntegrityTestTest(decor string, fstype common.FileSystemType) {
	app := k8stest.FioApplication{
		Decor:                   decor,
		VolSizeMb:               1024,
		OpenEbsEngine:           common.Hostpath,
		VolType:                 common.VolFileSystem,
		FsType:                  fstype,
		Loops:                   1,
		VolWaitForFirstConsumer: true,
		AddFioArgs:              []string{"--size=500m", "--status-interval=30"},
	}

	hostpathConfig, err := hostpath.GetHostPathAnnotationConfig()
	Expect(err).To(BeNil(), "failed to get hostpath annotation config")
	Expect(hostpathConfig).ToNot(BeNil(), "hostpath annotation config should not be nil")

	// setup sc parameters
	app.HostPath = k8stest.HostPathOptions{
		Annotations: hostpathConfig,
	}
	if app.FsType == common.BtrfsFsType {
		app.FsPercent = 60
	}
	// create sc, pvc, fio pod
	logf.Log.Info("create sc, pvc, fio pod")
	err = app.DeployApplication()
	Expect(err).To(BeNil(), "failed to deploy app")

	// Check fio pod status
	phase, podLogSysnopsis, err := k8stest.CheckFioPodCompleted(app.GetPodName(), common.NSDefault)
	Expect(err).To(BeNil(), "CheckPodComplete got error %s", err)
	Expect(phase).ShouldNot(Equal(coreV1.PodFailed), "fio pod phase is %s, %s", phase, podLogSysnopsis)

	logf.Log.Info("Waiting for fio to complete", "timeoutSecs", 240)
	err = app.WaitComplete(240)
	Expect(err).ToNot(HaveOccurred(), "failed to wait for fio to complete")

	// FIXME: Without sleep , PV deletion is getting stuck. Need to debug this issue further
	logf.Log.Info("sleep for 30 seconds")
	time.Sleep(30 * time.Second)

	// remove app pod, pvc,sc
	err = app.Cleanup()
	Expect(err).To(BeNil(), "failed to clean resources")
}

var _ = Describe("hostpath_data_integrity", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		after_err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(after_err).ToNot(HaveOccurred())

	})

	It("hostpath ext4: should verify a volume can be created, used and deleted", func() {
		dataIntegrityTestTest("ext4", common.Ext4FsType)
	})
	It("hostpath xfs: should verify a volume can be created, used and deleted", func() {
		dataIntegrityTestTest("xfs", common.XfsFsType)
	})
	It("hostpath btrfs: should verify a volume can be created, used and deleted", func() {
		dataIntegrityTestTest("btrfs", common.BtrfsFsType)
	})

})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

})

var _ = AfterSuite(func() {
	err := k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
