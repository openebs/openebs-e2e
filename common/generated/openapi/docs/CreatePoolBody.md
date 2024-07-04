# CreatePoolBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Disks** | **[]string** | disk device paths or URIs to be claimed by the pool | 
**Labels** | Pointer to **map[string]string** | labels to be set on the pools | [optional] 

## Methods

### NewCreatePoolBody

`func NewCreatePoolBody(disks []string, ) *CreatePoolBody`

NewCreatePoolBody instantiates a new CreatePoolBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePoolBodyWithDefaults

`func NewCreatePoolBodyWithDefaults() *CreatePoolBody`

NewCreatePoolBodyWithDefaults instantiates a new CreatePoolBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDisks

`func (o *CreatePoolBody) GetDisks() []string`

GetDisks returns the Disks field if non-nil, zero value otherwise.

### GetDisksOk

`func (o *CreatePoolBody) GetDisksOk() (*[]string, bool)`

GetDisksOk returns a tuple with the Disks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisks

`func (o *CreatePoolBody) SetDisks(v []string)`

SetDisks sets Disks field to given value.


### GetLabels

`func (o *CreatePoolBody) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CreatePoolBody) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CreatePoolBody) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CreatePoolBody) HasLabels() bool`

HasLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


