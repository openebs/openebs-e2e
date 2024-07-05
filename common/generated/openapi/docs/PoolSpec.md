# PoolSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Disks** | **[]string** | absolute disk paths claimed by the pool | 
**Id** | **string** | storage pool identifier | 
**Labels** | Pointer to **map[string]string** | labels to be set on the pools | [optional] 
**Node** | **string** | storage node identifier | 
**Status** | [**SpecStatus**](SpecStatus.md) |  | 

## Methods

### NewPoolSpec

`func NewPoolSpec(disks []string, id string, node string, status SpecStatus, ) *PoolSpec`

NewPoolSpec instantiates a new PoolSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolSpecWithDefaults

`func NewPoolSpecWithDefaults() *PoolSpec`

NewPoolSpecWithDefaults instantiates a new PoolSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDisks

`func (o *PoolSpec) GetDisks() []string`

GetDisks returns the Disks field if non-nil, zero value otherwise.

### GetDisksOk

`func (o *PoolSpec) GetDisksOk() (*[]string, bool)`

GetDisksOk returns a tuple with the Disks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisks

`func (o *PoolSpec) SetDisks(v []string)`

SetDisks sets Disks field to given value.


### GetId

`func (o *PoolSpec) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PoolSpec) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PoolSpec) SetId(v string)`

SetId sets Id field to given value.


### GetLabels

`func (o *PoolSpec) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *PoolSpec) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *PoolSpec) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *PoolSpec) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetNode

`func (o *PoolSpec) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *PoolSpec) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *PoolSpec) SetNode(v string)`

SetNode sets Node field to given value.


### GetStatus

`func (o *PoolSpec) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PoolSpec) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PoolSpec) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


