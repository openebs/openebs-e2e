# ReplicaSpecOperation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Operation** | **string** | Record of the operation | 
**Result** | Pointer to **bool** | Result of the operation | [optional] 

## Methods

### NewReplicaSpecOperation

`func NewReplicaSpecOperation(operation string, ) *ReplicaSpecOperation`

NewReplicaSpecOperation instantiates a new ReplicaSpecOperation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSpecOperationWithDefaults

`func NewReplicaSpecOperationWithDefaults() *ReplicaSpecOperation`

NewReplicaSpecOperationWithDefaults instantiates a new ReplicaSpecOperation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOperation

`func (o *ReplicaSpecOperation) GetOperation() string`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *ReplicaSpecOperation) GetOperationOk() (*string, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *ReplicaSpecOperation) SetOperation(v string)`

SetOperation sets Operation field to given value.


### GetResult

`func (o *ReplicaSpecOperation) GetResult() bool`

GetResult returns the Result field if non-nil, zero value otherwise.

### GetResultOk

`func (o *ReplicaSpecOperation) GetResultOk() (*bool, bool)`

GetResultOk returns a tuple with the Result field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResult

`func (o *ReplicaSpecOperation) SetResult(v bool)`

SetResult sets Result field to given value.

### HasResult

`func (o *ReplicaSpecOperation) HasResult() bool`

HasResult returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


