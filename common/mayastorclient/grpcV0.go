package mayastorclient

import (
	"fmt"

	v0 "github.com/openebs/openebs-e2e/common/mayastorclient/v0"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Implement the v0 grpc interface in the mayastorclient package
// here to avoid circular imports if implemented in v0
// It is inefficient in terms of invocation.

// Wrap v0NexusWrapper because we need to mutate
// [] v0.v0MayastorNexusChild to []MayastorNexusChild
// to avoid circular imports
type v0NexusWrapper struct {
	nexus v0.V0MayastorNexus
}

type v0RebuildStatsWrapper struct {
	rebuildStats v0.V0RebuildStatsReply
}

func (w v0NexusWrapper) GetString() string {
	return w.nexus.GetString()
}

func (w v0NexusWrapper) GetUuid() string {
	return w.nexus.GetUuid()
}

func (w v0NexusWrapper) GetSize() uint64 {
	return w.nexus.GetSize()
}

func (wh v0RebuildStatsWrapper) GetBlocksTotal() uint64 {
	return wh.rebuildStats.BlocksTotal
}

func (wh v0RebuildStatsWrapper) GetBlocksRecovered() uint64 {
	return wh.rebuildStats.BlocksRecovered
}

func (wh v0RebuildStatsWrapper) GetProgress() uint64 {
	return wh.rebuildStats.Progress
}

func (wh v0RebuildStatsWrapper) GetBlocksTransferred() uint64 {
	panic(fmt.Errorf("gRPC v0 does not have rebuild stats for BlocksTransferred"))
}

func (wh v0RebuildStatsWrapper) GetBlocksRemaining() uint64 {
	panic(fmt.Errorf("gRPC v0 does not have rebuild stats for BlocksRemaining"))
}

func (wh v0RebuildStatsWrapper) GetBlocksPerTask() uint64 {
	panic(fmt.Errorf("gRPC v0 does not have rebuild stats for BlocksPerTask"))
}

func (wh v0RebuildStatsWrapper) GetBlockSize() uint64 {
	return wh.rebuildStats.BlockSize
}

func (wh v0RebuildStatsWrapper) GetTasksTotal() uint64 {
	return wh.rebuildStats.TasksTotal
}

func (wh v0RebuildStatsWrapper) GetTasksActive() uint64 {
	return wh.rebuildStats.TasksActive
}

func (wh v0RebuildStatsWrapper) IsRebuildPartial() bool {
	panic(fmt.Errorf("gRPC v0 does not have rebuild stats for rebuild type"))
}

func (wh v0RebuildStatsWrapper) StartTime() *timestamppb.Timestamp {
	panic(fmt.Errorf("gRPC v0 does not have rebuild stats for rebuild StartTime"))
}

func (w v0NexusWrapper) GetChildren() []MayastorNexusChild {
	children0 := w.nexus.GetChildren()
	var children []MayastorNexusChild
	for _, child := range children0 {
		children = append(children, child)
	}
	return children
}

func (w v0NexusWrapper) GetStateString() string {
	return w.nexus.GetStateString()
}

type grpcV0 struct {
}

func (g grpcV0) Version() string {
	return "v0"
}

// ListNexuses list the set of nexuses found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of nexuses found, even if errors have occurred
func (g grpcV0) ListNexuses(addrs []string) ([]MayastorNexus, error) {
	v0nexuses, err := v0.ListNexuses(addrs)

	var nexuses []MayastorNexus
	for _, v0nexus := range v0nexuses {
		nexuses = append(nexuses, v0NexusWrapper{v0nexus})
	}
	return nexuses, err
}

func (g grpcV0) FaultNexusChild(address string, Uuid string, Uri string) error {
	return v0.FaultNexusChild(address, Uuid, Uri)
}

func (g grpcV0) FindNexus(uuid string, addrs []string) (*MayastorNexus, error) {
	v0nexus, err := v0.FindNexus(uuid, addrs)
	if err == nil && v0nexus != nil {
		var nexus MayastorNexus = v0NexusWrapper{*v0nexus}
		return &nexus, err
	}
	return nil, err
}

func (g grpcV0) ListNvmeControllers(addrs []string) ([]NvmeController, error) {
	var ncs []NvmeController
	v0ncs, err := v0.ListNvmeControllers(addrs)
	if err == nil {
		for _, v0nc := range v0ncs {
			ncs = append(ncs, v0nc)
		}
	}
	return ncs, err
}

func (g grpcV0) GetPool(name, addr string) (MayastorPool, error) {
	v0Pool, err := v0.GetPool(name, addr)
	if err == nil {
		var pool MayastorPool = v0Pool
		return pool, nil
	}
	return nil, err
}

// ListPools list the set of pools found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of pools found, even if errors have occurred
func (g grpcV0) ListPools(addrs []string) ([]MayastorPool, error) {
	var pools []MayastorPool
	v0pools, err := v0.ListPools(addrs)
	for _, v0Pool := range v0pools {
		pools = append(pools, v0Pool)
	}
	return pools, err
}

func (g grpcV0) DestroyAllPools(addrs []string) error {
	return v0.DestroyAllPools(addrs)
}

func (g grpcV0) DestroyPool(name, addr string) error {
	return v0.DestroyPool(name, addr)
}

func (g grpcV0) RmReplica(address string, uuid string) error {
	return v0.RmReplica(address, uuid)
}

func (g grpcV0) CreateReplicaExt(address string, uuid string, size uint64, pool string, thin bool) error {
	return v0.CreateReplicaExt(address, uuid, size, pool, thin)
}

func (g grpcV0) CreateReplica(address string, uuid string, size uint64, pool string) error {
	return v0.CreateReplica(address, uuid, size, pool)
}

// ListReplicas list the set of replicas found for the set of IP addresses
// The required semantics for this function are *always* return the list
// of replicas found, even if errors have occurred
func (g grpcV0) ListReplicas(addrs []string) ([]MayastorReplica, error) {
	var replicas []MayastorReplica
	v0Replicas, err := v0.ListReplicas(addrs)
	for _, v0Repl := range v0Replicas {
		replicas = append(replicas, v0Repl)
	}

	return replicas, err
}

func (g grpcV0) RmNodeReplicas(addrs []string) error {
	return v0.RmNodeReplicas(addrs)
}

func (g grpcV0) FindReplicas(uuid string, addrs []string) ([]MayastorReplica, error) {
	var replicas []MayastorReplica
	v0Replicas, err := v0.FindReplicas(uuid, addrs)
	if err == nil {
		for _, v0Repl := range v0Replicas {
			replicas = append(replicas, v0Repl)
		}
	}

	return replicas, err
}

func (g grpcV0) GetRebuildHistory(uuid string, addrs string) (RebuildHistory, error) {
	panic(fmt.Errorf("gRPC v0 does not have rebuild history rpc call"))
}

func (g grpcV0) GetRebuildStats(uuid string, dstUri string, addrs string) (RebuildStats, error) {
	var rebuildStats RebuildStats
	v0RebuildStats, err := v0.GetRebuildStats(uuid, dstUri, addrs)
	if err == nil {
		rebuildStats = v0RebuildStatsWrapper{v0RebuildStats}
	}
	return rebuildStats, err
}

func (g grpcV0) CheckAndSetConnect(nodes []string) error {
	return v0.CheckAndSetConnect(nodes)
}

func (g grpcV0) CanConnect() bool {
	return v0.CanConnect()
}

func (g grpcV0) WipeReplica(address string, replicaUUID string, poolName string) error {
	return fmt.Errorf("unsupported")
}

func (g grpcV0) ChecksumReplica(address string, replicaUUID string, poolName string) (uint32, error) {
	return 0, fmt.Errorf("unsupported")
}

func (g grpcV0) ShareBdev(address string, bdevUuid string) (string, error) {
	return v0.ShareBdev(address, bdevUuid)
}

func (g grpcV0) UnshareBdev(address string, bdevUuid string) error {
	return v0.UnshareBdev(address, bdevUuid)
}

func (g grpcV0) ResetIOStats(address string) error {
	panic(fmt.Errorf("gRPC v0 does not have io stats reset rpc call"))
}
