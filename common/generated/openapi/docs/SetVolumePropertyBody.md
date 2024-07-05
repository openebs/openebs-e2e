# SetVolumePropertyBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaxSnapshots** | Pointer to **int32** | Max Snapshots limit per volume. | [optional] 

## Methods

### NewSetVolumePropertyBody

`func NewSetVolumePropertyBody() *SetVolumePropertyBody`

NewSetVolumePropertyBody instantiates a new SetVolumePropertyBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSetVolumePropertyBodyWithDefaults

`func NewSetVolumePropertyBodyWithDefaults() *SetVolumePropertyBody`

NewSetVolumePropertyBodyWithDefaults instantiates a new SetVolumePropertyBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMaxSnapshots

`func (o *SetVolumePropertyBody) GetMaxSnapshots() int32`

GetMaxSnapshots returns the MaxSnapshots field if non-nil, zero value otherwise.

### GetMaxSnapshotsOk

`func (o *SetVolumePropertyBody) GetMaxSnapshotsOk() (*int32, bool)`

GetMaxSnapshotsOk returns a tuple with the MaxSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSnapshots

`func (o *SetVolumePropertyBody) SetMaxSnapshots(v int32)`

SetMaxSnapshots sets MaxSnapshots field to given value.

### HasMaxSnapshots

`func (o *SetVolumePropertyBody) HasMaxSnapshots() bool`

HasMaxSnapshots returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


