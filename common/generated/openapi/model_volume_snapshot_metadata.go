/*
IoEngine RESTful API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"time"
	"bytes"
	"fmt"
)

// checks if the VolumeSnapshotMetadata type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VolumeSnapshotMetadata{}

// VolumeSnapshotMetadata Volume Snapshot Metadata information.
type VolumeSnapshotMetadata struct {
	Status SpecStatus `json:"status"`
	// Timestamp when snapshot is taken on the storage system.
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// Size in bytes of the snapshot (which is equivalent to its source size).
	Size int64 `json:"size"`
	// Spec size in bytes of the snapshot (which is equivalent to its source spec size).
	SpecSize int64 `json:"spec_size"`
	// Size in bytes taken by the snapshot and its predecessors.
	TotalAllocatedSize int64 `json:"total_allocated_size"`
	TxnId string `json:"txn_id"`
	Transactions map[string][]ReplicaSnapshot `json:"transactions"`
	// Number of restores done from this snapshot.
	NumRestores int32 `json:"num_restores"`
	// Number of snapshot replicas for a volumesnapshot.
	NumSnapshotReplicas int32 `json:"num_snapshot_replicas"`
}

type _VolumeSnapshotMetadata VolumeSnapshotMetadata

// NewVolumeSnapshotMetadata instantiates a new VolumeSnapshotMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVolumeSnapshotMetadata(status SpecStatus, size int64, specSize int64, totalAllocatedSize int64, txnId string, transactions map[string][]ReplicaSnapshot, numRestores int32, numSnapshotReplicas int32) *VolumeSnapshotMetadata {
	this := VolumeSnapshotMetadata{}
	this.Status = status
	this.Size = size
	this.SpecSize = specSize
	this.TotalAllocatedSize = totalAllocatedSize
	this.TxnId = txnId
	this.Transactions = transactions
	this.NumRestores = numRestores
	this.NumSnapshotReplicas = numSnapshotReplicas
	return &this
}

// NewVolumeSnapshotMetadataWithDefaults instantiates a new VolumeSnapshotMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVolumeSnapshotMetadataWithDefaults() *VolumeSnapshotMetadata {
	this := VolumeSnapshotMetadata{}
	return &this
}

// GetStatus returns the Status field value
func (o *VolumeSnapshotMetadata) GetStatus() SpecStatus {
	if o == nil {
		var ret SpecStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetStatusOk() (*SpecStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *VolumeSnapshotMetadata) SetStatus(v SpecStatus) {
	o.Status = v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise.
func (o *VolumeSnapshotMetadata) GetTimestamp() time.Time {
	if o == nil || IsNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}
	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *VolumeSnapshotMetadata) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *VolumeSnapshotMetadata) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}

// GetSize returns the Size field value
func (o *VolumeSnapshotMetadata) GetSize() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetSizeOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *VolumeSnapshotMetadata) SetSize(v int64) {
	o.Size = v
}

// GetSpecSize returns the SpecSize field value
func (o *VolumeSnapshotMetadata) GetSpecSize() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.SpecSize
}

// GetSpecSizeOk returns a tuple with the SpecSize field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetSpecSizeOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SpecSize, true
}

// SetSpecSize sets field value
func (o *VolumeSnapshotMetadata) SetSpecSize(v int64) {
	o.SpecSize = v
}

// GetTotalAllocatedSize returns the TotalAllocatedSize field value
func (o *VolumeSnapshotMetadata) GetTotalAllocatedSize() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.TotalAllocatedSize
}

// GetTotalAllocatedSizeOk returns a tuple with the TotalAllocatedSize field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetTotalAllocatedSizeOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalAllocatedSize, true
}

// SetTotalAllocatedSize sets field value
func (o *VolumeSnapshotMetadata) SetTotalAllocatedSize(v int64) {
	o.TotalAllocatedSize = v
}

// GetTxnId returns the TxnId field value
func (o *VolumeSnapshotMetadata) GetTxnId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TxnId
}

// GetTxnIdOk returns a tuple with the TxnId field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetTxnIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TxnId, true
}

// SetTxnId sets field value
func (o *VolumeSnapshotMetadata) SetTxnId(v string) {
	o.TxnId = v
}

// GetTransactions returns the Transactions field value
func (o *VolumeSnapshotMetadata) GetTransactions() map[string][]ReplicaSnapshot {
	if o == nil {
		var ret map[string][]ReplicaSnapshot
		return ret
	}

	return o.Transactions
}

// GetTransactionsOk returns a tuple with the Transactions field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetTransactionsOk() (*map[string][]ReplicaSnapshot, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Transactions, true
}

// SetTransactions sets field value
func (o *VolumeSnapshotMetadata) SetTransactions(v map[string][]ReplicaSnapshot) {
	o.Transactions = v
}

// GetNumRestores returns the NumRestores field value
func (o *VolumeSnapshotMetadata) GetNumRestores() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NumRestores
}

// GetNumRestoresOk returns a tuple with the NumRestores field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetNumRestoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumRestores, true
}

// SetNumRestores sets field value
func (o *VolumeSnapshotMetadata) SetNumRestores(v int32) {
	o.NumRestores = v
}

// GetNumSnapshotReplicas returns the NumSnapshotReplicas field value
func (o *VolumeSnapshotMetadata) GetNumSnapshotReplicas() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NumSnapshotReplicas
}

// GetNumSnapshotReplicasOk returns a tuple with the NumSnapshotReplicas field value
// and a boolean to check if the value has been set.
func (o *VolumeSnapshotMetadata) GetNumSnapshotReplicasOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumSnapshotReplicas, true
}

// SetNumSnapshotReplicas sets field value
func (o *VolumeSnapshotMetadata) SetNumSnapshotReplicas(v int32) {
	o.NumSnapshotReplicas = v
}

func (o VolumeSnapshotMetadata) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VolumeSnapshotMetadata) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["status"] = o.Status
	if !IsNil(o.Timestamp) {
		toSerialize["timestamp"] = o.Timestamp
	}
	toSerialize["size"] = o.Size
	toSerialize["spec_size"] = o.SpecSize
	toSerialize["total_allocated_size"] = o.TotalAllocatedSize
	toSerialize["txn_id"] = o.TxnId
	toSerialize["transactions"] = o.Transactions
	toSerialize["num_restores"] = o.NumRestores
	toSerialize["num_snapshot_replicas"] = o.NumSnapshotReplicas
	return toSerialize, nil
}

func (o *VolumeSnapshotMetadata) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"status",
		"size",
		"spec_size",
		"total_allocated_size",
		"txn_id",
		"transactions",
		"num_restores",
		"num_snapshot_replicas",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varVolumeSnapshotMetadata := _VolumeSnapshotMetadata{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVolumeSnapshotMetadata)

	if err != nil {
		return err
	}

	*o = VolumeSnapshotMetadata(varVolumeSnapshotMetadata)

	return err
}

type NullableVolumeSnapshotMetadata struct {
	value *VolumeSnapshotMetadata
	isSet bool
}

func (v NullableVolumeSnapshotMetadata) Get() *VolumeSnapshotMetadata {
	return v.value
}

func (v *NullableVolumeSnapshotMetadata) Set(val *VolumeSnapshotMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableVolumeSnapshotMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableVolumeSnapshotMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVolumeSnapshotMetadata(val *VolumeSnapshotMetadata) *NullableVolumeSnapshotMetadata {
	return &NullableVolumeSnapshotMetadata{value: val, isSet: true}
}

func (v NullableVolumeSnapshotMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVolumeSnapshotMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


