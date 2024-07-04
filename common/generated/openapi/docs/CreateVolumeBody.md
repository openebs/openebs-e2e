# CreateVolumeBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Policy** | [**VolumePolicy**](VolumePolicy.md) |  | 
**Replicas** | **int32** | number of storage replicas | 
**Size** | **int64** | size of the volume in bytes | 
**Thin** | **bool** | flag indicating whether or not the volume is thin provisioned | 
**Topology** | Pointer to [**Topology**](Topology.md) |  | [optional] 
**Labels** | Pointer to **map[string]string** | Optionally used to store custom volume information | [optional] 
**AffinityGroup** | Pointer to [**AffinityGroup**](AffinityGroup.md) | Affinity Group related information. | [optional] 
**MaxSnapshots** | Pointer to **int32** | Max Snapshots limit per volume. | [optional] 

## Methods

### NewCreateVolumeBody

`func NewCreateVolumeBody(policy VolumePolicy, replicas int32, size int64, thin bool, ) *CreateVolumeBody`

NewCreateVolumeBody instantiates a new CreateVolumeBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateVolumeBodyWithDefaults

`func NewCreateVolumeBodyWithDefaults() *CreateVolumeBody`

NewCreateVolumeBodyWithDefaults instantiates a new CreateVolumeBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPolicy

`func (o *CreateVolumeBody) GetPolicy() VolumePolicy`

GetPolicy returns the Policy field if non-nil, zero value otherwise.

### GetPolicyOk

`func (o *CreateVolumeBody) GetPolicyOk() (*VolumePolicy, bool)`

GetPolicyOk returns a tuple with the Policy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicy

`func (o *CreateVolumeBody) SetPolicy(v VolumePolicy)`

SetPolicy sets Policy field to given value.


### GetReplicas

`func (o *CreateVolumeBody) GetReplicas() int32`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *CreateVolumeBody) GetReplicasOk() (*int32, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *CreateVolumeBody) SetReplicas(v int32)`

SetReplicas sets Replicas field to given value.


### GetSize

`func (o *CreateVolumeBody) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *CreateVolumeBody) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *CreateVolumeBody) SetSize(v int64)`

SetSize sets Size field to given value.


### GetThin

`func (o *CreateVolumeBody) GetThin() bool`

GetThin returns the Thin field if non-nil, zero value otherwise.

### GetThinOk

`func (o *CreateVolumeBody) GetThinOk() (*bool, bool)`

GetThinOk returns a tuple with the Thin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThin

`func (o *CreateVolumeBody) SetThin(v bool)`

SetThin sets Thin field to given value.


### GetTopology

`func (o *CreateVolumeBody) GetTopology() Topology`

GetTopology returns the Topology field if non-nil, zero value otherwise.

### GetTopologyOk

`func (o *CreateVolumeBody) GetTopologyOk() (*Topology, bool)`

GetTopologyOk returns a tuple with the Topology field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopology

`func (o *CreateVolumeBody) SetTopology(v Topology)`

SetTopology sets Topology field to given value.

### HasTopology

`func (o *CreateVolumeBody) HasTopology() bool`

HasTopology returns a boolean if a field has been set.

### GetLabels

`func (o *CreateVolumeBody) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CreateVolumeBody) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CreateVolumeBody) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CreateVolumeBody) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetAffinityGroup

`func (o *CreateVolumeBody) GetAffinityGroup() AffinityGroup`

GetAffinityGroup returns the AffinityGroup field if non-nil, zero value otherwise.

### GetAffinityGroupOk

`func (o *CreateVolumeBody) GetAffinityGroupOk() (*AffinityGroup, bool)`

GetAffinityGroupOk returns a tuple with the AffinityGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffinityGroup

`func (o *CreateVolumeBody) SetAffinityGroup(v AffinityGroup)`

SetAffinityGroup sets AffinityGroup field to given value.

### HasAffinityGroup

`func (o *CreateVolumeBody) HasAffinityGroup() bool`

HasAffinityGroup returns a boolean if a field has been set.

### GetMaxSnapshots

`func (o *CreateVolumeBody) GetMaxSnapshots() int32`

GetMaxSnapshots returns the MaxSnapshots field if non-nil, zero value otherwise.

### GetMaxSnapshotsOk

`func (o *CreateVolumeBody) GetMaxSnapshotsOk() (*int32, bool)`

GetMaxSnapshotsOk returns a tuple with the MaxSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSnapshots

`func (o *CreateVolumeBody) SetMaxSnapshots(v int32)`

SetMaxSnapshots sets MaxSnapshots field to given value.

### HasMaxSnapshots

`func (o *CreateVolumeBody) HasMaxSnapshots() bool`

HasMaxSnapshots returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


