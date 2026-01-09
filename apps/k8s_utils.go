package apps

import (
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/k8stest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	Standalone   Architecture = "standalone"
	Replicaset   Architecture = "replicaset"
	Replication  Architecture = "replication"
	replicaCount              = "replicaCount"
)

type Architecture string

func (a Architecture) String() string {
	return string(a)
}

func CreateStorageClass(mb *mongoBuilder) (string, error) {

	name := fmt.Sprintf(
		"mayastor-mongo-%d-%s-%s",
		mb.replicaCount,
		common.ThinProvisioning.String(),
		common.Ext4FsType,
	)

	scb := k8stest.NewScBuilder().
		WithNamespace(mb.namespace).
		WithName(name).
		WithReplicas(mb.replicaCount).
		WithProvisioningType(common.ThinProvisioning).
		WithFileSystemType(common.Ext4FsType).
		WithCloneFsIdAsVolumeId(common.CloneFsIdAsVolumeIdNone)

	scb.WithStsAffinityGroup(common.StsAffinityGroupEnable)

	if err := scb.BuildAndCreate(); err != nil {
		return "", err
	}

	logf.Log.Info("MongoDB StorageClass created",
		"name", name,
		"replicas", mb.replicaCount,
	)

	return name, nil
}

func CreatePostgresStorageClass(pb *postgresBuilder) (string, error) {
	var err error
	var poolsInCluster []common.MayastorPool
	const sleepTime = 3

	if pb.replicaCount == 0 && pb.architecture == Replication {
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
		pb.replicaCount = len(poolsInCluster)
		pb.values["readReplicas.replicaCount"] = pb.replicaCount
	} else if pb.replicaCount == 0 && pb.architecture == Standalone {
		pb.replicaCount = 1
	}

	name := fmt.Sprintf("mayastor-%s-%d-%s-%s", pb.architecture, pb.replicaCount, pb.provisioningType.String(), pb.filesystemType)
	scb := k8stest.NewScBuilder().
		WithNamespace(pb.namespace).
		WithReplicas(pb.replicaCount).
		WithProvisioningType(pb.provisioningType).
		WithFileSystemType(pb.filesystemType).
		WithCloneFsIdAsVolumeId(pb.CloneFsIdAsVolumeIdType)
	if pb.architecture == Replication {
		scb.WithStsAffinityGroup(common.StsAffinityGroupEnable)
	}
	if pb.filesystemType == common.BtrfsFsType {
		scb.WithMountOption("nodatacow")
	}
	err = scb.WithName(name).BuildAndCreate()
	if err != nil {
		return "", err
	}
	return name, nil
}
