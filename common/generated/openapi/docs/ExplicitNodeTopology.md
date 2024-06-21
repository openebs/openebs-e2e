# ExplicitNodeTopology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowedNodes** | **[]string** | replicas can only be placed on these nodes | 
**PreferredNodes** | **[]string** | preferred nodes to place the replicas | 

## Methods

### NewExplicitNodeTopology

`func NewExplicitNodeTopology(allowedNodes []string, preferredNodes []string, ) *ExplicitNodeTopology`

NewExplicitNodeTopology instantiates a new ExplicitNodeTopology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExplicitNodeTopologyWithDefaults

`func NewExplicitNodeTopologyWithDefaults() *ExplicitNodeTopology`

NewExplicitNodeTopologyWithDefaults instantiates a new ExplicitNodeTopology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAllowedNodes

`func (o *ExplicitNodeTopology) GetAllowedNodes() []string`

GetAllowedNodes returns the AllowedNodes field if non-nil, zero value otherwise.

### GetAllowedNodesOk

`func (o *ExplicitNodeTopology) GetAllowedNodesOk() (*[]string, bool)`

GetAllowedNodesOk returns a tuple with the AllowedNodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedNodes

`func (o *ExplicitNodeTopology) SetAllowedNodes(v []string)`

SetAllowedNodes sets AllowedNodes field to given value.


### GetPreferredNodes

`func (o *ExplicitNodeTopology) GetPreferredNodes() []string`

GetPreferredNodes returns the PreferredNodes field if non-nil, zero value otherwise.

### GetPreferredNodesOk

`func (o *ExplicitNodeTopology) GetPreferredNodesOk() (*[]string, bool)`

GetPreferredNodesOk returns a tuple with the PreferredNodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreferredNodes

`func (o *ExplicitNodeTopology) SetPreferredNodes(v []string)`

SetPreferredNodes sets PreferredNodes field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


