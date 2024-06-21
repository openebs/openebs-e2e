# Specs

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Nexuses** | [**[]NexusSpec**](NexusSpec.md) | Nexus Specs | 
**Pools** | [**[]PoolSpec**](PoolSpec.md) | Pool Specs | 
**Replicas** | [**[]ReplicaSpec**](ReplicaSpec.md) | Replica Specs | 
**Volumes** | [**[]VolumeSpec**](VolumeSpec.md) | Volume Specs | 

## Methods

### NewSpecs

`func NewSpecs(nexuses []NexusSpec, pools []PoolSpec, replicas []ReplicaSpec, volumes []VolumeSpec, ) *Specs`

NewSpecs instantiates a new Specs object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSpecsWithDefaults

`func NewSpecsWithDefaults() *Specs`

NewSpecsWithDefaults instantiates a new Specs object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNexuses

`func (o *Specs) GetNexuses() []NexusSpec`

GetNexuses returns the Nexuses field if non-nil, zero value otherwise.

### GetNexusesOk

`func (o *Specs) GetNexusesOk() (*[]NexusSpec, bool)`

GetNexusesOk returns a tuple with the Nexuses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNexuses

`func (o *Specs) SetNexuses(v []NexusSpec)`

SetNexuses sets Nexuses field to given value.


### GetPools

`func (o *Specs) GetPools() []PoolSpec`

GetPools returns the Pools field if non-nil, zero value otherwise.

### GetPoolsOk

`func (o *Specs) GetPoolsOk() (*[]PoolSpec, bool)`

GetPoolsOk returns a tuple with the Pools field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPools

`func (o *Specs) SetPools(v []PoolSpec)`

SetPools sets Pools field to given value.


### GetReplicas

`func (o *Specs) GetReplicas() []ReplicaSpec`

GetReplicas returns the Replicas field if non-nil, zero value otherwise.

### GetReplicasOk

`func (o *Specs) GetReplicasOk() (*[]ReplicaSpec, bool)`

GetReplicasOk returns a tuple with the Replicas field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicas

`func (o *Specs) SetReplicas(v []ReplicaSpec)`

SetReplicas sets Replicas field to given value.


### GetVolumes

`func (o *Specs) GetVolumes() []VolumeSpec`

GetVolumes returns the Volumes field if non-nil, zero value otherwise.

### GetVolumesOk

`func (o *Specs) GetVolumesOk() (*[]VolumeSpec, bool)`

GetVolumesOk returns a tuple with the Volumes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumes

`func (o *Specs) SetVolumes(v []VolumeSpec)`

SetVolumes sets Volumes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


