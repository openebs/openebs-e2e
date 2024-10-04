package hostpath_ci_smoke_test

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/src/tests/hostpath/hostpath_volume_provisioning"

	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var podName, pvcName string
var err error

func TestHostpathCiSmokeTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "hostpath_ci_smoke_test", "hostpath_ci_smoke_test")
}

var _ = Describe("hostpath_ci_smoke_test", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// cleanup k8s resources if exist
		logf.Log.Info("cleanup k8s resources")
		err := k8stest.CleanUpBusyboxResources([]string{podName}, pvcName)
		Expect(err).ToNot(HaveOccurred(), "failed delete busybox resource")

		// Check resource leakage.
		after_err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(after_err).ToNot(HaveOccurred())

	})

	It("hostpath ext4: should verify a volume can be created, used and deleted", func() {
		podName, pvcName, err = hostpath_volume_provisioning.VolumeProvisioningTest("ext4", common.Ext4FsType)
		Expect(err).ToNot(HaveOccurred())
	})

})

var _ = BeforeSuite(func() {
	err = e2e_ginkgo.SetupTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnv %v", err)

})

var _ = AfterSuite(func() {
	err = k8stest.TeardownTestEnv()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)
})
