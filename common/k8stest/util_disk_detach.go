package k8stest

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common/custom_resources"
	e2eagent "github.com/openebs/openebs-e2e/common/e2e_agent"
	"github.com/openebs/openebs-e2e/common/mayastor/disk_failures"
)

// InjectIOError injects IO error on the given pool device present on the given node.
func InjectIOError(nodeAddr string, poolDevice string) (string, error) {
	table := fmt.Sprintf(
		"0 50000 linear %s 0\n"+
			"50000 5000000 error\n"+
			"5050000 8143000 linear %s 5050000",
		poolDevice,
		poolDevice,
	)

	out, err := e2eagent.CreateFaultyDevice(nodeAddr, poolDevice, table)
	if err != nil {
		return out, fmt.Errorf(
			"failed to inject IO error on %s: %w (output=%s)",
			poolDevice,
			err,
			out,
		)
	}

	return out, nil
}

// RecoverIOError recovers from IO error on the given pool device present on the given node.
func RecoverIOError(nodeAddr string, poolDevice string) (string, error) {
	out, err := e2eagent.DeleteFaultyDevice(nodeAddr, poolDevice)
	if err != nil {
		return "", fmt.Errorf(
			"failed to recover IO error on %s: %w (output=%s)",
			poolDevice,
			err,
			out,
		)
	}

	return out, nil
}

// SetupTimeoutDevice sets up a timeout device on the given pool device present on the given node.
func SetupTimeoutDevice(nodeAddr, poolDevice string) (string, error) {
	resp, err := e2eagent.SetupTimeoutDevice(nodeAddr, poolDevice)
	if err != nil {
		return resp, fmt.Errorf(
			"setup timeout device failed on %s: %w",
			poolDevice, err,
		)
	}
	return resp, nil
}

// InjectIOTimeout injects IO timeout on the given pool device present on the given node.
func InjectIOTimeout(nodeAddr, poolDevice string) (string, error) {
	resp, err := e2eagent.InjectIOTimeout(nodeAddr, poolDevice)
	if err != nil {
		return resp, fmt.Errorf(
			"inject IO timeout failed on %s: %w",
			poolDevice, err,
		)
	}
	return resp, nil
}

// RecoverIOTimeout recovers from IO timeout on the given pool device present on the given node.
func RecoverIOTimeout(nodeAddr, poolDevice string) (string, error) {
	resp, err := e2eagent.RecoverIOTimeout(nodeAddr, poolDevice)
	if err != nil {
		return resp, fmt.Errorf(
			"recover IO timeout failed on %s: %w",
			poolDevice, err,
		)
	}
	return resp, nil
}

// CleanupTimeoutDevice cleans up the timeout device on the given pool device present on the given node.
func CleanupTimeoutDevice(nodeAddr, poolDevice string) (string, error) {
	resp, err := e2eagent.CleanupTimeoutDevice(nodeAddr, poolDevice)
	if err != nil {
		return resp, fmt.Errorf(
			"cleanup timeout device failed on %s: %w",
			poolDevice, err,
		)
	}
	return resp, nil
}

func GetPoolErrorThreshold(release, namespace, path string) (int, error) {

	val, err := GetUserHelmValueByPath(
		release,
		namespace,
		path,
	)
	if err != nil {
		return 0, err
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("unexpected type")
	}
}

func containsAlert(alerts []disk_failures.PoolAlert, expected disk_failures.PoolAlert) bool {
	for _, a := range alerts {
		if a == expected {
			return true
		}
	}
	return false
}

func PoolStateCheck(
	poolName string,
	expectedState string,
	expectedAlertStatus string,
	expectedAlert disk_failures.PoolAlert,
	alertType string,
) func() error {
	return func() error {
		poolCR, err := custom_resources.GetMsPool(poolName)
		if err != nil {
			return err
		}

		if poolCR.GetPoolStatus() != expectedState {
			return fmt.Errorf("expected state %s, got %s",
				expectedState, poolCR.GetPoolStatus())
		}

		if poolCR.GetPoolAlertStatus() != expectedAlertStatus {
			return fmt.Errorf("expected alert status %s, got %s",
				expectedAlertStatus, poolCR.GetPoolAlertStatus())
		}

		alerts := poolCR.GetAlerts()

		var alertList []disk_failures.PoolAlert
		switch alertType {
		case disk_failures.PoolAlertStatusCritical:
			alertList = alerts.Critical
		case disk_failures.PoolAlertStatusAttention:
			alertList = alerts.Attention
		case disk_failures.PoolAlertStatusWarning:
			alertList = alerts.Warning
		default:
			return fmt.Errorf("unknown alert type %s", alertType)
		}

		if !containsAlert(alertList, expectedAlert) {
			return fmt.Errorf("alert %s not found in %s list",
				expectedAlert, alertType)
		}

		return nil
	}
}
