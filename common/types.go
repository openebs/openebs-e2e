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
	ZfsFsType   FileSystemType = "zfs"
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

// PoolCordonConstraint represents the allowed constraint flags for pool cordon/uncordon operations.
// Using a typed enum improves type safety over free-form strings.
type PoolCordonConstraint int
type OfflinePoolDelete int

const (
	// CordonReplicas prevents scheduling new replicas on the pool
	CordonReplicas PoolCordonConstraint = iota
	// CordonSnapshots prevents creating snapshots on the pool
	CordonSnapshots PoolCordonConstraint = iota
	// CordonRestores prevents restore operations on the pool
	CordonRestores PoolCordonConstraint = iota
	// CordonImport prevents importing the pool after node restart
	CordonImport PoolCordonConstraint = iota
	// purgePool deletes the pool and all its data without further confirmation
	PurgePool OfflinePoolDelete = iota
	// ConfirmPoolDelete requires confirmation to delete the pool
	ConfirmPoolDelete OfflinePoolDelete = iota
	// AcceptDataLoss requires confirmation to delete the pool if data loss may occur
	AcceptDataLoss OfflinePoolDelete = iota
	// AcceptVolumeLoss requires confirmation to delete the pool if volume loss may occur
	AcceptVolumeLoss OfflinePoolDelete = iota
	// AcceptSnapshotLoss requires confirmation to delete the pool if snapshot loss may occur
	AcceptSnapshotLoss OfflinePoolDelete = iota
	// CleanupCr deletes the DiskPool CR (DiskPool/DSP) after pool deletion
	CleanupCr OfflinePoolDelete = iota
	// IgnoreNotFound makes pool delete succeed even if the pool is already deleted
	IgnoreNotFound OfflinePoolDelete = iota
)

// String returns the CLI flag name corresponding to the constraint.
func (c PoolCordonConstraint) String() string {
	switch c {
	case CordonReplicas:
		return "replicas"
	case CordonSnapshots:
		return "snapshots"
	case CordonRestores:
		return "restores"
	case CordonImport:
		return "import"
	default:
		return ""
	}
}

// String returns the CLI flag name corresponding to the constraint.
func (c OfflinePoolDelete) String() string {
	switch c {
	case PurgePool:
		return "purge"
	case ConfirmPoolDelete:
		return "yes"
	case AcceptDataLoss:
		return "accept-data-loss"
	case AcceptVolumeLoss:
		return "accept-volume-loss"
	case AcceptSnapshotLoss:
		return "accept-snapshot-loss"
	case CleanupCr:
		return "cleanup-dsp"
	case IgnoreNotFound:
		return "ignore-not-found"
	default:
		return ""
	}
}

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

type DiskPoolCrState int

const (
	DiskPoolCrCreated     DiskPoolCrState = iota
	DiskPoolCrCreating    DiskPoolCrState = iota
	DiskPoolCrTerminating DiskPoolCrState = iota
	DiskPoolCrDeleted     DiskPoolCrState = iota
)

func (poolCrState DiskPoolCrState) String() string {
	switch poolCrState {
	case DiskPoolCrCreated:
		return "Created"
	case DiskPoolCrCreating:
		return "Creating"
	case DiskPoolCrTerminating:
		return "Terminating"
	case DiskPoolCrDeleted:
		return "Deleted"
	default:
		return ""
	}
}

type DiskPoolStatus int

const (
	DiskPoolOnline   DiskPoolStatus = iota
	DiskPoolUnknown  DiskPoolStatus = iota
	DiskPoolOffline  DiskPoolStatus = iota
	DiskPoolDegraded DiskPoolStatus = iota
	DiskPoolFaulted  DiskPoolStatus = iota
)

func (poolStatus DiskPoolStatus) String() string {
	switch poolStatus {
	case DiskPoolOnline:
		return "Online"
	case DiskPoolUnknown:
		return "Unknown"
	case DiskPoolOffline:
		return "Offline"
	case DiskPoolDegraded:
		return "Degraded"
	case DiskPoolFaulted:
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

type OpenEbsEngine int

const (
	Lvm      OpenEbsEngine = iota
	Hostpath OpenEbsEngine = iota
	Zfs      OpenEbsEngine = iota
	Mayastor OpenEbsEngine = iota
	None     OpenEbsEngine = iota
)

func (Engine OpenEbsEngine) String() string {
	switch Engine {
	case Lvm:
		return "lvm"
	case Hostpath:
		return "hostpath"
	case Zfs:
		return "zfs"
	case Mayastor:
		return "mayastor"
	case None:
		return "none"
	default:
		return ""
	}
}

type YesNoVal int

const (
	Yes YesNoVal = iota
	No  YesNoVal = iota
)

func (Val YesNoVal) String() string {
	switch Val {
	case Yes:
		return "yes"
	case No:
		return "no"
	default:
		return ""
	}
}

type OnOffVal int

const (
	On  OnOffVal = iota
	Off OnOffVal = iota
)

func (val OnOffVal) String() string {
	switch val {
	case On:
		return "on"
	case Off:
		return "off"
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
	Encrypted     bool          `json:"encrypted"`
	AffinityGroup struct {
		Id string `json:"id"`
	} `json:"affinity_group"`
}

type Policy struct {
	Self_heal bool `json:"self_heal"`
}
type SpecTarget struct {
	Protocol      string          `json:"protocol"`
	Node          string          `json:"node"`
	FrontendNodes []FrontendNodes `json:"frontend_nodes"`
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
	Health          Health          `json:"health"`
	Replicas        []Replica       `json:"replicas"`
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
	Health      Health       `json:"health"`
	Encrypted   bool         `json:"encrypted"`
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

type FrontendNodes struct {
	Name string `json:"name"`
	Nqn  string `json:"nqn"`
}

type Health struct {
	CleanShutdown         bool  `json:"cleanShutdown"`
	HealthyReplicas       int32 `json:"healthyReplicas"`
	CleanReplicas         int32 `json:"cleanReplicas"`
	OnlineHealthyReplicas int32 `json:"onlineHealthyReplicas"`
	OnlineCleanReplicas   int32 `json:"onlineCleanReplicas"`
	LiveHealthyReplicas   int32 `json:"liveHealthyReplicas"`
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

type ZFSVolume struct {
	Spec   ZfsSpec   `json:"spec"`
	Status ZfsStatus `json:"status"`
}

type ZfsSpec struct {
	Capacity      string `json:"capacity"`
	Compression   string `json:"compression"`
	FsType        string `json:"fsType"`
	OwnerNodeID   string `json:"ownerNodeID"`
	PoolName      string `json:"poolName"`
	QuotaType     string `json:"quotaType"`
	Recordsize    string `json:"recordsize"`
	Shared        string `json:"shared"`
	ThinProvision string `json:"thinProvision"`
	VolumeType    string `json:"volumeType"`
}

type ZfsStatus struct {
	State string `json:"state"`
}

// ── Event enum types (PascalCase — values in JSON output) ──

type EventCategory string

const (
	EventCategoryPool             EventCategory = "Pool"
	EventCategoryVolume           EventCategory = "Volume"
	EventCategoryNexus            EventCategory = "Nexus"
	EventCategoryReplica          EventCategory = "Replica"
	EventCategoryNode             EventCategory = "Node"
	EventCategoryHighAvailability EventCategory = "HighAvailability"
	EventCategoryNvmePath         EventCategory = "NvmePath"
	EventCategoryHostInitiator    EventCategory = "HostInitiator"
	EventCategoryIoEngine         EventCategory = "IoEngineCategory"
	EventCategorySnapshot         EventCategory = "Snapshot"
	EventCategoryClone            EventCategory = "Clone"
)

type EventAction string

const (
	EventActionCreate               EventAction = "Create"
	EventActionDelete               EventAction = "Delete"
	EventActionStateChange          EventAction = "StateChange"
	EventActionRebuildBegin         EventAction = "RebuildBegin"
	EventActionRebuildEnd           EventAction = "RebuildEnd"
	EventActionSwitchOver           EventAction = "SwitchOver"
	EventActionAddChild             EventAction = "AddChild"
	EventActionRemoveChild          EventAction = "RemoveChild"
	EventActionOnlineChild          EventAction = "OnlineChild"
	EventActionNvmePathSuspect      EventAction = "NvmePathSuspect"
	EventActionNvmePathFail         EventAction = "NvmePathFail"
	EventActionNvmePathFix          EventAction = "NvmePathFix"
	EventActionNvmeConnect          EventAction = "NvmeConnect"
	EventActionNvmeDisconnect       EventAction = "NvmeDisconnect"
	EventActionNvmeKeepAliveTimeout EventAction = "NvmeKeepAliveTimeout"
	EventActionReactorFreeze        EventAction = "ReactorFreeze"
	EventActionReactorUnfreeze      EventAction = "ReactorUnfreeze"
	EventActionShutdown             EventAction = "Shutdown"
	EventActionStart                EventAction = "Start"
	EventActionStop                 EventAction = "Stop"
	EventActionSubsystemPause       EventAction = "SubsystemPause"
	EventActionSubsystemResume      EventAction = "SubsystemResume"
	EventActionInit                 EventAction = "Init"
	EventActionReconfiguring        EventAction = "Reconfiguring"
	EventActionNvmePathDeleting     EventAction = "NvmePathDeleting"
)

type EventComponent string

const (
	EventComponentCoreAgent      EventComponent = "CoreAgent"
	EventComponentIoEngine       EventComponent = "IoEngine"
	EventComponentHaClusterAgent EventComponent = "HaClusterAgent"
	EventComponentHaNodeAgent    EventComponent = "HaNodeAgent"
)

type RebuildStatus string

const (
	RebuildStatusStarted   RebuildStatus = "Started"
	RebuildStatusCompleted RebuildStatus = "Completed"
	RebuildStatusStopped   RebuildStatus = "Stopped"
	RebuildStatusFailed    RebuildStatus = "Failed"
)

type SwitchOverStatus string

const (
	SwitchOverStarted   SwitchOverStatus = "SwitchOverStarted"
	SwitchOverCompleted SwitchOverStatus = "SwitchOverCompleted"
	SwitchOverFailed    SwitchOverStatus = "SwitchOverFailed"
)

type EventVersion string

const (
	EventVersionV1 EventVersion = "V1"
)

// ── Event struct types (mirrors kubectl mayastor get events -o json) ──

type EventRecord struct {
	Category EventCategory `json:"category"`
	Action   EventAction   `json:"action"`
	Target   string        `json:"target,omitempty"`
	Metadata EventMeta     `json:"metadata"`
}

type EventMeta struct {
	Id             string       `json:"id"`
	Source         EventSource  `json:"source"`
	EventTimestamp string       `json:"timestamp"`
	Version        EventVersion `json:"version"`
}

type EventSource struct {
	Component    EventComponent `json:"component"`
	Node         string         `json:"node,omitempty"`
	EventDetails *EventDetails  `json:"eventDetails,omitempty"`
}

type EventDetails struct {
	RebuildDetails        *RebuildDetails        `json:"rebuildDetails,omitempty"`
	SwitchOverDetails     *SwitchOverDetails     `json:"switchOverDetails,omitempty"`
	NexusChildDetails     *NexusChildDetails     `json:"nexusChildDetails,omitempty"`
	NvmePathDetails       *NvmePathDetails       `json:"nvmePathDetails,omitempty"`
	HostInitiatorDetails  *HostInitiatorDetails  `json:"hostInitiatorDetails,omitempty"`
	StateChangeDetails    *StateChangeDetails    `json:"stateChangeDetails,omitempty"`
	ReplicaDetails        *ReplicaDetails        `json:"replicaDetails,omitempty"`
	SnapshotDetails       *SnapshotDetails       `json:"snapshotDetails,omitempty"`
	CloneDetails          *CloneDetails          `json:"cloneDetails,omitempty"`
	SubsystemPauseDetails *SubsystemPauseDetails `json:"subsystemPauseDetails,omitempty"`
	ActionDurationDetails *ActionDurationDetails  `json:"actionDurationDetails,omitempty"`
	ReactorDetails        *ReactorDetails        `json:"reactorDetails,omitempty"`
	ErrorDetails          *ErrorDetails          `json:"errorDetails,omitempty"`
}

type RebuildDetails struct {
	RebuildStatus      RebuildStatus `json:"rebuildStatus"`
	SourceReplica      string        `json:"sourceReplica"`
	DestinationReplica string        `json:"destinationReplica"`
	Error              string        `json:"error,omitempty"`
}

type SwitchOverDetails struct {
	SwitchOverStatus SwitchOverStatus `json:"switchOverStatus"`
	StartTime        string           `json:"startTime"`
	ExistingNqn      string           `json:"existingNqn"`
	NewPath          string           `json:"newPath,omitempty"`
	RetryCount       uint64           `json:"retryCount,omitempty"`
}

type NexusChildDetails struct {
	Uri string `json:"uri"`
}

type NvmePathDetails struct {
	Nqn  string `json:"nqn"`
	Path string `json:"path"`
}

type HostInitiatorDetails struct {
	HostNqn      string `json:"hostNqn"`
	SubsystemNqn string `json:"subsystemNqn"`
	Target       string `json:"target"`
	Uuid         string `json:"uuid"`
}

type StateChangeDetails struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type ReplicaDetails struct {
	PoolName    string `json:"poolName"`
	PoolUuid    string `json:"poolUuid"`
	ReplicaName string `json:"replicaName"`
}

type SnapshotDetails struct {
	ReplicaId  string `json:"replicaId"`
	CreateTime string `json:"createTime"`
	VolumeId   string `json:"volumeId"`
}

type CloneDetails struct {
	SourceUuid string `json:"sourceUuid"`
	CreateTime string `json:"createTime"`
}

type SubsystemPauseDetails struct {
	NexusPauseState string `json:"nexusPauseState"`
}

type ActionDurationDetails struct {
	TimeTaken TimeDuration `json:"timeTaken"`
}

type TimeDuration struct {
	Seconds uint64 `json:"seconds"`
	Nanos   uint64 `json:"nanos"`
}

type ReactorDetails struct {
	Lcore uint64 `json:"lcore"`
	State string `json:"state"`
}

type ErrorDetails struct {
	Error string `json:"error"`
}
