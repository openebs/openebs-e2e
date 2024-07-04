# CreateReplicaBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Share** | Pointer to [**ReplicaShareProtocol**](ReplicaShareProtocol.md) |  | [optional] 
**AllowedHosts** | Pointer to **[]string** |  | [optional] 
**Size** | **int64** | size of the replica in bytes | 
**Thin** | **bool** | thin provisioning | 

## Methods

### NewCreateReplicaBody

`func NewCreateReplicaBody(size int64, thin bool, ) *CreateReplicaBody`

NewCreateReplicaBody instantiates a new CreateReplicaBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateReplicaBodyWithDefaults

`func NewCreateReplicaBodyWithDefaults() *CreateReplicaBody`

NewCreateReplicaBodyWithDefaults instantiates a new CreateReplicaBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetShare

`func (o *CreateReplicaBody) GetShare() ReplicaShareProtocol`

GetShare returns the Share field if non-nil, zero value otherwise.

### GetShareOk

`func (o *CreateReplicaBody) GetShareOk() (*ReplicaShareProtocol, bool)`

GetShareOk returns a tuple with the Share field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShare

`func (o *CreateReplicaBody) SetShare(v ReplicaShareProtocol)`

SetShare sets Share field to given value.

### HasShare

`func (o *CreateReplicaBody) HasShare() bool`

HasShare returns a boolean if a field has been set.

### GetAllowedHosts

`func (o *CreateReplicaBody) GetAllowedHosts() []string`

GetAllowedHosts returns the AllowedHosts field if non-nil, zero value otherwise.

### GetAllowedHostsOk

`func (o *CreateReplicaBody) GetAllowedHostsOk() (*[]string, bool)`

GetAllowedHostsOk returns a tuple with the AllowedHosts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedHosts

`func (o *CreateReplicaBody) SetAllowedHosts(v []string)`

SetAllowedHosts sets AllowedHosts field to given value.

### HasAllowedHosts

`func (o *CreateReplicaBody) HasAllowedHosts() bool`

HasAllowedHosts returns a boolean if a field has been set.

### GetSize

`func (o *CreateReplicaBody) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *CreateReplicaBody) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *CreateReplicaBody) SetSize(v int64)`

SetSize sets Size field to given value.


### GetThin

`func (o *CreateReplicaBody) GetThin() bool`

GetThin returns the Thin field if non-nil, zero value otherwise.

### GetThinOk

`func (o *CreateReplicaBody) GetThinOk() (*bool, bool)`

GetThinOk returns a tuple with the Thin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThin

`func (o *CreateReplicaBody) SetThin(v bool)`

SetThin sets Thin field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


