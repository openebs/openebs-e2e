package v1

import (
	"context"
	"fmt"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type v1MayastorNexusChild struct {
	child *mayastorGrpc.Child
}

func (c v1MayastorNexusChild) Uri() string {
	return c.child.GetUri()
}

func (c v1MayastorNexusChild) IsOnline() bool {
	return c.child.State == mayastorGrpc.ChildState_CHILD_STATE_ONLINE
}

func (c v1MayastorNexusChild) IsDegraded() bool {
	return c.child.State == mayastorGrpc.ChildState_CHILD_STATE_DEGRADED
}

func (c v1MayastorNexusChild) GetState() int32 {
	return int32(c.child.State)
}

func (c v1MayastorNexusChild) RebuildProgress() int32 {
	return c.child.GetRebuildProgress()
}

type v1RebuildHistoryRecord struct {
	record *mayastorGrpc.RebuildHistoryRecord
}

// V1MayastorNexus Mayastor Nexus data
type V1MayastorNexus struct {
	Name      string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Uuid      string                  `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Size      uint64                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	State     mayastorGrpc.NexusState `protobuf:"varint,3,opt,name=state,proto3,enum=mayastor.NexusState" json:"state,omitempty"`
	Children  []*mayastorGrpc.Child   `protobuf:"bytes,4,rep,name=children,proto3" json:"children,omitempty"`
	DeviceUri string                  `protobuf:"bytes,5,opt,name=device_uri,json=deviceUri,proto3" json:"device_uri,omitempty"`
	Rebuilds  uint32                  `protobuf:"varint,6,opt,name=rebuilds,proto3" json:"rebuilds,omitempty"`
}

type V1RebuildHistory struct {
	Uuid    string                               `protobuf:"bytes,1,opt,name=nexus,proto3" json:"nexus,omitempty"`     // uuid of the nexus
	Records []*mayastorGrpc.RebuildHistoryRecord `protobuf:"bytes,2,rep,name=records,proto3" json:"records,omitempty"` // List of all the rebuild records
}

type V1RebuildStatsResponse struct {
	BlocksTotal       uint64                 `protobuf:"varint,1,opt,name=blocks_total,json=blocksTotal,proto3" json:"blocks_total,omitempty"`                   // total number of blocks to recover
	BlocksRecovered   uint64                 `protobuf:"varint,2,opt,name=blocks_recovered,json=blocksRecovered,proto3" json:"blocks_recovered,omitempty"`       // number of blocks already recovered
	BlocksTransferred uint64                 `protobuf:"varint,8,opt,name=blocks_transferred,json=blocksTransferred,proto3" json:"blocks_transferred,omitempty"` // number of blocks for which the actual data transfer already occurred
	BlocksRemaining   uint64                 `protobuf:"varint,9,opt,name=blocks_remaining,json=blocksRemaining,proto3" json:"blocks_remaining,omitempty"`       // number of blocks remaining to transfer
	Progress          uint64                 `protobuf:"varint,3,opt,name=progress,proto3" json:"progress,omitempty"`                                            // rebuild progress %
	BlocksPerTask     uint64                 `protobuf:"varint,4,opt,name=blocks_per_task,json=blocksPerTask,proto3" json:"blocks_per_task,omitempty"`           // granularity of each recovery task in blocks
	BlockSize         uint64                 `protobuf:"varint,5,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`                         // size in bytes of logical block
	TasksTotal        uint64                 `protobuf:"varint,6,opt,name=tasks_total,json=tasksTotal,proto3" json:"tasks_total,omitempty"`                      // total number of concurrent rebuild tasks
	TasksActive       uint64                 `protobuf:"varint,7,opt,name=tasks_active,json=tasksActive,proto3" json:"tasks_active,omitempty"`                   // number of current active tasks
	IsPartial         bool                   `protobuf:"varint,10,opt,name=is_partial,json=isPartial,proto3" json:"is_partial,omitempty"`                        // true for partial (only modified blocked transferred); false for the full rebuild (all blocks transferred)
	StartTime         *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`                         // start time of the rebuild (UTC)
}

func (msn V1MayastorNexus) GetString() string {
	descChildren := "["
	for _, child := range msn.Children {
		descChildren = fmt.Sprintf("%s(%v); ", descChildren, child)
	}
	descChildren += "]"
	return fmt.Sprintf("Uuid=%s; Size=%d; State=%v; DeviceUri=%s, Rebuilds=%d; Children=%v",
		msn.Uuid, msn.Size, msn.State, msn.DeviceUri, msn.Rebuilds, descChildren)
}

func (msn V1MayastorNexus) GetUuid() string {
	return msn.Uuid
}

func (msn V1MayastorNexus) GetSize() uint64 {
	return msn.Size
}

func (msn V1MayastorNexus) GetChildren() []v1MayastorNexusChild {
	var children []v1MayastorNexusChild
	for _, child := range msn.Children {
		children = append(children, v1MayastorNexusChild{
			child,
		})
	}
	return children
}

func listNexuses(address string) ([]V1MayastorNexus, error) {
	var nexusInfos []V1MayastorNexus
	var err error

	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("listNexuses", "error", err)
		return nexusInfos, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("ListNexuses", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewNexusRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var response *mayastorGrpc.ListNexusResponse
	retryBackoff(func() error {
		response, err = c.ListNexus(ctx, &mayastorGrpc.ListNexusOptions{})
		return err
	})

	if err == nil {
		if response != nil {
			for _, nexus := range response.NexusList {
				ni := V1MayastorNexus{
					Name:      nexus.Name,
					Uuid:      nexus.Uuid,
					Size:      nexus.Size,
					State:     nexus.State,
					Children:  nexus.Children,
					DeviceUri: nexus.DeviceUri,
					Rebuilds:  nexus.Rebuilds,
				}
				nexusInfos = append(nexusInfos, ni)
			}
		} else {
			err = fmt.Errorf("nil response for ListNexus on %s", address)
			logf.Log.Info("ListNexuses", "error", err)
		}
	} else {
		err = niceError(err)
		logf.Log.Info("ListNexuses", "error", err)
	}
	return nexusInfos, err
}

// ListNexuses given a list of node ip addresses, enumerate the set of nexuses on mayastor using gRPC on each of those nodes
// returns accumulated errors if gRPC communication failed.
func ListNexuses(addrs []string) ([]V1MayastorNexus, error) {
	var accErr error
	var nexusInfos []V1MayastorNexus
	for _, address := range addrs {
		nexusInfo, err := listNexuses(address)
		if err == nil {
			nexusInfos = append(nexusInfos, nexusInfo...)
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return nexusInfos, accErr
}

func FaultNexusChild(address string, Uuid string, Uri string) error {
	var err error
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("FaultNexusChild", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("FaultNexusChild", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewNexusRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	faultRequest := mayastorGrpc.FaultNexusChildRequest{
		Uuid: Uuid,
		Uri:  Uri,
	}
	var response *mayastorGrpc.FaultNexusChildResponse
	retryBackoff(func() error {
		response, err = c.FaultNexusChild(ctx, &faultRequest)
		return err
	})

	if err == nil {
		if response == nil {
			err = fmt.Errorf("nil response to FaultNexusChild")
		}
	} else {
		err = niceError(err)
		logf.Log.Info("FaultNexusChild", "error", err)
	}

	return err
}

// FindNexus given a list of node ip addresses, return the common.MayastorNexus with matching uuid
// returns accumulated errors if gRPC communication failed.
func FindNexus(uuid string, addrs []string) (*V1MayastorNexus, error) {
	var accErr error
	for _, address := range addrs {
		nexusInfos, err := listNexuses(address)
		if err == nil {
			for _, ni := range nexusInfos {
				if ni.Uuid == uuid {
					return &ni, nil
				}
			}
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, niceError(err))
			} else {
				accErr = err
			}
		}
	}
	return nil, accErr
}

// GetRebuildHistory given a node ip address, return the V1RebuildHistory with matching uuid
// returns accumulated errors if gRPC communication failed.
func GetRebuildHistory(uuid string, address string) (V1RebuildHistory, error) {
	var rebuildHistory V1RebuildHistory
	var err error

	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("GetRebuildHistory", "error", err)
		return rebuildHistory, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("GetRebuildHistory", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewNexusRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rebuildHistoryRequest := mayastorGrpc.RebuildHistoryRequest{
		Uuid: uuid,
	}
	var response *mayastorGrpc.RebuildHistoryResponse
	retryBackoff(func() error {
		response, err = c.GetRebuildHistory(ctx, &rebuildHistoryRequest)
		return err
	})

	if err == nil {
		if response == nil {
			err = fmt.Errorf("nil response to GetRebuildHistory")
		} else {
			rebuildHistory = V1RebuildHistory{
				Uuid:    response.Uuid,
				Records: response.Records,
			}
		}

	} else {
		err = niceError(err)
		logf.Log.Info("GetRebuildHistory", "error", err)
	}

	return rebuildHistory, err
}

// GetRebuildStats given a node ip address, return the V1RebuildStatsResponse with matching nexus uuid and destination uri
// returns accumulated errors if gRPC communication failed.
func GetRebuildStats(uuid string, dstUri string, address string) (V1RebuildStatsResponse, error) {
	var rebuildStats V1RebuildStatsResponse
	var err error

	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("GetRebuildStats", "error", err)
		return rebuildStats, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("GetRebuildStats", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewNexusRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rebuildStatsRequest := mayastorGrpc.RebuildStatsRequest{
		NexusUuid: uuid,
		Uri:       dstUri,
	}
	var response *mayastorGrpc.RebuildStatsResponse
	retryBackoff(func() error {
		response, err = c.GetRebuildStats(ctx, &rebuildStatsRequest)
		return err
	})

	if err == nil {
		if response == nil {
			err = fmt.Errorf("nil response to GetRebuildStats")
		} else {
			rebuildStats = V1RebuildStatsResponse{
				BlocksTotal:       response.BlocksTotal,
				BlocksRecovered:   response.BlocksRecovered,
				BlocksTransferred: response.BlocksTransferred,
				BlocksRemaining:   response.BlocksRemaining,
				BlocksPerTask:     response.BlocksPerTask,
				BlockSize:         response.BlockSize,
				Progress:          response.Progress,
				TasksTotal:        response.TasksTotal,
				TasksActive:       response.TasksActive,
				IsPartial:         response.IsPartial,
				StartTime:         response.StartTime,
			}
		}

	} else {
		err = niceError(err)
		logf.Log.Info("GetRebuildStats", "error", err)
	}

	return rebuildStats, err
}

func (h V1RebuildHistory) GetUuid() string {
	return h.Uuid
}

func (h v1RebuildHistoryRecord) GetChildUri() string {
	return h.record.ChildUri
}

func (h v1RebuildHistoryRecord) GetSrcUri() string {
	return h.record.SrcUri
}

func (h v1RebuildHistoryRecord) GetStateString() string {
	return h.record.GetState().String()
}

func (h v1RebuildHistoryRecord) GetBlocksTotal() uint64 {
	return h.record.BlocksTotal
}

func (h v1RebuildHistoryRecord) GetBlocksRecovered() uint64 {
	return h.record.BlocksRecovered
}

func (h v1RebuildHistoryRecord) GetBlocksTransferred() uint64 {
	return h.record.BlocksTransferred
}

func (h v1RebuildHistoryRecord) GetBlocksRemaining() uint64 {
	return h.record.BlocksRemaining
}

func (h v1RebuildHistoryRecord) GetBlocksPerTask() uint64 {
	return h.record.BlocksPerTask
}

func (h v1RebuildHistoryRecord) GetBlockSize() uint64 {
	return h.record.BlockSize
}

func (h v1RebuildHistoryRecord) IsPartial() bool {
	return h.record.IsPartial
}
func (h v1RebuildHistoryRecord) StartTime() *timestamppb.Timestamp {
	return h.record.StartTime
}
func (h v1RebuildHistoryRecord) EndTime() *timestamppb.Timestamp {
	return h.record.EndTime
}

func (h V1RebuildHistory) GetRecord() []v1RebuildHistoryRecord {
	var records []v1RebuildHistoryRecord
	for _, record := range h.Records {
		records = append(records, v1RebuildHistoryRecord{
			record,
		})
	}
	return records
}

func (h V1RebuildStatsResponse) GetBlocksTotal() uint64 {
	return h.BlocksTotal
}

func (h V1RebuildStatsResponse) GetBlocksRecovered() uint64 {
	return h.BlocksRecovered
}

func (h V1RebuildStatsResponse) GetBlocksTransferred() uint64 {
	return h.BlocksTransferred
}

func (h V1RebuildStatsResponse) GetBlocksRemaining() uint64 {
	return h.BlocksRemaining
}

func (h V1RebuildStatsResponse) GetBlocksPerTask() uint64 {
	return h.BlocksPerTask
}

func (h V1RebuildStatsResponse) GetBlockSize() uint64 {
	return h.BlockSize
}

func (h V1RebuildStatsResponse) GetTasksTotal() uint64 {
	return h.TasksTotal
}

func (h V1RebuildStatsResponse) GetTasksActive() uint64 {
	return h.TasksActive
}

func (h V1RebuildStatsResponse) IsRebuildPartial() bool {
	return h.IsPartial
}

func (h V1RebuildStatsResponse) GetStartTime() *timestamppb.Timestamp {
	return h.StartTime
}

func (h V1RebuildStatsResponse) GetProgress() uint64 {
	return h.Progress
}

func (msn V1MayastorNexus) GetStateString() string {
	switch msn.State {
	case mayastorGrpc.NexusState_NEXUS_ONLINE:
		return "Online"
	case mayastorGrpc.NexusState_NEXUS_DEGRADED:
		return "Degraded"
	case mayastorGrpc.NexusState_NEXUS_FAULTED:
		return "Faulted"
	case mayastorGrpc.NexusState_NEXUS_UNKNOWN:
		return "Unknown"
	}
	return "?"
}
