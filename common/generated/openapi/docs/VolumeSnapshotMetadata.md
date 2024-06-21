# VolumeSnapshotMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | [**SpecStatus**](SpecStatus.md) |  | 
**Timestamp** | Pointer to **time.Time** | Timestamp when snapshot is taken on the storage system. | [optional] 
**Size** | **int64** | Size in bytes of the snapshot (which is equivalent to its source size). | 
**SpecSize** | **int64** | Spec size in bytes of the snapshot (which is equivalent to its source spec size). | 
**TotalAllocatedSize** | **int64** | Size in bytes taken by the snapshot and its predecessors. | 
**TxnId** | **string** |  | 
**Transactions** | [**map[string][]ReplicaSnapshot**](array.md) |  | 
**NumRestores** | **int32** | Number of restores done from this snapshot. | 
**NumSnapshotReplicas** | **int32** | Number of snapshot replicas for a volumesnapshot. | 

## Methods

### NewVolumeSnapshotMetadata

`func NewVolumeSnapshotMetadata(status SpecStatus, size int64, specSize int64, totalAllocatedSize int64, txnId string, transactions map[string][]ReplicaSnapshot, numRestores int32, numSnapshotReplicas int32, ) *VolumeSnapshotMetadata`

NewVolumeSnapshotMetadata instantiates a new VolumeSnapshotMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeSnapshotMetadataWithDefaults

`func NewVolumeSnapshotMetadataWithDefaults() *VolumeSnapshotMetadata`

NewVolumeSnapshotMetadataWithDefaults instantiates a new VolumeSnapshotMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *VolumeSnapshotMetadata) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VolumeSnapshotMetadata) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VolumeSnapshotMetadata) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.


### GetTimestamp

`func (o *VolumeSnapshotMetadata) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *VolumeSnapshotMetadata) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *VolumeSnapshotMetadata) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *VolumeSnapshotMetadata) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.

### GetSize

`func (o *VolumeSnapshotMetadata) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *VolumeSnapshotMetadata) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *VolumeSnapshotMetadata) SetSize(v int64)`

SetSize sets Size field to given value.


### GetSpecSize

`func (o *VolumeSnapshotMetadata) GetSpecSize() int64`

GetSpecSize returns the SpecSize field if non-nil, zero value otherwise.

### GetSpecSizeOk

`func (o *VolumeSnapshotMetadata) GetSpecSizeOk() (*int64, bool)`

GetSpecSizeOk returns a tuple with the SpecSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpecSize

`func (o *VolumeSnapshotMetadata) SetSpecSize(v int64)`

SetSpecSize sets SpecSize field to given value.


### GetTotalAllocatedSize

`func (o *VolumeSnapshotMetadata) GetTotalAllocatedSize() int64`

GetTotalAllocatedSize returns the TotalAllocatedSize field if non-nil, zero value otherwise.

### GetTotalAllocatedSizeOk

`func (o *VolumeSnapshotMetadata) GetTotalAllocatedSizeOk() (*int64, bool)`

GetTotalAllocatedSizeOk returns a tuple with the TotalAllocatedSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalAllocatedSize

`func (o *VolumeSnapshotMetadata) SetTotalAllocatedSize(v int64)`

SetTotalAllocatedSize sets TotalAllocatedSize field to given value.


### GetTxnId

`func (o *VolumeSnapshotMetadata) GetTxnId() string`

GetTxnId returns the TxnId field if non-nil, zero value otherwise.

### GetTxnIdOk

`func (o *VolumeSnapshotMetadata) GetTxnIdOk() (*string, bool)`

GetTxnIdOk returns a tuple with the TxnId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTxnId

`func (o *VolumeSnapshotMetadata) SetTxnId(v string)`

SetTxnId sets TxnId field to given value.


### GetTransactions

`func (o *VolumeSnapshotMetadata) GetTransactions() map[string][]ReplicaSnapshot`

GetTransactions returns the Transactions field if non-nil, zero value otherwise.

### GetTransactionsOk

`func (o *VolumeSnapshotMetadata) GetTransactionsOk() (*map[string][]ReplicaSnapshot, bool)`

GetTransactionsOk returns a tuple with the Transactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactions

`func (o *VolumeSnapshotMetadata) SetTransactions(v map[string][]ReplicaSnapshot)`

SetTransactions sets Transactions field to given value.


### GetNumRestores

`func (o *VolumeSnapshotMetadata) GetNumRestores() int32`

GetNumRestores returns the NumRestores field if non-nil, zero value otherwise.

### GetNumRestoresOk

`func (o *VolumeSnapshotMetadata) GetNumRestoresOk() (*int32, bool)`

GetNumRestoresOk returns a tuple with the NumRestores field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumRestores

`func (o *VolumeSnapshotMetadata) SetNumRestores(v int32)`

SetNumRestores sets NumRestores field to given value.


### GetNumSnapshotReplicas

`func (o *VolumeSnapshotMetadata) GetNumSnapshotReplicas() int32`

GetNumSnapshotReplicas returns the NumSnapshotReplicas field if non-nil, zero value otherwise.

### GetNumSnapshotReplicasOk

`func (o *VolumeSnapshotMetadata) GetNumSnapshotReplicasOk() (*int32, bool)`

GetNumSnapshotReplicasOk returns a tuple with the NumSnapshotReplicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumSnapshotReplicas

`func (o *VolumeSnapshotMetadata) SetNumSnapshotReplicas(v int32)`

SetNumSnapshotReplicas sets NumSnapshotReplicas field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


