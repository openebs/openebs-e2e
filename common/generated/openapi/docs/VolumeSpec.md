# VolumeSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Labels** | Pointer to **map[string]string** | Optionally used to store custom volume information | [optional] 
**NumReplicas** | **int32** | Number of children the volume should have. | 
**Operation** | Pointer to [**VolumeSpecOperation**](VolumeSpecOperation.md) |  | [optional] 
**Size** | **int64** | Size that the volume should be. | 
**Status** | [**SpecStatus**](SpecStatus.md) |  | 
**Target** | Pointer to [**VolumeTarget**](VolumeTarget.md) |  | [optional] 
**Uuid** | **string** | Volume Id | 
**Topology** | Pointer to [**Topology**](Topology.md) |  | [optional] 
**Policy** | [**VolumePolicy**](VolumePolicy.md) |  | 
**Thin** | **bool** | Thin provisioning flag. | 
**AsThin** | Pointer to **bool** | Volume converted to thin provisioned. | [optional] 
**AffinityGroup** | Pointer to [**AffinityGroup**](AffinityGroup.md) |  | [optional] 
**ContentSource** | Pointer to [**NullableVolumeContentSource**](VolumeContentSource.md) |  | [optional] 
**NumSnapshots** | **int32** | Number of snapshots taken on this volume. | 
**MaxSnapshots** | Pointer to **int32** | Max snapshots to limit per volume. | [optional] 

## Methods

### NewVolumeSpec

`func NewVolumeSpec(numReplicas int32, size int64, status SpecStatus, uuid string, policy VolumePolicy, thin bool, numSnapshots int32, ) *VolumeSpec`

NewVolumeSpec instantiates a new VolumeSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeSpecWithDefaults

`func NewVolumeSpecWithDefaults() *VolumeSpec`

NewVolumeSpecWithDefaults instantiates a new VolumeSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabels

`func (o *VolumeSpec) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *VolumeSpec) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *VolumeSpec) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *VolumeSpec) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetNumReplicas

`func (o *VolumeSpec) GetNumReplicas() int32`

GetNumReplicas returns the NumReplicas field if non-nil, zero value otherwise.

### GetNumReplicasOk

`func (o *VolumeSpec) GetNumReplicasOk() (*int32, bool)`

GetNumReplicasOk returns a tuple with the NumReplicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumReplicas

`func (o *VolumeSpec) SetNumReplicas(v int32)`

SetNumReplicas sets NumReplicas field to given value.


### GetOperation

`func (o *VolumeSpec) GetOperation() VolumeSpecOperation`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *VolumeSpec) GetOperationOk() (*VolumeSpecOperation, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *VolumeSpec) SetOperation(v VolumeSpecOperation)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *VolumeSpec) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetSize

`func (o *VolumeSpec) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *VolumeSpec) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *VolumeSpec) SetSize(v int64)`

SetSize sets Size field to given value.


### GetStatus

`func (o *VolumeSpec) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *VolumeSpec) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *VolumeSpec) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.


### GetTarget

`func (o *VolumeSpec) GetTarget() VolumeTarget`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *VolumeSpec) GetTargetOk() (*VolumeTarget, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *VolumeSpec) SetTarget(v VolumeTarget)`

SetTarget sets Target field to given value.

### HasTarget

`func (o *VolumeSpec) HasTarget() bool`

HasTarget returns a boolean if a field has been set.

### GetUuid

`func (o *VolumeSpec) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *VolumeSpec) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *VolumeSpec) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetTopology

`func (o *VolumeSpec) GetTopology() Topology`

GetTopology returns the Topology field if non-nil, zero value otherwise.

### GetTopologyOk

`func (o *VolumeSpec) GetTopologyOk() (*Topology, bool)`

GetTopologyOk returns a tuple with the Topology field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopology

`func (o *VolumeSpec) SetTopology(v Topology)`

SetTopology sets Topology field to given value.

### HasTopology

`func (o *VolumeSpec) HasTopology() bool`

HasTopology returns a boolean if a field has been set.

### GetPolicy

`func (o *VolumeSpec) GetPolicy() VolumePolicy`

GetPolicy returns the Policy field if non-nil, zero value otherwise.

### GetPolicyOk

`func (o *VolumeSpec) GetPolicyOk() (*VolumePolicy, bool)`

GetPolicyOk returns a tuple with the Policy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicy

`func (o *VolumeSpec) SetPolicy(v VolumePolicy)`

SetPolicy sets Policy field to given value.


### GetThin

`func (o *VolumeSpec) GetThin() bool`

GetThin returns the Thin field if non-nil, zero value otherwise.

### GetThinOk

`func (o *VolumeSpec) GetThinOk() (*bool, bool)`

GetThinOk returns a tuple with the Thin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThin

`func (o *VolumeSpec) SetThin(v bool)`

SetThin sets Thin field to given value.


### GetAsThin

`func (o *VolumeSpec) GetAsThin() bool`

GetAsThin returns the AsThin field if non-nil, zero value otherwise.

### GetAsThinOk

`func (o *VolumeSpec) GetAsThinOk() (*bool, bool)`

GetAsThinOk returns a tuple with the AsThin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsThin

`func (o *VolumeSpec) SetAsThin(v bool)`

SetAsThin sets AsThin field to given value.

### HasAsThin

`func (o *VolumeSpec) HasAsThin() bool`

HasAsThin returns a boolean if a field has been set.

### GetAffinityGroup

`func (o *VolumeSpec) GetAffinityGroup() AffinityGroup`

GetAffinityGroup returns the AffinityGroup field if non-nil, zero value otherwise.

### GetAffinityGroupOk

`func (o *VolumeSpec) GetAffinityGroupOk() (*AffinityGroup, bool)`

GetAffinityGroupOk returns a tuple with the AffinityGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAffinityGroup

`func (o *VolumeSpec) SetAffinityGroup(v AffinityGroup)`

SetAffinityGroup sets AffinityGroup field to given value.

### HasAffinityGroup

`func (o *VolumeSpec) HasAffinityGroup() bool`

HasAffinityGroup returns a boolean if a field has been set.

### GetContentSource

`func (o *VolumeSpec) GetContentSource() VolumeContentSource`

GetContentSource returns the ContentSource field if non-nil, zero value otherwise.

### GetContentSourceOk

`func (o *VolumeSpec) GetContentSourceOk() (*VolumeContentSource, bool)`

GetContentSourceOk returns a tuple with the ContentSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContentSource

`func (o *VolumeSpec) SetContentSource(v VolumeContentSource)`

SetContentSource sets ContentSource field to given value.

### HasContentSource

`func (o *VolumeSpec) HasContentSource() bool`

HasContentSource returns a boolean if a field has been set.

### SetContentSourceNil

`func (o *VolumeSpec) SetContentSourceNil(b bool)`

 SetContentSourceNil sets the value for ContentSource to be an explicit nil

### UnsetContentSource
`func (o *VolumeSpec) UnsetContentSource()`

UnsetContentSource ensures that no value is present for ContentSource, not even an explicit nil
### GetNumSnapshots

`func (o *VolumeSpec) GetNumSnapshots() int32`

GetNumSnapshots returns the NumSnapshots field if non-nil, zero value otherwise.

### GetNumSnapshotsOk

`func (o *VolumeSpec) GetNumSnapshotsOk() (*int32, bool)`

GetNumSnapshotsOk returns a tuple with the NumSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumSnapshots

`func (o *VolumeSpec) SetNumSnapshots(v int32)`

SetNumSnapshots sets NumSnapshots field to given value.


### GetMaxSnapshots

`func (o *VolumeSpec) GetMaxSnapshots() int32`

GetMaxSnapshots returns the MaxSnapshots field if non-nil, zero value otherwise.

### GetMaxSnapshotsOk

`func (o *VolumeSpec) GetMaxSnapshotsOk() (*int32, bool)`

GetMaxSnapshotsOk returns a tuple with the MaxSnapshots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSnapshots

`func (o *VolumeSpec) SetMaxSnapshots(v int32)`

SetMaxSnapshots sets MaxSnapshots field to given value.

### HasMaxSnapshots

`func (o *VolumeSpec) HasMaxSnapshots() bool`

HasMaxSnapshots returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


