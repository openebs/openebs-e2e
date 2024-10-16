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

// checks if the NodeAccessInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NodeAccessInfo{}

// NodeAccessInfo Frontend Node access information.
type NodeAccessInfo struct {
	// The nodename of the node.
	Name string `json:"name"`
	// The Nvme Nqn of the node's initiator.
	Nqn string `json:"nqn"`
}

type _NodeAccessInfo NodeAccessInfo

// NewNodeAccessInfo instantiates a new NodeAccessInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodeAccessInfo(name string, nqn string) *NodeAccessInfo {
	this := NodeAccessInfo{}
	this.Name = name
	this.Nqn = nqn
	return &this
}

// NewNodeAccessInfoWithDefaults instantiates a new NodeAccessInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodeAccessInfoWithDefaults() *NodeAccessInfo {
	this := NodeAccessInfo{}
	return &this
}

// GetName returns the Name field value
func (o *NodeAccessInfo) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *NodeAccessInfo) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *NodeAccessInfo) SetName(v string) {
	o.Name = v
}

// GetNqn returns the Nqn field value
func (o *NodeAccessInfo) GetNqn() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Nqn
}

// GetNqnOk returns a tuple with the Nqn field value
// and a boolean to check if the value has been set.
func (o *NodeAccessInfo) GetNqnOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Nqn, true
}

// SetNqn sets field value
func (o *NodeAccessInfo) SetNqn(v string) {
	o.Nqn = v
}

func (o NodeAccessInfo) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NodeAccessInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["nqn"] = o.Nqn
	return toSerialize, nil
}

func (o *NodeAccessInfo) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"nqn",
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

	varNodeAccessInfo := _NodeAccessInfo{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varNodeAccessInfo)

	if err != nil {
		return err
	}

	*o = NodeAccessInfo(varNodeAccessInfo)

	return err
}

type NullableNodeAccessInfo struct {
	value *NodeAccessInfo
	isSet bool
}

func (v NullableNodeAccessInfo) Get() *NodeAccessInfo {
	return v.value
}

func (v *NullableNodeAccessInfo) Set(val *NodeAccessInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableNodeAccessInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableNodeAccessInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodeAccessInfo(val *NodeAccessInfo) *NullableNodeAccessInfo {
	return &NullableNodeAccessInfo{value: val, isSet: true}
}

func (v NullableNodeAccessInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodeAccessInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


