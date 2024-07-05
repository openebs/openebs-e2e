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

// Protocol Common Protocol
type Protocol string

// List of Protocol
const (
	PROTOCOL_NONE  Protocol = "none"
	PROTOCOL_NVMF  Protocol = "nvmf"
	PROTOCOL_ISCSI Protocol = "iscsi"
	PROTOCOL_NBD   Protocol = "nbd"
)

// All allowed values of Protocol enum
var AllowedProtocolEnumValues = []Protocol{
	"none",
	"nvmf",
	"iscsi",
	"nbd",
}

func (v *Protocol) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Protocol(value)
	for _, existing := range AllowedProtocolEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Protocol", value)
}

// NewProtocolFromValue returns a pointer to a valid Protocol
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewProtocolFromValue(v string) (*Protocol, error) {
	ev := Protocol(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Protocol: valid values are %v", v, AllowedProtocolEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Protocol) IsValid() bool {
	for _, existing := range AllowedProtocolEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Protocol value
func (v Protocol) Ptr() *Protocol {
	return &v
}

type NullableProtocol struct {
	value *Protocol
	isSet bool
}

func (v NullableProtocol) Get() *Protocol {
	return v.value
}

func (v *NullableProtocol) Set(val *Protocol) {
	v.value = val
	v.isSet = true
}

func (v NullableProtocol) IsSet() bool {
	return v.isSet
}

func (v *NullableProtocol) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProtocol(val *Protocol) *NullableProtocol {
	return &NullableProtocol{value: val, isSet: true}
}

func (v NullableProtocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProtocol) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

