package v1

import (
	"context"
	"fmt"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/runtime/schema"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type v1MayastorPool struct {
	Name     string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Disks    []string               `protobuf:"bytes,2,rep,name=disks,proto3" json:"disks,omitempty"`
	State    mayastorGrpc.PoolState `protobuf:"varint,3,opt,name=state,proto3,enum=mayastor.PoolState" json:"state,omitempty"`
	Capacity uint64                 `protobuf:"varint,4,opt,name=capacity,proto3" json:"capacity,omitempty"`
	Used     uint64                 `protobuf:"varint,5,opt,name=used,proto3" json:"used,omitempty"`
	Uuid     string                 `protobuf:"bytes,6,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (msp v1MayastorPool) GetString() string {
	return fmt.Sprintf("Name=%s; Disks=%v; State=%v; Used=%d, Capacity=%d;",
		msp.Name, msp.Disks, msp.State, msp.Used, msp.Capacity)
}

func (msp v1MayastorPool) GetName() string {
	return msp.Name
}

func (msp v1MayastorPool) GetDisks() []string {
	return msp.Disks
}

func (msp v1MayastorPool) GetState() mayastorGrpc.PoolState {
	return msp.State
}

func (msp v1MayastorPool) GetCapacity() uint64 {
	return msp.Capacity
}

func (msp v1MayastorPool) GetUsed() uint64 {
	return msp.Used
}

func (msp v1MayastorPool) GetStateString() string {
	return msp.State.String()
}

func (msp v1MayastorPool) IsPoolOnline() bool {
	return msp.State == mayastorGrpc.PoolState_POOL_ONLINE
}

func (msp v1MayastorPool) GetUuid() string {
	return msp.Uuid
}

func (msp v1MayastorPool) ToCrdState() string {
	switch msp.State {
	case mayastorGrpc.PoolState_POOL_UNKNOWN:
		return "Pending"
	case mayastorGrpc.PoolState_POOL_ONLINE:
		return "Online"
	case mayastorGrpc.PoolState_POOL_DEGRADED:
		return "Degraded"
	case mayastorGrpc.PoolState_POOL_FAULTED:
		return "Faulted"
	default:
		return "Offline"
	}
}

func listPool(address string) ([]v1MayastorPool, error) {
	var poolInfos []v1MayastorPool
	var err error
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("listPool", "error", err)
		return poolInfos, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("listPool", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewPoolRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var response *mayastorGrpc.ListPoolsResponse
	retryBackoff(func() error {
		response, err = c.ListPools(ctx, &mayastorGrpc.ListPoolOptions{})
		return err
	})

	if err == nil {
		if response != nil {
			for _, pool := range response.Pools {
				pi := v1MayastorPool{
					Name:     pool.Name,
					Disks:    pool.Disks,
					State:    pool.State,
					Capacity: pool.Capacity,
					Used:     pool.Used,
					Uuid:     pool.Uuid,
				}
				poolInfos = append(poolInfos, pi)
			}
		} else {
			err = fmt.Errorf("nil response for ListPools on %s", address)
			logf.Log.Info("listPool", "error", err)
		}
	} else {
		err = niceError(err)
		logf.Log.Info("listPool", "error", err)
	}
	return poolInfos, err
}

func GetPool(name, addr string) (*v1MayastorPool, error) {
	poolInfo, err := listPool(addr)
	if err != nil {
		return nil, err
	}
	for _, pool := range poolInfo {
		if pool.Name == name {
			return &pool, nil
		}
	}
	return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
}

// ListPools given a list of node ip addresses, enumerate the set of pools on mayastor using gRPC on each of those nodes
// returns accumulated errors if gRPC communication failed.
func ListPools(addrs []string) ([]v1MayastorPool, error) {
	var accErr error
	var poolInfos []v1MayastorPool
	for _, address := range addrs {
		poolInfo, err := listPool(address)
		if err == nil {
			poolInfos = append(poolInfos, poolInfo...)
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return poolInfos, accErr
}

func DestroyAllPools(addrs []string) error {

	for _, addr := range addrs {
		poolInfo, err := listPool(addr)
		if err != nil {
			return err
		}
		if len(poolInfo) == 0 {
			continue
		}
		for _, pool := range poolInfo {
			err = DestroyPool(pool.Name, addr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DestroyPool(name, address string) error {
	var err error
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("destroyPool", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("destroyPool", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewPoolRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.DestroyPool(ctx, &mayastorGrpc.DestroyPoolRequest{Name: name})
	return err
}
