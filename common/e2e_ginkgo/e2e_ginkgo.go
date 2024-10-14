package e2e_ginkgo

import (
	"fmt"
	"os"
	"testing"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/event"
	"github.com/openebs/openebs-e2e/common/k8stest"

	"github.com/openebs/openebs-e2e/common/loki"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"gopkg.in/yaml.v3"

	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// flag that records that the test suite has failed test cases
var haveFailedTestCases = false

// InitTesting initialise testing and setup class name + report filename.
func InitTesting(t *testing.T, classname string, reportname string) {
	// Set the product environment variable if not set and the cluster
	// has one of the products installed.
	// panic otherwise
	if _, ok := os.LookupEnv("e2e_product"); !ok {
		product, err := k8stest.DiscoverProduct()
		if err == nil {
			err = os.Setenv("e2e_product", product)
		}
		if err != nil {
			panic(err)
		}
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, classname)
	loki.SendLokiMarker("Start of test " + classname)
}

func SetupTestEnvBasic() error {
	e2e_config.SetContext(e2e_config.E2eTesting)
	log.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(ginkgo.GinkgoWriter)))
	fmt.Printf("Mayastor namespace is \"%s\"\n", common.NSMayastor())

	ginkgo.By("bootstrapping test environment")
	return k8stest.SetupK8sEnvBasic()
}

func SetupTestEnv() error {
	err := SetupTestEnvBasic()

	if err != nil {
		return fmt.Errorf("failed to setup test environment : SetupTestEnvBasic: %v", err)
	}
	return k8stest.SetupK8sEnv()
}

var resourceCheckError error

// BeforeEachCheck asserts that the state of mayastor resources is fit for the test to run
func BeforeEachCheck() error {
	testDesc := ginkgo.CurrentSpecReport()
	common.SetTestCaseLogsPath(testDesc.FullText())

	log.Log.Info("BeforeEachCheck",
		"FailQuick", e2e_config.GetConfig().FailQuick,
		"Restart", e2e_config.GetConfig().BeforeEachCheckAndRestart,
		"resourceCheckError", resourceCheckError,
	)

	if e2e_config.GetConfig().FailQuick && resourceCheckError != nil {
		return fmt.Errorf("FailQuick: prior ResourceCheck failed")
	}

	if e2e_config.GetConfig().FailQuick && haveFailedTestCases {
		return fmt.Errorf("FailQuick: a prior testcase has failed")
	}

	if e2e_config.GetConfig().BeforeEachCheckAndRestart {
		if resourceCheckError == nil {
			// no previous failure, check resources
			resourceCheckError = k8stest.ResourceCheck(false)
		}

		if resourceCheckError != nil {
			// previous failure or resource check failed, restart
			log.Log.Info("BeforeEachCheck: restarting Mayastor", "err", resourceCheckError)
			_ = k8stest.RestartMayastor(120, 120, 120)
			_ = k8stest.RestoreConfiguredPools()
			log.Log.Info("BeforeEachCheck: restart complete")
		} else {
			// resource check succeeded
			return resourceCheckError
		}
	}

	if resourceCheckError = k8stest.ResourceCheck(false); resourceCheckError != nil {
		log.Log.Info("BeforeEachCheck failed", "error", resourceCheckError)
		resourceCheckError = fmt.Errorf("%w; not running test case, k8s cluster is not \"clean\"!!! ", resourceCheckError)
	} else {
		podNames, err := k8stest.ListRunningMayastorPods(nil)
		if err != nil {
			err = fmt.Errorf("%w; not running test case, not able to get pod list", err)
			return err
		}
		k8stest.SetMayastorInitialPodCount(len(podNames))
	}

	return resourceCheckError
}

func BeforeEachK8sCheck() error {

	testDesc := ginkgo.CurrentSpecReport()
	common.SetTestCaseLogsPath(testDesc.FullText())

	log.Log.Info("BeforeEachCheck",
		"FailQuick", e2e_config.GetConfig().FailQuick,
		"Restart", e2e_config.GetConfig().BeforeEachCheckAndRestart,
		"resourceCheckError", resourceCheckError,
	)

	if e2e_config.GetConfig().FailQuick && resourceCheckError != nil {
		return fmt.Errorf("FailQuick: prior ResourceCheck failed")
	}

	if e2e_config.GetConfig().FailQuick && haveFailedTestCases {
		return fmt.Errorf("FailQuick: a prior testcase has failed")
	}

	if e2e_config.GetConfig().BeforeEachCheckAndRestart {
		if resourceCheckError == nil {
			// no previous failure, check resources
			resourceCheckError = k8stest.ResourceK8sCheck()
		}
	}
	return resourceCheckError
}

// GenerateSupportBundleIfTestFailed Support function to generate a support bundle if the testcase has failed
// This function exists for tests which are disruptive and need more control when performing
// checks after a test case has completed
func GenerateSupportBundleIfTestFailed() {
	testDesc := ginkgo.CurrentSpecReport()
	if testDesc.Failed() {
		logsPath, err := common.GetTestCaseLogsPath()
		if err == nil {
			k8stest.GenerateSupportBundle(logsPath)
		} else {
			log.Log.Info("test case logs path was not set")
		}
	}
}

func afterEachCheckResources(canGenSupportBundle bool) error {
	// resourceCheckError is set if BeforeEachCheck fails
	// so test case starting conditions are invalid - do nothing.
	if e2e_config.GetConfig().FailQuick && resourceCheckError != nil {
		return fmt.Errorf("prior ResourceCheck failed")
	}
	var waitForPools bool
	if !ginkgo.CurrentSpecReport().Failed() {
		waitForPools = true
	}
	err := k8stest.ResourceCheck(waitForPools)
	if err != nil && canGenSupportBundle {
		haveFailedTestCases = true
		logsPath, err := common.GetTestCaseLogsPath()
		if err == nil {
			k8stest.GenerateSupportBundle(logsPath + "-AfterEach")
		} else {
			log.Log.Info("test case logs path was not set")
		}
	}
	eventMessages, err := event.GetAllEventMessagesSymbolic()
	if err != nil {
		log.Log.Info("failed to retrieve events", "error", err)
	} else {
		yml, err := yaml.Marshal(eventMessages)
		if err != nil {
			log.Log.Info("failed to deserialise events", "error", err)
		} else {
			logsPath, err := common.GetTestCaseLogsPath()
			if err != nil {
				log.Log.Info("failed to retrieve logs path", "error", err)
			} else {
				err = os.MkdirAll(logsPath, 0755)
				if err != nil {
					log.Log.Info("Failed to create path", logsPath, err)
				} else {
					err = os.WriteFile(logsPath+"/events.yml", yml, 0644)
					if err != nil {
						log.Log.Info("failed to write", logsPath+"/events.yml", err)
					} else {
						log.Log.Info("events collected", "at", logsPath+"/events.yml")
					}
				}
			}
		}
	}
	common.ResetTestCaseLogsPath()
	return err
}

// AfterEachCheckResourceOnly asserts that the state of mayastor resources has been restored.
// This function exists for tests which are disruptive and need more control when performing
// checks after a test case has completed
func AfterEachCheckResourceOnly() error {
	return afterEachCheckResources(true)
}

// AfterEachCheck asserts that the state of mayastor resources has been restored.
// and generates a support bundle if the test failed.
func AfterEachCheck() error {
	canGenSupportBundle := true
	log.Log.Info("AfterEachCheck")
	testDesc := ginkgo.CurrentSpecReport()
	if testDesc.Failed() {
		haveFailedTestCases = true
		logsPath, err := common.GetTestCaseLogsPath()
		if err == nil {
			k8stest.GenerateSupportBundle(logsPath)
			canGenSupportBundle = false
		} else {
			log.Log.Info("test case logs path was not set!")
		}
	}

	// it is very likely that if the test case failed, resource check will fail,
	// suppress generation of another support bundle since we've just generated
	// the more useful support bundle.
	return afterEachCheckResources(canGenSupportBundle)
}

func AfterEachK8sCheck() error {
	return k8stest.ResourceK8sCheck()
}
