package v1_rest_api

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/generated/openapi"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func oaPoolToMsp(oaPool openapi.Pool) common.MayastorPool {
	var poolSpec openapi.PoolSpec = oaPool.GetSpec()
	var poolState openapi.PoolState = oaPool.GetState()

	return common.MayastorPool{
		Name: oaPool.Id,
		Spec: common.MayastorPoolSpec{
			Node:  poolSpec.GetNode(),
			Disks: poolSpec.GetDisks(),
		},
		Status: common.MayastorPoolStatus{
			Capacity: uint64(poolState.GetCapacity()),
			Used:     uint64(poolState.GetUsed()),
			Disks:    poolState.GetDisks(),
			Spec: common.MayastorPoolSpec{
				Disks: poolSpec.GetDisks(),
				Node:  poolSpec.GetNode(),
			},
			State:  string(poolState.GetStatus()),
			Avail:  uint64(poolState.GetCapacity() - poolState.GetUsed()),
			Reason: "",
		},
	}
}

// GetMsPool Get pointer to a mayastor control plane pool
func (cp CPv1RestApi) GetMsPool(poolName string) (*common.MayastorPool, error) {
	oaPool, err, statusCode := cp.oa.getPool(poolName)

	if err != nil {
		logf.Log.Info("getPool failed", "statusCode", statusCode)
		return nil, err
	}

	msp := oaPoolToMsp(oaPool)
	return &msp, nil
}

func (cp CPv1RestApi) ListMsPools() ([]common.MayastorPool, error) {
	var msPools []common.MayastorPool
	oaPools, err, statusCode := cp.oa.getPools()

	if err != nil {
		logf.Log.Info("getPools failed", "statusCode", statusCode)
	} else {
		for _, item := range oaPools {
			msPools = append(msPools, oaPoolToMsp(item))
		}
	}

	return msPools, err
}
