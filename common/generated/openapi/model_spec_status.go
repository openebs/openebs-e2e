/*
IoEngine RESTful API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// SpecStatus Common base state for a resource
type SpecStatus string

// List of SpecStatus
const (
	CREATING SpecStatus = "Creating"
	CREATED SpecStatus = "Created"
	DELETING SpecStatus = "Deleting"
	DELETED SpecStatus = "Deleted"
)

// All allowed values of SpecStatus enum
var AllowedSpecStatusEnumValues = []SpecStatus{
	"Creating",
	"Created",
	"Deleting",
	"Deleted",
}

func (v *SpecStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := SpecStatus(value)
	for _, existing := range AllowedSpecStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid SpecStatus", value)
}

// NewSpecStatusFromValue returns a pointer to a valid SpecStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewSpecStatusFromValue(v string) (*SpecStatus, error) {
	ev := SpecStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for SpecStatus: valid values are %v", v, AllowedSpecStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v SpecStatus) IsValid() bool {
	for _, existing := range AllowedSpecStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to SpecStatus value
func (v SpecStatus) Ptr() *SpecStatus {
	return &v
}

type NullableSpecStatus struct {
	value *SpecStatus
	isSet bool
}

func (v NullableSpecStatus) Get() *SpecStatus {
	return v.value
}

func (v *NullableSpecStatus) Set(val *SpecStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableSpecStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableSpecStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSpecStatus(val *SpecStatus) *NullableSpecStatus {
	return &NullableSpecStatus{value: val, isSet: true}
}

func (v NullableSpecStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSpecStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

