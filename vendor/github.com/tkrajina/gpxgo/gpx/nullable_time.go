// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import "time"

//NullableTime implements a nullable time
type NullableTime struct {
	data    time.Time
	notNull bool
}

//Null checks if value is null
func (n *NullableTime) Null() bool {
	return !n.notNull
}

//NotNull checks if value is not null
func (n *NullableTime) NotNull() bool {
	return n.notNull
}

//Value returns the value
func (n *NullableTime) Value() time.Time {
	return n.data
}

//SetValue sets the value
func (n *NullableTime) SetValue(data time.Time) {
	n.data = data
	n.notNull = true
}

//SetNull sets the value to null
func (n *NullableTime) SetNull() {
	var defaultValue time.Time
	n.data = defaultValue
	n.notNull = false
}

//NewNullableTime creates a new NullableTime
func NewNullableTime(data time.Time) *NullableTime {
	result := new(NullableTime)
	result.data = data
	result.notNull = true
	return result
}
