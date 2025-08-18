package dsp_cluster_size

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func SetMayastorDspClusterSize(size, helmChart, helmRelease, helmVersion string) error {
	var values map[string]interface{}
	prefix := "agents.core.poolClusterSize"
	if e2e_config.GetConfig().Product.UseUmbrellaOpenEBSChart {
		prefix = fmt.Sprintf("%s.%s", e2e_config.GetConfig().Product.ChartName, prefix)
	}

	logf.Log.Info("Setting DSP cluster size for Mayastor",
		"size", size,
		"helmChart", helmChart,
		"helmRelease", helmRelease,
		"helmVersion", helmVersion,
		"prefix", prefix)
	values = map[string]interface{}{
		prefix: size,
	}

	_, err := k8stest.UpgradeHelmChart(helmChart,
		common.NSMayastor(),
		helmRelease,
		helmVersion,
		values,
	)
	if err != nil {
		return fmt.Errorf("failed to set dsp cluster size via helm, error: %v", err)
	}

	return nil
}
