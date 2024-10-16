# AppNodeState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | App node identifier. | 
**Endpoint** | **string** | gRPC server endpoint of the app node. | 
**Status** | **string** | Deemed Status of the app node. | 

## Methods

### NewAppNodeState

`func NewAppNodeState(id string, endpoint string, status string, ) *AppNodeState`

NewAppNodeState instantiates a new AppNodeState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAppNodeStateWithDefaults

`func NewAppNodeStateWithDefaults() *AppNodeState`

NewAppNodeStateWithDefaults instantiates a new AppNodeState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AppNodeState) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AppNodeState) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AppNodeState) SetId(v string)`

SetId sets Id field to given value.


### GetEndpoint

`func (o *AppNodeState) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *AppNodeState) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *AppNodeState) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.


### GetStatus

`func (o *AppNodeState) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *AppNodeState) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *AppNodeState) SetStatus(v string)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


