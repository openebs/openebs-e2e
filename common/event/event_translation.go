package event

// This file contains code to convert event structures to forms that can
// be serialised with human-legible symbols.

import (
	"fmt"
)

// The following structures have enums values converted to corresponding enum strings

type RebuildDetailsSym struct {
	DestinationReplica string `json:"destination_replica"`
	Error              string `json:"error"`
	RebuildStatus      string `json:"rebuild_status"`
	SourceReplica      string `json:"source_replica"`
}

type SwitchOverEventDetailsSym struct {
	ExistingNqn      string `json:"existing_nqn"`
	NewPath          string `json:"new_path"`
	RetryCount       uint64 `json:"retry_count"`
	StartTime        string `json:"start_time"`
	SwitchOverStatus string `json:"switch_over_status"`
}

type EventDetailsSym struct {
	CloneEventDetails         *CloneEventDetails         `json:"clone_details,omitempty" yaml:"clone_details,omitempty"`
	ErrorDetails              *ErrorDetails              `json:"error_details,omitempty" yaml:"error_details,omitempty"`
	HostInitiatorEventDetails *HostInitiatorEventDetails `json:"host_initiator_details,omitempty" yaml:"host_initiator_details,omitempty"`
	NexusChildEventDetails    *NexusChildEventDetails    `json:"nexus_child_details,omitempty" yaml:"nexus_child_details,omitempty"`
	NvmePathEventDetails      *NvmePathEventDetails      `json:"nvme_path_details,omitempty" yaml:"nvme_path_details,omitempty"`
	ReactorEventDetails       *ReactorEventDetails       `json:"reactor_details,omitempty" yaml:"reactor_details,omitempty"`
	RebuildDetails            *RebuildDetailsSym         `json:"rebuild_details,omitempty" yaml:"rebuild_details,omitempty"`
	ReplicaEventDetails       *ReplicaEventDetails       `json:"replica_details,omitempty" yaml:"replica_details,omitempty"`
	SnapshotEventDetails      *SnapshotEventDetails      `json:"snapshot_details,omitempty" yaml:"snapshot_details,omitempty"`
	StateChangeEventDetails   *StateChangeEventDetails   `json:"state_change_details,omitempty" yaml:"state_change_details,omitempty"`
	SubsystemPauseDetails     *SubsystemPauseDetails     `json:"subsystem_pause_details,omitempty" yaml:"subsystem_pause_details,omitempty"`
	SwitchOverEventDetails    *SwitchOverEventDetailsSym `json:"switch_over_details,omitempty" yaml:"switch_over_details,omitempty"`
}

type EventSourceSym struct {
	// Io-engine or core-agent.
	Component string `json:"component"`
	// Node name
	Node         string           `json:"node"`
	EventDetails *EventDetailsSym `json:"event_details,omitempty" yaml:"event_details,omitempty"`
}

type EventMetaSym struct {
	// Something that uniquely identifies events.
	// UUIDv4.
	// GUID.
	Id     string         `json:"id"`
	Source EventSourceSym `json:"source"`
	// Event timestamp.
	EventTimestamp string `json:"timestamp"`
	// Version of the event message.
	Version int `json:"version"`
}

type EventMessageSym struct {
	// Event Category.
	Category string `json:"category"`
	// Event Action.
	Action string `json:"action"`
	// Target id for the category against which action is performed.
	Target string `json:"target"`
	// Event meta data.
	Metadata EventMetaSym `json:"metadata"`
}

func (category Category) String() string {
	switch category {
	case UnknownCategory:
		return "UnknownCategory"
	case CategoryPool:
		return "Pool"
	case CategoryVolume:
		return "Volume"
	case CategoryNexus:
		return "Nexus"
	case CategoryReplica:
		return "Replica"
	case CategoryNode:
		return "Node"
	case CategoryHighAvailability:
		return "HighAvailability"
	case CategoryNvmePath:
		return "NvmePath"
	case CategoryHostInitiator:
		return "HostInitiator"
	case CategoryIoEngine:
		return "IoEngine"
	case CategorySnapshot:
		return "Snapshot"
	case CategoryClone:
		return "Clone"
	default:
		return fmt.Sprintf("UnrecognisedCategory_%d", category)
	}
}

func (action Action) String() string {
	switch action {
	case UnknownAction:
		return "UnknownAction"
	case ActionCreate:
		return "Create"
	case ActionDelete:
		return "Delete"
	case ActionStateChange:
		return "StateChange"
	case ActionRebuildBegin:
		return "RebuildBegin"
	case ActionRebuildEnd:
		return "RebuildEnd"
	case ActionSwitchOver:
		return "SwitchOver"
	case ActionAddChild:
		return "AddChild"
	case ActionRemoveChild:
		return "RemoveChild"
	case ActionNvmePathSuspect:
		return "NvmeSuspect"
	case ActionNvmePathFail:
		return "NvmePathFail"
	case ActionNvmePathFix:
		return "NvmePathFix"
	case ActionOnlineChild:
		return "OnlineChild"
	case ActionNvmeConnect:
		return "NvmeConnect"
	case ActionNvmeDisconnect:
		return "NvmeDisconnect"
	case ActionNvmeKeepAliveTimeout:
		return "NvmeKeepAliveTimeout"
	case ActionReactorFreeze:
		return "ReactorFreeze"
	case ActionReactorUnfreeze:
		return "ReactorUnfreeze"
	case ActionShutdown:
		return "Shutdown"
	case ActionStart:
		return "Start"
	case ActionStop:
		return "Stop"
	case ActionSubsystemPause:
		return "SubsystemPause"
	case ActionSubsystemResume:
		return "SubsystemResume"
	case ActionInit:
		return "Init"
	case ActionReconfiguring:
		return "Reconfiguring"
	default:
		return fmt.Sprintf("UnrecognisedAction_%d", action)
	}
}

func (component Component) String() string {
	switch component {
	case UnknownComponent:
		return "UnknownComponent"
	case ComponentCoreAgent:
		return "CoreAgent"
	case ComponentIoEngine:
		return "IoEngine"
	case ComponentHaClusterAgent:
		return "HaClusterAgent"
	case ComponentHaNodeAgent:
		return "HaNodeAgent"
	default:
		return fmt.Sprintf("UnrecognisedComponent_%d", component)
	}
}

func (rebuildstatus RebuildStatus) String() string {
	switch rebuildstatus {
	case RebuildStatusUnknown:
		return "Unknown"
	case RebuildStatusStarted:
		return "Started"
	case RebuildStatusCompleted:
		return "Completed"
	case RebuildStatusStopped:
		return "Stopped"
	case RebuildStatusFailed:
		return "Failed"
	default:
		return fmt.Sprintf("UnrecognisedStatus_%d", rebuildstatus)
	}
}

func (switchoverstatus SwitchOverStatus) String() string {
	switch switchoverstatus {
	case UnknownSwitchOverStatus:
		return "Unknown"
	case SwitchOverStarted:
		return "Started"
	case SwitchOverCompleted:
		return "Completed"
	case SwitchOverFailed:
		return "Failed"
	default:
		return fmt.Sprintf("UnrecognisedStatus_%d", switchoverstatus)
	}
}

// Convert an EventMessage into an EventMessageSym which has
// enum numeric values converted to the corresponding Rust symbol names
func ToSymbolic(msg *EventMessage) (EventMessageSym, error) {
	eventMessageSym := EventMessageSym{}
	var err error

	eventMessageSym.Category = msg.Category.String()
	eventMessageSym.Action = msg.Action.String()
	eventMessageSym.Target = msg.Target
	eventMessageSym.Metadata.Source.Component = msg.Metadata.Source.Component.String()
	eventMessageSym.Metadata.Source.Node = msg.Metadata.Source.Node
	eventMessageSym.Metadata.Source.Node = msg.Metadata.Source.Node

	if msg.Metadata.Source.EventDetails != nil {
		eventMessageSym.Metadata.Source.EventDetails = &EventDetailsSym{}

		eventMessageSym.Metadata.Source.EventDetails.CloneEventDetails = msg.Metadata.Source.EventDetails.CloneEventDetails
		eventMessageSym.Metadata.Source.EventDetails.ErrorDetails = msg.Metadata.Source.EventDetails.ErrorDetails
		eventMessageSym.Metadata.Source.EventDetails.HostInitiatorEventDetails = msg.Metadata.Source.EventDetails.HostInitiatorEventDetails
		eventMessageSym.Metadata.Source.EventDetails.NexusChildEventDetails = msg.Metadata.Source.EventDetails.NexusChildEventDetails
		eventMessageSym.Metadata.Source.EventDetails.NvmePathEventDetails = msg.Metadata.Source.EventDetails.NvmePathEventDetails
		eventMessageSym.Metadata.Source.EventDetails.ReactorEventDetails = msg.Metadata.Source.EventDetails.ReactorEventDetails

		if msg.Metadata.Source.EventDetails.RebuildDetails != nil {
			rebuild_details := RebuildDetailsSym{
				SourceReplica:      msg.Metadata.Source.EventDetails.RebuildDetails.SourceReplica,
				DestinationReplica: msg.Metadata.Source.EventDetails.RebuildDetails.DestinationReplica,
				Error:              msg.Metadata.Source.EventDetails.RebuildDetails.Error,
				RebuildStatus:      msg.Metadata.Source.EventDetails.RebuildDetails.RebuildStatus.String(),
			}
			eventMessageSym.Metadata.Source.EventDetails.RebuildDetails = &rebuild_details
		}

		eventMessageSym.Metadata.Source.EventDetails.ReplicaEventDetails = msg.Metadata.Source.EventDetails.ReplicaEventDetails
		eventMessageSym.Metadata.Source.EventDetails.SnapshotEventDetails = msg.Metadata.Source.EventDetails.SnapshotEventDetails
		eventMessageSym.Metadata.Source.EventDetails.StateChangeEventDetails = msg.Metadata.Source.EventDetails.StateChangeEventDetails
		eventMessageSym.Metadata.Source.EventDetails.SubsystemPauseDetails = msg.Metadata.Source.EventDetails.SubsystemPauseDetails

		if msg.Metadata.Source.EventDetails.SwitchOverEventDetails != nil {
			switchover_event_details := SwitchOverEventDetailsSym{
				SwitchOverStatus: msg.Metadata.Source.EventDetails.SwitchOverEventDetails.SwitchOverStatus.String(),
				StartTime:        msg.Metadata.Source.EventDetails.SwitchOverEventDetails.StartTime,
				ExistingNqn:      msg.Metadata.Source.EventDetails.SwitchOverEventDetails.ExistingNqn,
				NewPath:          msg.Metadata.Source.EventDetails.SwitchOverEventDetails.NewPath,
				RetryCount:       msg.Metadata.Source.EventDetails.SwitchOverEventDetails.RetryCount,
			}
			eventMessageSym.Metadata.Source.EventDetails.SwitchOverEventDetails = &switchover_event_details
		}
	}

	eventMessageSym.Metadata.EventTimestamp = msg.Metadata.EventTimestamp
	eventMessageSym.Metadata.Id = msg.Metadata.Id
	eventMessageSym.Metadata.Version = msg.Metadata.Version

	return eventMessageSym, err
}
