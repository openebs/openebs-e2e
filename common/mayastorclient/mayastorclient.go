package mayastorclient

import (
	"fmt"
	"sync"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type GrpcInterface interface {
	Version() string

	// Nexus abstraction
	ListNexuses(addrs []string) ([]MayastorNexus, error)
	FaultNexusChild(address string, Uuid string, Uri string) error
	FindNexus(uuid string, addrs []string) (*MayastorNexus, error)
	GetRebuildHistory(uuid string, addrs string) (RebuildHistory, error)            // uuid of the nexus
	GetRebuildStats(uuid string, dstUri string, addrs string) (RebuildStats, error) // uuid of the nexus

	// Nvme controller abstraction
	ListNvmeControllers(addrs []string) ([]NvmeController, error)

	// Disk pool abstraction
	GetPool(name, addr string) (MayastorPool, error)
	ListPools(addrs []string) ([]MayastorPool, error)
	DestroyAllPools(addrs []string) error
	DestroyPool(name, addr string) error

	// Replica abstraction
	RmReplica(address string, uuid string) error
	CreateReplicaExt(address string, uuid string, size uint64, pool string, thin bool) error
	CreateReplica(address string, uuid string, size uint64, pool string) error
	ListReplicas(addrs []string) ([]MayastorReplica, error)
	RmNodeReplicas(addrs []string) error
	FindReplicas(uuid string, addrs []string) ([]MayastorReplica, error)
	WipeReplica(address string, replicaUUID string, poolName string) error
	ChecksumReplica(address string, replicaUUID string, poolName string) (uint32, error)

	//bdev abstraction
	ShareBdev(address string, bdevUuid string) (string, error)
	UnshareBdev(address string, bdevUuid string) error

	// grpc connect abstraction
	CheckAndSetConnect(nodes []string) error
	CanConnect() bool

	// io stats
	ResetIOStats(address string) error
}

// The default grpc interface
var defaultGrpcIfc GrpcInterface
var once sync.Once

// Initialise should be called before other functions which use the
// default grpc interface namely all public functions in this package
// except GetGrpcIfc
func Initialise(nodes []string) {
	once.Do(func() {
		var grpcVers = []string{"v1", "v0"}
		for _, ver := range grpcVers {
			grpcIface, err := GetGrpcIfc(ver)
			if err == nil {
				err = grpcIface.CheckAndSetConnect(nodes)
				if err == nil {
					logf.Log.Info("*** Using gRPC ", "version", ver)
					defaultGrpcIfc = grpcIface
					return
				}
			}
		}
	})
}

// GetGrpcIfc returns a mayastor grpc interface object for the requested version or nil
func GetGrpcIfc(version string) (GrpcInterface, error) {
	switch version {
	case "V0", "v0":
		return grpcV0{}, nil
	case "V1", "v1":
		return grpcV1{}, nil
	default:
		return nil, fmt.Errorf("unknown version")
	}
}

func Version() string {
	if defaultGrpcIfc == nil {
		return ""
	}
	return defaultGrpcIfc.Version()
}

func ListNexuses(addrs []string) ([]MayastorNexus, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.ListNexuses(addrs)
}

func FaultNexusChild(address string, Uuid string, Uri string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.FaultNexusChild(address, Uuid, Uri)
}

func FindNexus(uuid string, addrs []string) (*MayastorNexus, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.FindNexus(uuid, addrs)
}

func ListNvmeControllers(addrs []string) ([]NvmeController, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.ListNvmeControllers(addrs)
}

func GetPool(name, addr string) (MayastorPool, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.GetPool(name, addr)
}

func ListPools(addrs []string) ([]MayastorPool, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.ListPools(addrs)
}

/*
func DestroyAllPools(addrs []string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.DestroyAllPools(addrs)
}

func DestroyPool(name, addr string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.DestroyPool(name, addr)
}
*/

func RmReplica(address string, uuid string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.RmReplica(address, uuid)
}

func CreateReplicaExt(address string, uuid string, size uint64, pool string, thin bool) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.CreateReplicaExt(address, uuid, size, pool, thin)
}

func CreateReplica(address string, uuid string, size uint64, pool string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.CreateReplica(address, uuid, size, pool)
}

func ListReplicas(addrs []string) ([]MayastorReplica, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.ListReplicas(addrs)
}

func RmNodeReplicas(addrs []string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.RmNodeReplicas(addrs)
}

func FindReplicas(uuid string, addrs []string) ([]MayastorReplica, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.FindReplicas(uuid, addrs)
}

func GetRebuildHistory(uuid string, addrs string) (RebuildHistory, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.GetRebuildHistory(uuid, addrs)
}

func GetRebuildStats(uuid string, dstUri string, addrs string) (RebuildStats, error) {
	if defaultGrpcIfc == nil {
		return nil, fmt.Errorf("mayastor client package has not been initialised")
	}
	return defaultGrpcIfc.GetRebuildStats(uuid, dstUri, addrs)
}

func CanConnect() bool {
	if defaultGrpcIfc != nil {
		return defaultGrpcIfc.CanConnect()
	}
	return false
}

func WipeReplica(address string, replicaUUID string, poolName string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	logf.Log.Info("WipeReplica ", "address", address, "replicaUUID", replicaUUID, "poolName", poolName)
	return defaultGrpcIfc.WipeReplica(address, replicaUUID, poolName)
}

func ChecksumReplica(address string, replicaUUID string, poolName string) (uint32, error) {
	if defaultGrpcIfc == nil {
		return 0, fmt.Errorf("mayastor client package has not been initialised")
	}
	logf.Log.Info("ChecksumReplica ", "address", address, "replicaUUID", replicaUUID, "poolName", poolName)
	return defaultGrpcIfc.ChecksumReplica(address, replicaUUID, poolName)
}

func ShareBdev(address string, bdevUuid string) (string, error) {
	if defaultGrpcIfc == nil {
		return "", fmt.Errorf("mayastor client package has not been initialised")
	}
	logf.Log.Info("ShareBdev ", "address", address, "bdev", bdevUuid)
	return defaultGrpcIfc.ShareBdev(address, bdevUuid)
}

func UnshareBdev(address string, bdevUuid string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	logf.Log.Info("UnshareBdev ", "address", address, "bdevUuid", bdevUuid)
	return defaultGrpcIfc.UnshareBdev(address, bdevUuid)
}

func ResetIOStats(address string) error {
	if defaultGrpcIfc == nil {
		return fmt.Errorf("mayastor client package has not been initialised")
	}
	logf.Log.Info("reset io stats")
	return defaultGrpcIfc.ResetIOStats(address)
}
