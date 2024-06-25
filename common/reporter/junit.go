package reporter

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/openebs/openebs-e2e/common/e2e_config"
)

func GetReporters(name string) []Reporter {
	cfg := e2e_config.GetConfig()

	if cfg.ReportsDir == "" {
		return []Reporter{}
	}
	testGroupPrefix := "e2e."
	xmlFileSpec := cfg.ReportsDir + "/" + testGroupPrefix + name + "-junit.xml"
	junitReporter := reporters.NewJUnitReporter(xmlFileSpec)
	return []Reporter{junitReporter}
}
