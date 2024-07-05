# PoolTopology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Labelled** | Pointer to [**LabelledTopology**](LabelledTopology.md) | volume pool topology definition through labels | [optional] 

## Methods

### NewPoolTopology

`func NewPoolTopology() *PoolTopology`

NewPoolTopology instantiates a new PoolTopology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPoolTopologyWithDefaults

`func NewPoolTopologyWithDefaults() *PoolTopology`

NewPoolTopologyWithDefaults instantiates a new PoolTopology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabelled

`func (o *PoolTopology) GetLabelled() LabelledTopology`

GetLabelled returns the Labelled field if non-nil, zero value otherwise.

### GetLabelledOk

`func (o *PoolTopology) GetLabelledOk() (*LabelledTopology, bool)`

GetLabelledOk returns a tuple with the Labelled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabelled

`func (o *PoolTopology) SetLabelled(v LabelledTopology)`

SetLabelled sets Labelled field to given value.

### HasLabelled

`func (o *PoolTopology) HasLabelled() bool`

HasLabelled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


