package v1_rest_api

import (
	"fmt"

	"github.com/openebs/openebs-e2e/common"
)

func (cp CPv1RestApi) GetSnapshots() ([]common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting snapshots"))
}

func (cp CPv1RestApi) GetSnapshot(snapshotId string) (common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting snapshots"))
}

func (cp CPv1RestApi) GetVolumeSnapshot(volUuid string, snapshotId string) (common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting snapshots"))
}

func (cp CPv1RestApi) GetVolumeSnapshots(volUuid string) ([]common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting snapshots"))
}

func (cp CPv1RestApi) GetVolumeSnapshotTopology() ([]common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting volume-snapshot-topology"))
}

func (cp CPv1RestApi) GetPerSnapshotVolumeSnapshotTopology(snapshotId string) (common.SnapshotSchema, error) {
	// #TODO implement REST api for snapshot
	panic(fmt.Errorf("not implemented REST api for getting volume-snapshot-topology"))
}
