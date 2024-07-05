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

// checks if the AppNodeSpec type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AppNodeSpec{}

// AppNodeSpec App node attributes.
type AppNodeSpec struct {
	// App node identifier.
	Id string `json:"id"`
	// gRPC server endpoint of the app node.
	Endpoint string `json:"endpoint"`
	// Labels to be set on the app node.
	Labels *map[string]string `json:"labels,omitempty"`
}

type _AppNodeSpec AppNodeSpec

// NewAppNodeSpec instantiates a new AppNodeSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAppNodeSpec(id string, endpoint string) *AppNodeSpec {
	this := AppNodeSpec{}
	this.Id = id
	this.Endpoint = endpoint
	return &this
}

// NewAppNodeSpecWithDefaults instantiates a new AppNodeSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAppNodeSpecWithDefaults() *AppNodeSpec {
	this := AppNodeSpec{}
	return &this
}

// GetId returns the Id field value
func (o *AppNodeSpec) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AppNodeSpec) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AppNodeSpec) SetId(v string) {
	o.Id = v
}

// GetEndpoint returns the Endpoint field value
func (o *AppNodeSpec) GetEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value
// and a boolean to check if the value has been set.
func (o *AppNodeSpec) GetEndpointOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Endpoint, true
}

// SetEndpoint sets field value
func (o *AppNodeSpec) SetEndpoint(v string) {
	o.Endpoint = v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *AppNodeSpec) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AppNodeSpec) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *AppNodeSpec) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *AppNodeSpec) SetLabels(v map[string]string) {
	o.Labels = &v
}

func (o AppNodeSpec) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AppNodeSpec) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["endpoint"] = o.Endpoint
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	return toSerialize, nil
}

func (o *AppNodeSpec) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
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

	varAppNodeSpec := _AppNodeSpec{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAppNodeSpec)

	if err != nil {
		return err
	}

	*o = AppNodeSpec(varAppNodeSpec)

	return err
}

type NullableAppNodeSpec struct {
	value *AppNodeSpec
	isSet bool
}

func (v NullableAppNodeSpec) Get() *AppNodeSpec {
	return v.value
}

func (v *NullableAppNodeSpec) Set(val *AppNodeSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableAppNodeSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableAppNodeSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAppNodeSpec(val *AppNodeSpec) *NullableAppNodeSpec {
	return &NullableAppNodeSpec{value: val, isSet: true}
}

func (v NullableAppNodeSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAppNodeSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


