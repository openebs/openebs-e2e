# VolumePolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SelfHeal** | **bool** | If true the control plane will attempt to heal the volume by itself | 

## Methods

### NewVolumePolicy

`func NewVolumePolicy(selfHeal bool, ) *VolumePolicy`

NewVolumePolicy instantiates a new VolumePolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumePolicyWithDefaults

`func NewVolumePolicyWithDefaults() *VolumePolicy`

NewVolumePolicyWithDefaults instantiates a new VolumePolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSelfHeal

`func (o *VolumePolicy) GetSelfHeal() bool`

GetSelfHeal returns the SelfHeal field if non-nil, zero value otherwise.

### GetSelfHealOk

`func (o *VolumePolicy) GetSelfHealOk() (*bool, bool)`

GetSelfHealOk returns a tuple with the SelfHeal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSelfHeal

`func (o *VolumePolicy) SetSelfHeal(v bool)`

SetSelfHeal sets SelfHeal field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


