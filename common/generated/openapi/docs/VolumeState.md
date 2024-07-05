# VolumeState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Target** | Pointer to [**Nexus**](Nexus.md) | target exposed via a Nexus | [optional] 
**Size** | **int64** | size of the volume in bytes | 
**Status** | [**VolumeStatus**](VolumeStatus.md) |  | 
**Uuid** | **string** | name of the volume | 
**ReplicaTopology** | [**map[string]ReplicaTopology**](ReplicaTopology.md) | replica topology information | 
**Usage** | Pointer to [**VolumeUsage**](VolumeUsage.md) |  | [optional] 

## Methods

### NewVolumeState

`func NewVolumeState(size int64, status VolumeStatus, uuid string, replicaTopology map[string]ReplicaTopology, ) *VolumeState`

NewVolumeState instantiates a new VolumeState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeStateWithDefaults

`func NewVolumeStateWithDefaults() *VolumeState`

NewVolumeStateWithDefaults instantiates a new VolumeState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTarget

`func (o *VolumeState) GetTarget() Nexus`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *VolumeState) GetTargetOk() (*Nexus, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *VolumeState) SetTarget(v Nexus)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *VolumeState) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### GetSize

`func (o *VolumeState) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *VolumeState) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *VolumeState) SetSize(v int64)`

SetSize sets Size field to given value.


### GetStatus

`func (o *VolumeState) GetStatus() VolumeStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VolumeState) GetStatusOk() (*VolumeStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VolumeState) SetStatus(v VolumeStatus)`

SetStatus sets Status field to given value.


### GetUuid

`func (o *VolumeState) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *VolumeState) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *VolumeState) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetReplicaTopology

`func (o *VolumeState) GetReplicaTopology() map[string]ReplicaTopology`

GetReplicaTopology returns the ReplicaTopology field if non-nil, zero value otherwise.

### GetReplicaTopologyOk

`func (o *VolumeState) GetReplicaTopologyOk() (*map[string]ReplicaTopology, bool)`

GetReplicaTopologyOk returns a tuple with the ReplicaTopology field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicaTopology

`func (o *VolumeState) SetReplicaTopology(v map[string]ReplicaTopology)`

SetReplicaTopology sets ReplicaTopology field to given value.


### GetUsage

`func (o *VolumeState) GetUsage() VolumeUsage`

GetUsage returns the Usage field if non-nil, zero value otherwise.

### GetUsageOk

`func (o *VolumeState) GetUsageOk() (*VolumeUsage, bool)`

GetUsageOk returns a tuple with the Usage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsage

`func (o *VolumeState) SetUsage(v VolumeUsage)`

SetUsage sets Usage field to given value.

### HasUsage

`func (o *VolumeState) HasUsage() bool`

HasUsage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


