# ReplicaSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Managed** | **bool** | Managed by our control plane | 
**Operation** | Pointer to [**ReplicaSpecOperation**](ReplicaSpecOperation.md) |  | [optional] 
**Owners** | [**ReplicaSpecOwners**](ReplicaSpecOwners.md) |  | 
**Pool** | **string** | The pool that the replica should live on. | 
**PoolUuid** | Pointer to **string** | storage pool unique identifier | [optional] 
**Share** | [**Protocol**](Protocol.md) |  | 
**Size** | **int64** | The size that the replica should be. | 
**Status** | [**SpecStatus**](SpecStatus.md) |  | 
**Thin** | **bool** | Thin provisioning. | 
**Uuid** | **string** | uuid of the replica | 
**Kind** | Pointer to [**ReplicaKind**](ReplicaKind.md) |  | [optional] [default to REGULAR]

## Methods

### NewReplicaSpec

`func NewReplicaSpec(managed bool, owners ReplicaSpecOwners, pool string, share Protocol, size int64, status SpecStatus, thin bool, uuid string, ) *ReplicaSpec`

NewReplicaSpec instantiates a new ReplicaSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaSpecWithDefaults

`func NewReplicaSpecWithDefaults() *ReplicaSpec`

NewReplicaSpecWithDefaults instantiates a new ReplicaSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetManaged

`func (o *ReplicaSpec) GetManaged() bool`

GetManaged returns the Managed field if non-nil, zero value otherwise.

### GetManagedOk

`func (o *ReplicaSpec) GetManagedOk() (*bool, bool)`

GetManagedOk returns a tuple with the Managed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManaged

`func (o *ReplicaSpec) SetManaged(v bool)`

SetManaged sets Managed field to given value.


### GetOperation

`func (o *ReplicaSpec) GetOperation() ReplicaSpecOperation`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *ReplicaSpec) GetOperationOk() (*ReplicaSpecOperation, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *ReplicaSpec) SetOperation(v ReplicaSpecOperation)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *ReplicaSpec) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetOwners

`func (o *ReplicaSpec) GetOwners() ReplicaSpecOwners`

GetOwners returns the Owners field if non-nil, zero value otherwise.

### GetOwnersOk

`func (o *ReplicaSpec) GetOwnersOk() (*ReplicaSpecOwners, bool)`

GetOwnersOk returns a tuple with the Owners field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwners

`func (o *ReplicaSpec) SetOwners(v ReplicaSpecOwners)`

SetOwners sets Owners field to given value.


### GetPool

`func (o *ReplicaSpec) GetPool() string`

GetPool returns the Pool field if non-nil, zero value otherwise.

### GetPoolOk

`func (o *ReplicaSpec) GetPoolOk() (*string, bool)`

GetPoolOk returns a tuple with the Pool field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPool

`func (o *ReplicaSpec) SetPool(v string)`

SetPool sets Pool field to given value.


### GetPoolUuid

`func (o *ReplicaSpec) GetPoolUuid() string`

GetPoolUuid returns the PoolUuid field if non-nil, zero value otherwise.

### GetPoolUuidOk

`func (o *ReplicaSpec) GetPoolUuidOk() (*string, bool)`

GetPoolUuidOk returns a tuple with the PoolUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolUuid

`func (o *ReplicaSpec) SetPoolUuid(v string)`

SetPoolUuid sets PoolUuid field to given value.

### HasPoolUuid

`func (o *ReplicaSpec) HasPoolUuid() bool`

HasPoolUuid returns a boolean if a field has been set.

### GetShare

`func (o *ReplicaSpec) GetShare() Protocol`

GetShare returns the Share field if non-nil, zero value otherwise.

### GetShareOk

`func (o *ReplicaSpec) GetShareOk() (*Protocol, bool)`

GetShareOk returns a tuple with the Share field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShare

`func (o *ReplicaSpec) SetShare(v Protocol)`

SetShare sets Share field to given value.


### GetSize

`func (o *ReplicaSpec) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *ReplicaSpec) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *ReplicaSpec) SetSize(v int64)`

SetSize sets Size field to given value.


### GetStatus

`func (o *ReplicaSpec) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ReplicaSpec) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ReplicaSpec) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.


### GetThin

`func (o *ReplicaSpec) GetThin() bool`

GetThin returns the Thin field if non-nil, zero value otherwise.

### GetThinOk

`func (o *ReplicaSpec) GetThinOk() (*bool, bool)`

GetThinOk returns a tuple with the Thin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThin

`func (o *ReplicaSpec) SetThin(v bool)`

SetThin sets Thin field to given value.


### GetUuid

`func (o *ReplicaSpec) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *ReplicaSpec) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *ReplicaSpec) SetUuid(v string)`

SetUuid sets Uuid field to given value.


### GetKind

`func (o *ReplicaSpec) GetKind() ReplicaKind`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *ReplicaSpec) GetKindOk() (*ReplicaKind, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *ReplicaSpec) SetKind(v ReplicaKind)`

SetKind sets Kind field to given value.

### HasKind

`func (o *ReplicaSpec) HasKind() bool`

HasKind returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


