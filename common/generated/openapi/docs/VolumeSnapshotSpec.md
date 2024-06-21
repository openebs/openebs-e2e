# VolumeSnapshotSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uuid** | **string** |  | 
**SourceVolume** | **string** |  | 

## Methods

### NewVolumeSnapshotSpec

`func NewVolumeSnapshotSpec(uuid string, sourceVolume string, ) *VolumeSnapshotSpec`

NewVolumeSnapshotSpec instantiates a new VolumeSnapshotSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeSnapshotSpecWithDefaults

`func NewVolumeSnapshotSpecWithDefaults() *VolumeSnapshotSpec`

NewVolumeSnapshotSpecWithDefaults instantiates a new VolumeSnapshotSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUuid

`func (o *VolumeSnapshotSpec) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *VolumeSnapshotSpec) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *VolumeSnapshotSpec) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetSourceVolume

`func (o *VolumeSnapshotSpec) GetSourceVolume() string`

GetSourceVolume returns the SourceVolume field if non-nil, zero value otherwise.

### GetSourceVolumeOk

`func (o *VolumeSnapshotSpec) GetSourceVolumeOk() (*string, bool)`

GetSourceVolumeOk returns a tuple with the SourceVolume field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceVolume

`func (o *VolumeSnapshotSpec) SetSourceVolume(v string)`

SetSourceVolume sets SourceVolume field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


