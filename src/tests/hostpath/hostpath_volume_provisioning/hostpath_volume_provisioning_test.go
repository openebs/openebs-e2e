package hostpath_volume_provisioning

import (
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/openebs/openebs-e2e/common/k8stest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var podName, pvcName string
var err error

func TestHostpathVolumeProvisioningTest(t *testing.T) {
	// Initialise test and set class and file names for reports
	e2e_ginkgo.InitTesting(t, "hostpath_volume_provisioning", "hostpath_volume_provisioning")
}

var _ = Describe("hostpath_volume_provisioning", func() {

	BeforeEach(func() {
		// Check ready to run
		err := e2e_ginkgo.BeforeEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {

		// cleanup k8s resources
		logf.Log.Info("cleanup k8s resources ")
		err = k8stest.CleanUpBusyboxResources([]string{podName}, pvcName)
		Expect(err).ToNot(HaveOccurred(), "failed to k8s resource")

		// Check resource leakage.
		err := e2e_ginkgo.AfterEachK8sCheck()
		Expect(err).ToNot(HaveOccurred())

	})

	It("hostpath ext4: should verify a volume can be created, used and deleted", func() {
		podName, pvcName, err = VolumeProvisioningTest("ext4", common.Ext4FsType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("hostpath xfs: should verify a volume can be created, used and deleted", func() {
		podName, pvcName, err = VolumeProvisioningTest("xfs", common.XfsFsType)
		Expect(err).ToNot(HaveOccurred())
	})
	It("hostpath btrfs: should verify a volume can be created, used and deleted", func() {
		podName, pvcName, err = VolumeProvisioningTest("btrfs", common.BtrfsFsType)
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
