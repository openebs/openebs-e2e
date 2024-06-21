# CreateNexusBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Children** | **[]string** | replica can be iscsi and nvmf remote targets or a local spdk bdev  (i.e. bdev:///name-of-the-bdev).   uris to the targets we connect to | 
**Size** | **int64** | size of the device in bytes | 

## Methods

### NewCreateNexusBody

`func NewCreateNexusBody(children []string, size int64, ) *CreateNexusBody`

NewCreateNexusBody instantiates a new CreateNexusBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateNexusBodyWithDefaults

`func NewCreateNexusBodyWithDefaults() *CreateNexusBody`

NewCreateNexusBodyWithDefaults instantiates a new CreateNexusBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChildren

`func (o *CreateNexusBody) GetChildren() []string`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *CreateNexusBody) GetChildrenOk() (*[]string, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *CreateNexusBody) SetChildren(v []string)`

SetChildren sets Children field to given value.


### GetSize

`func (o *CreateNexusBody) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *CreateNexusBody) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *CreateNexusBody) SetSize(v int64)`

SetSize sets Size field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


