package apps

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func CreateStorageClass(mb *mongoBuilder) (string, error) {
	var err error
	var poolsInCluster []common.MayastorPool
	const sleepTime = 3

	if mb.replicaCount == 0 && mb.architecture == Replicaset {
		for ix := 0; ix < (k8stest.DefTimeoutSecs+sleepTime-1)/sleepTime; ix++ {
			poolsInCluster, err = k8stest.ListMsPools()
			if err != nil {
				logf.Log.Info("ListMsPools", "Error", err)
				time.Sleep(sleepTime * time.Second)
				continue
			}
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to list disk pools, error %v", err)
		}
		mb.replicaCount = len(poolsInCluster)
		mb.values["replicaCount"] = mb.replicaCount
	} else if mb.replicaCount == 0 && mb.architecture == Standalone {
		mb.replicaCount = 1
	}

	name := fmt.Sprintf("mayastor-%s-%d-%s-%s", mb.architecture, mb.replicaCount, mb.provisioningType.String(), mb.filesystemType)
	scb := k8stest.NewScBuilder().
		WithNamespace(mb.namespace).
		WithReplicas(mb.replicaCount).
		WithProvisioningType(mb.provisioningType).
		WithFileSystemType(mb.filesystemType).
		WithCloneFsIdAsVolumeId(mb.CloneFsIdAsVolumeIdType)
	if mb.architecture == Replicaset {
		scb.WithStsAffinityGroup(common.StsAffinityGroupEnable)
	}
	if mb.filesystemType == common.BtrfsFsType {
		scb.WithMountOption("nodatacow")
	}
	err = scb.WithName(name).BuildAndCreate()
	if err != nil {
		return "", err
	}
	return name, nil
}
