package v0

import (
	"context"
	"fmt"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v0/protobuf"

	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcStatus "google.golang.org/grpc/status"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func isDeadlineExceeded(err error) bool {
	status, ok := grpcStatus.FromError(err)
	if !ok {
		return false
	}
	if status.Code() == grpcCodes.DeadlineExceeded {
		return true
	}
	return false
}

func niceError(err error) error {
	if err != nil {
		if isDeadlineExceeded(err) {
			// stop huge print out of error on deadline exceeded
			return grpcStatus.Error(grpcCodes.DeadlineExceeded, fmt.Sprintf("%v", context.DeadlineExceeded))
		}
	}
	return err
}

var canConnect bool = false

func mayastorInfo(address string) (*mayastorGrpc.MayastorInfoRequest, error) {
	var err error
	addrPort := fmt.Sprintf("%s:%d", address, mayastorPort)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	c := mayastorGrpc.NewMayastorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.GetMayastorInfo(ctx, &null)
	return info, err
}

// CheckAndSetConnect call to cache connectable state to Mayastor instances on the cluster under test
// Just dialing does not work, we need to make a simple gRPC call  (GetMayastorInfo)
func CheckAndSetConnect(nodes []string) error {
	var connErr error
	logf.Log.Info("Checking gRPC connections to Mayastor on", "nodes", nodes)
	if len(nodes) != 0 {
		for _, node := range nodes {
			info, err := mayastorInfo(node)
			logf.Log.Info("", "mayastorInfo", info)
			if err != nil || info == nil {
				connErr = fmt.Errorf("e:%v, i:%v", err, info)
				logf.Log.Info("gRPC connect failed", "node", node, "err", connErr)
				break
			}
		}
		canConnect = connErr == nil
	}
	return connErr
}

// CanConnect retrieve the cached connectable state to Mayastor instances on the cluster under test
func CanConnect() bool {
	return canConnect
}

const ctxTimeout = 30 * time.Second

// retry a function upto 6 times with exponential backoff,
// starting at 5 seconds if the error(s) returned are
// is deadline_exceeded.
func retryBackoff(f func() (err error)) {
	timeout := 5 * time.Second
	for i := 0; i < 6; i++ {
		if !isDeadlineExceeded(f()) {
			return
		}
		logf.Log.Info("retrying gRPC call", "after", timeout)
		time.Sleep(timeout)
		timeout *= 2
	}
}
