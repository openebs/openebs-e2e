# AppNodeSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | App node identifier. | 
**Endpoint** | **string** | gRPC server endpoint of the app node. | 
**Labels** | Pointer to **map[string]string** | Labels to be set on the app node. | [optional] 

## Methods

### NewAppNodeSpec

`func NewAppNodeSpec(id string, endpoint string, ) *AppNodeSpec`

NewAppNodeSpec instantiates a new AppNodeSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAppNodeSpecWithDefaults

`func NewAppNodeSpecWithDefaults() *AppNodeSpec`

NewAppNodeSpecWithDefaults instantiates a new AppNodeSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *AppNodeSpec) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AppNodeSpec) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AppNodeSpec) SetId(v string)`

SetId sets Id field to given value.


### GetEndpoint

`func (o *AppNodeSpec) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *AppNodeSpec) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *AppNodeSpec) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.


### GetLabels

`func (o *AppNodeSpec) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *AppNodeSpec) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *AppNodeSpec) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *AppNodeSpec) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


