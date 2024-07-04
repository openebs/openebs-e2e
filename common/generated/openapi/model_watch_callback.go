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

// checks if the WatchCallback type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WatchCallback{}

// WatchCallback Watch Callbacks
type WatchCallback struct {
	Uri *string `json:"uri,omitempty"`
}

// NewWatchCallback instantiates a new WatchCallback object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWatchCallback() *WatchCallback {
	this := WatchCallback{}
	return &this
}

// NewWatchCallbackWithDefaults instantiates a new WatchCallback object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWatchCallbackWithDefaults() *WatchCallback {
	this := WatchCallback{}
	return &this
}

// GetUri returns the Uri field value if set, zero value otherwise.
func (o *WatchCallback) GetUri() string {
	if o == nil || IsNil(o.Uri) {
		var ret string
		return ret
	}
	return *o.Uri
}

// GetUriOk returns a tuple with the Uri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WatchCallback) GetUriOk() (*string, bool) {
	if o == nil || IsNil(o.Uri) {
		return nil, false
	}
	return o.Uri, true
}

// HasUri returns a boolean if a field has been set.
func (o *WatchCallback) HasUri() bool {
	if o != nil && !IsNil(o.Uri) {
		return true
	}

	return false
}

// SetUri gets a reference to the given string and assigns it to the Uri field.
func (o *WatchCallback) SetUri(v string) {
	o.Uri = &v
}

func (o WatchCallback) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WatchCallback) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Uri) {
		toSerialize["uri"] = o.Uri
	}
	return toSerialize, nil
}

type NullableWatchCallback struct {
	value *WatchCallback
	isSet bool
}

func (v NullableWatchCallback) Get() *WatchCallback {
	return v.value
}

func (v *NullableWatchCallback) Set(val *WatchCallback) {
	v.value = val
	v.isSet = true
}

func (v NullableWatchCallback) IsSet() bool {
	return v.isSet
}

func (v *NullableWatchCallback) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWatchCallback(val *WatchCallback) *NullableWatchCallback {
	return &NullableWatchCallback{value: val, isSet: true}
}

func (v NullableWatchCallback) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWatchCallback) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


