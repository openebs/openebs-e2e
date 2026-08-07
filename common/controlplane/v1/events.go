package v1

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// extractJSON strips any non-JSON prefix the plugin may write to stdout
// (e.g. "Loki is not found, continuing...") before the actual JSON array/object.
func extractJSON(raw []byte) []byte {
	for i, b := range raw {
		if b == '[' || b == '{' {
			return raw[i:]
		}
	}
	return raw
}

func (cp CPv1) GetEvents(flags ...string) ([]common.EventRecord, error) {
	args := append([]string{"get", "events", "-n", common.NSMayastor(), "-o", "json"}, flags...)
	cmd := GetMayastorPluginCmd(args...)
	jsonInput, err := cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var records []common.EventRecord
	if err := json.Unmarshal(extractJSON(jsonInput), &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events JSON: %w, output: %s", err, string(jsonInput))
	}
	return records, nil
}

// GetEventsWithStderr runs the plugin and captures stdout (JSON events) and stderr (warnings) separately.
// Non-JSON prefix lines in stdout are captured in the returned stderr string.
func (cp CPv1) GetEventsWithStderr(flags ...string) ([]common.EventRecord, string, error) {
	args := append([]string{"get", "events", "-n", common.NSMayastor(), "-o", "json"}, flags...)
	cmd := GetMayastorPluginCmd(args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stdoutBytes := stdout.Bytes()
	stderrStr := stderr.String()
	// plugin may write informational messages to stdout before JSON;
	// capture those as part of stderrStr so the caller can check them
	jsonBytes := extractJSON(stdoutBytes)
	if len(jsonBytes) < len(stdoutBytes) {
		prefix := string(stdoutBytes[:len(stdoutBytes)-len(jsonBytes)])
		stderrStr = prefix + stderrStr
	}
	if err != nil {
		return nil, stderrStr, fmt.Errorf("plugin command failed: %w, stderr: %s", err, stderrStr)
	}
	var records []common.EventRecord
	if err := json.Unmarshal(jsonBytes, &records); err != nil {
		return nil, stderrStr, fmt.Errorf("failed to unmarshal events JSON: %w, output: %s", err, stdout.String())
	}
	return records, stderrStr, nil
}

func (cp CPv1) CountEvents(flags ...string) int {
	records, err := cp.GetEvents(flags...)
	if err != nil {
		logf.Log.Info("CountEvents: failed to get events", "error", err)
		return 0
	}
	return len(records)
}
