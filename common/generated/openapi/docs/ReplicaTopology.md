# ReplicaTopology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Node** | Pointer to **string** | storage node identifier | [optional] 
**Pool** | Pointer to **string** | storage pool identifier | [optional] 
**State** | [**ReplicaState**](ReplicaState.md) |  | 
**ChildStatus** | Pointer to [**ChildState**](ChildState.md) |  | [optional] 
**ChildStatusReason** | Pointer to [**ChildStateReason**](ChildStateReason.md) |  | [optional] 
**Usage** | Pointer to [**ReplicaUsage**](ReplicaUsage.md) |  | [optional] 
**RebuildProgress** | Pointer to **int32** | current rebuild progress (%) | [optional] 

## Methods

### NewReplicaTopology

`func NewReplicaTopology(state ReplicaState, ) *ReplicaTopology`

NewReplicaTopology instantiates a new ReplicaTopology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaTopologyWithDefaults

`func NewReplicaTopologyWithDefaults() *ReplicaTopology`

NewReplicaTopologyWithDefaults instantiates a new ReplicaTopology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNode

`func (o *ReplicaTopology) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *ReplicaTopology) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *ReplicaTopology) SetNode(v string)`

SetNode sets Node field to given value.

### HasNode

`func (o *ReplicaTopology) HasNode() bool`

HasNode returns a boolean if a field has been set.

### GetPool

`func (o *ReplicaTopology) GetPool() string`

GetPool returns the Pool field if non-nil, zero value otherwise.

### GetPoolOk

`func (o *ReplicaTopology) GetPoolOk() (*string, bool)`

GetPoolOk returns a tuple with the Pool field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPool

`func (o *ReplicaTopology) SetPool(v string)`

SetPool sets Pool field to given value.

### HasPool

`func (o *ReplicaTopology) HasPool() bool`

HasPool returns a boolean if a field has been set.

### GetState

`func (o *ReplicaTopology) GetState() ReplicaState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ReplicaTopology) GetStateOk() (*ReplicaState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ReplicaTopology) SetState(v ReplicaState)`

SetState sets State field to given value.


### GetChildStatus

`func (o *ReplicaTopology) GetChildStatus() ChildState`

GetChildStatus returns the ChildStatus field if non-nil, zero value otherwise.

### GetChildStatusOk

`func (o *ReplicaTopology) GetChildStatusOk() (*ChildState, bool)`

GetChildStatusOk returns a tuple with the ChildStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildStatus

`func (o *ReplicaTopology) SetChildStatus(v ChildState)`

SetChildStatus sets ChildStatus field to given value.

### HasChildStatus

`func (o *ReplicaTopology) HasChildStatus() bool`

HasChildStatus returns a boolean if a field has been set.

### GetChildStatusReason

`func (o *ReplicaTopology) GetChildStatusReason() ChildStateReason`

GetChildStatusReason returns the ChildStatusReason field if non-nil, zero value otherwise.

### GetChildStatusReasonOk

`func (o *ReplicaTopology) GetChildStatusReasonOk() (*ChildStateReason, bool)`

GetChildStatusReasonOk returns a tuple with the ChildStatusReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildStatusReason

`func (o *ReplicaTopology) SetChildStatusReason(v ChildStateReason)`

SetChildStatusReason sets ChildStatusReason field to given value.

### HasChildStatusReason

`func (o *ReplicaTopology) HasChildStatusReason() bool`

HasChildStatusReason returns a boolean if a field has been set.

### GetUsage

`func (o *ReplicaTopology) GetUsage() ReplicaUsage`

GetUsage returns the Usage field if non-nil, zero value otherwise.

### GetUsageOk

`func (o *ReplicaTopology) GetUsageOk() (*ReplicaUsage, bool)`

GetUsageOk returns a tuple with the Usage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsage

`func (o *ReplicaTopology) SetUsage(v ReplicaUsage)`

SetUsage sets Usage field to given value.

### HasUsage

`func (o *ReplicaTopology) HasUsage() bool`

HasUsage returns a boolean if a field has been set.

### GetRebuildProgress

`func (o *ReplicaTopology) GetRebuildProgress() int32`

GetRebuildProgress returns the RebuildProgress field if non-nil, zero value otherwise.

### GetRebuildProgressOk

`func (o *ReplicaTopology) GetRebuildProgressOk() (*int32, bool)`

GetRebuildProgressOk returns a tuple with the RebuildProgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRebuildProgress

`func (o *ReplicaTopology) SetRebuildProgress(v int32)`

SetRebuildProgress sets RebuildProgress field to given value.

### HasRebuildProgress

`func (o *ReplicaTopology) HasRebuildProgress() bool`

HasRebuildProgress returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


