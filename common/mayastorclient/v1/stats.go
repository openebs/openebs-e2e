package v1

import (
	"context"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func ResetIOStats(address string) error {

	logf.Log.Info("reset io stats", "address", address)
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("ResetIOStats", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("ResetIOStats", "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewStatsRpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	retryBackoff(func() error {
		_, err = c.ResetIoStats(ctx, &emptypb.Empty{})
		return err
	})

	return niceError(err)
}
