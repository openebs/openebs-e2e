# BlockDeviceFilesystem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Fstype** | **string** | filesystem type: ext3, ntfs, ... | 
**Label** | **string** | volume label | 
**Mountpoint** | **string** | path where filesystem is currently mounted | 
**Uuid** | **string** | UUID identifying the volume (filesystem) | 

## Methods

### NewBlockDeviceFilesystem

`func NewBlockDeviceFilesystem(fstype string, label string, mountpoint string, uuid string, ) *BlockDeviceFilesystem`

NewBlockDeviceFilesystem instantiates a new BlockDeviceFilesystem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockDeviceFilesystemWithDefaults

`func NewBlockDeviceFilesystemWithDefaults() *BlockDeviceFilesystem`

NewBlockDeviceFilesystemWithDefaults instantiates a new BlockDeviceFilesystem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFstype

`func (o *BlockDeviceFilesystem) GetFstype() string`

GetFstype returns the Fstype field if non-nil, zero value otherwise.

### GetFstypeOk

`func (o *BlockDeviceFilesystem) GetFstypeOk() (*string, bool)`

GetFstypeOk returns a tuple with the Fstype field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFstype

`func (o *BlockDeviceFilesystem) SetFstype(v string)`

SetFstype sets Fstype field to given value.


### GetLabel

`func (o *BlockDeviceFilesystem) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *BlockDeviceFilesystem) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *BlockDeviceFilesystem) SetLabel(v string)`

SetLabel sets Label field to given value.


### GetMountpoint

`func (o *BlockDeviceFilesystem) GetMountpoint() string`

GetMountpoint returns the Mountpoint field if non-nil, zero value otherwise.

### GetMountpointOk

`func (o *BlockDeviceFilesystem) GetMountpointOk() (*string, bool)`

GetMountpointOk returns a tuple with the Mountpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMountpoint

`func (o *BlockDeviceFilesystem) SetMountpoint(v string)`

SetMountpoint sets Mountpoint field to given value.


### GetUuid

`func (o *BlockDeviceFilesystem) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *BlockDeviceFilesystem) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *BlockDeviceFilesystem) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


