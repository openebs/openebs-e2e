package zfs_ci_smoke_test

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/zfs"
	"github.com/openebs/openebs-e2e/src/tests/zfs/zfs_volume_provisioning"

	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var nodeConfig zfs.ZfsNodesDevicePoolConfig
var app k8stest.FioApplication
var err error

func TestZfsCiSmokeTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "zfs_ci_smoke_test", "zfs_ci_smoke_test")
}

var _ = Describe("zfs_ci_smoke_test", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err = app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		// Check resource leakage.
		after_err := e2e_ginkgo.AfterEachK8sCheck()

		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources if exist")
		err := app.Cleanup()
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		Expect(after_err).ToNot(HaveOccurred())

	})

	It("zfs smoke test: should verify a volume can be created, used and deleted", func() {
		app, err = zfs_volume_provisioning.VolumeProvisioningTest("ext4", common.Zfs, common.VolFileSystem, common.Ext4FsType, false, nodeConfig)
		Expect(err).ToNot(HaveOccurred())
	})

})

var _ = BeforeSuite(func() {
	err = e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

	//setup nodes with zfs pool
	nodeConfig, err = zfs.SetupZfsNodes("zfspv-pool", 10737418240)
	Expect(err).ToNot(HaveOccurred(), "failed to setup zfs pool")

})

var _ = AfterSuite(func() {
	logf.Log.Info("remove node with device and zpool", "node config", nodeConfig)
	err = nodeConfig.RemoveConfiguredDeviceZfsPool()
	Expect(err).ToNot(HaveOccurred(), "failed to cleanup node with device")
	// NB This only tears down the local structures for talking to the cluster,
	// not the kubernetes cluster itself.	By("tearing down the test environment")
	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
