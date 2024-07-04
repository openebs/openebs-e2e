# NexusSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Children** | **[]string** | List of children. | 
**Managed** | **bool** | Managed by our control plane | 
**Node** | **string** | Node where the nexus should live. | 
**Operation** | Pointer to [**NexusSpecOperation**](NexusSpecOperation.md) |  | [optional] 
**Owner** | Pointer to **string** | Volume which owns this nexus, if any | [optional] 
**Share** | [**Protocol**](Protocol.md) |  | 
**Size** | **int64** | Size of the nexus. | 
**Status** | [**SpecStatus**](SpecStatus.md) |  | 
**Uuid** | **string** | Nexus Id | 

## Methods

### NewNexusSpec

`func NewNexusSpec(children []string, managed bool, node string, share Protocol, size int64, status SpecStatus, uuid string, ) *NexusSpec`

NewNexusSpec instantiates a new NexusSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNexusSpecWithDefaults

`func NewNexusSpecWithDefaults() *NexusSpec`

NewNexusSpecWithDefaults instantiates a new NexusSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChildren

`func (o *NexusSpec) GetChildren() []string`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *NexusSpec) GetChildrenOk() (*[]string, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *NexusSpec) SetChildren(v []string)`

SetChildren sets Children field to given value.


### GetManaged

`func (o *NexusSpec) GetManaged() bool`

GetManaged returns the Managed field if non-nil, zero value otherwise.

### GetManagedOk

`func (o *NexusSpec) GetManagedOk() (*bool, bool)`

GetManagedOk returns a tuple with the Managed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManaged

`func (o *NexusSpec) SetManaged(v bool)`

SetManaged sets Managed field to given value.


### GetNode

`func (o *NexusSpec) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *NexusSpec) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *NexusSpec) SetNode(v string)`

SetNode sets Node field to given value.


### GetOperation

`func (o *NexusSpec) GetOperation() NexusSpecOperation`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *NexusSpec) GetOperationOk() (*NexusSpecOperation, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *NexusSpec) SetOperation(v NexusSpecOperation)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *NexusSpec) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetOwner

`func (o *NexusSpec) GetOwner() string`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *NexusSpec) GetOwnerOk() (*string, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *NexusSpec) SetOwner(v string)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *NexusSpec) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetShare

`func (o *NexusSpec) GetShare() Protocol`

GetShare returns the Share field if non-nil, zero value otherwise.

### GetShareOk

`func (o *NexusSpec) GetShareOk() (*Protocol, bool)`

GetShareOk returns a tuple with the Share field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShare

`func (o *NexusSpec) SetShare(v Protocol)`

SetShare sets Share field to given value.


### GetSize

`func (o *NexusSpec) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *NexusSpec) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *NexusSpec) SetSize(v int64)`

SetSize sets Size field to given value.


### GetStatus

`func (o *NexusSpec) GetStatus() SpecStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NexusSpec) GetStatusOk() (*SpecStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NexusSpec) SetStatus(v SpecStatus)`

SetStatus sets Status field to given value.


### GetUuid

`func (o *NexusSpec) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *NexusSpec) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *NexusSpec) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


