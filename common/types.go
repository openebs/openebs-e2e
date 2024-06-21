package common

import "fmt"

type ShareProto string

const (
	ShareProtoNvmf  ShareProto = "nvmf"
	ShareProtoIscsi ShareProto = "iscsi"
)

type FileSystemType string

const (
	NoneFsType  FileSystemType = ""
	Ext4FsType  FileSystemType = "ext4"
	XfsFsType   FileSystemType = "xfs"
	BtrfsFsType FileSystemType = "btrfs"
)

type VolumeType int

const (
	VolFileSystem VolumeType = iota
	VolRawBlock   VolumeType = iota
	VolTypeNone   VolumeType = iota
)

func (volType VolumeType) String() string {
	switch volType {
	case VolFileSystem:
		return "FileSystem"
	case VolRawBlock:
		return "RawBlock"
	default:
		return "Unknown"
	}
}

type StsAffinityGroup string

const (
	StsAffinityGroupDisable StsAffinityGroup = "false"
	StsAffinityGroupEnable  StsAffinityGroup = "true"
)

type ProvisioningType int

const (
	ThinProvisioning  ProvisioningType = iota
	ThickProvisioning ProvisioningType = iota
)

func (provisioningType ProvisioningType) String() string {
	switch provisioningType {
	case ThickProvisioning:
		return "thick"
	case ThinProvisioning:
		return "thin"
	default:
		return "thick"
	}
}

type CloneFsIdAsVolumeIdType int

const (
	CloneFsIdAsVolumeIdNone    CloneFsIdAsVolumeIdType = iota
	CloneFsIdAsVolumeIdEnable  CloneFsIdAsVolumeIdType = iota
	CloneFsIdAsVolumeIdDisable CloneFsIdAsVolumeIdType = iota
)

func (CloneFsId CloneFsIdAsVolumeIdType) String() string {
	switch CloneFsId {
	case CloneFsIdAsVolumeIdEnable:
		return "enable"
	case CloneFsIdAsVolumeIdDisable:
		return "disable"
	case CloneFsIdAsVolumeIdNone:
		return "none"
	default:
		return ""
	}
}

type ReplicaTopologyChildState int

const (
	ChildStateOnline   ReplicaTopologyChildState = iota
	ChildStateUnknown  ReplicaTopologyChildState = iota
	ChildStateDegraded ReplicaTopologyChildState = iota
	ChildStateFaulted  ReplicaTopologyChildState = iota
)

func (replicaTopologyChildState ReplicaTopologyChildState) String() string {
	switch replicaTopologyChildState {
	case ChildStateOnline:
		return "Online"
	case ChildStateUnknown:
		return "Unknown"
	case ChildStateDegraded:
		return "Degraded"
	case ChildStateFaulted:
		return "Faulted"
	default:
		return ""
	}
}

type AllowVolumeExpansion int

const (
	AllowVolumeExpansionNone    AllowVolumeExpansion = iota
	AllowVolumeExpansionEnable  AllowVolumeExpansion = iota
	AllowVolumeExpansionDisable AllowVolumeExpansion = iota
)

func (CloneFsId AllowVolumeExpansion) String() string {
	switch CloneFsId {
	case AllowVolumeExpansionEnable:
		return "enable"
	case AllowVolumeExpansionDisable:
		return "disable"
	case AllowVolumeExpansionNone:
		return "none"
	default:
		return ""
	}
}

const MsvStatusStateOnline = "Online"

type MayastorVolume struct {
	Spec  MsvSpec  `json:"spec"`
	State MsvState `json:"state"`
}

type MsvSpec struct {
	Num_replicas  int           `json:"num_replicas"`
	Size          int64         `json:"size"`
	Status        string        `json:"status"`
	Target        SpecTarget    `json:"target"`
	Uuid          string        `json:"uuid"`
	Topology      Topology      `json:"topology"`
	Policy        Policy        `json:"policy"`
	Thin          bool          `json:"thin"`
	AsThin        bool          `json:"as_thin"`
	NumSnapshots  int32         `json:"num_snapshots"`
	ContentSource ContentSource `json:"content_source"`
	MaxSnapshots  int32         `json:"max_snapshots"`
}

type Policy struct {
	Self_heal bool `json:"self_heal"`
}
type SpecTarget struct {
	Protocol string `json:"protocol"`
	Node     string `json:"node"`
}

type Topology struct {
	NodeTopology Node_topology `json:"node_topology"`
	PoolTopology Pool_topology `json:"pool_topology"`
}
type Node_topology struct {
	Explicit Explicit `json:"explicit"`
}
type Pool_topology struct {
	Labelled Labelled `json:"labelled"`
}
type Labelled struct {
	Inclusion map[string]interface{} `json:"inclusion"`
	Exclusion map[string]interface{} `json:"exclusion"`
}

type Explicit struct {
	AllowedNodes   []string `json:"allowed_nodes"`
	PreferredNodes []string `json:"preferred_nodes"`
}

type ContentSource struct {
	Snapshot Snapshot `json:"snapshot"`
}

type Snapshot struct {
	Snapshot string `json:"snapshot"`
	Volume   string `json:"volume"`
}

type MsvState struct {
	Target          StateTarget     `json:"target"`
	Size            int64           `json:"size"`
	Status          string          `json:"status"`
	Uuid            string          `json:"uuid"`
	ReplicaTopology ReplicaTopology `json:"replica_topology"`
	Usage           Usage           `json:"usage"`
}

type ReplicaTopology map[string]Replica

// MsvReplica contains replica Uri along with uuid and replica details. In older mayastor volume schema, replica
// Uri was present which is not the case in current version
type MsvReplica struct {
	Uuid    string
	Uri     string
	Replica Replica
}

type ReplicaUsage struct {
	Capacity              int64 `json:"capacity"`
	Allocated             int64 `json:"allocated"`
	AllocatedSnapshots    int64 `json:"allocated_snapshots"`
	AllocatedAllSnapshots int64 `json:"allocated_all_snapshots"`
}

type Replica struct {
	Node        string       `json:"node"`
	Pool        string       `json:"pool"`
	State       string       `json:"state"`
	ChildStatus string       `json:"child-status"`
	Usage       ReplicaUsage `json:"usage"`
}

type StateTarget struct {
	Children  []TargetChild `json:"children"`
	DeviceUri string        `json:"deviceUri"`
	Node      string        `json:"node"`
	Rebuilds  int32         `json:"rebuilds"`
	Protocol  string        `json:"protocol"`
	Size      int64         `json:"size"`
	State     string        `json:"state"`
	Uuid      string        `json:"uuid"`
}

type TargetChild struct {
	State           string `json:"state"`
	Uri             string `json:"uri"`
	RebuildProgress *int32 `json:"rebuildProgress"`
}

type Usage struct {
	Capacity                int64 `json:"capacity"`
	Allocated               int64 `json:"allocated"`
	AllocatedReplica        int64 `json:"allocated_replica"`
	AllocatedSnapshots      int64 `json:"allocated_snapshots"`
	AllocatedAllSnapshots   int64 `json:"allocated_all_snapshots"`
	TotalAllocated          int64 `json:"total_allocated"`
	TotalAllocatedReplicas  int64 `json:"total_allocated_replicas"`
	TotalAllocatedSnapshots int64 `json:"total_allocated_snapshots"`
}

type MayastorVolumeInterface interface {
	GetMSV(uuid string) (*MayastorVolume, error)
	GetMsvNodes(uuid string) (string, []string)
	DeleteMsv(volName string) error
	ListMsvs() ([]MayastorVolume, error)
	SetMsvReplicaCount(uuid string, replicaCount int) error
	GetMsvState(uuid string) (string, error)
	GetMsvReplicas(volName string) ([]Replica, error)
	GetMsvReplicaTopology(volUuid, replicaUuid string) (Replica, error)
	GetMsvReplicaTopologies(volUuid string) (ReplicaTopology, error)
	GetMsvNexusChildren(volName string) ([]TargetChild, error)
	GetMsvNexusState(uuid string) (string, error)
	IsMsvPublished(uuid string) bool
	IsMsvDeleted(uuid string) bool
	CheckForMsvs() (bool, error)
	CheckAllMsvsAreHealthy() error
}

type MayastorNodeInterface interface {
	GetMSN(node string) (*MayastorNode, error)
	ListMsns() ([]MayastorNode, error)
}

type MayastorNode struct {
	Name  string            `json:"name"`
	Spec  MayastorNodeSpec  `json:"spec"`
	State MayastorNodeState `json:"state"`
}

type MayastorNodeSpec struct {
	GrpcEndpoint string `json:"grpcEndpoint"`
	ID           string `json:"id"`
	Node_nqn     string `json:"node_nqn"`
}

type MayastorNodeState struct {
	GrpcEndpoint string `json:"grpcEndpoint"`
	ID           string `json:"id"`
	Status       string `json:"status"`
	Node_nqn     string `json:"node_nqn"`
}

type MayastorPool struct {
	Name   string             `json:"name"`
	Spec   MayastorPoolSpec   `json:"spec"`
	Status MayastorPoolStatus `json:"status"`
}

type MayastorPoolSpec struct {
	Disks []string `json:"disks"`
	Node  string   `json:"node"`
}

type MayastorPoolStatus struct {
	Avail     uint64           `json:"avail"`
	Capacity  uint64           `json:"capacity"`
	Disks     []string         `json:"disks"`
	Reason    string           `json:"reason"`
	Spec      MayastorPoolSpec `json:"spec"`
	State     string           `json:"state"`
	Used      uint64           `json:"used"`
	Committed uint64           `json:"committed"`
}

type SnapshotMetadata struct {
	Status              string                           `json:"status"`
	Timestamp           string                           `json:"timestamp"`
	TxnID               string                           `json:"txn_id"`
	Transactions        map[string][]SnapshotTransaction `json:"transactions"`
	NumSnapshotReplicas int                              `json:"num_snapshot_replicas"`
	NumRestores         int                              `json:"num_restores"`
}

type SnapshotTransaction struct {
	UUID     string `json:"uuid"`
	SourceID string `json:"source_id"`
	Status   string `json:"status"`
}

type SnapshotSpec struct {
	UUID         string `json:"uuid"`
	SourceVolume string `json:"source_volume"`
}

type SnapshotState struct {
	UUID             string            `json:"uuid"`
	AllocatedSize    int64             `json:"allocated_size"`
	SourceVolume     string            `json:"source_volume"`
	Timestamp        string            `json:"timestamp"`
	ReadyAsSource    bool              `json:"ready_as_source"`
	ReplicaSnapshots []ReplicaSnapshot `json:"replica_snapshots"`
}

type ReplicaSnapshot struct {
	Online  OnlineSnapshot  `json:"online"`
	Offline OfflineSnapshot `json:"offline"`
}

type OnlineSnapshot struct {
	UUID                 string `json:"uuid"`
	SourceID             string `json:"source_id"`
	PoolID               string `json:"pool_id"`
	PoolUUID             string `json:"pool_uuid"`
	Timestamp            string `json:"timestamp"`
	Size                 int64  `json:"size"`
	AllocatedSize        int64  `json:"allocated_size"`
	PredecessorAllocSize int64  `json:"predecessor_alloc_size"`
}

type OfflineSnapshot struct {
	UUID     string `json:"uuid"`
	SourceID string `json:"source_id"`
	PoolID   string `json:"pool_id"`
	PoolUUID string `json:"pool_uuid"`
}

type SnapshotSchema struct {
	Definition SnapshotDefinition `json:"definition"`
	State      SnapshotState      `json:"state"`
}

type SnapshotDefinition struct {
	Metadata SnapshotMetadata `json:"metadata"`
	Spec     SnapshotSpec     `json:"spec"`
}

type MayastorPoolInterface interface {
	GetMsPool(poolName string) (*MayastorPool, error)
	ListMsPools() ([]MayastorPool, error)
}

type ErrorAccumulator struct {
	errs []error
}

func (acc *ErrorAccumulator) Accumulate(err error) {
	if err != nil {
		acc.errs = append(acc.errs, err)
	}
}

func (acc *ErrorAccumulator) GetError() error {
	var err error
	for _, e := range acc.errs {
		if err != nil {
			err = fmt.Errorf("%w; %v", err, e)
		} else {
			err = e
		}
	}
	return err
}

type CmpReplicas int

const (
	CmpReplicasMatch    CmpReplicas = iota
	CmpReplicasMismatch CmpReplicas = iota
	CmpReplicasFailed   CmpReplicas = iota
)

type ReplicasComparison struct {
	Result      CmpReplicas
	Description string
	Err         error
}
