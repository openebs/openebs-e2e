# ReplicaSnapshot

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**SourceId** | **string** |  | 
**Status** | [**SpecStatus**](SpecStatus.md) |  | 

## Methods

### NewReplicaSnapshot

`func NewReplicaSnapshot(uuid string, sourceId string, status SpecStatus, ) *ReplicaSnapshot`

NewReplicaSnapshot instantiates a new ReplicaSnapshot object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSnapshotWithDefaults

`func NewReplicaSnapshotWithDefaults() *ReplicaSnapshot`

NewReplicaSnapshotWithDefaults instantiates a new ReplicaSnapshot object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *ReplicaSnapshot) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ReplicaSnapshot) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ReplicaSnapshot) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetSourceId

`func (o *ReplicaSnapshot) GetSourceId() string`

GetSourceId returns the SourceId field if non-nil, zero value otherwise.

### GetSourceIdOk

`func (o *ReplicaSnapshot) GetSourceIdOk() (*string, bool)`

GetSourceIdOk returns a tuple with the SourceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceId

`func (o *ReplicaSnapshot) SetSourceId(v string)`

SetSourceId sets SourceId field to given value.


### GetStatus

`func (o *ReplicaSnapshot) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ReplicaSnapshot) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ReplicaSnapshot) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


