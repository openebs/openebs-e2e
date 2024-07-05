# Nexus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Children** | [**[]Child**](Child.md) | Array of Nexus Children | 
**DeviceUri** | **string** | URI of the device for the volume (missing if not published).  Missing property and empty string are treated the same. | 
**Node** | **string** | id of the io-engine instance | 
**Rebuilds** | **int32** | total number of rebuild tasks | 
**Protocol** | [**Protocol**](Protocol.md) |  | 
**Size** | **int64** | size of the volume in bytes | 
**State** | [**NexusState**](NexusState.md) |  | 
**Uuid** | **string** | uuid of the nexus | 

## Methods

### NewNexus

`func NewNexus(children []Child, deviceUri string, node string, rebuilds int32, protocol Protocol, size int64, state NexusState, uuid string, ) *Nexus`

NewNexus instantiates a new Nexus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNexusWithDefaults

`func NewNexusWithDefaults() *Nexus`

NewNexusWithDefaults instantiates a new Nexus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChildren

`func (o *Nexus) GetChildren() []Child`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *Nexus) GetChildrenOk() (*[]Child, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *Nexus) SetChildren(v []Child)`

SetChildren sets Children field to given value.


### GetDeviceUri

`func (o *Nexus) GetDeviceUri() string`

GetDeviceUri returns the DeviceUri field if non-nil, zero value otherwise.

### GetDeviceUriOk

`func (o *Nexus) GetDeviceUriOk() (*string, bool)`

GetDeviceUriOk returns a tuple with the DeviceUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceUri

`func (o *Nexus) SetDeviceUri(v string)`

SetDeviceUri sets DeviceUri field to given value.


### GetNode

`func (o *Nexus) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *Nexus) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *Nexus) SetNode(v string)`

SetNode sets Node field to given value.


### GetRebuilds

`func (o *Nexus) GetRebuilds() int32`

GetRebuilds returns the Rebuilds field if non-nil, zero value otherwise.

### GetRebuildsOk

`func (o *Nexus) GetRebuildsOk() (*int32, bool)`

GetRebuildsOk returns a tuple with the Rebuilds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRebuilds

`func (o *Nexus) SetRebuilds(v int32)`

SetRebuilds sets Rebuilds field to given value.


### GetProtocol

`func (o *Nexus) GetProtocol() Protocol`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *Nexus) GetProtocolOk() (*Protocol, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *Nexus) SetProtocol(v Protocol)`

SetProtocol sets Protocol field to given value.


### GetSize

`func (o *Nexus) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Nexus) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Nexus) SetSize(v int64)`

SetSize sets Size field to given value.


### GetState

`func (o *Nexus) GetState() NexusState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Nexus) GetStateOk() (*NexusState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Nexus) SetState(v NexusState)`

SetState sets State field to given value.


### GetUuid

`func (o *Nexus) GetUuid() string`

GetUuid returns the Uuid field if non-nil, zero value otherwise.

### GetUuidOk

`func (o *Nexus) GetUuidOk() (*string, bool)`

GetUuidOk returns a tuple with the Uuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUuid

`func (o *Nexus) SetUuid(v string)`

SetUuid sets Uuid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


