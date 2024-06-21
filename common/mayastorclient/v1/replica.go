package v1

import (
	"context"
	"fmt"
	"io"
	"time"

	mayastorGrpc "github.com/openebs/openebs-e2e/common/mayastorclient/v1/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type v1MayastorReplica struct {
	Name     string                     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                                   // name of the replica
	Uuid     string                     `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`                                   // uuid of the replica
	Pooluuid string                     `protobuf:"bytes,3,opt,name=pooluuid,proto3" json:"pooluuid,omitempty"`                           // uuid of the pool on which replica is present
	Size     uint64                     `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`                                  // size of the replica in bytes
	Thin     bool                       `protobuf:"varint,5,opt,name=thin,proto3" json:"thin,omitempty"`                                  // thin provisioning
	Share    mayastorGrpc.ShareProtocol `protobuf:"varint,6,opt,name=share,proto3,enum=mayastor.v1.ShareProtocol" json:"share,omitempty"` // protocol used for exposing the replica
	Uri      string                     `protobuf:"bytes,7,opt,name=uri,proto3" json:"uri,omitempty"`                                     // uri under which the replica is accessible by nexus
	Poolname string                     `protobuf:"bytes,8,opt,name=poolname,proto3" json:"poolname,omitempty"`                           // name of the pool on which replica is present
}

func (msr v1MayastorReplica) Description() string {
	return fmt.Sprintf("Uuid=%s; Pool=%s; Thin=%v; Size=%d; Share=%s; Uri=%s;",
		msr.Uuid, msr.Poolname, msr.Thin, msr.Size, msr.Share, msr.Uri)
}

func (msr v1MayastorReplica) GetUuid() string {
	return msr.Uuid
}

func (msr v1MayastorReplica) GetPool() string {
	return msr.Poolname
}

func (msr v1MayastorReplica) GetThin() bool {
	return msr.Thin
}

func (msr v1MayastorReplica) GetSize() uint64 {
	return msr.Size
}
func (msr v1MayastorReplica) GetShareString() string {
	return msr.Share.String()
}

func (msr v1MayastorReplica) GetUri() string {
	return msr.Uri
}

func (msr v1MayastorReplica) GetName() string {
	return msr.Name
}

func (msr v1MayastorReplica) GetPoolUuid() string {
	return msr.Pooluuid
}

func listReplica(address string) ([]v1MayastorReplica, error) {
	var replicaInfos []v1MayastorReplica
	var err error
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("listReplica connect failure", "address", address, "addrPort", addrPort, "error", err)
		return replicaInfos, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("listReplicas", "error on close", err)
		}
	}(conn)

	c := mayastorGrpc.NewReplicaRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var response *mayastorGrpc.ListReplicasResponse
	retryBackoffOnUnavailable(func() error {
		response, err = c.ListReplicas(ctx, &mayastorGrpc.ListReplicaOptions{})
		return err
	})

	if err == nil {
		if response != nil {
			for _, replica := range response.Replicas {
				ri := v1MayastorReplica{
					Name:     replica.Name,
					Uuid:     replica.Uuid,
					Poolname: replica.Poolname,
					Thin:     replica.Thin,
					Size:     replica.Size,
					Share:    replica.Share,
					Uri:      replica.Uri,
					Pooluuid: replica.Uuid,
				}
				replicaInfos = append(replicaInfos, ri)
			}
		} else {
			err = fmt.Errorf("nil response for ListReplicas on %s", address)
			logf.Log.Info("listReplicas", "error", err)
		}
	} else {
		err = niceError(err)
		logf.Log.Info("listReplicas", "error", err)
	}
	return replicaInfos, err
}

// RmReplica remove a replica identified by node and uuid
func RmReplica(address string, uuid string) error {
	logf.Log.Info("RmReplica", "address", address, "UUID", uuid)
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("rmReplicas", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("RmReplicas", "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewReplicaRpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := mayastorGrpc.DestroyReplicaRequest{Uuid: uuid}
	retryBackoff(func() error {
		_, err = c.DestroyReplica(ctx, &req)
		return err
	})

	return niceError(err)
}

// CreateReplicaExt create a replica on a mayastor node
func CreateReplicaExt(address string, uuid string, size uint64, pool string, thin bool) error {
	shareProto := mayastorGrpc.ShareProtocol_NVMF
	logf.Log.Info("CreateReplica", "address", address, "UUID", uuid, "size", size, "pool", pool, "Thin", thin, "Share", shareProto)
	addrPort := getAddrPort(address)
	var err error

	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("createReplica", "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("CreateReplicaExt", "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewReplicaRpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := mayastorGrpc.CreateReplicaRequest{
		Name:     fmt.Sprintf("%s-%s", pool, uuid),
		Uuid:     uuid,
		Size:     size,
		Thin:     thin,
		Pooluuid: pool,
		Share:    shareProto,
	}

	retryBackoff(func() error {
		_, err = c.CreateReplica(ctx, &req)
		return err
	})

	return niceError(err)
}

// CreateReplica create a replica on a mayastor node, with parameters
//
//	 thin fixed to false and share fixed to NVMF.
//	Other parameters must be specified
func CreateReplica(address string, uuid string, size uint64, pool string) error {
	return CreateReplicaExt(address, uuid, size, pool, false)
}

// ListReplicas given a list of node ip addresses, enumerate the set of replicas on mayastor using gRPC on each of those nodes
// returns accumulated errors if gRPC communication failed.
func ListReplicas(addrs []string) ([]v1MayastorReplica, error) {
	var accErr error
	var replicaInfos []v1MayastorReplica
	for _, address := range addrs {
		replicaInfo, err := listReplica(address)
		if err == nil {
			replicaInfos = append(replicaInfos, replicaInfo...)
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return replicaInfos, accErr
}

// RmNodeReplicas given a list of node ip addresses, delete the set of replicas on mayastor using gRPC on each of those nodes
// returns errors if gRPC communication failed.
func RmNodeReplicas(addrs []string) error {
	var accErr error
	for _, address := range addrs {
		replicaInfos, err := listReplica(address)
		if err == nil {
			for _, replicaInfo := range replicaInfos {
				err = RmReplica(address, replicaInfo.Uuid)
			}
		}
		if err != nil {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return accErr
}

// FindReplicas given a list of node ip addresses, enumerate the set of replicas on mayastor using gRPC on each of those nodes
// returns accumulated errors if gRPC communication failed.
func FindReplicas(uuid string, addrs []string) ([]v1MayastorReplica, error) {
	var accErr error
	var replicaInfos []v1MayastorReplica
	for _, address := range addrs {
		replicaInfo, err := listReplica(address)
		if err == nil {
			for _, repl := range replicaInfo {
				if repl.Uuid == uuid {
					replicaInfos = append(replicaInfos, repl)
				}
			}
		} else {
			if accErr != nil {
				accErr = fmt.Errorf("%v;%v", accErr, err)
			} else {
				accErr = err
			}
		}
	}
	return replicaInfos, accErr
}

// WipeReplica fill a replica with zeroes
func WipeReplica(address string, replicaUuid string, poolName string) error {
	desc := fmt.Sprintf("addr:%s uuid:%s pool:%s", address, replicaUuid, poolName)
	logf.Log.Info("WipeReplica", "address", address, "UUID", replicaUuid, "poolName", poolName)
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("WipeReplica", "target", desc, "error", err)
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("WipeReplica", "target", desc, "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewTestRpcClient(conn)

	reqPoolName := mayastorGrpc.WipeReplicaRequest_PoolName{
		PoolName: poolName,
	}
	wipeOpts := mayastorGrpc.WipeOptions{
		WipeMethod: mayastorGrpc.WipeOptions_WRITE_ZEROES,
		//		WritePattern: ,
	}
	streamWipeOpts := mayastorGrpc.StreamWipeOptions{
		Options:   &wipeOpts,
		ChunkSize: 0, // 0 => chunk size defaults tp replica size
	}
	var req = mayastorGrpc.WipeReplicaRequest{
		Uuid:        replicaUuid,
		Pool:        &reqPoolName,
		WipeOptions: &streamWipeOpts,
	}
	startTime := time.Now()
	logf.Log.Info("Start WipeReplica .............", "target", desc)
	stream, err := c.WipeReplica(context.TODO(), &req)
	if err == nil {
		for {
			resp, rcvErr := stream.Recv()
			if rcvErr == io.EOF {
				break
			}
			if err != nil {
				logf.Log.Info("unexpected error", "target", desc, "error", err)
				return err
			} else {
				if resp == nil {
					logf.Log.Info("unexpected nil return, parameters may be invalid", "resp", resp)
					return fmt.Errorf("gRPC call returned nil response")
				} else {
					logf.Log.Info("got", "target", desc, "resp", resp, "elapsed", time.Since(startTime))
				}
			}
		}
	}
	logf.Log.Info("WipeReplica Finished", "target", desc, "elapsed", time.Since(startTime))

	return niceError(err)
}

// ChecksumReplica calculate the checksum of a replica
func ChecksumReplica(address string, replicaUuid string, poolName string) (uint32, error) {
	var cksum uint32
	var cksumSet = false
	desc := fmt.Sprintf("addr:%s uuid:%s pool:%s", address, replicaUuid, poolName)
	logf.Log.Info("ChecksumReplica", "address", address, "UUID", replicaUuid, "poolName", poolName)
	addrPort := getAddrPort(address)
	conn, err := grpc.Dial(addrPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logf.Log.Info("ChecksumReplica", "target", desc, "error", err)
		return cksum, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logf.Log.Info("ChecksumReplica", "target", desc, "error on close", err)
		}
	}(conn)
	c := mayastorGrpc.NewTestRpcClient(conn)

	var features *mayastorGrpc.TestFeatures
	features, err = c.GetFeatures(context.TODO(), &emptypb.Empty{})
	if err != nil {
		logf.Log.Info("failed to retrieve test features set", "err", err)
		return cksum, err
	}
	if features == nil || features.CksumAlgs == nil {
		return cksum, fmt.Errorf("unsupported")
	}

	var req = mayastorGrpc.WipeReplicaRequest{
		Uuid: replicaUuid,
		Pool: &mayastorGrpc.WipeReplicaRequest_PoolName{
			PoolName: poolName,
		},
		WipeOptions: &mayastorGrpc.StreamWipeOptions{
			Options: &mayastorGrpc.WipeOptions{
				WipeMethod:   mayastorGrpc.WipeOptions_CHECKSUM,
				WritePattern: nil,
				CksumAlg:     mayastorGrpc.WipeOptions_Crc32c,
			},
			ChunkSize: 0,
		},
	}
	startTime := time.Now()
	logf.Log.Info("ChecksumReplica ", "req.WipeOptions", req.WipeOptions, "req.WipeOptions.Options", req.WipeOptions.Options)
	logf.Log.Info("Start ChecksumReplica .............", "target", desc, "wipeM", req.GetWipeOptions().GetOptions().GetWipeMethod(),
		"cksa", req.GetWipeOptions().GetOptions().GetCksumAlg())
	stream, err := c.WipeReplica(context.TODO(), &req)
	if err == nil {
		for {
			resp, rcvErr := stream.Recv()
			if rcvErr == io.EOF {
				break
			}
			if err != nil {
				logf.Log.Info("unexpected error", "target", desc, "error", err)
				return cksum, err
			} else {
				if resp == nil {
					logf.Log.Info("unexpected nil return, parameters may be invalid", "resp", resp)
					return cksum, fmt.Errorf("gRPC call returned nil response")
				} else {
					logf.Log.Info("got", "resp.Checksum", resp.Checksum, "target", desc, "resp", resp, "elapsed", time.Since(startTime))
					if resp.Checksum != nil && resp.RemainingBytes == 0 {
						cksum = resp.GetCrc32()
						cksumSet = true
					}
				}
			}
		}
	}
	if !cksumSet {
		return cksum, fmt.Errorf("failed to retrieve replica checksum")
	}
	logf.Log.Info("ChecksumReplica Finished", "target", desc, "elapsed", time.Since(startTime))

	return cksum, niceError(err)
}
