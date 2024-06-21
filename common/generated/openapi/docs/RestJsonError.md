# RestJsonError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Details** | **string** | detailed error information | 
**Message** | **string** | last reported error information | 
**Kind** | **string** | error kind | 

## Methods

### NewRestJsonError

`func NewRestJsonError(details string, message string, kind string, ) *RestJsonError`

NewRestJsonError instantiates a new RestJsonError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRestJsonErrorWithDefaults

`func NewRestJsonErrorWithDefaults() *RestJsonError`

NewRestJsonErrorWithDefaults instantiates a new RestJsonError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDetails

`func (o *RestJsonError) GetDetails() string`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *RestJsonError) GetDetailsOk() (*string, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *RestJsonError) SetDetails(v string)`

SetDetails sets Details field to given value.


### GetMessage

`func (o *RestJsonError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *RestJsonError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *RestJsonError) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetKind

`func (o *RestJsonError) GetKind() string`

GetKind returns the Kind field if non-nil, zero value otherwise.

### GetKindOk

`func (o *RestJsonError) GetKindOk() (*string, bool)`

GetKindOk returns a tuple with the Kind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKind

`func (o *RestJsonError) SetKind(v string)`

SetKind sets Kind field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


