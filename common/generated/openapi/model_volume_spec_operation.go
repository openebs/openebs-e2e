/*
IoEngine RESTful API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the VolumeSpecOperation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VolumeSpecOperation{}

// VolumeSpecOperation Record of the operation in progress
type VolumeSpecOperation struct {
	// Record of the operation
	Operation string `json:"operation"`
	// Result of the operation
	Result *bool `json:"result,omitempty"`
}

type _VolumeSpecOperation VolumeSpecOperation

// NewVolumeSpecOperation instantiates a new VolumeSpecOperation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVolumeSpecOperation(operation string) *VolumeSpecOperation {
	this := VolumeSpecOperation{}
	this.Operation = operation
	return &this
}

// NewVolumeSpecOperationWithDefaults instantiates a new VolumeSpecOperation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVolumeSpecOperationWithDefaults() *VolumeSpecOperation {
	this := VolumeSpecOperation{}
	return &this
}

// GetOperation returns the Operation field value
func (o *VolumeSpecOperation) GetOperation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Operation
}

// GetOperationOk returns a tuple with the Operation field value
// and a boolean to check if the value has been set.
func (o *VolumeSpecOperation) GetOperationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Operation, true
}

// SetOperation sets field value
func (o *VolumeSpecOperation) SetOperation(v string) {
	o.Operation = v
}

// GetResult returns the Result field value if set, zero value otherwise.
func (o *VolumeSpecOperation) GetResult() bool {
	if o == nil || IsNil(o.Result) {
		var ret bool
		return ret
	}
	return *o.Result
}

// GetResultOk returns a tuple with the Result field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VolumeSpecOperation) GetResultOk() (*bool, bool) {
	if o == nil || IsNil(o.Result) {
		return nil, false
	}
	return o.Result, true
}

// HasResult returns a boolean if a field has been set.
func (o *VolumeSpecOperation) HasResult() bool {
	if o != nil && !IsNil(o.Result) {
		return true
	}

	return false
}

// SetResult gets a reference to the given bool and assigns it to the Result field.
func (o *VolumeSpecOperation) SetResult(v bool) {
	o.Result = &v
}

func (o VolumeSpecOperation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VolumeSpecOperation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["operation"] = o.Operation
	if !IsNil(o.Result) {
		toSerialize["result"] = o.Result
	}
	return toSerialize, nil
}

func (o *VolumeSpecOperation) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"operation",
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

	varVolumeSpecOperation := _VolumeSpecOperation{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVolumeSpecOperation)

	if err != nil {
		return err
	}

	*o = VolumeSpecOperation(varVolumeSpecOperation)

	return err
}

type NullableVolumeSpecOperation struct {
	value *VolumeSpecOperation
	isSet bool
}

func (v NullableVolumeSpecOperation) Get() *VolumeSpecOperation {
	return v.value
}

func (v *NullableVolumeSpecOperation) Set(val *VolumeSpecOperation) {
	v.value = val
	v.isSet = true
}

func (v NullableVolumeSpecOperation) IsSet() bool {
	return v.isSet
}

func (v *NullableVolumeSpecOperation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVolumeSpecOperation(val *VolumeSpecOperation) *NullableVolumeSpecOperation {
	return &NullableVolumeSpecOperation{value: val, isSet: true}
}

func (v NullableVolumeSpecOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVolumeSpecOperation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


