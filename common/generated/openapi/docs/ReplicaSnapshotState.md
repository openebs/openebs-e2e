# ReplicaSnapshotState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Online** | Pointer to [**OnlineReplicaSnapshotState**](OnlineReplicaSnapshotState.md) |  | [optional] 
**Offline** | Pointer to [**OfflineReplicaSnapshotState**](OfflineReplicaSnapshotState.md) |  | [optional] 

## Methods

### NewReplicaSnapshotState

`func NewReplicaSnapshotState() *ReplicaSnapshotState`

NewReplicaSnapshotState instantiates a new ReplicaSnapshotState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSnapshotStateWithDefaults

`func NewReplicaSnapshotStateWithDefaults() *ReplicaSnapshotState`

NewReplicaSnapshotStateWithDefaults instantiates a new ReplicaSnapshotState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOnline

`func (o *ReplicaSnapshotState) GetOnline() OnlineReplicaSnapshotState`

GetOnline returns the Online field if non-nil, zero value otherwise.

### GetOnlineOk

`func (o *ReplicaSnapshotState) GetOnlineOk() (*OnlineReplicaSnapshotState, bool)`

GetOnlineOk returns a tuple with the Online field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnline

`func (o *ReplicaSnapshotState) SetOnline(v OnlineReplicaSnapshotState)`

SetOnline sets Online field to given value.

### HasOnline

`func (o *ReplicaSnapshotState) HasOnline() bool`

HasOnline returns a boolean if a field has been set.

### GetOffline

`func (o *ReplicaSnapshotState) GetOffline() OfflineReplicaSnapshotState`

GetOffline returns the Offline field if non-nil, zero value otherwise.

### GetOfflineOk

`func (o *ReplicaSnapshotState) GetOfflineOk() (*OfflineReplicaSnapshotState, bool)`

GetOfflineOk returns a tuple with the Offline field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOffline

`func (o *ReplicaSnapshotState) SetOffline(v OfflineReplicaSnapshotState)`

SetOffline sets Offline field to given value.

### HasOffline

`func (o *ReplicaSnapshotState) HasOffline() bool`

HasOffline returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


