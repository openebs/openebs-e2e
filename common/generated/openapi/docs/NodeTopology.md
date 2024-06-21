# NodeTopology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Explicit** | Pointer to [**ExplicitNodeTopology**](ExplicitNodeTopology.md) | volume topology, explicitly selected | [optional] 
**Labelled** | Pointer to [**LabelledTopology**](LabelledTopology.md) | volume topology definition through labels | [optional] 

## Methods

### NewNodeTopology

`func NewNodeTopology() *NodeTopology`

NewNodeTopology instantiates a new NodeTopology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeTopologyWithDefaults

`func NewNodeTopologyWithDefaults() *NodeTopology`

NewNodeTopologyWithDefaults instantiates a new NodeTopology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExplicit

`func (o *NodeTopology) GetExplicit() ExplicitNodeTopology`

GetExplicit returns the Explicit field if non-nil, zero value otherwise.

### GetExplicitOk

`func (o *NodeTopology) GetExplicitOk() (*ExplicitNodeTopology, bool)`

GetExplicitOk returns a tuple with the Explicit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExplicit

`func (o *NodeTopology) SetExplicit(v ExplicitNodeTopology)`

SetExplicit sets Explicit field to given value.

### HasExplicit

`func (o *NodeTopology) HasExplicit() bool`

HasExplicit returns a boolean if a field has been set.

### GetLabelled

`func (o *NodeTopology) GetLabelled() LabelledTopology`

GetLabelled returns the Labelled field if non-nil, zero value otherwise.

### GetLabelledOk

`func (o *NodeTopology) GetLabelledOk() (*LabelledTopology, bool)`

GetLabelledOk returns a tuple with the Labelled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabelled

`func (o *NodeTopology) SetLabelled(v LabelledTopology)`

SetLabelled sets Labelled field to given value.

### HasLabelled

`func (o *NodeTopology) HasLabelled() bool`

HasLabelled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


