package v0

import (
	"context"
	"fmt"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v0/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ShareBdev share a bdev with uuid
func ShareBdev(address string, bdevUuid string) (string, error) {
	logf.Log.Info("ShareBdev", "address", address, "bdevUuid", bdevUuid)
	var bdevUri string
	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("ShareBdev", "error", err)
		return bdevUri, err
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
		Name: bdevUuid,
	}
	var response *mayastorGrpc.BdevShareReply
	retryBackoff(func() error {
		response, err = c.Share(ctx, &req)
		return err
	})

	if response == nil {
		return bdevUri, fmt.Errorf("failed to share bdev, response: %v", response)
	}
	bdevUri = response.GetUri()

	return bdevUri, niceError(err)
}

// UnshareBdev unshare a bdev with uuid
func UnshareBdev(address string, bdevUuid string) error {
	logf.Log.Info("UnshareBdev", "address", address, "bdevUuid", bdevUuid)
	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
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

	req := mayastorGrpc.CreateReply{
		Name: bdevUuid,
	}
	retryBackoff(func() error {
		_, err = c.Unshare(ctx, &req)
		return err
	})

	return niceError(err)
}
