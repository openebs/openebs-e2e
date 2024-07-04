# Pool

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | storage pool identifier | 
**Spec** | Pointer to [**PoolSpec**](PoolSpec.md) |  | [optional] 
**State** | Pointer to [**PoolState**](PoolState.md) |  | [optional] 

## Methods

### NewPool

`func NewPool(id string, ) *Pool`

NewPool instantiates a new Pool object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolWithDefaults

`func NewPoolWithDefaults() *Pool`

NewPoolWithDefaults instantiates a new Pool object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Pool) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Pool) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Pool) SetId(v string)`

SetId sets Id field to given value.


### GetSpec

`func (o *Pool) GetSpec() PoolSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *Pool) GetSpecOk() (*PoolSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *Pool) SetSpec(v PoolSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *Pool) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetState

`func (o *Pool) GetState() PoolState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Pool) GetStateOk() (*PoolState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Pool) SetState(v PoolState)`

SetState sets State field to given value.

### HasState

`func (o *Pool) HasState() bool`

HasState returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


