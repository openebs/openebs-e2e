# NodeState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GrpcEndpoint** | **string** | gRPC endpoint of the io-engine instance | 
**Id** | **string** | storage node identifier | 
**Status** | [**NodeStatus**](NodeStatus.md) |  | 
**NodeNqn** | Pointer to **string** | NVMe Qualified Names (NQNs) are used to uniquely describe a host or NVM subsystem for the purposes of identification and authentication | [optional] 

## Methods

### NewNodeState

`func NewNodeState(grpcEndpoint string, id string, status NodeStatus, ) *NodeState`

NewNodeState instantiates a new NodeState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeStateWithDefaults

`func NewNodeStateWithDefaults() *NodeState`

NewNodeStateWithDefaults instantiates a new NodeState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGrpcEndpoint

`func (o *NodeState) GetGrpcEndpoint() string`

GetGrpcEndpoint returns the GrpcEndpoint field if non-nil, zero value otherwise.

### GetGrpcEndpointOk

`func (o *NodeState) GetGrpcEndpointOk() (*string, bool)`

GetGrpcEndpointOk returns a tuple with the GrpcEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrpcEndpoint

`func (o *NodeState) SetGrpcEndpoint(v string)`

SetGrpcEndpoint sets GrpcEndpoint field to given value.


### GetId

`func (o *NodeState) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NodeState) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NodeState) SetId(v string)`

SetId sets Id field to given value.


### GetStatus

`func (o *NodeState) GetStatus() NodeStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *NodeState) GetStatusOk() (*NodeStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *NodeState) SetStatus(v NodeStatus)`

SetStatus sets Status field to given value.


### GetNodeNqn

`func (o *NodeState) GetNodeNqn() string`

GetNodeNqn returns the NodeNqn field if non-nil, zero value otherwise.

### GetNodeNqnOk

`func (o *NodeState) GetNodeNqnOk() (*string, bool)`

GetNodeNqnOk returns a tuple with the NodeNqn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeNqn

`func (o *NodeState) SetNodeNqn(v string)`

SetNodeNqn sets NodeNqn field to given value.

### HasNodeNqn

`func (o *NodeState) HasNodeNqn() bool`

HasNodeNqn returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


