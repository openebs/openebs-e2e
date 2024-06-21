# NodeSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GrpcEndpoint** | **string** | gRPC endpoint of the io-engine instance | 
**Id** | **string** | storage node identifier | 
**Labels** | Pointer to **map[string]string** | labels to be set on the node | [optional] 
**Cordondrainstate** | Pointer to [**NullableCordonDrainState**](CordonDrainState.md) | the drain state | [optional] 
**NodeNqn** | Pointer to **string** | NVMe Qualified Names (NQNs) are used to uniquely describe a host or NVM subsystem for the purposes of identification and authentication | [optional] 

## Methods

### NewNodeSpec

`func NewNodeSpec(grpcEndpoint string, id string, ) *NodeSpec`

NewNodeSpec instantiates a new NodeSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeSpecWithDefaults

`func NewNodeSpecWithDefaults() *NodeSpec`

NewNodeSpecWithDefaults instantiates a new NodeSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGrpcEndpoint

`func (o *NodeSpec) GetGrpcEndpoint() string`

GetGrpcEndpoint returns the GrpcEndpoint field if non-nil, zero value otherwise.

### GetGrpcEndpointOk

`func (o *NodeSpec) GetGrpcEndpointOk() (*string, bool)`

GetGrpcEndpointOk returns a tuple with the GrpcEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrpcEndpoint

`func (o *NodeSpec) SetGrpcEndpoint(v string)`

SetGrpcEndpoint sets GrpcEndpoint field to given value.


### GetId

`func (o *NodeSpec) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NodeSpec) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NodeSpec) SetId(v string)`

SetId sets Id field to given value.


### GetLabels

`func (o *NodeSpec) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *NodeSpec) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *NodeSpec) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *NodeSpec) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetCordondrainstate

`func (o *NodeSpec) GetCordondrainstate() CordonDrainState`

GetCordondrainstate returns the Cordondrainstate field if non-nil, zero value otherwise.

### GetCordondrainstateOk

`func (o *NodeSpec) GetCordondrainstateOk() (*CordonDrainState, bool)`

GetCordondrainstateOk returns a tuple with the Cordondrainstate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCordondrainstate

`func (o *NodeSpec) SetCordondrainstate(v CordonDrainState)`

SetCordondrainstate sets Cordondrainstate field to given value.

### HasCordondrainstate

`func (o *NodeSpec) HasCordondrainstate() bool`

HasCordondrainstate returns a boolean if a field has been set.

### SetCordondrainstateNil

`func (o *NodeSpec) SetCordondrainstateNil(b bool)`

 SetCordondrainstateNil sets the value for Cordondrainstate to be an explicit nil

### UnsetCordondrainstate
`func (o *NodeSpec) UnsetCordondrainstate()`

UnsetCordondrainstate ensures that no value is present for Cordondrainstate, not even an explicit nil
### GetNodeNqn

`func (o *NodeSpec) GetNodeNqn() string`

GetNodeNqn returns the NodeNqn field if non-nil, zero value otherwise.

### GetNodeNqnOk

`func (o *NodeSpec) GetNodeNqnOk() (*string, bool)`

GetNodeNqnOk returns a tuple with the NodeNqn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeNqn

`func (o *NodeSpec) SetNodeNqn(v string)`

SetNodeNqn sets NodeNqn field to given value.

### HasNodeNqn

`func (o *NodeSpec) HasNodeNqn() bool`

HasNodeNqn returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


