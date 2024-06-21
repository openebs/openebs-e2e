# VolumeUsage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Capacity** | **int64** | Capacity of the volume in bytes. | 
**Allocated** | **int64** | -| Allocated size in bytes, related the largest healthy replica, including snapshots. For example, if a volume has 2 replicas, each with 1MiB allocated space, then this field will be 1MiB. | 
**AllocatedReplica** | **int64** | -| Allocated size in bytes, related to the largest healthy replica, excluding snapshots. | 
**AllocatedSnapshots** | **int64** | -| Allocated size in bytes, related the healthy replica with the highest snapshot usage. | 
**AllocatedAllSnapshots** | **int64** | -| For a restored/cloned volume, allocated size in bytes, related to the healthy replica with largest parent snapshot allocation. | 
**TotalAllocated** | **int64** | -| Allocated size in bytes, accrued from all the replicas, including snapshots. For example, if a volume has 2 replicas, each with 1MiB allocated space, then this field will be 2MiB. | 
**TotalAllocatedReplicas** | **interface{}** | -| Allocated size in bytes, accrued from all the replicas, excluding snapshots. | 
**TotalAllocatedSnapshots** | **int64** | -| Allocated size in bytes, accrued from all the replica&#39;s snapshots. | 

## Methods

### NewVolumeUsage

`func NewVolumeUsage(capacity int64, allocated int64, allocatedReplica int64, allocatedSnapshots int64, allocatedAllSnapshots int64, totalAllocated int64, totalAllocatedReplicas interface{}, totalAllocatedSnapshots int64, ) *VolumeUsage`

NewVolumeUsage instantiates a new VolumeUsage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeUsageWithDefaults

`func NewVolumeUsageWithDefaults() *VolumeUsage`

NewVolumeUsageWithDefaults instantiates a new VolumeUsage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacity

`func (o *VolumeUsage) GetCapacity() int64`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *VolumeUsage) GetCapacityOk() (*int64, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *VolumeUsage) SetCapacity(v int64)`

SetCapacity sets Capacity field to given value.


### GetAllocated

`func (o *VolumeUsage) GetAllocated() int64`

GetAllocated returns the Allocated field if non-nil, zero value otherwise.

### GetAllocatedOk

`func (o *VolumeUsage) GetAllocatedOk() (*int64, bool)`

GetAllocatedOk returns a tuple with the Allocated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocated

`func (o *VolumeUsage) SetAllocated(v int64)`

SetAllocated sets Allocated field to given value.


### GetAllocatedReplica

`func (o *VolumeUsage) GetAllocatedReplica() int64`

GetAllocatedReplica returns the AllocatedReplica field if non-nil, zero value otherwise.

### GetAllocatedReplicaOk

`func (o *VolumeUsage) GetAllocatedReplicaOk() (*int64, bool)`

GetAllocatedReplicaOk returns a tuple with the AllocatedReplica field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedReplica

`func (o *VolumeUsage) SetAllocatedReplica(v int64)`

SetAllocatedReplica sets AllocatedReplica field to given value.


### GetAllocatedSnapshots

`func (o *VolumeUsage) GetAllocatedSnapshots() int64`

GetAllocatedSnapshots returns the AllocatedSnapshots field if non-nil, zero value otherwise.

### GetAllocatedSnapshotsOk

`func (o *VolumeUsage) GetAllocatedSnapshotsOk() (*int64, bool)`

GetAllocatedSnapshotsOk returns a tuple with the AllocatedSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedSnapshots

`func (o *VolumeUsage) SetAllocatedSnapshots(v int64)`

SetAllocatedSnapshots sets AllocatedSnapshots field to given value.


### GetAllocatedAllSnapshots

`func (o *VolumeUsage) GetAllocatedAllSnapshots() int64`

GetAllocatedAllSnapshots returns the AllocatedAllSnapshots field if non-nil, zero value otherwise.

### GetAllocatedAllSnapshotsOk

`func (o *VolumeUsage) GetAllocatedAllSnapshotsOk() (*int64, bool)`

GetAllocatedAllSnapshotsOk returns a tuple with the AllocatedAllSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedAllSnapshots

`func (o *VolumeUsage) SetAllocatedAllSnapshots(v int64)`

SetAllocatedAllSnapshots sets AllocatedAllSnapshots field to given value.


### GetTotalAllocated

`func (o *VolumeUsage) GetTotalAllocated() int64`

GetTotalAllocated returns the TotalAllocated field if non-nil, zero value otherwise.

### GetTotalAllocatedOk

`func (o *VolumeUsage) GetTotalAllocatedOk() (*int64, bool)`

GetTotalAllocatedOk returns a tuple with the TotalAllocated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalAllocated

`func (o *VolumeUsage) SetTotalAllocated(v int64)`

SetTotalAllocated sets TotalAllocated field to given value.


### GetTotalAllocatedReplicas

`func (o *VolumeUsage) GetTotalAllocatedReplicas() interface{}`

GetTotalAllocatedReplicas returns the TotalAllocatedReplicas field if non-nil, zero value otherwise.

### GetTotalAllocatedReplicasOk

`func (o *VolumeUsage) GetTotalAllocatedReplicasOk() (*interface{}, bool)`

GetTotalAllocatedReplicasOk returns a tuple with the TotalAllocatedReplicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalAllocatedReplicas

`func (o *VolumeUsage) SetTotalAllocatedReplicas(v interface{})`

SetTotalAllocatedReplicas sets TotalAllocatedReplicas field to given value.


### SetTotalAllocatedReplicasNil

`func (o *VolumeUsage) SetTotalAllocatedReplicasNil(b bool)`

 SetTotalAllocatedReplicasNil sets the value for TotalAllocatedReplicas to be an explicit nil

### UnsetTotalAllocatedReplicas
`func (o *VolumeUsage) UnsetTotalAllocatedReplicas()`

UnsetTotalAllocatedReplicas ensures that no value is present for TotalAllocatedReplicas, not even an explicit nil
### GetTotalAllocatedSnapshots

`func (o *VolumeUsage) GetTotalAllocatedSnapshots() int64`

GetTotalAllocatedSnapshots returns the TotalAllocatedSnapshots field if non-nil, zero value otherwise.

### GetTotalAllocatedSnapshotsOk

`func (o *VolumeUsage) GetTotalAllocatedSnapshotsOk() (*int64, bool)`

GetTotalAllocatedSnapshotsOk returns a tuple with the TotalAllocatedSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalAllocatedSnapshots

`func (o *VolumeUsage) SetTotalAllocatedSnapshots(v int64)`

SetTotalAllocatedSnapshots sets TotalAllocatedSnapshots field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


