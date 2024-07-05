# AppNodes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Entries** | [**[]AppNode**](AppNode.md) |  | 
**NextToken** | Pointer to **int32** |  | [optional] 

## Methods

### NewAppNodes

`func NewAppNodes(entries []AppNode, ) *AppNodes`

NewAppNodes instantiates a new AppNodes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAppNodesWithDefaults

`func NewAppNodesWithDefaults() *AppNodes`

NewAppNodesWithDefaults instantiates a new AppNodes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEntries

`func (o *AppNodes) GetEntries() []AppNode`

GetEntries returns the Entries field if non-nil, zero value otherwise.

### GetEntriesOk

`func (o *AppNodes) GetEntriesOk() (*[]AppNode, bool)`

GetEntriesOk returns a tuple with the Entries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntries

`func (o *AppNodes) SetEntries(v []AppNode)`

SetEntries sets Entries field to given value.


### GetNextToken

`func (o *AppNodes) GetNextToken() int32`

GetNextToken returns the NextToken field if non-nil, zero value otherwise.

### GetNextTokenOk

`func (o *AppNodes) GetNextTokenOk() (*int32, bool)`

GetNextTokenOk returns a tuple with the NextToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextToken

`func (o *AppNodes) SetNextToken(v int32)`

SetNextToken sets NextToken field to given value.

### HasNextToken

`func (o *AppNodes) HasNextToken() bool`

HasNextToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


