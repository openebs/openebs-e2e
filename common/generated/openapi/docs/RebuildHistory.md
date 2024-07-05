# RebuildHistory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TargetUuid** | **string** | Id of the volume target | 
**Records** | [**[]RebuildRecord**](RebuildRecord.md) | Array of rebuild record | 

## Methods

### NewRebuildHistory

`func NewRebuildHistory(targetUuid string, records []RebuildRecord, ) *RebuildHistory`

NewRebuildHistory instantiates a new RebuildHistory object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRebuildHistoryWithDefaults

`func NewRebuildHistoryWithDefaults() *RebuildHistory`

NewRebuildHistoryWithDefaults instantiates a new RebuildHistory object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTargetUuid

`func (o *RebuildHistory) GetTargetUuid() string`

GetTargetUuid returns the TargetUuid field if non-nil, zero value otherwise.

### GetTargetUuidOk

`func (o *RebuildHistory) GetTargetUuidOk() (*string, bool)`

GetTargetUuidOk returns a tuple with the TargetUuid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetUuid

`func (o *RebuildHistory) SetTargetUuid(v string)`

SetTargetUuid sets TargetUuid field to given value.


### GetRecords

`func (o *RebuildHistory) GetRecords() []RebuildRecord`

GetRecords returns the Records field if non-nil, zero value otherwise.

### GetRecordsOk

`func (o *RebuildHistory) GetRecordsOk() (*[]RebuildRecord, bool)`

GetRecordsOk returns a tuple with the Records field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecords

`func (o *RebuildHistory) SetRecords(v []RebuildRecord)`

SetRecords sets Records field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


