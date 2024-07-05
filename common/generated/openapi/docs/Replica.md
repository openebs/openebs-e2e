# Replica

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Node** | **string** | storage node identifier | 
**Pool** | **string** | storage pool identifier | 
**PoolUuid** | Pointer to **string** | storage pool unique identifier | [optional] 
**Share** | [**Protocol**](Protocol.md) |  | 
**Size** | **int64** | size of the replica in bytes | 
**Space** | Pointer to [**ReplicaSpaceUsage**](ReplicaSpaceUsage.md) |  | [optional] 
**State** | [**ReplicaState**](ReplicaState.md) |  | 
**Thin** | **bool** | thin provisioning | 
**Uri** | **string** | uri usable by nexus to access it | 
**Uuid** | **string** | uuid of the replica | 
**AllowedHosts** | Pointer to **[]string** | NQNs of hosts allowed to connect to this replica | [optional] 
**Kind** | [**ReplicaKind**](ReplicaKind.md) |  | [default to REGULAR]

## Methods

### NewReplica

`func NewReplica(node string, pool string, share Protocol, size int64, state ReplicaState, thin bool, uri string, uuid string, kind ReplicaKind, ) *Replica`

NewReplica instantiates a new Replica object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaWithDefaults

`func NewReplicaWithDefaults() *Replica`

NewReplicaWithDefaults instantiates a new Replica object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNode

`func (o *Replica) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *Replica) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *Replica) SetNode(v string)`

SetNode sets Node field to given value.


### GetPool

`func (o *Replica) GetPool() string`

GetPool returns the Pool field if non-nil, zero value otherwise.

### GetPoolOk

`func (o *Replica) GetPoolOk() (*string, bool)`

GetPoolOk returns a tuple with the Pool field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPool

`func (o *Replica) SetPool(v string)`

SetPool sets Pool field to given value.


### GetPoolUuid

`func (o *Replica) GetPoolUuid() string`

GetPoolUuid returns the PoolUuid field if non-nil, zero value otherwise.

### GetPoolUuidOk

`func (o *Replica) GetPoolUuidOk() (*string, bool)`

GetPoolUuidOk returns a tuple with the PoolUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolUuid

`func (o *Replica) SetPoolUuid(v string)`

SetPoolUuid sets PoolUuid field to given value.

### HasPoolUuid

`func (o *Replica) HasPoolUuid() bool`

HasPoolUuid returns a boolean if a field has been set.

### GetShare

`func (o *Replica) GetShare() Protocol`

GetShare returns the Share field if non-nil, zero value otherwise.

### GetShareOk

`func (o *Replica) GetShareOk() (*Protocol, bool)`

GetShareOk returns a tuple with the Share field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShare

`func (o *Replica) SetShare(v Protocol)`

SetShare sets Share field to given value.


### GetSize

`func (o *Replica) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Replica) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Replica) SetSize(v int64)`

SetSize sets Size field to given value.


### GetSpace

`func (o *Replica) GetSpace() ReplicaSpaceUsage`

GetSpace returns the Space field if non-nil, zero value otherwise.

### GetSpaceOk

`func (o *Replica) GetSpaceOk() (*ReplicaSpaceUsage, bool)`

GetSpaceOk returns a tuple with the Space field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpace

`func (o *Replica) SetSpace(v ReplicaSpaceUsage)`

SetSpace sets Space field to given value.

### HasSpace

`func (o *Replica) HasSpace() bool`

HasSpace returns a boolean if a field has been set.

### GetState

`func (o *Replica) GetState() ReplicaState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Replica) GetStateOk() (*ReplicaState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Replica) SetState(v ReplicaState)`

SetState sets State field to given value.


### GetThin

`func (o *Replica) GetThin() bool`

GetThin returns the Thin field if non-nil, zero value otherwise.

### GetThinOk

`func (o *Replica) GetThinOk() (*bool, bool)`

GetThinOk returns a tuple with the Thin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThin

`func (o *Replica) SetThin(v bool)`

SetThin sets Thin field to given value.


### GetUri

`func (o *Replica) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *Replica) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *Replica) SetUri(v string)`

SetUri sets Uri field to given value.


### GetUuid

`func (o *Replica) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *Replica) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *Replica) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetAllowedHosts

`func (o *Replica) GetAllowedHosts() []string`

GetAllowedHosts returns the AllowedHosts field if non-nil, zero value otherwise.

### GetAllowedHostsOk

`func (o *Replica) GetAllowedHostsOk() (*[]string, bool)`

GetAllowedHostsOk returns a tuple with the AllowedHosts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedHosts

`func (o *Replica) SetAllowedHosts(v []string)`

SetAllowedHosts sets AllowedHosts field to given value.

### HasAllowedHosts

`func (o *Replica) HasAllowedHosts() bool`

HasAllowedHosts returns a boolean if a field has been set.

### GetKind

`func (o *Replica) GetKind() ReplicaKind`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *Replica) GetKindOk() (*ReplicaKind, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *Replica) SetKind(v ReplicaKind)`

SetKind sets Kind field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


