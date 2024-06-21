# VolumeTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Node** | **string** | The node where front-end IO will be sent to | 
**Protocol** | Pointer to [**VolumeShareProtocol**](VolumeShareProtocol.md) |  | [optional] 
**FrontendNodes** | Pointer to [**[]NodeAccessInfo**](NodeAccessInfo.md) | The nodes where the front-end workload resides. If the workload moves then the volume must be republished. | [optional] 

## Methods

### NewVolumeTarget

`func NewVolumeTarget(node string, ) *VolumeTarget`

NewVolumeTarget instantiates a new VolumeTarget object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeTargetWithDefaults

`func NewVolumeTargetWithDefaults() *VolumeTarget`

NewVolumeTargetWithDefaults instantiates a new VolumeTarget object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNode

`func (o *VolumeTarget) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *VolumeTarget) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *VolumeTarget) SetNode(v string)`

SetNode sets Node field to given value.


### GetProtocol

`func (o *VolumeTarget) GetProtocol() VolumeShareProtocol`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *VolumeTarget) GetProtocolOk() (*VolumeShareProtocol, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *VolumeTarget) SetProtocol(v VolumeShareProtocol)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *VolumeTarget) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetFrontendNodes

`func (o *VolumeTarget) GetFrontendNodes() []NodeAccessInfo`

GetFrontendNodes returns the FrontendNodes field if non-nil, zero value otherwise.

### GetFrontendNodesOk

`func (o *VolumeTarget) GetFrontendNodesOk() (*[]NodeAccessInfo, bool)`

GetFrontendNodesOk returns a tuple with the FrontendNodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontendNodes

`func (o *VolumeTarget) SetFrontendNodes(v []NodeAccessInfo)`

SetFrontendNodes sets FrontendNodes field to given value.

### HasFrontendNodes

`func (o *VolumeTarget) HasFrontendNodes() bool`

HasFrontendNodes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


