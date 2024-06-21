package v0

import (
	"context"
	"fmt"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v0/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type v0MayastorNexusChild struct {
	child *mayastorGrpc.Child
}

func (c v0MayastorNexusChild) Uri() string {
	return c.child.GetUri()
}

func (c v0MayastorNexusChild) IsOnline() bool {
	return c.child.State == mayastorGrpc.ChildState_CHILD_ONLINE
}

func (c v0MayastorNexusChild) IsDegraded() bool {
	return c.child.State == mayastorGrpc.ChildState_CHILD_DEGRADED
}

func (c v0MayastorNexusChild) GetState() int32 {
	return int32(c.child.State)
}

func (c v0MayastorNexusChild) RebuildProgress() int32 {
	return c.child.GetRebuildProgress()
}

// V0MayastorNexus Mayastor Nexus data
type V0MayastorNexus struct {
	Name      string                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Uuid      string                  `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Size      uint64                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	State     mayastorGrpc.NexusState `protobuf:"varint,3,opt,name=state,proto3,enum=mayastor.NexusState" json:"state,omitempty"`
	Children  []*mayastorGrpc.Child   `protobuf:"bytes,4,rep,name=children,proto3" json:"children,omitempty"`
	DeviceUri string                  `protobuf:"bytes,5,opt,name=device_uri,json=deviceUri,proto3" json:"device_uri,omitempty"`
	Rebuilds  uint32                  `protobuf:"varint,6,opt,name=rebuilds,proto3" json:"rebuilds,omitempty"`
}

type V0RebuildStatsReply struct {
	BlocksTotal     uint64 `protobuf:"varint,1,opt,name=blocks_total,json=blocksTotal,proto3" json:"blocks_total,omitempty"`               // total number of blocks to recover
	BlocksRecovered uint64 `protobuf:"varint,2,opt,name=blocks_recovered,json=blocksRecovered,proto3" json:"blocks_recovered,omitempty"`   // number of blocks recovered
	Progress        uint64 `protobuf:"varint,3,opt,name=progress,proto3" json:"progress,omitempty"`                                        // rebuild progress %
	SegmentSizeBlks uint64 `protobuf:"varint,4,opt,name=segment_size_blks,json=segmentSizeBlks,proto3" json:"segment_size_blks,omitempty"` // granularity of each recovery copy in blocks
	BlockSize       uint64 `protobuf:"varint,5,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`                     // size in bytes of each block
	TasksTotal      uint64 `protobuf:"varint,6,opt,name=tasks_total,json=tasksTotal,proto3" json:"tasks_total,omitempty"`                  // total number of concurrent rebuild tasks
	TasksActive     uint64 `protobuf:"varint,7,opt,name=tasks_active,json=tasksActive,proto3" json:"tasks_active,omitempty"`               // number of current active tasks
}

func (msn V0MayastorNexus) GetString() string {
	descChildren := "["
	for _, child := range msn.Children {
		descChildren = fmt.Sprintf("%s(%v); ", descChildren, child)
	}
	descChildren += "]"
	return fmt.Sprintf("Uuid=%s; Size=%d; State=%v; DeviceUri=%s, Rebuilds=%d; Children=%v",
		msn.Uuid, msn.Size, msn.State, msn.DeviceUri, msn.Rebuilds, descChildren)
}

func (msn V0MayastorNexus) GetUuid() string {
	return msn.Uuid
}

func (msn V0MayastorNexus) GetSize() uint64 {
	return msn.Size
}

func (msn V0MayastorNexus) GetChildren() []v0MayastorNexusChild {
	var children []v0MayastorNexusChild
	for _, child := range msn.Children {
		children = append(children, v0MayastorNexusChild{
			child,
		})
	}
	return children
}

func listNexuses(address string) ([]V0MayastorNexus, error) {
	var nexusInfos []V0MayastorNexus
	var err error

	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
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

	c := mayastorGrpc.NewMayastorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var response *mayastorGrpc.ListNexusV2Reply
	retryBackoff(func() error {
		response, err = c.ListNexusV2(ctx, &null)
		return err
	})

	if err == nil {
		if response != nil {
			for _, nexus := range response.NexusList {
				ni := V0MayastorNexus{
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
func ListNexuses(addrs []string) ([]V0MayastorNexus, error) {
	var accErr error
	var nexusInfos []V0MayastorNexus
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
	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
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

	c := mayastorGrpc.NewMayastorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	faultRequest := mayastorGrpc.FaultNexusChildRequest{
		Uuid: Uuid,
		Uri:  Uri,
	}
	var response *mayastorGrpc.Null
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

// FindNexus given a list of node ip addresses, return the V0MayastorNexus with matching uuid
// returns accumulated errors if gRPC communication failed.
func FindNexus(uuid string, addrs []string) (*V0MayastorNexus, error) {
	var accErr error
	for _, address := range addrs {
		nexusInfos, err := listNexuses(address)
		if err == nil {
			for _, ni := range nexusInfos {
				if ni.GetUuid() == uuid {
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

// GetRebuildStats given a node ip address, return the V0RebuildHistory with matching uuid and uri
// returns accumulated errors if gRPC communication failed.
func GetRebuildStats(uuid string, dstUri string, addrs string) (V0RebuildStatsReply, error) {
	var err error
	addrPort := fmt.Sprintf("%s:%d", addrs, mayastorPort)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	var rebuildStats V0RebuildStatsReply
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

	c := mayastorGrpc.NewMayastorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rebuildStatsRequest := mayastorGrpc.RebuildStatsRequest{
		Uuid: uuid,
		Uri:  dstUri,
	}
	var response *mayastorGrpc.RebuildStatsReply
	retryBackoff(func() error {
		response, err = c.GetRebuildStats(ctx, &rebuildStatsRequest)
		return err
	})

	if err == nil {
		if response == nil {
			err = fmt.Errorf("nil response to GetRebuildStats")
		} else {
			rebuildStats = V0RebuildStatsReply{
				BlocksTotal:     response.BlocksTotal,
				BlocksRecovered: response.BlocksRecovered,
				Progress:        response.Progress,
				SegmentSizeBlks: response.SegmentSizeBlks,
				BlockSize:       response.BlockSize,
				TasksTotal:      response.TasksTotal,
				TasksActive:     response.TasksActive,
			}
		}

	} else {
		err = niceError(err)
		logf.Log.Info("GetRebuildStats", "error", err)
	}

	return rebuildStats, err
}

func (msn V0MayastorNexus) GetStateString() string {
	switch msn.State {
	case mayastorGrpc.NexusState_NEXUS_ONLINE:
		return "Online"
	case mayastorGrpc.NexusState_NEXUS_DEGRADED:
		return "Degraded"
	case mayastorGrpc.NexusState_NEXUS_FAULTED:
		return "Faulted"
	case mayastorGrpc.NexusState_NEXUS_UNKNOWN:
		return "Unknown"
	case mayastorGrpc.NexusState_NEXUS_SHUTDOWN:
		return "Shutdown"
	case mayastorGrpc.NexusState_NEXUS_SHUTTING_DOWN:
		return "Shutting_Down"
	}
	return "?"
}
