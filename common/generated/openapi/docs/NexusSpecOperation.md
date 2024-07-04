# NexusSpecOperation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Operation** | **string** | Record of the operation | 
**Result** | Pointer to **bool** | Result of the operation | [optional] 

## Methods

### NewNexusSpecOperation

`func NewNexusSpecOperation(operation string, ) *NexusSpecOperation`

NewNexusSpecOperation instantiates a new NexusSpecOperation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNexusSpecOperationWithDefaults

`func NewNexusSpecOperationWithDefaults() *NexusSpecOperation`

NewNexusSpecOperationWithDefaults instantiates a new NexusSpecOperation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOperation

`func (o *NexusSpecOperation) GetOperation() string`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *NexusSpecOperation) GetOperationOk() (*string, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *NexusSpecOperation) SetOperation(v string)`

SetOperation sets Operation field to given value.


### GetResult

`func (o *NexusSpecOperation) GetResult() bool`

GetResult returns the Result field if non-nil, zero value otherwise.

### GetResultOk

`func (o *NexusSpecOperation) GetResultOk() (*bool, bool)`

GetResultOk returns a tuple with the Result field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResult

`func (o *NexusSpecOperation) SetResult(v bool)`

SetResult sets Result field to given value.

### HasResult

`func (o *NexusSpecOperation) HasResult() bool`

HasResult returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


