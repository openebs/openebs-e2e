# BlockDevice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Available** | **bool** | identifies if device is available for use (ie. is not \&quot;currently\&quot; in  use) | 
**ConnectionType** | **string** | the type of bus through which the device is connected to the system | 
**Devlinks** | **[]string** | list of udev generated symlinks by which device may be identified | 
**Devmajor** | **int32** | major device number | 
**Devminor** | **int32** | minor device number | 
**Devname** | **string** | entry in /dev associated with device | 
**Devpath** | **string** | official device path | 
**Devtype** | **string** | currently \&quot;disk\&quot; or \&quot;partition\&quot; | 
**Filesystem** | Pointer to [**BlockDeviceFilesystem**](BlockDeviceFilesystem.md) |  | [optional] 
**IsRotational** | Pointer to **bool** | indicates whether the device is rotational or non-rotational | [optional] 
**Model** | **string** | device model - useful for identifying devices | 
**Partition** | Pointer to [**BlockDevicePartition**](BlockDevicePartition.md) |  | [optional] 
**Size** | **int64** | size of device in (512 byte) blocks | 

## Methods

### NewBlockDevice

`func NewBlockDevice(available bool, connectionType string, devlinks []string, devmajor int32, devminor int32, devname string, devpath string, devtype string, model string, size int64, ) *BlockDevice`

NewBlockDevice instantiates a new BlockDevice object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockDeviceWithDefaults

`func NewBlockDeviceWithDefaults() *BlockDevice`

NewBlockDeviceWithDefaults instantiates a new BlockDevice object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvailable

`func (o *BlockDevice) GetAvailable() bool`

GetAvailable returns the Available field if non-nil, zero value otherwise.

### GetAvailableOk

`func (o *BlockDevice) GetAvailableOk() (*bool, bool)`

GetAvailableOk returns a tuple with the Available field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailable

`func (o *BlockDevice) SetAvailable(v bool)`

SetAvailable sets Available field to given value.


### GetConnectionType

`func (o *BlockDevice) GetConnectionType() string`

GetConnectionType returns the ConnectionType field if non-nil, zero value otherwise.

### GetConnectionTypeOk

`func (o *BlockDevice) GetConnectionTypeOk() (*string, bool)`

GetConnectionTypeOk returns a tuple with the ConnectionType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionType

`func (o *BlockDevice) SetConnectionType(v string)`

SetConnectionType sets ConnectionType field to given value.


### GetDevlinks

`func (o *BlockDevice) GetDevlinks() []string`

GetDevlinks returns the Devlinks field if non-nil, zero value otherwise.

### GetDevlinksOk

`func (o *BlockDevice) GetDevlinksOk() (*[]string, bool)`

GetDevlinksOk returns a tuple with the Devlinks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevlinks

`func (o *BlockDevice) SetDevlinks(v []string)`

SetDevlinks sets Devlinks field to given value.


### GetDevmajor

`func (o *BlockDevice) GetDevmajor() int32`

GetDevmajor returns the Devmajor field if non-nil, zero value otherwise.

### GetDevmajorOk

`func (o *BlockDevice) GetDevmajorOk() (*int32, bool)`

GetDevmajorOk returns a tuple with the Devmajor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevmajor

`func (o *BlockDevice) SetDevmajor(v int32)`

SetDevmajor sets Devmajor field to given value.


### GetDevminor

`func (o *BlockDevice) GetDevminor() int32`

GetDevminor returns the Devminor field if non-nil, zero value otherwise.

### GetDevminorOk

`func (o *BlockDevice) GetDevminorOk() (*int32, bool)`

GetDevminorOk returns a tuple with the Devminor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevminor

`func (o *BlockDevice) SetDevminor(v int32)`

SetDevminor sets Devminor field to given value.


### GetDevname

`func (o *BlockDevice) GetDevname() string`

GetDevname returns the Devname field if non-nil, zero value otherwise.

### GetDevnameOk

`func (o *BlockDevice) GetDevnameOk() (*string, bool)`

GetDevnameOk returns a tuple with the Devname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevname

`func (o *BlockDevice) SetDevname(v string)`

SetDevname sets Devname field to given value.


### GetDevpath

`func (o *BlockDevice) GetDevpath() string`

GetDevpath returns the Devpath field if non-nil, zero value otherwise.

### GetDevpathOk

`func (o *BlockDevice) GetDevpathOk() (*string, bool)`

GetDevpathOk returns a tuple with the Devpath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevpath

`func (o *BlockDevice) SetDevpath(v string)`

SetDevpath sets Devpath field to given value.


### GetDevtype

`func (o *BlockDevice) GetDevtype() string`

GetDevtype returns the Devtype field if non-nil, zero value otherwise.

### GetDevtypeOk

`func (o *BlockDevice) GetDevtypeOk() (*string, bool)`

GetDevtypeOk returns a tuple with the Devtype field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevtype

`func (o *BlockDevice) SetDevtype(v string)`

SetDevtype sets Devtype field to given value.


### GetFilesystem

`func (o *BlockDevice) GetFilesystem() BlockDeviceFilesystem`

GetFilesystem returns the Filesystem field if non-nil, zero value otherwise.

### GetFilesystemOk

`func (o *BlockDevice) GetFilesystemOk() (*BlockDeviceFilesystem, bool)`

GetFilesystemOk returns a tuple with the Filesystem field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilesystem

`func (o *BlockDevice) SetFilesystem(v BlockDeviceFilesystem)`

SetFilesystem sets Filesystem field to given value.

### HasFilesystem

`func (o *BlockDevice) HasFilesystem() bool`

HasFilesystem returns a boolean if a field has been set.

### GetIsRotational

`func (o *BlockDevice) GetIsRotational() bool`

GetIsRotational returns the IsRotational field if non-nil, zero value otherwise.

### GetIsRotationalOk

`func (o *BlockDevice) GetIsRotationalOk() (*bool, bool)`

GetIsRotationalOk returns a tuple with the IsRotational field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsRotational

`func (o *BlockDevice) SetIsRotational(v bool)`

SetIsRotational sets IsRotational field to given value.

### HasIsRotational

`func (o *BlockDevice) HasIsRotational() bool`

HasIsRotational returns a boolean if a field has been set.

### GetModel

`func (o *BlockDevice) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *BlockDevice) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *BlockDevice) SetModel(v string)`

SetModel sets Model field to given value.


### GetPartition

`func (o *BlockDevice) GetPartition() BlockDevicePartition`

GetPartition returns the Partition field if non-nil, zero value otherwise.

### GetPartitionOk

`func (o *BlockDevice) GetPartitionOk() (*BlockDevicePartition, bool)`

GetPartitionOk returns a tuple with the Partition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartition

`func (o *BlockDevice) SetPartition(v BlockDevicePartition)`

SetPartition sets Partition field to given value.

### HasPartition

`func (o *BlockDevice) HasPartition() bool`

HasPartition returns a boolean if a field has been set.

### GetSize

`func (o *BlockDevice) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *BlockDevice) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *BlockDevice) SetSize(v int64)`

SetSize sets Size field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


