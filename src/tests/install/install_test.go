package install

import (
	"testing"

	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/e2e_ginkgo"
	"github.com/openebs/openebs-e2e/common/k8stest"
	"github.com/openebs/openebs-e2e/src/common/k8sinstall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInstall(t *testing.T) {
	e2e_config.SetContext(e2e_config.E2eTesting)

	e2e_ginkgo.InitTesting(t, "Install Test", "install")
}

var _ = Describe("Install test", func() {
	AfterEach(func() {
		// Check resource leakage.
		e2e_ginkgo.GenerateSupportBundleIfTestFailed()
	})

	It("should install using yaml files", func() {
		Expect(k8sinstall.InstallProduct()).ToNot(HaveOccurred(), "install failed")
	})
})

var _ = BeforeSuite(func() {
	err := e2e_ginkgo.SetupTestEnvBasic()
	Expect(err).ToNot(HaveOccurred(), "failed to setup test environment in BeforeSuite : SetupTestEnvBasic %v", err)

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := k8stest.TeardownTestEnvNoCleanup()
	Expect(err).ToNot(HaveOccurred(), "failed to tear down test environment in AfterSuite : TeardownTestEnv %v", err)

})
