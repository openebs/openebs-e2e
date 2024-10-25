package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openebs/openebs-e2e/common/e2e_config"
)

// NSMayastor return the name of the namespace in which Mayastor is installed
func NSMayastor() string {
	return e2e_config.GetConfig().Product.ProductNamespace
}

// NSMayastor return the name of the namespace in which Mayastor is installed
func NSOpenEBS() string {
	return e2e_config.GetConfig().Product.OpenEBSProductNamespace
}

// default fio arguments for E2E fio runs
var fioArgs = []string{
	"--name=benchtest",
	"--numjobs=1",
}

var fioParams = []string{
	"--direct=1",
	"--random_generator=tausworthe64",
	"--rw=randrw",
	"--ioengine=libaio",
	"--bs=4k",
	"--iodepth=16",
	"--verify_fatal=1",
	"--verify=crc32",
	"--verify_async=2",
}

// GetFioArgs return the default command line for fio - for use with Mayastor,
// for single volume
func GetFioArgs() []string {
	return append(fioArgs, fioParams...)
}

// GetDefaultFioArguments return the default settings (arguments) for fio - for use with Mayastor
func GetDefaultFioArguments() []string {
	return fioParams
}

func GetFioImage() string {
	return e2e_config.GetConfig().E2eFioImage
}

func DefaultReplicaCount() int {
	return e2e_config.GetConfig().DefaultReplicaCount
}

// SanitizePathname map helper function for runes used to create directories
// only allow A-Z, a-z, 0-9 and replace ' ' with '_'
func SanitizePathname(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r
	case r >= 'a' && r <= 'z':
		return r
	case r >= '0' && r <= '9':
		return r
	case r == ' ':
		return '_'
	}
	return -1
}

var testcaseLogsPath string
var currentTestCase string

// SetTestCaseLogsPath call at the start of test case
func SetTestCaseLogsPath(testcase string) {
	logRoot, ok := os.LookupEnv("e2etestlogdir")
	if !ok {
		logRoot = "/tmp/e2e/logs"
	}
	t0 := time.Now().UTC()
	ts := fmt.Sprintf("%v%02d%02d%v%v%v", t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), t0.Second())
	testcaseLogsPath = fmt.Sprintf("%s/%s/%s", logRoot, strings.Map(SanitizePathname, testcase), ts)
	currentTestCase = testcase
}

// GetTestCaseLogsPath get the path to the logs directory for the current test case instance
func GetTestCaseLogsPath() (string, error) {
	if currentTestCase == "" {
		return "", fmt.Errorf("test case has not been set")
	}
	return testcaseLogsPath, nil
}

// ResetTestCaseLogsPath  call when the test case has completed to clear state
func ResetTestCaseLogsPath() {
	currentTestCase = ""
	testcaseLogsPath = ""
}

// GetTestSuiteLogsPath get the path to the logs directory for the current test suite
func GetTestSuiteLogsPath(testsuite string) (string, error) {
	if len(testsuite) > 1 {
		logRoot, ok := os.LookupEnv("e2etestlogdir")
		if !ok {
			logRoot = "/tmp/e2e/logs"
		}
		t0 := time.Now().UTC()
		ts := fmt.Sprintf("%v%02d%02d%v%v%v", t0.Year(), t0.Month(), t0.Day(), t0.Hour(), t0.Minute(), t0.Second())
		tsLogsPath := fmt.Sprintf("%s/%s/%s", logRoot, strings.Map(SanitizePathname, testsuite), ts)
		return tsLogsPath, nil
	} else {
		return "", fmt.Errorf("zero length testsuite name")
	}
}
