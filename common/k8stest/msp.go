package k8stest

import (
	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/controlplane"
)

func GetMsPool(poolName string) (*common.MayastorPool, error) {
	return controlplane.GetMsPool(poolName)
}

func ListMsPools() ([]common.MayastorPool, error) {
	return controlplane.ListMsPools()
}
