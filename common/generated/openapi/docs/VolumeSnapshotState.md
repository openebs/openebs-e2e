# VolumeSnapshotState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**AllocatedSize** | **int64** | Runtime size in bytes of the snapshot. Equal to the volume allocation at the time of the snapshot creation. It may grow larger if any of its predecessors are deleted. | 
**SourceVolume** | **string** |  | 
**Timestamp** | Pointer to **time.Time** | Timestamp when snapshot is taken on the storage system. | [optional] 
**ReadyAsSource** | **bool** | Indicates if a snapshot is ready to be used as a new volume source. | [default to false]
**ReplicaSnapshots** | [**[]ReplicaSnapshotState**](ReplicaSnapshotState.md) | List of individual ReplicaSnapshotStates. | 

## Methods

### NewVolumeSnapshotState

`func NewVolumeSnapshotState(uuid string, allocatedSize int64, sourceVolume string, readyAsSource bool, replicaSnapshots []ReplicaSnapshotState, ) *VolumeSnapshotState`

NewVolumeSnapshotState instantiates a new VolumeSnapshotState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeSnapshotStateWithDefaults

`func NewVolumeSnapshotStateWithDefaults() *VolumeSnapshotState`

NewVolumeSnapshotStateWithDefaults instantiates a new VolumeSnapshotState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *VolumeSnapshotState) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *VolumeSnapshotState) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *VolumeSnapshotState) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetAllocatedSize

`func (o *VolumeSnapshotState) GetAllocatedSize() int64`

GetAllocatedSize returns the AllocatedSize field if non-nil, zero value otherwise.

### GetAllocatedSizeOk

`func (o *VolumeSnapshotState) GetAllocatedSizeOk() (*int64, bool)`

GetAllocatedSizeOk returns a tuple with the AllocatedSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedSize

`func (o *VolumeSnapshotState) SetAllocatedSize(v int64)`

SetAllocatedSize sets AllocatedSize field to given value.


### GetSourceVolume

`func (o *VolumeSnapshotState) GetSourceVolume() string`

GetSourceVolume returns the SourceVolume field if non-nil, zero value otherwise.

### GetSourceVolumeOk

`func (o *VolumeSnapshotState) GetSourceVolumeOk() (*string, bool)`

GetSourceVolumeOk returns a tuple with the SourceVolume field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceVolume

`func (o *VolumeSnapshotState) SetSourceVolume(v string)`

SetSourceVolume sets SourceVolume field to given value.


### GetTimestamp

`func (o *VolumeSnapshotState) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *VolumeSnapshotState) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *VolumeSnapshotState) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *VolumeSnapshotState) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetReadyAsSource

`func (o *VolumeSnapshotState) GetReadyAsSource() bool`

GetReadyAsSource returns the ReadyAsSource field if non-nil, zero value otherwise.

### GetReadyAsSourceOk

`func (o *VolumeSnapshotState) GetReadyAsSourceOk() (*bool, bool)`

GetReadyAsSourceOk returns a tuple with the ReadyAsSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadyAsSource

`func (o *VolumeSnapshotState) SetReadyAsSource(v bool)`

SetReadyAsSource sets ReadyAsSource field to given value.


### GetReplicaSnapshots

`func (o *VolumeSnapshotState) GetReplicaSnapshots() []ReplicaSnapshotState`

GetReplicaSnapshots returns the ReplicaSnapshots field if non-nil, zero value otherwise.

### GetReplicaSnapshotsOk

`func (o *VolumeSnapshotState) GetReplicaSnapshotsOk() (*[]ReplicaSnapshotState, bool)`

GetReplicaSnapshotsOk returns a tuple with the ReplicaSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicaSnapshots

`func (o *VolumeSnapshotState) SetReplicaSnapshots(v []ReplicaSnapshotState)`

SetReplicaSnapshots sets ReplicaSnapshots field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


