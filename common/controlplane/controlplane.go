package controlplane

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/openebs/openebs-e2e/common"
	"github.com/openebs/openebs-e2e/common/e2e_config"

	v1 "github.com/openebs/openebs-e2e/common/controlplane/v1"
	v1rest "github.com/openebs/openebs-e2e/common/controlplane/v1-rest-api"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type ControlPlaneInterface interface {
	// Version

	MajorVersion() int
	Version() string

	IsTimeoutError(error) bool

	// Resource state strings abstraction

	VolStateHealthy() string
	VolStateDegraded() string
	VolStateUnknown() string
	VolStateFaulted() string
	ChildStateUnknown() string
	ChildStateOnline() string
	ChildStateDegraded() string
	ChildStateFaulted() string
	NexusStateUnknown() string
	NexusStateOnline() string
	NexusStateDegraded() string
	NexusStateFaulted() string
	MspStateOnline() string
	ReplicaStateUnknown() string
	ReplicaStateOnline() string
	ReplicaStateDegraded() string
	ReplicaStateFaulted() string

	// MSV abstraction

	GetMSV(uuid string) (*common.MayastorVolume, error)
	GetMsvNodes(uuid string) (string, []string)
	CanDeleteMsv() bool
	DeleteMsv(volName string) error
	ListMsvs() ([]common.MayastorVolume, error)
	SetMsvReplicaCount(uuid string, replicaCount int) error
	GetMsvState(uuid string) (string, error)
	GetMsvReplicas(volName string) ([]common.MsvReplica, error)
	GetMsvReplicaTopology(volUuid string) (common.ReplicaTopology, error)
	GetMsvNexusChildren(volName string) ([]common.TargetChild, error)
	GetMsvNexusState(uuid string) (string, error)
	IsMsvPublished(uuid string) bool
	IsMsvDeleted(uuid string) bool
	CheckForMsvs() (bool, error)
	CheckAllMsvsAreHealthy() error
	GetMsvTargetNode(uuid string) (string, error)
	GetMsvTargetUuid(uuid string) (string, error)
	ListRestoredMsvs() ([]common.MayastorVolume, error)
	GetMsvSize(uuid string) (int64, error)
	SetVolumeMaxSnapshotCount(uuid string, maxSnapshotCount int32) error
	GetMsvMaxSnapshotCount(uuid string) (int32, error)
	GetMsvDeviceUri(uuid string) (string, error)

	// Mayastor Node abstraction

	GetMSN(node string) (*common.MayastorNode, error)
	ListMsns() ([]common.MayastorNode, error)
	GetMsNodeStatus(node string) (string, error)
	UpdateNodeLabel(nodeName string, labelKey, labelValue string) error

	NodeStateOffline() string
	NodeStateOnline() string
	NodeStateUnknown() string
	NodeStateEmpty() string

	// Mayastor Pool abstraction
	CreatePoolOnInstall() bool
	GetMsPool(poolName string) (*common.MayastorPool, error)
	ListMsPools() ([]common.MayastorPool, error)

	//cordon
	CordonNode(nodeName string, cordonLabel string) error
	GetCordonNodeLabels(nodeName string) ([]string, error)
	UnCordonNode(nodeName string, cordonLabel string) error

	//drain
	DrainNode(nodeName string, drainLabel string, drainTimeOut int) error
	GetDrainNodeLabels(nodeName string) ([]string, []string, error)

	//upgrade
	Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error)
	UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error
	UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error
	UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error
	UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error
	GetUpgradeStatus() (string, error)
	GetToUpgradeVersion() (string, error)
	DeleteUpgrade() error

	//snapshot
	GetVolumeSnapshots(volUuid string) ([]common.SnapshotSchema, error)
	GetVolumeSnapshot(volUuid string, snapshotId string) (common.SnapshotSchema, error)
	GetSnapshot(snapshotId string) (common.SnapshotSchema, error)
	GetSnapshots() ([]common.SnapshotSchema, error)
	GetVolumeSnapshotTopology() ([]common.SnapshotSchema, error)
	GetPerSnapshotVolumeSnapshotTopology(snapshotId string) (common.SnapshotSchema, error)
}

var ifc ControlPlaneInterface

var once sync.Once

func getControlPlane() ControlPlaneInterface {
	once.Do(func() {
		version := e2e_config.GetConfig().MayastorVersion
		verComponents := strings.Split(version, ".")
		major, err := strconv.Atoi(verComponents[0])
		if err == nil {
			switch major {
			case 1:
				if e2e_config.GetConfig().KubectlPluginDir == "" || os.Getenv("e2e_control_plane_rest_api") == `true` {
					logf.Log.Info("*** Using REST API for control plane communication ***")
					ifc, err = v1rest.MakeCP()
				} else {
					logf.Log.Info("*** Using kubectl plugin for control plane communication ***")
					ifc, err = v1.MakeCP()
				}
			default:
				panic(fmt.Errorf("unsupported control plane version %v", version))
			}
		}
		if err != nil {
			panic(err)
		}
		if ifc == nil {
			panic("failed to set control plane object")
		}
	})
	return ifc
}

func VolStateHealthy() string {
	return getControlPlane().VolStateHealthy()
}

func VolStateDegraded() string {
	return getControlPlane().VolStateDegraded()
}

func VolStateUnknown() string {
	return getControlPlane().VolStateUnknown()
}

func VolStateFaulted() string {
	return getControlPlane().VolStateFaulted()
}

func ChildStateUnknown() string {
	return getControlPlane().ChildStateUnknown()
}

func ChildStateOnline() string {
	return getControlPlane().ChildStateOnline()
}

func ChildStateDegraded() string {
	return getControlPlane().ChildStateDegraded()
}

func ChildStateFaulted() string {
	return getControlPlane().ChildStateFaulted()
}

func NexusStateUnknown() string {
	return getControlPlane().NexusStateUnknown()
}

func NexusStateOnline() string {
	return getControlPlane().NexusStateOnline()
}

func NexusStateDegraded() string {
	return getControlPlane().NexusStateDegraded()
}

func NexusStateFaulted() string {
	return getControlPlane().NexusStateFaulted()
}

func MspStateOnline() string {
	return getControlPlane().MspStateOnline()
}

func ReplicaStateUnknown() string {
	return getControlPlane().ReplicaStateUnknown()
}

func ReplicaStateOnline() string {
	return getControlPlane().ReplicaStateOnline()
}

func ReplicaStateDegraded() string {
	return getControlPlane().ReplicaStateDegraded()
}

func ReplicaStateFaulted() string {
	return getControlPlane().ReplicaStateFaulted()
}

//FIXME: MSV These functions are only guaranteed to
// work correctly if invoked from k8stest/msv.go
// which ensures that necessary setup functions
// have been called
// The issue is that for control plane v1 we need
// node addresses and the k8stest pkg provides that.

// GetMSV Get pointer to a mayastor volume custom resource
// returns nil and no error if the msv is in pending state.
func GetMSV(uuid string) (*common.MayastorVolume, error) {
	return getControlPlane().GetMSV(uuid)
}

// GetMsvNodes Retrieve the nexus node hosting the Mayastor Volume,
// and the names of the replica nodes
// function asserts if the volume CR is not found.
func GetMsvNodes(uuid string) (string, []string) {
	return getControlPlane().GetMsvNodes(uuid)
}

func CanDeleteMsv() bool {
	return getControlPlane().CanDeleteMsv()
}

func DeleteMsv(volName string) error {
	return getControlPlane().DeleteMsv(volName)
}

func ListMsvs() ([]common.MayastorVolume, error) {
	return getControlPlane().ListMsvs()
}

func ListRestoredMsvs() ([]common.MayastorVolume, error) {
	return getControlPlane().ListRestoredMsvs()
}

func SetMsvReplicaCount(uuid string, replicaCount int) error {
	return getControlPlane().SetMsvReplicaCount(uuid, replicaCount)
}

func GetMsvState(uuid string) (string, error) {
	return getControlPlane().GetMsvState(uuid)
}

func GetMsvReplicas(volName string) ([]common.MsvReplica, error) {
	return getControlPlane().GetMsvReplicas(volName)
}

func GetMsvReplicaTopology(volUuid string) (common.ReplicaTopology, error) {
	return getControlPlane().GetMsvReplicaTopology(volUuid)
}

func GetMsvNexusChildren(volName string) ([]common.TargetChild, error) {
	return getControlPlane().GetMsvNexusChildren(volName)
}

func GetMsvNexusState(uuid string) (string, error) {
	return getControlPlane().GetMsvNexusState(uuid)
}

func IsMsvPublished(uuid string) bool {
	return getControlPlane().IsMsvPublished(uuid)
}

func IsMsvDeleted(uuid string) bool {
	return getControlPlane().IsMsvDeleted(uuid)
}

func CheckForMsvs() (bool, error) {
	return getControlPlane().CheckForMsvs()
}

func CheckAllMsvsAreHealthy() error {
	return getControlPlane().CheckAllMsvsAreHealthy()
}

func GetMsvTargetNode(uuid string) (string, error) {
	return getControlPlane().GetMsvTargetNode(uuid)
}

func GetMsvTargetUuid(uuid string) (string, error) {
	return getControlPlane().GetMsvTargetUuid(uuid)
}

func GetMsvSize(uuid string) (int64, error) {
	return getControlPlane().GetMsvSize(uuid)
}

func GetMsvDeviceUri(uuid string) (string, error) {
	return getControlPlane().GetMsvDeviceUri(uuid)
}

func SetVolumeMaxSnapshotCount(uuid string, maxSnapshotCount int32) error {
	return getControlPlane().SetVolumeMaxSnapshotCount(uuid, maxSnapshotCount)
}

func GetMsvMaxSnapshotCount(uuid string) (int32, error) {
	return getControlPlane().GetMsvMaxSnapshotCount(uuid)
}

//FIXME: MSN These functions are only guaranteed to
// work correctly if invoked from k8stest/msn.go
// which ensures that necessary setup functions
// have been called
// The issue is that for control plane v1 we need
// node addresses and the k8stest pkg provides that.

// GetMSN Get pointer to a mayastor node custom resource
// returns nil and no error if the msn is in pending state.
func GetMSN(nodeName string) (*common.MayastorNode, error) {
	return getControlPlane().GetMSN(nodeName)
}

func ListMsns() ([]common.MayastorNode, error) {
	return getControlPlane().ListMsns()
}
func GetMsNodeStatus(nodeName string) (string, error) {
	return getControlPlane().GetMsNodeStatus(nodeName)
}

func UpdateNodeLabel(nodeName string, labelKey, labelValue string) error {
	return getControlPlane().UpdateNodeLabel(nodeName, labelKey, labelValue)
}

func NodeStateOnline() string {
	return getControlPlane().NodeStateOnline()
}

// NodeStateOffline is set when the node misses its watchdog deadline
func NodeStateOffline() string {
	return getControlPlane().NodeStateOffline()
}

// NodeStateUnknown is set if the mayastor instance deregisters itself (when the pod goes down gracefully),
// or if there's an error when we're issuing issuing "list" requests
func NodeStateUnknown() string {
	return getControlPlane().NodeStateUnknown()
}

// NodeStateEmpty i.e. no state at all if the control plane restarts and the node is not available at that time
func NodeStateEmpty() string {
	return getControlPlane().NodeStateEmpty()
}

func Version() string {
	return getControlPlane().Version()
}

func MajorVersion() int {
	return getControlPlane().MajorVersion()
}

func IsTimeoutError(err error) bool {
	return getControlPlane().IsTimeoutError(err)
}

func CreatePoolOnInstall() bool {
	return getControlPlane().CreatePoolOnInstall()
}

func GetMsPool(poolName string) (*common.MayastorPool, error) {
	return getControlPlane().GetMsPool(poolName)
}

func ListMsPools() ([]common.MayastorPool, error) {
	return getControlPlane().ListMsPools()
}

func CordonNode(nodeName string, cordonLabel string) error {
	return getControlPlane().CordonNode(nodeName, cordonLabel)
}

func GetCordonNodeLabels(nodeName string) ([]string, error) {
	return getControlPlane().GetCordonNodeLabels(nodeName)
}

func UnCordonNode(nodeName string, cordonLabel string) error {
	return getControlPlane().UnCordonNode(nodeName, cordonLabel)
}

func DrainNode(nodeName string, drainLabel string, drainTimeOut int) error {
	return getControlPlane().DrainNode(nodeName, drainLabel, drainTimeOut)
}

func GetDrainNodeLabels(nodeName string) ([]string, []string, error) {
	return getControlPlane().GetDrainNodeLabels(nodeName)
}

func Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) (string, error) {
	return getControlPlane().Upgrade(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded)
}

func UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	return getControlPlane().UpgradeWithSkipDataPlaneRestart(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded)
}

func UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	return getControlPlane().UpgradeWithSkipReplicaRebuild(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded)
}

func UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	return getControlPlane().UpgradeWithSkipSingleReplicaValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded)
}

func UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded bool) error {
	return getControlPlane().UpgradeWithSkipCordonNodeValidation(isUpgradingToUnstableBranch, isPartialRebuildDisableNeeded)
}

func GetUpgradeStatus() (string, error) {
	return getControlPlane().GetUpgradeStatus()
}

func GetToUpgradeVersion() (string, error) {
	return getControlPlane().GetToUpgradeVersion()
}

func DeleteUpgrade() error {
	return getControlPlane().DeleteUpgrade()
}

func GetVolumeSnapshots(volUuid string) ([]common.SnapshotSchema, error) {
	return getControlPlane().GetVolumeSnapshots(volUuid)
}

func GetVolumeSnapshot(volUuid string, snapshotId string) (common.SnapshotSchema, error) {
	return getControlPlane().GetVolumeSnapshot(volUuid, snapshotId)
}

func GetSnapshot(snapshotId string) (common.SnapshotSchema, error) {
	return getControlPlane().GetSnapshot(snapshotId)
}

func GetSnapshots() ([]common.SnapshotSchema, error) {
	return getControlPlane().GetSnapshots()
}

func GetVolumeSnapshotTopology() ([]common.SnapshotSchema, error) {
	return getControlPlane().GetVolumeSnapshotTopology()
}

func GetPerSnapshotVolumeSnapshotTopology(snapshotId string) (common.SnapshotSchema, error) {
	return getControlPlane().GetPerSnapshotVolumeSnapshotTopology(snapshotId)
}
