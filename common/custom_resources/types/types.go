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
	GetPoolReadyStatus() string
	GetPoolReadyReason() string
	GetPoolErrorCount() uint64
	GetPoolAlertStatus() string
	GetPoolErrorCode() string
	GetPoolStatus() string
	GetClusterSize() string
	SetClusterSize(clusterSize string) (DiskPool, error)
	GetSpecEncryptionSecret() string
	SetSpecEncryptionSecret(secretName string) (DiskPool, error)
	IsPoolEncrypted() bool
	// Advanced features (v1beta3+)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string) (DiskPool, error)
	GetMaxExpansion() string
	SetMaxExpansion(maxExpansion string) (DiskPool, error)
	// Optional: returns status maxExpandableSize as a human-readable string (e.g., "255.8 GiB").
	// May return empty string for CR versions that don't expose it.
	GetStatusMaxExpandableSize() string
}

// DiskPoolFunctions interface to implement support for a DiskPool CRD version
type DiskPoolFunctions interface {
	CreateMsPool(poolName string, node string, disks []string) (DiskPool, error)
	CreateMsPoolWithTopologySpec(poolName string, node string, disks []string, labels map[string]string) (DiskPool, error)
	CreateMsPoolWithEncryption(poolName string, node string, disks []string, encryptionSecretName string) (DiskPool, error)
	CreateMsPoolWithClusterSize(poolName string, node string, disks []string, clusterSize string) (DiskPool, error)
	GetMsPool(poolName string) (DiskPool, error)
	DeleteMsPool(poolName string) error
	ListMsPoolCrs() ([]DiskPool, error)
	// Advanced features (v1beta3+)
	CreateMsPoolWithMaxSize(poolName string, node string, disks []string, maxSize string) (DiskPool, error)
	CreateMsPoolWithMaxSizeAndClusterSize(poolName string, node string, disks []string, maxSize string, clusterSize string) (DiskPool, error)
	CreateMsPoolWithEncryptionAndMaxSize(poolName string, node string, disks []string, encryptionSecretName string, maxSize string, clusterSize string) (DiskPool, error)
	AnnotatePoolForExpansion(poolName string) error
	VerifyPoolCapacityAndMaxExpansion(poolName string, expectedCapacity uint64, expectedMaxExpansion string) error
	AnnotateOfflinePoolForDelete(poolName string, opts ...string) error
}
