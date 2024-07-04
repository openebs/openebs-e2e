# RestWatch

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Callback** | **string** | callback used to notify the watch of a change | 
**Resource** | **string** | id of the resource to watch on | 

## Methods

### NewRestWatch

`func NewRestWatch(callback string, resource string, ) *RestWatch`

NewRestWatch instantiates a new RestWatch object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRestWatchWithDefaults

`func NewRestWatchWithDefaults() *RestWatch`

NewRestWatchWithDefaults instantiates a new RestWatch object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCallback

`func (o *RestWatch) GetCallback() string`

GetCallback returns the Callback field if non-nil, zero value otherwise.

### GetCallbackOk

`func (o *RestWatch) GetCallbackOk() (*string, bool)`

GetCallbackOk returns a tuple with the Callback field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallback

`func (o *RestWatch) SetCallback(v string)`

SetCallback sets Callback field to given value.


### GetResource

`func (o *RestWatch) GetResource() string`

GetResource returns the Resource field if non-nil, zero value otherwise.

### GetResourceOk

`func (o *RestWatch) GetResourceOk() (*string, bool)`

GetResourceOk returns a tuple with the Resource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResource

`func (o *RestWatch) SetResource(v string)`

SetResource sets Resource field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


