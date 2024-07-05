# Volume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Spec** | [**VolumeSpec**](VolumeSpec.md) |  | 
**State** | [**VolumeState**](VolumeState.md) |  | 

## Methods

### NewVolume

`func NewVolume(spec VolumeSpec, state VolumeState, ) *Volume`

NewVolume instantiates a new Volume object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeWithDefaults

`func NewVolumeWithDefaults() *Volume`

NewVolumeWithDefaults instantiates a new Volume object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSpec

`func (o *Volume) GetSpec() VolumeSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *Volume) GetSpecOk() (*VolumeSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *Volume) SetSpec(v VolumeSpec)`

SetSpec sets Spec field to given value.


### GetState

`func (o *Volume) GetState() VolumeState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Volume) GetStateOk() (*VolumeState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Volume) SetState(v VolumeState)`

SetState sets State field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


