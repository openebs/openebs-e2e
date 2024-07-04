# VolumeSnapshot

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Definition** | [**VolumeSnapshotDefinition**](VolumeSnapshotDefinition.md) |  | 
**State** | [**VolumeSnapshotState**](VolumeSnapshotState.md) |  | 

## Methods

### NewVolumeSnapshot

`func NewVolumeSnapshot(definition VolumeSnapshotDefinition, state VolumeSnapshotState, ) *VolumeSnapshot`

NewVolumeSnapshot instantiates a new VolumeSnapshot object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeSnapshotWithDefaults

`func NewVolumeSnapshotWithDefaults() *VolumeSnapshot`

NewVolumeSnapshotWithDefaults instantiates a new VolumeSnapshot object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefinition

`func (o *VolumeSnapshot) GetDefinition() VolumeSnapshotDefinition`

GetDefinition returns the Definition field if non-nil, zero value otherwise.

### GetDefinitionOk

`func (o *VolumeSnapshot) GetDefinitionOk() (*VolumeSnapshotDefinition, bool)`

GetDefinitionOk returns a tuple with the Definition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefinition

`func (o *VolumeSnapshot) SetDefinition(v VolumeSnapshotDefinition)`

SetDefinition sets Definition field to given value.


### GetState

`func (o *VolumeSnapshot) GetState() VolumeSnapshotState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *VolumeSnapshot) GetStateOk() (*VolumeSnapshotState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *VolumeSnapshot) SetState(v VolumeSnapshotState)`

SetState sets State field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


