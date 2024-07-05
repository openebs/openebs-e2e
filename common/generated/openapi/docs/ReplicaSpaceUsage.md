# ReplicaSpaceUsage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CapacityBytes** | **int64** | Replica capacity in bytes. | 
**AllocatedBytes** | **int64** | Amount of actually allocated disk space for this replica in bytes. | 
**AllocatedBytesSnapshots** | **int64** | Amount of actually allocated disk space for this replica&#39;s snapshots in bytes. | [default to 0]
**AllocatedBytesAllSnapshots** | Pointer to **int64** | Amount of actually allocated disk space for this replica&#39;s snapshots and its predecessors in bytes. For a restored/cloned replica this includes snapshots from the parent source.  | [optional] [default to 0]
**ClusterSize** | **int64** | Cluster size in bytes. | 
**Clusters** | **int64** | Total number of clusters. | 
**AllocatedClusters** | **int64** | Number of actually used clusters. | 

## Methods

### NewReplicaSpaceUsage

`func NewReplicaSpaceUsage(capacityBytes int64, allocatedBytes int64, allocatedBytesSnapshots int64, clusterSize int64, clusters int64, allocatedClusters int64, ) *ReplicaSpaceUsage`

NewReplicaSpaceUsage instantiates a new ReplicaSpaceUsage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSpaceUsageWithDefaults

`func NewReplicaSpaceUsageWithDefaults() *ReplicaSpaceUsage`

NewReplicaSpaceUsageWithDefaults instantiates a new ReplicaSpaceUsage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacityBytes

`func (o *ReplicaSpaceUsage) GetCapacityBytes() int64`

GetCapacityBytes returns the CapacityBytes field if non-nil, zero value otherwise.

### GetCapacityBytesOk

`func (o *ReplicaSpaceUsage) GetCapacityBytesOk() (*int64, bool)`

GetCapacityBytesOk returns a tuple with the CapacityBytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacityBytes

`func (o *ReplicaSpaceUsage) SetCapacityBytes(v int64)`

SetCapacityBytes sets CapacityBytes field to given value.


### GetAllocatedBytes

`func (o *ReplicaSpaceUsage) GetAllocatedBytes() int64`

GetAllocatedBytes returns the AllocatedBytes field if non-nil, zero value otherwise.

### GetAllocatedBytesOk

`func (o *ReplicaSpaceUsage) GetAllocatedBytesOk() (*int64, bool)`

GetAllocatedBytesOk returns a tuple with the AllocatedBytes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedBytes

`func (o *ReplicaSpaceUsage) SetAllocatedBytes(v int64)`

SetAllocatedBytes sets AllocatedBytes field to given value.


### GetAllocatedBytesSnapshots

`func (o *ReplicaSpaceUsage) GetAllocatedBytesSnapshots() int64`

GetAllocatedBytesSnapshots returns the AllocatedBytesSnapshots field if non-nil, zero value otherwise.

### GetAllocatedBytesSnapshotsOk

`func (o *ReplicaSpaceUsage) GetAllocatedBytesSnapshotsOk() (*int64, bool)`

GetAllocatedBytesSnapshotsOk returns a tuple with the AllocatedBytesSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedBytesSnapshots

`func (o *ReplicaSpaceUsage) SetAllocatedBytesSnapshots(v int64)`

SetAllocatedBytesSnapshots sets AllocatedBytesSnapshots field to given value.


### GetAllocatedBytesAllSnapshots

`func (o *ReplicaSpaceUsage) GetAllocatedBytesAllSnapshots() int64`

GetAllocatedBytesAllSnapshots returns the AllocatedBytesAllSnapshots field if non-nil, zero value otherwise.

### GetAllocatedBytesAllSnapshotsOk

`func (o *ReplicaSpaceUsage) GetAllocatedBytesAllSnapshotsOk() (*int64, bool)`

GetAllocatedBytesAllSnapshotsOk returns a tuple with the AllocatedBytesAllSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedBytesAllSnapshots

`func (o *ReplicaSpaceUsage) SetAllocatedBytesAllSnapshots(v int64)`

SetAllocatedBytesAllSnapshots sets AllocatedBytesAllSnapshots field to given value.

### HasAllocatedBytesAllSnapshots

`func (o *ReplicaSpaceUsage) HasAllocatedBytesAllSnapshots() bool`

HasAllocatedBytesAllSnapshots returns a boolean if a field has been set.

### GetClusterSize

`func (o *ReplicaSpaceUsage) GetClusterSize() int64`

GetClusterSize returns the ClusterSize field if non-nil, zero value otherwise.

### GetClusterSizeOk

`func (o *ReplicaSpaceUsage) GetClusterSizeOk() (*int64, bool)`

GetClusterSizeOk returns a tuple with the ClusterSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterSize

`func (o *ReplicaSpaceUsage) SetClusterSize(v int64)`

SetClusterSize sets ClusterSize field to given value.


### GetClusters

`func (o *ReplicaSpaceUsage) GetClusters() int64`

GetClusters returns the Clusters field if non-nil, zero value otherwise.

### GetClustersOk

`func (o *ReplicaSpaceUsage) GetClustersOk() (*int64, bool)`

GetClustersOk returns a tuple with the Clusters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusters

`func (o *ReplicaSpaceUsage) SetClusters(v int64)`

SetClusters sets Clusters field to given value.


### GetAllocatedClusters

`func (o *ReplicaSpaceUsage) GetAllocatedClusters() int64`

GetAllocatedClusters returns the AllocatedClusters field if non-nil, zero value otherwise.

### GetAllocatedClustersOk

`func (o *ReplicaSpaceUsage) GetAllocatedClustersOk() (*int64, bool)`

GetAllocatedClustersOk returns a tuple with the AllocatedClusters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedClusters

`func (o *ReplicaSpaceUsage) SetAllocatedClusters(v int64)`

SetAllocatedClusters sets AllocatedClusters field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


