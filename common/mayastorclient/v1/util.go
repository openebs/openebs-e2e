package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/openebs/openebs-e2e/common/k8s_portforward"
	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"

	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func isRetryErr(err error) bool {
	status, ok := grpcStatus.FromError(err)
	if ok {
		switch status.Code() {
		case grpcCodes.DeadlineExceeded, grpcCodes.Unavailable:
			return true
		default:
			return false
		}
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

func getAddrPort(address string) string {
	return k8s_portforward.TryPortForwardNode(address, mayastorPort)
}

func mayastorInfo(address string) (*mayastorGrpc.MayastorInfoResponse, error) {
	var err error
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	c := mayastorGrpc.NewHostRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.GetMayastorInfo(ctx, &emptypb.Empty{})
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
	if !isDeadlineExceeded(f()) {
		return
	}
	timeout := 5 * time.Second
	for i := 0; i < 6; i++ {
		logf.Log.Info("retrying gRPC call", "after", timeout)
		time.Sleep(timeout)
		timeout *= 2
		if !isDeadlineExceeded(f()) {
			return
		}
	}
}

// retry a function upto 6 times with incremental backoff,
// starting at 1 seconds upt 32 seconds (63 seconds total)
// if the error(s) returned are
// is deadline_exceeded or unavailable.
func retryBackoffOnUnavailable(f func() (err error)) {
	if !isRetryErr(f()) {
		return
	}
	timeout := 5 * time.Second
	for i := 0; i < 6; i++ {
		logf.Log.Info("retrying gRPC call", "after", timeout)
		time.Sleep(timeout)
		timeout *= 2
		if !isRetryErr(f()) {
			return
		}
	}
}
