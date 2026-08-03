package v1

import (
	"encoding/json"
	"fmt"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (cp CPv1) GetEvents(flags ...string) ([]common.EventRecord, error) {
	args := append([]string{"get", "events", "-n", common.NSMayastor(), "-o", "json"}, flags...)
	cmd := GetMayastorPluginCmd(args...)
	jsonInput, err := cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var records []common.EventRecord
	if err := json.Unmarshal(jsonInput, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events JSON: %w, output: %s", err, string(jsonInput))
	}
	return records, nil
}

func (cp CPv1) CountEvents(flags ...string) int {
	records, err := cp.GetEvents(flags...)
	if err != nil {
		logf.Log.Info("CountEvents: failed to get events", "error", err)
		return 0
	}
	return len(records)
}
