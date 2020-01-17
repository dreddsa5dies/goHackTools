// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

//NullableString implements a nullable string
type NullableString struct {
	data    string
	notNull bool
}

//Null checks if value is null
func (n *NullableString) Null() bool {
	return !n.notNull
}

//NotNull checks if value is not null
func (n *NullableString) NotNull() bool {
	return n.notNull
}

//Value returns the value
func (n *NullableString) Value() string {
	return n.data
}

//SetValue sets the value
func (n *NullableString) SetValue(data string) {
	n.data = data
	n.notNull = true
}

//SetNull sets the value to null
func (n *NullableString) SetNull() {
	var defaultValue string
	n.data = defaultValue
	n.notNull = false
}

//NewNullableString creates a new NullableString
func NewNullableString(data string) *NullableString {
	result := new(NullableString)
	result.data = data
	result.notNull = true
	return result
}
