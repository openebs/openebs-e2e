package v1

import (
	"context"
	"fmt"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ShareBdev share a bdev with uuid
func ShareBdev(address string, bdevUuid string) (string, error) {
	logf.Log.Info("ShareBdev", "address", address, "bdevUuid", bdevUuid)
	var bdevShareUri string
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("ShareBdev", "error", err)
		return bdevShareUri, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("ShareBdev", "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewBdevRpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := mayastorGrpc.BdevShareRequest{
		Name:     bdevUuid,
		Protocol: mayastorGrpc.ShareProtocol_NVMF,
	}
	var response *mayastorGrpc.BdevShareResponse
	retryBackoff(func() error {
		response, err = c.Share(ctx, &req)
		return err
	})
	if response == nil {
		return bdevShareUri, fmt.Errorf("failed to share bdev, response: %v", response)
	}
	bdevShareUri = response.Bdev.GetShareUri()

	logf.Log.Info("bdev", "uuid", bdevUuid, "sharedUri", bdevShareUri)
	return bdevShareUri, niceError(err)
}

// UnshareBdev unshare a bdev with uuid
func UnshareBdev(address string, bdevUuid string) error {
	logf.Log.Info("UnshareBdev", "address", address, "bdevUuid", bdevUuid)
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("UnshareBdev", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("UnshareBdev", "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewBdevRpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := mayastorGrpc.BdevUnshareRequest{
		Name: bdevUuid,
	}
	retryBackoff(func() error {
		_, err = c.Unshare(ctx, &req)
		return err
	})

	return niceError(err)
}
