package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/openebs/openebs-e2e/common"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type MayastorCpPool struct {
	Id    string   `json:"id"`
	Spec  mspSpec  `json:"spec"`
	State mspState `json:"state"`
}

type mspSpec struct {
	Disks  []string          `json:"disks"`
	Id     string            `json:"id"`
	Labels map[string]string `json:"labels"`
	Node   string            `json:"node"`
	Status string            `json:"status"`
}

type mspState struct {
	Capacity  uint64   `json:"capacity"`
	Disks     []string `json:"disks"`
	ID        string   `json:"id"`
	Node      string   `json:"node"`
	Status    string   `json:"status"`
	Used      uint64   `json:"used"`
	Committed uint64   `json:"committed"`
}

func (cp CPv1) CreatePoolOnInstall() bool {
	no_pool_install := os.Getenv("no_pool_install")
	return no_pool_install != "true"
}

func GetMayastorCpPool(name string) (*MayastorCpPool, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "pool", name)
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

func ListMayastorCpPools() ([]MayastorCpPool, error) {
	pluginpath := GetPluginPath()

	var jsonInput []byte
	var err error
	cmd := exec.Command(pluginpath, "-n", common.NSMayastor(), "-ojson", "get", "pools")
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
