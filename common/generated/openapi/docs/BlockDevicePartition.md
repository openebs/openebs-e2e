# BlockDevicePartition

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | partition name | 
**Number** | **int32** | partition number | 
**Parent** | **string** | devname of parent device to which this partition belongs | 
**Scheme** | **string** | partition scheme: gpt, dos, ... | 
**Typeid** | **string** | partition type identifier | 
**Uuid** | **string** | UUID identifying partition | 

## Methods

### NewBlockDevicePartition

`func NewBlockDevicePartition(name string, number int32, parent string, scheme string, typeid string, uuid string, ) *BlockDevicePartition`

NewBlockDevicePartition instantiates a new BlockDevicePartition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockDevicePartitionWithDefaults

`func NewBlockDevicePartitionWithDefaults() *BlockDevicePartition`

NewBlockDevicePartitionWithDefaults instantiates a new BlockDevicePartition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BlockDevicePartition) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BlockDevicePartition) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BlockDevicePartition) SetName(v string)`

SetName sets Name field to given value.


### GetNumber

`func (o *BlockDevicePartition) GetNumber() int32`

GetNumber returns the Number field if non-nil, zero value otherwise.

### GetNumberOk

`func (o *BlockDevicePartition) GetNumberOk() (*int32, bool)`

GetNumberOk returns a tuple with the Number field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumber

`func (o *BlockDevicePartition) SetNumber(v int32)`

SetNumber sets Number field to given value.


### GetParent

`func (o *BlockDevicePartition) GetParent() string`

GetParent returns the Parent field if non-nil, zero value otherwise.

### GetParentOk

`func (o *BlockDevicePartition) GetParentOk() (*string, bool)`

GetParentOk returns a tuple with the Parent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParent

`func (o *BlockDevicePartition) SetParent(v string)`

SetParent sets Parent field to given value.


### GetScheme

`func (o *BlockDevicePartition) GetScheme() string`

GetScheme returns the Scheme field if non-nil, zero value otherwise.

### GetSchemeOk

`func (o *BlockDevicePartition) GetSchemeOk() (*string, bool)`

GetSchemeOk returns a tuple with the Scheme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheme

`func (o *BlockDevicePartition) SetScheme(v string)`

SetScheme sets Scheme field to given value.


### GetTypeid

`func (o *BlockDevicePartition) GetTypeid() string`

GetTypeid returns the Typeid field if non-nil, zero value otherwise.

### GetTypeidOk

`func (o *BlockDevicePartition) GetTypeidOk() (*string, bool)`

GetTypeidOk returns a tuple with the Typeid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypeid

`func (o *BlockDevicePartition) SetTypeid(v string)`

SetTypeid sets Typeid field to given value.


### GetUuid

`func (o *BlockDevicePartition) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *BlockDevicePartition) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *BlockDevicePartition) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


