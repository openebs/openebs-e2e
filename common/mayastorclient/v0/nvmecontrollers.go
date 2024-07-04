package v0

import (
	"context"
	"fmt"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v0/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type v0NvmeController struct {
	Name    string                           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	State   mayastorGrpc.NvmeControllerState `protobuf:"varint,2,opt,name=state,proto3,enum=mayastor.NvmeControllerState" json:"state,omitempty"`
	Size    uint64                           `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	BlkSize uint32                           `protobuf:"varint,4,opt,name=blk_size,json=blkSize,proto3" json:"blk_size,omitempty"`
}

func (nvmectlr v0NvmeController) GetString() string {
	return fmt.Sprintf("Name=%s; State=%s; ; Size=%d; BlkSize=%d;",
		nvmectlr.Name, nvmectlr.State, nvmectlr.Size, nvmectlr.BlkSize)
}

func (nvmectlr v0NvmeController) GetName() string {
	return nvmectlr.Name
}

func (nvmectlr v0NvmeController) GetStateString() string {
	return nvmectlr.State.String()
}

func (nvmectlr v0NvmeController) GetSize() uint64 {
	return nvmectlr.Size
}

func (nvmectlr v0NvmeController) GetBlkSize() uint32 {
	return nvmectlr.BlkSize
}

func listNvmeController(address string) ([]v0NvmeController, error) {
	var nvmeControllers []v0NvmeController
	var err error
	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("listReplica", "error", err)
		return nvmeControllers, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("listReplicas", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewMayastorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var response *mayastorGrpc.ListNvmeControllersReply
	retryBackoff(func() error {
		response, err = c.ListNvmeControllers(ctx, &null)
		return err
	})

	if err == nil {
		if response != nil {
			for _, nvmeController := range response.Controllers {
				nc := v0NvmeController{
					Name:    nvmeController.Name,
					State:   nvmeController.State,
					Size:    nvmeController.Size,
					BlkSize: nvmeController.BlkSize,
				}
				nvmeControllers = append(nvmeControllers, nc)
			}
		} else {
			err = fmt.Errorf("nil response for ListReplicas on %s", address)
			logf.Log.Info("listReplicas", "error", err)
		}
	} else {
		logf.Log.Info("listReplicas", "error", err)
	}
	return nvmeControllers, err
}

// ListNvmeControllers given a list of node ip addresses, enumerate the set of nvmeControllers on mayastor using gRPC on each of those nodes
// returns accumulated errors if gRPC communication failed.
func ListNvmeControllers(addrs []string) ([]v0NvmeController, error) {
	var accErr error
	var nvmeControllers []v0NvmeController
	for _, address := range addrs {
		nvmeController, err := listNvmeController(address)
		if err == nil {
			nvmeControllers = append(nvmeControllers, nvmeController...)
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return nvmeControllers, accErr
}
