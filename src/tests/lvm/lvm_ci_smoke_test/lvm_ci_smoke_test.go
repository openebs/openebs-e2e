package lvm_ci_smoke_test

import (
	"testing"

	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/lvm"
	"github.com/openebs/openebs-e2e/src/tests/lvm/lvm_volume_provisioning"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var err error
var nodeConfig lvm.LvmNodesDevicePvVgConfig
var app k8stest.FioApplication

func TestLvmCiSmokeTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "lvm_ci_smoke_test", "lvm_ci_smoke_test")
}

var _ = Describe("lvm_ci_smoke_test", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Check resource leakage.
		after_err := e2e_ginkgo.AfterEachK8sCheck()
		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err := app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		Expect(after_err).ToNot(HaveOccurred())
	})

	It("lvm ci smoke test: should verify a volume can be created, used and deleted", func() {
		app, err = lvm_volume_provisioning.VolumeProvisioningTest("lvm-ext4", common.Lvm, common.VolFileSystem, common.Ext4FsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
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
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	logf.Log.Info("remove node with device and vg", "node config", nodeConfig)
	err := nodeConfig.RemoveConfiguredLvmNodesWithDeviceAndVg()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")

	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
