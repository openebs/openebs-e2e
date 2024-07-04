# Child

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RebuildProgress** | Pointer to **int32** | current rebuild progress (%) | [optional] 
**State** | [**ChildState**](ChildState.md) |  | 
**StateReason** | Pointer to [**ChildStateReason**](ChildStateReason.md) |  | [optional] 
**Uri** | **string** | uri of the child device | 

## Methods

### NewChild

`func NewChild(state ChildState, uri string, ) *Child`

NewChild instantiates a new Child object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChildWithDefaults

`func NewChildWithDefaults() *Child`

NewChildWithDefaults instantiates a new Child object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRebuildProgress

`func (o *Child) GetRebuildProgress() int32`

GetRebuildProgress returns the RebuildProgress field if non-nil, zero value otherwise.

### GetRebuildProgressOk

`func (o *Child) GetRebuildProgressOk() (*int32, bool)`

GetRebuildProgressOk returns a tuple with the RebuildProgress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRebuildProgress

`func (o *Child) SetRebuildProgress(v int32)`

SetRebuildProgress sets RebuildProgress field to given value.

### HasRebuildProgress

`func (o *Child) HasRebuildProgress() bool`

HasRebuildProgress returns a boolean if a field has been set.

### GetState

`func (o *Child) GetState() ChildState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Child) GetStateOk() (*ChildState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Child) SetState(v ChildState)`

SetState sets State field to given value.


### GetStateReason

`func (o *Child) GetStateReason() ChildStateReason`

GetStateReason returns the StateReason field if non-nil, zero value otherwise.

### GetStateReasonOk

`func (o *Child) GetStateReasonOk() (*ChildStateReason, bool)`

GetStateReasonOk returns a tuple with the StateReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateReason

`func (o *Child) SetStateReason(v ChildStateReason)`

SetStateReason sets StateReason field to given value.

### HasStateReason

`func (o *Child) HasStateReason() bool`

HasStateReason returns a boolean if a field has been set.

### GetUri

`func (o *Child) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *Child) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *Child) SetUri(v string)`

SetUri sets Uri field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


