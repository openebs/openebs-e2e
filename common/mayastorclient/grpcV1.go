package mayastorclient

import (
	v1 "github.com/openebs/openebs-e2e/common/mayastorclient/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcV1 struct {
}

func (g grpcV1) Version() string {
	return "v1"
}

// Implement the v1 grpc interface in the mayastorclient package
// here to avoid circular imports if implemented in v1
// It is inefficient in terms of invocation.

// Wrap v1NexusWrapper because we need to mutate
// [] v1.v1MayastorNexusChild to []MayastorNexusChild
// to avoid circular imports
type v1NexusWrapper struct {
	nexus v1.V1MayastorNexus
}

type v1RebuildHistoryWrapper struct {
	rebuildHistory v1.V1RebuildHistory
}

type v1RebuildStatsWrapper struct {
	rebuildStats v1.V1RebuildStatsResponse
}

func (w v1NexusWrapper) GetString() string {
	return w.nexus.GetString()
}

func (w v1NexusWrapper) GetUuid() string {
	return w.nexus.GetUuid()
}

func (w v1NexusWrapper) GetSize() uint64 {
	return w.nexus.GetSize()
}

func (wh v1RebuildHistoryWrapper) GetUuid() string {
	return wh.rebuildHistory.GetUuid()
}

func (wh v1RebuildHistoryWrapper) GetRecords() []RebuildHistoryRecord {
	recordsV1 := wh.rebuildHistory.GetRecord()
	var records []RebuildHistoryRecord
	for _, record := range recordsV1 {
		records = append(records, record)
	}
	return records
}

func (rs v1RebuildStatsWrapper) GetBlocksTotal() uint64 {
	return rs.rebuildStats.BlocksTotal
}

func (rs v1RebuildStatsWrapper) GetBlocksRecovered() uint64 {
	return rs.rebuildStats.BlocksRecovered
}

func (rs v1RebuildStatsWrapper) GetBlocksTransferred() uint64 {
	return rs.rebuildStats.BlocksTransferred
}

func (rs v1RebuildStatsWrapper) GetBlocksRemaining() uint64 {
	return rs.rebuildStats.BlocksRemaining
}

func (rs v1RebuildStatsWrapper) GetBlocksPerTask() uint64 {
	return rs.rebuildStats.BlocksPerTask
}

func (rs v1RebuildStatsWrapper) GetProgress() uint64 {
	return rs.rebuildStats.Progress
}

func (rs v1RebuildStatsWrapper) GetBlockSize() uint64 {
	return rs.rebuildStats.BlockSize
}

func (rs v1RebuildStatsWrapper) GetTasksTotal() uint64 {
	return rs.rebuildStats.TasksTotal
}

func (rs v1RebuildStatsWrapper) GetTasksActive() uint64 {
	return rs.rebuildStats.TasksActive
}

func (rs v1RebuildStatsWrapper) IsRebuildPartial() bool {
	return rs.rebuildStats.IsPartial
}

func (rs v1RebuildStatsWrapper) StartTime() *timestamppb.Timestamp {
	return rs.rebuildStats.StartTime
}

func (w v1NexusWrapper) GetChildren() []MayastorNexusChild {
	children0 := w.nexus.GetChildren()
	var children []MayastorNexusChild
	for _, child := range children0 {
		children = append(children, child)
	}
	return children
}

func (w v1NexusWrapper) GetStateString() string {
	return w.nexus.GetStateString()
}

// ListNexuses list the set of nexuses found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of nexuses found, even if errors have occurred
func (g grpcV1) ListNexuses(addrs []string) ([]MayastorNexus, error) {
	v1nexuses, err := v1.ListNexuses(addrs)
	var nexuses []MayastorNexus
	for _, v1nexus := range v1nexuses {
		nexuses = append(nexuses, v1NexusWrapper{v1nexus})
	}
	return nexuses, err
}

func (g grpcV1) FaultNexusChild(address string, Uuid string, Uri string) error {
	return v1.FaultNexusChild(address, Uuid, Uri)
}

func (g grpcV1) FindNexus(uuid string, addrs []string) (*MayastorNexus, error) {
	v1nexus, err := v1.FindNexus(uuid, addrs)
	if err == nil && v1nexus != nil {
		var nexus MayastorNexus = v1NexusWrapper{*v1nexus}
		return &nexus, err
	}
	return nil, err
}

func (g grpcV1) ListNvmeControllers(addrs []string) ([]NvmeController, error) {
	var ncs []NvmeController
	v1ncs, err := v1.ListNvmeControllers(addrs)
	if err == nil {
		for _, v1nc := range v1ncs {
			ncs = append(ncs, v1nc)
		}
	}
	return ncs, err
}

func (g grpcV1) GetPool(name, addr string) (MayastorPool, error) {
	v1Pool, err := v1.GetPool(name, addr)
	if err == nil {
		var pool MayastorPool = v1Pool
		return pool, nil
	}
	return nil, err
}

// ListPools list the set of pools found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of pools found, even if errors have occurred
func (g grpcV1) ListPools(addrs []string) ([]MayastorPool, error) {
	var pools []MayastorPool
	v1pools, err := v1.ListPools(addrs)
	for _, v1Pool := range v1pools {
		pools = append(pools, v1Pool)
	}
	return pools, err
}

func (g grpcV1) DestroyAllPools(addrs []string) error {
	return v1.DestroyAllPools(addrs)
}

func (g grpcV1) DestroyPool(name, addr string) error {
	return v1.DestroyPool(name, addr)
}

func (g grpcV1) RmReplica(address string, uuid string) error {
	return v1.RmReplica(address, uuid)
}

func (g grpcV1) CreateReplicaExt(address string, uuid string, size uint64, pool string, thin bool) error {
	return v1.CreateReplicaExt(address, uuid, size, pool, thin)
}

func (g grpcV1) CreateReplica(address string, uuid string, size uint64, pool string) error {
	return v1.CreateReplica(address, uuid, size, pool)
}

// ListReplicas list the set of replicas found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of replicas found, even if errors have occurred
func (g grpcV1) ListReplicas(addrs []string) ([]MayastorReplica, error) {
	var replicas []MayastorReplica
	v1Replicas, err := v1.ListReplicas(addrs)

	for _, v1Repl := range v1Replicas {
		replicas = append(replicas, v1Repl)
	}

	return replicas, err
}

func (g grpcV1) RmNodeReplicas(addrs []string) error {
	return v1.RmNodeReplicas(addrs)
}

func (g grpcV1) FindReplicas(uuid string, addrs []string) ([]MayastorReplica, error) {
	var replicas []MayastorReplica
	v1Replicas, err := v1.FindReplicas(uuid, addrs)
	if err == nil {
		for _, v1Repl := range v1Replicas {
			replicas = append(replicas, v1Repl)
		}
	}

	return replicas, err
}

func (g grpcV1) GetRebuildHistory(uuid string, addrs string) (RebuildHistory, error) {
	var rebuildHistory RebuildHistory
	v1RebuildHistory, err := v1.GetRebuildHistory(uuid, addrs)
	if err == nil {
		rebuildHistory = v1RebuildHistoryWrapper{v1RebuildHistory}
	}
	return rebuildHistory, err
}

func (g grpcV1) GetRebuildStats(uuid string, dstUri string, addrs string) (RebuildStats, error) {
	var rebuildStats RebuildStats
	v1RebuildStats, err := v1.GetRebuildStats(uuid, dstUri, addrs)
	if err == nil {
		rebuildStats = v1RebuildStatsWrapper{v1RebuildStats}
	}
	return rebuildStats, err
}

func (g grpcV1) CheckAndSetConnect(nodes []string) error {
	return v1.CheckAndSetConnect(nodes)
}

func (g grpcV1) CanConnect() bool {
	return v1.CanConnect()
}

func (g grpcV1) WipeReplica(address string, replicaUUID string, poolName string) error {
	return v1.WipeReplica(address, replicaUUID, poolName)
}

func (g grpcV1) ChecksumReplica(address string, replicaUUID string, poolName string) (uint32, error) {
	return v1.ChecksumReplica(address, replicaUUID, poolName)
}

func (g grpcV1) ShareBdev(address string, bdevUuid string) (string, error) {
	return v1.ShareBdev(address, bdevUuid)
}

func (g grpcV1) UnshareBdev(address string, bdevUuid string) error {
	return v1.UnshareBdev(address, bdevUuid)
}

func (g grpcV1) ResetIOStats(address string) error {
	return v1.ResetIOStats(address)
}
