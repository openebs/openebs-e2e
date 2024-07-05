# Topology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NodeTopology** | Pointer to [**NullableNodeTopology**](NodeTopology.md) |  | [optional] 
**PoolTopology** | Pointer to [**NullablePoolTopology**](PoolTopology.md) |  | [optional] 

## Methods

### NewTopology

`func NewTopology() *Topology`

NewTopology instantiates a new Topology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTopologyWithDefaults

`func NewTopologyWithDefaults() *Topology`

NewTopologyWithDefaults instantiates a new Topology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNodeTopology

`func (o *Topology) GetNodeTopology() NodeTopology`

GetNodeTopology returns the NodeTopology field if non-nil, zero value otherwise.

### GetNodeTopologyOk

`func (o *Topology) GetNodeTopologyOk() (*NodeTopology, bool)`

GetNodeTopologyOk returns a tuple with the NodeTopology field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodeTopology

`func (o *Topology) SetNodeTopology(v NodeTopology)`

SetNodeTopology sets NodeTopology field to given value.

### HasNodeTopology

`func (o *Topology) HasNodeTopology() bool`

HasNodeTopology returns a boolean if a field has been set.

### SetNodeTopologyNil

`func (o *Topology) SetNodeTopologyNil(b bool)`

 SetNodeTopologyNil sets the value for NodeTopology to be an explicit nil

### UnsetNodeTopology
`func (o *Topology) UnsetNodeTopology()`

UnsetNodeTopology ensures that no value is present for NodeTopology, not even an explicit nil
### GetPoolTopology

`func (o *Topology) GetPoolTopology() PoolTopology`

GetPoolTopology returns the PoolTopology field if non-nil, zero value otherwise.

### GetPoolTopologyOk

`func (o *Topology) GetPoolTopologyOk() (*PoolTopology, bool)`

GetPoolTopologyOk returns a tuple with the PoolTopology field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoolTopology

`func (o *Topology) SetPoolTopology(v PoolTopology)`

SetPoolTopology sets PoolTopology field to given value.

### HasPoolTopology

`func (o *Topology) HasPoolTopology() bool`

HasPoolTopology returns a boolean if a field has been set.

### SetPoolTopologyNil

`func (o *Topology) SetPoolTopologyNil(b bool)`

 SetPoolTopologyNil sets the value for PoolTopology to be an explicit nil

### UnsetPoolTopology
`func (o *Topology) UnsetPoolTopology()`

UnsetPoolTopology ensures that no value is present for PoolTopology, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


