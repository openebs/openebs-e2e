# OnlineReplicaSnapshotState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**SourceId** | **string** |  | 
**PoolId** | **string** | storage pool identifier | 
**PoolUuid** | **string** | storage pool unique identifier | 
**Timestamp** | **time.Time** | Timestamp when the replica snapshot is taken on the storage system. | 
**Size** | **int64** | Replica snapshot size. | 
**AllocatedSize** | **int64** | Runtime size in bytes of the snapshot. Equal to the volume allocation at the time of the snapshot creation. It may grow larger if any of its predecessors are deleted. | 
**PredecessorAllocSize** | **int64** | Total allocated size of all the snapshot predecessors. | 

## Methods

### NewOnlineReplicaSnapshotState

`func NewOnlineReplicaSnapshotState(uuid string, sourceId string, poolId string, poolUuid string, timestamp time.Time, size int64, allocatedSize int64, predecessorAllocSize int64, ) *OnlineReplicaSnapshotState`

NewOnlineReplicaSnapshotState instantiates a new OnlineReplicaSnapshotState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOnlineReplicaSnapshotStateWithDefaults

`func NewOnlineReplicaSnapshotStateWithDefaults() *OnlineReplicaSnapshotState`

NewOnlineReplicaSnapshotStateWithDefaults instantiates a new OnlineReplicaSnapshotState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *OnlineReplicaSnapshotState) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *OnlineReplicaSnapshotState) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *OnlineReplicaSnapshotState) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetSourceId

`func (o *OnlineReplicaSnapshotState) GetSourceId() string`

GetSourceId returns the SourceId field if non-nil, zero value otherwise.

### GetSourceIdOk

`func (o *OnlineReplicaSnapshotState) GetSourceIdOk() (*string, bool)`

GetSourceIdOk returns a tuple with the SourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceId

`func (o *OnlineReplicaSnapshotState) SetSourceId(v string)`

SetSourceId sets SourceId field to given value.


### GetPoolId

`func (o *OnlineReplicaSnapshotState) GetPoolId() string`

GetPoolId returns the PoolId field if non-nil, zero value otherwise.

### GetPoolIdOk

`func (o *OnlineReplicaSnapshotState) GetPoolIdOk() (*string, bool)`

GetPoolIdOk returns a tuple with the PoolId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolId

`func (o *OnlineReplicaSnapshotState) SetPoolId(v string)`

SetPoolId sets PoolId field to given value.


### GetPoolUuid

`func (o *OnlineReplicaSnapshotState) GetPoolUuid() string`

GetPoolUuid returns the PoolUuid field if non-nil, zero value otherwise.

### GetPoolUuidOk

`func (o *OnlineReplicaSnapshotState) GetPoolUuidOk() (*string, bool)`

GetPoolUuidOk returns a tuple with the PoolUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolUuid

`func (o *OnlineReplicaSnapshotState) SetPoolUuid(v string)`

SetPoolUuid sets PoolUuid field to given value.


### GetTimestamp

`func (o *OnlineReplicaSnapshotState) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *OnlineReplicaSnapshotState) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *OnlineReplicaSnapshotState) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.


### GetSize

`func (o *OnlineReplicaSnapshotState) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *OnlineReplicaSnapshotState) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *OnlineReplicaSnapshotState) SetSize(v int64)`

SetSize sets Size field to given value.


### GetAllocatedSize

`func (o *OnlineReplicaSnapshotState) GetAllocatedSize() int64`

GetAllocatedSize returns the AllocatedSize field if non-nil, zero value otherwise.

### GetAllocatedSizeOk

`func (o *OnlineReplicaSnapshotState) GetAllocatedSizeOk() (*int64, bool)`

GetAllocatedSizeOk returns a tuple with the AllocatedSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedSize

`func (o *OnlineReplicaSnapshotState) SetAllocatedSize(v int64)`

SetAllocatedSize sets AllocatedSize field to given value.


### GetPredecessorAllocSize

`func (o *OnlineReplicaSnapshotState) GetPredecessorAllocSize() int64`

GetPredecessorAllocSize returns the PredecessorAllocSize field if non-nil, zero value otherwise.

### GetPredecessorAllocSizeOk

`func (o *OnlineReplicaSnapshotState) GetPredecessorAllocSizeOk() (*int64, bool)`

GetPredecessorAllocSizeOk returns a tuple with the PredecessorAllocSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPredecessorAllocSize

`func (o *OnlineReplicaSnapshotState) SetPredecessorAllocSize(v int64)`

SetPredecessorAllocSize sets PredecessorAllocSize field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


