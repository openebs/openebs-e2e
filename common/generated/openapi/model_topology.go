/*
IoEngine RESTful API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: v0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the Topology type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Topology{}

// Topology node and pool topology for volumes
type Topology struct {
	NodeTopology NullableNodeTopology `json:"node_topology,omitempty"`
	PoolTopology NullablePoolTopology `json:"pool_topology,omitempty"`
}

// NewTopology instantiates a new Topology object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTopology() *Topology {
	this := Topology{}
	return &this
}

// NewTopologyWithDefaults instantiates a new Topology object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTopologyWithDefaults() *Topology {
	this := Topology{}
	return &this
}

// GetNodeTopology returns the NodeTopology field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Topology) GetNodeTopology() NodeTopology {
	if o == nil || IsNil(o.NodeTopology.Get()) {
		var ret NodeTopology
		return ret
	}
	return *o.NodeTopology.Get()
}

// GetNodeTopologyOk returns a tuple with the NodeTopology field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Topology) GetNodeTopologyOk() (*NodeTopology, bool) {
	if o == nil {
		return nil, false
	}
	return o.NodeTopology.Get(), o.NodeTopology.IsSet()
}

// HasNodeTopology returns a boolean if a field has been set.
func (o *Topology) HasNodeTopology() bool {
	if o != nil && o.NodeTopology.IsSet() {
		return true
	}

	return false
}

// SetNodeTopology gets a reference to the given NullableNodeTopology and assigns it to the NodeTopology field.
func (o *Topology) SetNodeTopology(v NodeTopology) {
	o.NodeTopology.Set(&v)
}
// SetNodeTopologyNil sets the value for NodeTopology to be an explicit nil
func (o *Topology) SetNodeTopologyNil() {
	o.NodeTopology.Set(nil)
}

// UnsetNodeTopology ensures that no value is present for NodeTopology, not even an explicit nil
func (o *Topology) UnsetNodeTopology() {
	o.NodeTopology.Unset()
}

// GetPoolTopology returns the PoolTopology field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Topology) GetPoolTopology() PoolTopology {
	if o == nil || IsNil(o.PoolTopology.Get()) {
		var ret PoolTopology
		return ret
	}
	return *o.PoolTopology.Get()
}

// GetPoolTopologyOk returns a tuple with the PoolTopology field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Topology) GetPoolTopologyOk() (*PoolTopology, bool) {
	if o == nil {
		return nil, false
	}
	return o.PoolTopology.Get(), o.PoolTopology.IsSet()
}

// HasPoolTopology returns a boolean if a field has been set.
func (o *Topology) HasPoolTopology() bool {
	if o != nil && o.PoolTopology.IsSet() {
		return true
	}

	return false
}

// SetPoolTopology gets a reference to the given NullablePoolTopology and assigns it to the PoolTopology field.
func (o *Topology) SetPoolTopology(v PoolTopology) {
	o.PoolTopology.Set(&v)
}
// SetPoolTopologyNil sets the value for PoolTopology to be an explicit nil
func (o *Topology) SetPoolTopologyNil() {
	o.PoolTopology.Set(nil)
}

// UnsetPoolTopology ensures that no value is present for PoolTopology, not even an explicit nil
func (o *Topology) UnsetPoolTopology() {
	o.PoolTopology.Unset()
}

func (o Topology) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Topology) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.NodeTopology.IsSet() {
		toSerialize["node_topology"] = o.NodeTopology.Get()
	}
	if o.PoolTopology.IsSet() {
		toSerialize["pool_topology"] = o.PoolTopology.Get()
	}
	return toSerialize, nil
}

type NullableTopology struct {
	value *Topology
	isSet bool
}

func (v NullableTopology) Get() *Topology {
	return v.value
}

func (v *NullableTopology) Set(val *Topology) {
	v.value = val
	v.isSet = true
}

func (v NullableTopology) IsSet() bool {
	return v.isSet
}

func (v *NullableTopology) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTopology(val *Topology) *NullableTopology {
	return &NullableTopology{value: val, isSet: true}
}

func (v NullableTopology) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTopology) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


