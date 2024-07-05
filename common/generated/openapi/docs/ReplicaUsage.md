# ReplicaUsage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Capacity** | **int64** | Replica capacity in bytes. | [default to 0]
**Allocated** | **int64** | Amount of actually allocated disk space for this replica in bytes. | [default to 0]
**AllocatedSnapshots** | **int64** | Amount of actually allocated disk space for this replica&#39;s snapshots in bytes. | [default to 0]
**AllocatedAllSnapshots** | **int64** | Amount of actually allocated disk space for this replica&#39;s snapshots and its predecessors in bytes. For a restored/cloned replica this includes snapshots from the parent source.  | [default to 0]

## Methods

### NewReplicaUsage

`func NewReplicaUsage(capacity int64, allocated int64, allocatedSnapshots int64, allocatedAllSnapshots int64, ) *ReplicaUsage`

NewReplicaUsage instantiates a new ReplicaUsage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaUsageWithDefaults

`func NewReplicaUsageWithDefaults() *ReplicaUsage`

NewReplicaUsageWithDefaults instantiates a new ReplicaUsage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacity

`func (o *ReplicaUsage) GetCapacity() int64`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *ReplicaUsage) GetCapacityOk() (*int64, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *ReplicaUsage) SetCapacity(v int64)`

SetCapacity sets Capacity field to given value.


### GetAllocated

`func (o *ReplicaUsage) GetAllocated() int64`

GetAllocated returns the Allocated field if non-nil, zero value otherwise.

### GetAllocatedOk

`func (o *ReplicaUsage) GetAllocatedOk() (*int64, bool)`

GetAllocatedOk returns a tuple with the Allocated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocated

`func (o *ReplicaUsage) SetAllocated(v int64)`

SetAllocated sets Allocated field to given value.


### GetAllocatedSnapshots

`func (o *ReplicaUsage) GetAllocatedSnapshots() int64`

GetAllocatedSnapshots returns the AllocatedSnapshots field if non-nil, zero value otherwise.

### GetAllocatedSnapshotsOk

`func (o *ReplicaUsage) GetAllocatedSnapshotsOk() (*int64, bool)`

GetAllocatedSnapshotsOk returns a tuple with the AllocatedSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedSnapshots

`func (o *ReplicaUsage) SetAllocatedSnapshots(v int64)`

SetAllocatedSnapshots sets AllocatedSnapshots field to given value.


### GetAllocatedAllSnapshots

`func (o *ReplicaUsage) GetAllocatedAllSnapshots() int64`

GetAllocatedAllSnapshots returns the AllocatedAllSnapshots field if non-nil, zero value otherwise.

### GetAllocatedAllSnapshotsOk

`func (o *ReplicaUsage) GetAllocatedAllSnapshotsOk() (*int64, bool)`

GetAllocatedAllSnapshotsOk returns a tuple with the AllocatedAllSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedAllSnapshots

`func (o *ReplicaUsage) SetAllocatedAllSnapshots(v int64)`

SetAllocatedAllSnapshots sets AllocatedAllSnapshots field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


