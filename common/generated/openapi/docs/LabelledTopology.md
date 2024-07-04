# LabelledTopology

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Exclusion** | **map[string]string** | Excludes resources with the same $label name, eg:  \&quot;Zone\&quot; would not allow for resources with the same \&quot;Zone\&quot; value  to be used for a certain operation, eg:  A node with \&quot;Zone: A\&quot; would not be paired up with a node with \&quot;Zone: A\&quot;,  but it could be paired up with a node with \&quot;Zone: B\&quot;  exclusive label NAME in the form \&quot;NAME\&quot;, and not \&quot;NAME: VALUE\&quot; | 
**Inclusion** | **map[string]string** | Includes resources with the same $label or $label:$value eg:  if label is \&quot;Zone: A\&quot;:  A resource with \&quot;Zone: A\&quot; would be paired up with a resource with \&quot;Zone: A\&quot;,  but not with a resource with \&quot;Zone: B\&quot;  if label is \&quot;Zone\&quot;:  A resource with \&quot;Zone: A\&quot; would be paired up with a resource with \&quot;Zone: B\&quot;,  but not with a resource with \&quot;OtherLabel: B\&quot;  inclusive label key value in the form \&quot;NAME: VALUE\&quot; | 

## Methods

### NewLabelledTopology

`func NewLabelledTopology(exclusion map[string]string, inclusion map[string]string, ) *LabelledTopology`

NewLabelledTopology instantiates a new LabelledTopology object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLabelledTopologyWithDefaults

`func NewLabelledTopologyWithDefaults() *LabelledTopology`

NewLabelledTopologyWithDefaults instantiates a new LabelledTopology object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExclusion

`func (o *LabelledTopology) GetExclusion() map[string]string`

GetExclusion returns the Exclusion field if non-nil, zero value otherwise.

### GetExclusionOk

`func (o *LabelledTopology) GetExclusionOk() (*map[string]string, bool)`

GetExclusionOk returns a tuple with the Exclusion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExclusion

`func (o *LabelledTopology) SetExclusion(v map[string]string)`

SetExclusion sets Exclusion field to given value.


### GetInclusion

`func (o *LabelledTopology) GetInclusion() map[string]string`

GetInclusion returns the Inclusion field if non-nil, zero value otherwise.

### GetInclusionOk

`func (o *LabelledTopology) GetInclusionOk() (*map[string]string, bool)`

GetInclusionOk returns a tuple with the Inclusion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInclusion

`func (o *LabelledTopology) SetInclusion(v map[string]string)`

SetInclusion sets Inclusion field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


