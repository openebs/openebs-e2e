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

// checks if the RegisterAppNode type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RegisterAppNode{}

// RegisterAppNode struct for RegisterAppNode
type RegisterAppNode struct {
	// gRPC server endpoint of the app node.
	Endpoint string `json:"endpoint"`
	// Labels to be set on the app node.
	Labels *map[string]string `json:"labels,omitempty"`
}

type _RegisterAppNode RegisterAppNode

// NewRegisterAppNode instantiates a new RegisterAppNode object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegisterAppNode(endpoint string) *RegisterAppNode {
	this := RegisterAppNode{}
	this.Endpoint = endpoint
	return &this
}

// NewRegisterAppNodeWithDefaults instantiates a new RegisterAppNode object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegisterAppNodeWithDefaults() *RegisterAppNode {
	this := RegisterAppNode{}
	return &this
}

// GetEndpoint returns the Endpoint field value
func (o *RegisterAppNode) GetEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value
// and a boolean to check if the value has been set.
func (o *RegisterAppNode) GetEndpointOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Endpoint, true
}

// SetEndpoint sets field value
func (o *RegisterAppNode) SetEndpoint(v string) {
	o.Endpoint = v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *RegisterAppNode) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegisterAppNode) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *RegisterAppNode) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *RegisterAppNode) SetLabels(v map[string]string) {
	o.Labels = &v
}

func (o RegisterAppNode) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RegisterAppNode) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["endpoint"] = o.Endpoint
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	return toSerialize, nil
}

func (o *RegisterAppNode) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"endpoint",
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

	varRegisterAppNode := _RegisterAppNode{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRegisterAppNode)

	if err != nil {
		return err
	}

	*o = RegisterAppNode(varRegisterAppNode)

	return err
}

type NullableRegisterAppNode struct {
	value *RegisterAppNode
	isSet bool
}

func (v NullableRegisterAppNode) Get() *RegisterAppNode {
	return v.value
}

func (v *NullableRegisterAppNode) Set(val *RegisterAppNode) {
	v.value = val
	v.isSet = true
}

func (v NullableRegisterAppNode) IsSet() bool {
	return v.isSet
}

func (v *NullableRegisterAppNode) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegisterAppNode(val *RegisterAppNode) *NullableRegisterAppNode {
	return &NullableRegisterAppNode{value: val, isSet: true}
}

func (v NullableRegisterAppNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegisterAppNode) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


