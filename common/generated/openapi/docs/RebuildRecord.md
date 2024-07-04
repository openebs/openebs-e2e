# RebuildRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ChildUri** | **string** | Uri of the rebuilding child | 
**SrcUri** | **string** | Uri of source child for rebuild job | 
**RebuildJobState** | [**RebuildJobState**](RebuildJobState.md) |  | 
**BlocksTotal** | **int32** | Total blocks to rebuild | 
**BlocksRecovered** | **int32** | Number of blocks processed | 
**BlocksTransferred** | **int32** | Number of blocks to transferred | 
**BlocksRemaining** | **int32** | Number of blocks remaining | 
**BlockSize** | **int32** | Size of each block in the task | 
**IsPartial** | **bool** | True means its Partial rebuild job. If false, its Full rebuild job | 
**StartTime** | **time.Time** | Start time of the rebuild job (UTC) | 
**EndTime** | **time.Time** | End time of the rebuild job (UTC) | 

## Methods

### NewRebuildRecord

`func NewRebuildRecord(childUri string, srcUri string, rebuildJobState RebuildJobState, blocksTotal int32, blocksRecovered int32, blocksTransferred int32, blocksRemaining int32, blockSize int32, isPartial bool, startTime time.Time, endTime time.Time, ) *RebuildRecord`

NewRebuildRecord instantiates a new RebuildRecord object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRebuildRecordWithDefaults

`func NewRebuildRecordWithDefaults() *RebuildRecord`

NewRebuildRecordWithDefaults instantiates a new RebuildRecord object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChildUri

`func (o *RebuildRecord) GetChildUri() string`

GetChildUri returns the ChildUri field if non-nil, zero value otherwise.

### GetChildUriOk

`func (o *RebuildRecord) GetChildUriOk() (*string, bool)`

GetChildUriOk returns a tuple with the ChildUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildUri

`func (o *RebuildRecord) SetChildUri(v string)`

SetChildUri sets ChildUri field to given value.


### GetSrcUri

`func (o *RebuildRecord) GetSrcUri() string`

GetSrcUri returns the SrcUri field if non-nil, zero value otherwise.

### GetSrcUriOk

`func (o *RebuildRecord) GetSrcUriOk() (*string, bool)`

GetSrcUriOk returns a tuple with the SrcUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrcUri

`func (o *RebuildRecord) SetSrcUri(v string)`

SetSrcUri sets SrcUri field to given value.


### GetRebuildJobState

`func (o *RebuildRecord) GetRebuildJobState() RebuildJobState`

GetRebuildJobState returns the RebuildJobState field if non-nil, zero value otherwise.

### GetRebuildJobStateOk

`func (o *RebuildRecord) GetRebuildJobStateOk() (*RebuildJobState, bool)`

GetRebuildJobStateOk returns a tuple with the RebuildJobState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRebuildJobState

`func (o *RebuildRecord) SetRebuildJobState(v RebuildJobState)`

SetRebuildJobState sets RebuildJobState field to given value.


### GetBlocksTotal

`func (o *RebuildRecord) GetBlocksTotal() int32`

GetBlocksTotal returns the BlocksTotal field if non-nil, zero value otherwise.

### GetBlocksTotalOk

`func (o *RebuildRecord) GetBlocksTotalOk() (*int32, bool)`

GetBlocksTotalOk returns a tuple with the BlocksTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlocksTotal

`func (o *RebuildRecord) SetBlocksTotal(v int32)`

SetBlocksTotal sets BlocksTotal field to given value.


### GetBlocksRecovered

`func (o *RebuildRecord) GetBlocksRecovered() int32`

GetBlocksRecovered returns the BlocksRecovered field if non-nil, zero value otherwise.

### GetBlocksRecoveredOk

`func (o *RebuildRecord) GetBlocksRecoveredOk() (*int32, bool)`

GetBlocksRecoveredOk returns a tuple with the BlocksRecovered field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlocksRecovered

`func (o *RebuildRecord) SetBlocksRecovered(v int32)`

SetBlocksRecovered sets BlocksRecovered field to given value.


### GetBlocksTransferred

`func (o *RebuildRecord) GetBlocksTransferred() int32`

GetBlocksTransferred returns the BlocksTransferred field if non-nil, zero value otherwise.

### GetBlocksTransferredOk

`func (o *RebuildRecord) GetBlocksTransferredOk() (*int32, bool)`

GetBlocksTransferredOk returns a tuple with the BlocksTransferred field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlocksTransferred

`func (o *RebuildRecord) SetBlocksTransferred(v int32)`

SetBlocksTransferred sets BlocksTransferred field to given value.


### GetBlocksRemaining

`func (o *RebuildRecord) GetBlocksRemaining() int32`

GetBlocksRemaining returns the BlocksRemaining field if non-nil, zero value otherwise.

### GetBlocksRemainingOk

`func (o *RebuildRecord) GetBlocksRemainingOk() (*int32, bool)`

GetBlocksRemainingOk returns a tuple with the BlocksRemaining field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlocksRemaining

`func (o *RebuildRecord) SetBlocksRemaining(v int32)`

SetBlocksRemaining sets BlocksRemaining field to given value.


### GetBlockSize

`func (o *RebuildRecord) GetBlockSize() int32`

GetBlockSize returns the BlockSize field if non-nil, zero value otherwise.

### GetBlockSizeOk

`func (o *RebuildRecord) GetBlockSizeOk() (*int32, bool)`

GetBlockSizeOk returns a tuple with the BlockSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockSize

`func (o *RebuildRecord) SetBlockSize(v int32)`

SetBlockSize sets BlockSize field to given value.


### GetIsPartial

`func (o *RebuildRecord) GetIsPartial() bool`

GetIsPartial returns the IsPartial field if non-nil, zero value otherwise.

### GetIsPartialOk

`func (o *RebuildRecord) GetIsPartialOk() (*bool, bool)`

GetIsPartialOk returns a tuple with the IsPartial field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPartial

`func (o *RebuildRecord) SetIsPartial(v bool)`

SetIsPartial sets IsPartial field to given value.


### GetStartTime

`func (o *RebuildRecord) GetStartTime() time.Time`

GetStartTime returns the StartTime field if non-nil, zero value otherwise.

### GetStartTimeOk

`func (o *RebuildRecord) GetStartTimeOk() (*time.Time, bool)`

GetStartTimeOk returns a tuple with the StartTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartTime

`func (o *RebuildRecord) SetStartTime(v time.Time)`

SetStartTime sets StartTime field to given value.


### GetEndTime

`func (o *RebuildRecord) GetEndTime() time.Time`

GetEndTime returns the EndTime field if non-nil, zero value otherwise.

### GetEndTimeOk

`func (o *RebuildRecord) GetEndTimeOk() (*time.Time, bool)`

GetEndTimeOk returns a tuple with the EndTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndTime

`func (o *RebuildRecord) SetEndTime(v time.Time)`

SetEndTime sets EndTime field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


