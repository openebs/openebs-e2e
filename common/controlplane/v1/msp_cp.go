package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type MayastorCpPool struct {
	Id    string   `json:"id"`
	Spec  mspSpec  `json:"spec"`
	State mspState `json:"state"`
}

type mspSpec struct {
	Disks       []string          `json:"disks"`
	Id          string            `json:"id"`
	Labels      map[string]string `json:"labels"`
	Node        string            `json:"node"`
	Status      string            `json:"status"`
	CordonDrain *CordonDrainSpec  `json:"cordonDrain,omitempty"`
}

// CordonDrainSpec represents the cordon drain specification structure
type CordonDrainSpec struct {
	Cordoned *PoolCordonedState `json:"cordoned,omitempty"`
}

// PoolCordonedState represents the cordoned state with constraints for pools
type PoolCordonedState struct {
	Replicas  bool `json:"replicas"`
	Snapshots bool `json:"snapshots"`
	Restores  bool `json:"restores"`
	Import    bool `json:"import"`
}

// IsCordoned returns true if any constraint is enabled
func (c *PoolCordonedState) IsCordoned() bool {
	return c.Replicas || c.Snapshots || c.Restores || c.Import
}

type mspState struct {
	Capacity  uint64   `json:"capacity"`
	Disks     []string `json:"disks"`
	ID        string   `json:"id"`
	Node      string   `json:"node"`
	Status    string   `json:"status"`
	Used      uint64   `json:"used"`
	Committed uint64   `json:"committed"`
	// Additional fields exposed by plugin JSON
	DiskCapacityBytes  uint64 `json:"diskCapacity"`
	MaxExpandableBytes uint64 `json:"maxExpandableSize"`
}

func (cp CPv1) CreatePoolOnInstall() bool {
	no_pool_install := os.Getenv("no_pool_install")
	return no_pool_install != "true"
}

func GetMayastorCpPool(name string) (*MayastorCpPool, error) {
	var jsonInput []byte
	var err error
	cmd := GetMayastorPluginCmd("-n", common.NSMayastor(), "-ojson", "get", "pool", name)
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}

	var response MayastorCpPool
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		msg := string(jsonInput)
		if !HasNotFoundRestJsonError(msg) {
			logf.Log.Info("Failed to unmarshal (get pool)", "string", msg)
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return &response, nil
}

// Expose disk capacity and max expandable size in bytes for tests
func (cp CPv1) GetPoolDiskCapacityAndMaxExpandable(name string) (uint64, uint64, error) {
	p, err := GetMayastorCpPool(name)
	if err != nil {
		return 0, 0, err
	}
	return p.State.DiskCapacityBytes, p.State.MaxExpandableBytes, nil
}

func ListMayastorCpPools() ([]MayastorCpPool, error) {
	var jsonInput []byte
	var err error
	cmd := GetMayastorPluginCmd("-n", common.NSMayastor(), "-ojson", "get", "pools")
	jsonInput, err = cmd.CombinedOutput()
	err = CheckPluginError(jsonInput, err)
	if err != nil {
		return nil, err
	}
	var response []MayastorCpPool
	err = json.Unmarshal(jsonInput, &response)
	if err != nil {
		errMsg := string(jsonInput)
		logf.Log.Info("Failed to unmarshal (get pools)", "string", string(jsonInput))
		return []MayastorCpPool{}, fmt.Errorf("%s", errMsg)
	}
	return response, nil
}

func cpMspToMsp(cpMsp *MayastorCpPool) common.MayastorPool {
	return common.MayastorPool{
		Name: cpMsp.Id,
		Spec: common.MayastorPoolSpec{
			Node:  cpMsp.Spec.Node,
			Disks: cpMsp.Spec.Disks,
		},
		Status: common.MayastorPoolStatus{
			Capacity:  cpMsp.State.Capacity,
			Used:      cpMsp.State.Used,
			Committed: cpMsp.State.Committed,
			Disks:     cpMsp.State.Disks,
			Spec: common.MayastorPoolSpec{
				Disks: cpMsp.Spec.Disks,
				Node:  cpMsp.Spec.Node,
			},
			State:  cpMsp.State.Status,
			Avail:  cpMsp.State.Capacity - cpMsp.State.Used,
			Reason: "",
		},
	}
}

// GetMsPool Get pointer to a mayastor control plane pool
func (cp CPv1) GetMsPool(poolName string) (*common.MayastorPool, error) {
	cpMsp, err := GetMayastorCpPool(poolName)
	if err != nil {
		return nil, fmt.Errorf("GetMsPool: %v", err)
	}

	if cpMsp == nil {
		logf.Log.Info("Msp not found", "pool", poolName)
		return nil, nil
	}

	msp := cpMspToMsp(cpMsp)
	return &msp, nil
}

func (cp CPv1) ListMsPools() ([]common.MayastorPool, error) {
	var msps []common.MayastorPool
	list, err := ListMayastorCpPools()
	if err == nil {
		for _, item := range list {
			msps = append(msps, cpMspToMsp(&item))
		}
	}
	return msps, err
}

// CordonPool cordons a pool with the specified constraints
func (cp CPv1) CordonPool(poolID string, constraints ...common.PoolCordonConstraint) error {
	constraintStrings := make([]string, 0, len(constraints))
	for _, c := range constraints {
		if s := c.String(); s != "" {
			constraintStrings = append(constraintStrings, s)
		}
	}
	logf.Log.Info("Executing cordon pool command", "pool", poolID, "constraints", constraintStrings)

	args := []string{"-n", common.NSMayastor(), "cordon", "pool"}

	// Add constraint flags if specified
	for _, constraint := range constraints {
		switch constraint {
		case common.CordonReplicas:
			args = append(args, "--replicas")
		case common.CordonSnapshots:
			args = append(args, "--snapshots")
		case common.CordonRestores:
			args = append(args, "--restores")
		case common.CordonImport:
			args = append(args, "--import")
		}
	}

	args = append(args, poolID)

	// Log the actual command being executed
	logf.Log.Info("Executing command", "command", "mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to cordon pool %s with constraints %v, error %v, output: %s", poolID, constraintStrings, err, out.String())
	}

	logf.Log.Info("Successfully cordoned pool", "pool", poolID, "constraints", constraintStrings, "output", out.String())
	return nil
}

// UnCordonPool uncordons a pool by removing the specified constraints
func (cp CPv1) UnCordonPool(poolID string, constraints ...common.PoolCordonConstraint) error {
	constraintStrings := make([]string, 0, len(constraints))
	for _, c := range constraints {
		if s := c.String(); s != "" {
			constraintStrings = append(constraintStrings, s)
		}
	}
	logf.Log.Info("Executing uncordon pool command", "pool", poolID, "constraints", constraintStrings)

	args := []string{"-n", common.NSMayastor(), "uncordon", "pool"}

	// Add constraint flags if specified
	for _, constraint := range constraints {
		switch constraint {
		case common.CordonReplicas:
			args = append(args, "--replicas")
		case common.CordonSnapshots:
			args = append(args, "--snapshots")
		case common.CordonRestores:
			args = append(args, "--restores")
		case common.CordonImport:
			args = append(args, "--import")
		}
	}

	args = append(args, poolID)

	// Log the actual command being executed
	logf.Log.Info("Executing command", "command", "mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to uncordon pool %s with constraints %v, error %v, output: %s", poolID, constraintStrings, err, out.String())
	}

	logf.Log.Info("Successfully uncordoned pool", "pool", poolID, "constraints", constraintStrings, "output", out.String())
	return nil
}

// GetPoolCordonStatus gets the current cordon status of a pool
func (cp CPv1) GetPoolCordonStatus(poolID string) (*PoolCordonStatus, error) {
	logf.Log.Info("Getting cordon status for pool", "pool", poolID)

	args := []string{"-n", common.NSMayastor(), "-ojson", "get", "pool", poolID}

	// Log the actual command being executed
	logf.Log.Info("Executing command", "command", "mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get cordon status for pool %s, error %v", poolID, err)
	}

	outputString := out.String()
	var poolInfo MayastorCpPool
	err = json.Unmarshal([]byte(outputString), &poolInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal command output for pool %s, error %v", outputString, err)
	}

	status := &PoolCordonStatus{
		PoolID:     poolID,
		IsCordoned: poolInfo.Spec.CordonDrain != nil && poolInfo.Spec.CordonDrain.Cordoned != nil && poolInfo.Spec.CordonDrain.Cordoned.IsCordoned(),
	}

	if poolInfo.Spec.CordonDrain != nil && poolInfo.Spec.CordonDrain.Cordoned != nil {
		// Parse the cordon constraints from the YAML structure
		// The cordonDrain field contains the constraint information
		status.Constraints = parseCordonConstraints(poolInfo.Spec.CordonDrain.Cordoned)
	}

	return status, nil
}

// PoolCordonStatus represents the cordon status of a pool
type PoolCordonStatus struct {
	PoolID      string   `json:"pool_id"`
	IsCordoned  bool     `json:"is_cordoned"`
	Constraints []string `json:"constraints,omitempty"`
}

// parseCordonConstraints parses the cordon constraints from the cordonDrain field
// This is a helper function to extract constraint information
func parseCordonConstraints(cordoned *PoolCordonedState) []string {
	if cordoned == nil {
		return []string{}
	}

	// Parse cordonDrain field for constraints
	constraints := []string{}

	if cordoned.Replicas {
		constraints = append(constraints, "replicas")
	}
	if cordoned.Snapshots {
		constraints = append(constraints, "snapshots")
	}
	if cordoned.Restores {
		constraints = append(constraints, "restores")
	}
	if cordoned.Import {
		constraints = append(constraints, "import")
	}

	return constraints
}

// ExpandPoolViaPlugin uses kubectl mayastor plugin to expand the pool
func (cp CPv1) ExpandPoolViaPlugin(poolName string) error {
	logf.Log.Info("Expanding pool via plugin", "pool", poolName)

	args := []string{"-n", common.NSMayastor(), "expand", "pool", poolName}

	// Log the actual command being executed
	logf.Log.Info("Executing command", "command", "kubectl mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to expand pool %s, error %v, output: %s", poolName, err, out.String())
	}

	logf.Log.Info("Successfully expanded pool via plugin", "pool", poolName, "output", out.String())
	return nil
}

// Common Utilities for Offline Pool Deletion Tests

// DeletePoolViaPlugin uses kubectl mayastor plugin to delete the pool
func (cp CPv1) DeleteOfflinePoolViaPlugin(poolName string, flags ...common.OfflinePoolDelete) error {
	logf.Log.Info("Deleting pool via plugin", "pool", poolName)
	args := []string{"-n", common.NSMayastor(), "delete", "pool", poolName}

	// Add constraint flags if specified
	for _, flag := range flags {
		switch flag {
		case common.PurgePool:
			args = append(args, "--purge")
		case common.ConfirmPoolDelete:
			args = append(args, "--yes")
		case common.AcceptDataLoss:
			args = append(args, "--accept-data-loss")
		case common.AcceptVolumeLoss:
			args = append(args, "--accept-volume-loss")
		case common.AcceptSnapshotLoss:
			args = append(args, "--accept-snapshot-loss")
		case common.CleanupCr:
			args = append(args, "--cleanup-cr")
		case common.IgnoreNotFound:
			args = append(args, "--ignore-not-found")
		}
	}

	// Log the actual command being executed
	logf.Log.Info("Executing command", "command", "kubectl mayastor", "args", args)

	cmd := GetMayastorPluginCmd(args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("plugin failed to delete pool %s, error %v, output: %s", poolName, err, out.String())
	}

	logf.Log.Info("Successfully deleted pool via plugin", "pool", poolName, "output", out.String())
	return nil
}
