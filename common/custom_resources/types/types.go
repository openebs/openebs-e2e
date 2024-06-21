package types

import "reflect"

// Types to support DiskPool CRD abstraction to support multiple
// CRD versions

// DiskPool interface to access all DiskPool CRs
type DiskPool interface {
	String() string
	GetType() reflect.Type
	GetName() string
	GetStatusCapacity() uint64
	GetStatusUsed() uint64
	CompareStatus(other *DiskPool) bool
	GetFinalizers() []string
	SetFinalizers(finalizers []string) (DiskPool, error)
	GetSpecDisks() []string
	SetSpecDisks(disks []string) (DiskPool, error)
	GetSpecNode() string
	SetSpecNode(node string) (DiskPool, error)
	GetCRStatus() string
	GetPoolStatus() string
}

// DiskPoolFunctions interface to implement support for a DiskPool CRD version
type DiskPoolFunctions interface {
	CreateMsPool(poolName string, node string, disks []string) (DiskPool, error)
	CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (DiskPool, error)
	GetMsPool(poolName string) (DiskPool, error)
	DeleteMsPool(poolName string) error
	ListMsPoolCrs() ([]DiskPool, error)
}
