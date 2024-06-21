# RegisterAppNode

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Endpoint** | **string** | gRPC server endpoint of the app node. | 
**Labels** | Pointer to **map[string]string** | Labels to be set on the app node. | [optional] 

## Methods

### NewRegisterAppNode

`func NewRegisterAppNode(endpoint string, ) *RegisterAppNode`

NewRegisterAppNode instantiates a new RegisterAppNode object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterAppNodeWithDefaults

`func NewRegisterAppNodeWithDefaults() *RegisterAppNode`

NewRegisterAppNodeWithDefaults instantiates a new RegisterAppNode object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEndpoint

`func (o *RegisterAppNode) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *RegisterAppNode) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *RegisterAppNode) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.


### GetLabels

`func (o *RegisterAppNode) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *RegisterAppNode) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *RegisterAppNode) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *RegisterAppNode) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


