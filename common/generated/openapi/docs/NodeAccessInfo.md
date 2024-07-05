# NodeAccessInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The nodename of the node. | 
**Nqn** | **string** | The Nvme Nqn of the node&#39;s initiator. | 

## Methods

### NewNodeAccessInfo

`func NewNodeAccessInfo(name string, nqn string, ) *NodeAccessInfo`

NewNodeAccessInfo instantiates a new NodeAccessInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeAccessInfoWithDefaults

`func NewNodeAccessInfoWithDefaults() *NodeAccessInfo`

NewNodeAccessInfoWithDefaults instantiates a new NodeAccessInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NodeAccessInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NodeAccessInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NodeAccessInfo) SetName(v string)`

SetName sets Name field to given value.


### GetNqn

`func (o *NodeAccessInfo) GetNqn() string`

GetNqn returns the Nqn field if non-nil, zero value otherwise.

### GetNqnOk

`func (o *NodeAccessInfo) GetNqnOk() (*string, bool)`

GetNqnOk returns a tuple with the Nqn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNqn

`func (o *NodeAccessInfo) SetNqn(v string)`

SetNqn sets Nqn field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


