# PoolState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Capacity** | **int64** | size of the pool in bytes | 
**Disks** | **[]string** | absolute disk paths claimed by the pool | 
**Id** | **string** | storage pool identifier | 
**Node** | **string** | storage node identifier | 
**Status** | [**PoolStatus**](PoolStatus.md) |  | 
**Used** | **int64** | used bytes from the pool | 
**Committed** | Pointer to **int64** | accrued size of all replicas contained in this pool | [optional] 

## Methods

### NewPoolState

`func NewPoolState(capacity int64, disks []string, id string, node string, status PoolStatus, used int64, ) *PoolState`

NewPoolState instantiates a new PoolState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolStateWithDefaults

`func NewPoolStateWithDefaults() *PoolState`

NewPoolStateWithDefaults instantiates a new PoolState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacity

`func (o *PoolState) GetCapacity() int64`

GetCapacity returns the Capacity field if non-nil, zero value otherwise.

### GetCapacityOk

`func (o *PoolState) GetCapacityOk() (*int64, bool)`

GetCapacityOk returns a tuple with the Capacity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacity

`func (o *PoolState) SetCapacity(v int64)`

SetCapacity sets Capacity field to given value.


### GetDisks

`func (o *PoolState) GetDisks() []string`

GetDisks returns the Disks field if non-nil, zero value otherwise.

### GetDisksOk

`func (o *PoolState) GetDisksOk() (*[]string, bool)`

GetDisksOk returns a tuple with the Disks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisks

`func (o *PoolState) SetDisks(v []string)`

SetDisks sets Disks field to given value.


### GetId

`func (o *PoolState) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PoolState) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PoolState) SetId(v string)`

SetId sets Id field to given value.


### GetNode

`func (o *PoolState) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *PoolState) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *PoolState) SetNode(v string)`

SetNode sets Node field to given value.


### GetStatus

`func (o *PoolState) GetStatus() PoolStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PoolState) GetStatusOk() (*PoolStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PoolState) SetStatus(v PoolStatus)`

SetStatus sets Status field to given value.


### GetUsed

`func (o *PoolState) GetUsed() int64`

GetUsed returns the Used field if non-nil, zero value otherwise.

### GetUsedOk

`func (o *PoolState) GetUsedOk() (*int64, bool)`

GetUsedOk returns a tuple with the Used field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsed

`func (o *PoolState) SetUsed(v int64)`

SetUsed sets Used field to given value.


### GetCommitted

`func (o *PoolState) GetCommitted() int64`

GetCommitted returns the Committed field if non-nil, zero value otherwise.

### GetCommittedOk

`func (o *PoolState) GetCommittedOk() (*int64, bool)`

GetCommittedOk returns a tuple with the Committed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitted

`func (o *PoolState) SetCommitted(v int64)`

SetCommitted sets Committed field to given value.

### HasCommitted

`func (o *PoolState) HasCommitted() bool`

HasCommitted returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


