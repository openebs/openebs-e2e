package mayastorclient

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MayastorNexus interface {
	GetString() string
	GetUuid() string
	GetSize() uint64
	GetChildren() []MayastorNexusChild
	GetStateString() string
}

type MayastorNexusChild interface {
	Uri() string
	RebuildProgress() int32
	IsOnline() bool
	IsDegraded() bool
}

// NvmeController  Nvme controller data
type NvmeController interface {
	GetName() string
	GetStateString() string
	GetSize() uint64
	GetBlkSize() uint32
	GetString() string
}

// MayastorPool Mayastor Pool data
type MayastorPool interface {
	GetName() string
	GetDisks() []string
	GetCapacity() uint64
	GetUsed() uint64
	GetString() string
	GetStateString() string
	IsPoolOnline() bool
	ToCrdState() string
}

// MayastorReplica Mayastor Replica data
type MayastorReplica interface {
	GetUuid() string
	GetPool() string
	GetThin() bool
	GetSize() uint64
	GetShareString() string
	GetUri() string
	Description() string
}

type RebuildHistory interface {
	GetUuid() string
	GetRecords() []RebuildHistoryRecord
}

type RebuildHistoryRecord interface {
	GetChildUri() string
	GetSrcUri() string
	GetStateString() string
	GetBlocksTotal() uint64
	GetBlocksRecovered() uint64
	GetBlocksTransferred() uint64
	GetBlocksRemaining() uint64
	GetBlocksPerTask() uint64
	GetBlockSize() uint64
	IsPartial() bool
	StartTime() *timestamppb.Timestamp
	EndTime() *timestamppb.Timestamp
}

type RebuildStats interface {
	GetBlocksTotal() uint64
	GetBlocksRecovered() uint64
	GetBlocksTransferred() uint64
	GetBlocksRemaining() uint64
	GetBlocksPerTask() uint64
	GetProgress() uint64
	GetBlockSize() uint64
	GetTasksTotal() uint64
	GetTasksActive() uint64
	IsRebuildPartial() bool
	StartTime() *timestamppb.Timestamp
}

type MayastorReplicaArray []MayastorReplica

func (msr MayastorReplicaArray) Len() int           { return len(msr) }
func (msr MayastorReplicaArray) Less(i, j int) bool { return msr[i].GetPool() < msr[j].GetPool() }
func (msr MayastorReplicaArray) Swap(i, j int)      { msr[i], msr[j] = msr[j], msr[i] }
