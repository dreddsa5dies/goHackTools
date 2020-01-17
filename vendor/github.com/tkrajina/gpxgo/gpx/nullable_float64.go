// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

//NullableFloat64 implements a nullable float64
type NullableFloat64 struct {
	data    float64
	notNull bool
}

//Null checks if value is null
func (n *NullableFloat64) Null() bool {
	return !n.notNull
}

//NotNull checks if value is not null
func (n *NullableFloat64) NotNull() bool {
	return n.notNull
}

//Value returns the value
func (n *NullableFloat64) Value() float64 {
	return n.data
}

//SetValue sets the value
func (n *NullableFloat64) SetValue(data float64) {
	n.data = data
	n.notNull = true
}

//SetNull sets the value to null
func (n *NullableFloat64) SetNull() {
	var defaultValue float64
	n.data = defaultValue
	n.notNull = false
}

//NewNullableFloat64 creates a new NullableFloat64
func NewNullableFloat64(data float64) *NullableFloat64 {
	result := new(NullableFloat64)
	result.data = data
	result.notNull = true
	return result
}

//UnmarshalXML implements xml unmarshalling
func (n *NullableFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		n.SetNull()
		return nil
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		value, err := strconv.ParseFloat(strData, 64)
		if err != nil {
			n.SetNull()
			return nil
		}
		n.SetValue(value)
	}
	d.Skip()
	return nil
}

//UnmarshalXMLAttr implements xml attribute unmarshalling
func (n *NullableFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	strData := strings.Trim(string(attr.Value), " ")
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		n.SetNull()
		return nil
	}
	n.SetValue(value)
	return nil
}

//MarshalXML implements xml marshalling
func (n NullableFloat64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Null() {
		return nil
	}
	xmlName := xml.Name{Local: start.Name.Local}
	if err := e.EncodeToken(xml.StartElement{Name: xmlName}); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData([]byte(fmt.Sprintf("%g", n.Value())))); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.EndElement{Name: xmlName}); err != nil {
		return err
	}
	return nil
}

//MarshalXMLAttr implements xml attribute marshalling
func (n NullableFloat64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var result xml.Attr
	if n.Null() {
		return result, nil
	}
	return xml.Attr{
			Name:  xml.Name{Local: name.Local},
			Value: fmt.Sprintf("%g", n.Value())},
		nil
}
