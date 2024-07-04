# Node

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | storage node identifier | 
**Spec** | Pointer to [**NodeSpec**](NodeSpec.md) |  | [optional] 
**State** | Pointer to [**NodeState**](NodeState.md) |  | [optional] 

## Methods

### NewNode

`func NewNode(id string, ) *Node`

NewNode instantiates a new Node object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeWithDefaults

`func NewNodeWithDefaults() *Node`

NewNodeWithDefaults instantiates a new Node object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Node) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Node) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Node) SetId(v string)`

SetId sets Id field to given value.


### GetSpec

`func (o *Node) GetSpec() NodeSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *Node) GetSpecOk() (*NodeSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *Node) SetSpec(v NodeSpec)`

SetSpec sets Spec field to given value.

### HasSpec

`func (o *Node) HasSpec() bool`

HasSpec returns a boolean if a field has been set.

### GetState

`func (o *Node) GetState() NodeState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Node) GetStateOk() (*NodeState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Node) SetState(v NodeState)`

SetState sets State field to given value.

### HasState

`func (o *Node) HasState() bool`

HasState returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


