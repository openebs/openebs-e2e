# ReplicaSpecOwners

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Nexuses** | **[]string** |  | 
**Volume** | Pointer to **string** |  | [optional] 

## Methods

### NewReplicaSpecOwners

`func NewReplicaSpecOwners(nexuses []string, ) *ReplicaSpecOwners`

NewReplicaSpecOwners instantiates a new ReplicaSpecOwners object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSpecOwnersWithDefaults

`func NewReplicaSpecOwnersWithDefaults() *ReplicaSpecOwners`

NewReplicaSpecOwnersWithDefaults instantiates a new ReplicaSpecOwners object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNexuses

`func (o *ReplicaSpecOwners) GetNexuses() []string`

GetNexuses returns the Nexuses field if non-nil, zero value otherwise.

### GetNexusesOk

`func (o *ReplicaSpecOwners) GetNexusesOk() (*[]string, bool)`

GetNexusesOk returns a tuple with the Nexuses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNexuses

`func (o *ReplicaSpecOwners) SetNexuses(v []string)`

SetNexuses sets Nexuses field to given value.


### GetVolume

`func (o *ReplicaSpecOwners) GetVolume() string`

GetVolume returns the Volume field if non-nil, zero value otherwise.

### GetVolumeOk

`func (o *ReplicaSpecOwners) GetVolumeOk() (*string, bool)`

GetVolumeOk returns a tuple with the Volume field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolume

`func (o *ReplicaSpecOwners) SetVolume(v string)`

SetVolume sets Volume field to given value.

### HasVolume

`func (o *ReplicaSpecOwners) HasVolume() bool`

HasVolume returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


