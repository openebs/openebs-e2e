# PublishVolumeBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PublishContext** | **map[string]string** | Controller Volume Publish context | 
**ReuseExisting** | Pointer to **bool** | Allows reusing of the current target. | [optional] 
**Node** | Pointer to **string** | The node where the target will reside in. It may be moved elsewhere during volume republish. | [optional] 
**Protocol** | [**VolumeShareProtocol**](VolumeShareProtocol.md) | The protocol used to connect to the front-end node. | 
**Republish** | Pointer to **bool** | Allows republishing the volume on the node by shutting down the existing target first. | [optional] 
**FrontendNode** | Pointer to **string** | The node where the front-end workload resides. If the workload moves then the volume must be republished. | [optional] 

## Methods

### NewPublishVolumeBody

`func NewPublishVolumeBody(publishContext map[string]string, protocol VolumeShareProtocol, ) *PublishVolumeBody`

NewPublishVolumeBody instantiates a new PublishVolumeBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPublishVolumeBodyWithDefaults

`func NewPublishVolumeBodyWithDefaults() *PublishVolumeBody`

NewPublishVolumeBodyWithDefaults instantiates a new PublishVolumeBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPublishContext

`func (o *PublishVolumeBody) GetPublishContext() map[string]string`

GetPublishContext returns the PublishContext field if non-nil, zero value otherwise.

### GetPublishContextOk

`func (o *PublishVolumeBody) GetPublishContextOk() (*map[string]string, bool)`

GetPublishContextOk returns a tuple with the PublishContext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishContext

`func (o *PublishVolumeBody) SetPublishContext(v map[string]string)`

SetPublishContext sets PublishContext field to given value.


### GetReuseExisting

`func (o *PublishVolumeBody) GetReuseExisting() bool`

GetReuseExisting returns the ReuseExisting field if non-nil, zero value otherwise.

### GetReuseExistingOk

`func (o *PublishVolumeBody) GetReuseExistingOk() (*bool, bool)`

GetReuseExistingOk returns a tuple with the ReuseExisting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReuseExisting

`func (o *PublishVolumeBody) SetReuseExisting(v bool)`

SetReuseExisting sets ReuseExisting field to given value.

### HasReuseExisting

`func (o *PublishVolumeBody) HasReuseExisting() bool`

HasReuseExisting returns a boolean if a field has been set.

### GetNode

`func (o *PublishVolumeBody) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *PublishVolumeBody) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *PublishVolumeBody) SetNode(v string)`

SetNode sets Node field to given value.

### HasNode

`func (o *PublishVolumeBody) HasNode() bool`

HasNode returns a boolean if a field has been set.

### GetProtocol

`func (o *PublishVolumeBody) GetProtocol() VolumeShareProtocol`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *PublishVolumeBody) GetProtocolOk() (*VolumeShareProtocol, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *PublishVolumeBody) SetProtocol(v VolumeShareProtocol)`

SetProtocol sets Protocol field to given value.


### GetRepublish

`func (o *PublishVolumeBody) GetRepublish() bool`

GetRepublish returns the Republish field if non-nil, zero value otherwise.

### GetRepublishOk

`func (o *PublishVolumeBody) GetRepublishOk() (*bool, bool)`

GetRepublishOk returns a tuple with the Republish field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRepublish

`func (o *PublishVolumeBody) SetRepublish(v bool)`

SetRepublish sets Republish field to given value.

### HasRepublish

`func (o *PublishVolumeBody) HasRepublish() bool`

HasRepublish returns a boolean if a field has been set.

### GetFrontendNode

`func (o *PublishVolumeBody) GetFrontendNode() string`

GetFrontendNode returns the FrontendNode field if non-nil, zero value otherwise.

### GetFrontendNodeOk

`func (o *PublishVolumeBody) GetFrontendNodeOk() (*string, bool)`

GetFrontendNodeOk returns a tuple with the FrontendNode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontendNode

`func (o *PublishVolumeBody) SetFrontendNode(v string)`

SetFrontendNode sets FrontendNode field to given value.

### HasFrontendNode

`func (o *PublishVolumeBody) HasFrontendNode() bool`

HasFrontendNode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


