// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
)

/*

The GPX XML hierarchy:

gpx
    - attr: version (xsd:string) required
    - attr: creator (xsd:string) required
    name
    desc
    author
    email
    url
    urlname
    time
    keywords
    bounds
    wpt
        - attr: lat (gpx:latitudeType) required
        - attr: lon (gpx:longitudeType) required
        ele
        time
        magvar
        geoidheight
        name
        cmt
        desc
        src
        url
        urlname
        sym
        type
        fix
        sat
        hdop
        vdop
        pdop
        ageofdgpsdata
        dgpsid
    rte
        name
        cmt
        desc
        src
        url
        urlname
        number
        rtept
            - attr: lat (gpx:latitudeType) required
            - attr: lon (gpx:longitudeType) required
            ele
            time
            magvar
            geoidheight
            name
            cmt
            desc
            src
            url
            urlname
            sym
            type
            fix
            sat
            hdop
            vdop
            pdop
            ageofdgpsdata
            dgpsid
    trk
        name
        cmt
        desc
        src
        url
        urlname
        number
        trkseg
            trkpt
                - attr: lat (gpx:latitudeType) required
                - attr: lon (gpx:longitudeType) required
                ele
                time
                course
                speed
                magvar
                geoidheight
                name
                cmt
                desc
                src
                url
                urlname
                sym
                type
                fix
                sat
                hdop
                vdop
                pdop
                ageofdgpsdata
                dgpsid
*/

type gpx10Gpx struct {
	XMLName      xml.Name `xml:"gpx"`
	XMLNs        string   `xml:"xmlns,attr,omitempty"`
	XmlNsXsi     string   `xml:"xmlns:xsi,attr,omitempty"`
	XmlSchemaLoc string   `xml:"xsi:schemaLocation,attr,omitempty"`

	Version   string           `xml:"version,attr"`
	Creator   string           `xml:"creator,attr"`
	Name      string           `xml:"name,omitempty"`
	Desc      string           `xml:"desc,omitempty"`
	Author    string           `xml:"author,omitempty"`
	Email     string           `xml:"email,omitempty"`
	Url       string           `xml:"url,omitempty"`
	UrlName   string           `xml:"urlname,omitempty"`
	Time      string           `xml:"time,omitempty"`
	Keywords  string           `xml:"keywords,omitempty"`
	Bounds    *GpxBounds       `xml:"bounds"`
	Waypoints []*gpx10GpxPoint `xml:"wpt"`
	Routes    []*gpx10GpxRte   `xml:"rte"`
	Tracks    []*gpx10GpxTrk   `xml:"trk"`
}

//type gpx10GpxBounds struct {
//	//XMLName xml.Name `xml:"bounds"`
//	MinLat float64 `xml:"minlat,attr"`
//	MaxLat float64 `xml:"maxlat,attr"`
//	MinLon float64 `xml:"minlon,attr"`
//	MaxLon float64 `xml:"maxlon,attr"`
//}

//type gpx10GpxAuthor struct {
//	Name  string        `xml:"name,omitempty"`
//	Email string        `xml:"email,omitempty"`
//	Link  *gpx10GpxLink `xml:"link"`
//}

//type gpx10GpxEmail struct {
//	Id     string `xml:"id,attr"`
//	Domain string `xml:"domain,attr"`
//}

type gpx10GpxLink struct {
	Href string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

//type gpx10GpxMetadata struct {
//	XMLName xml.Name        `xml:"metadata"`
//	Name    string          `xml:"name,omitempty"`
//	Desc    string          `xml:"desc,omitempty"`
//	Author  *gpx10GpxAuthor `xml:"author,omitempty"`
//	//	Links     []GpxLink     `xml:"link"`
//	Timestamp string `xml:"time,omitempty"`
//	Keywords  string `xml:"keywords,omitempty"`
//	//	Bounds    *GpxBounds    `xml:"bounds"`
//}

//type gpx10GpxExtensions struct {
//	Bytes []byte `xml:",innerxml"`
//}

/**
 * Common struct fields for all points
 */
type gpx10GpxPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	// Position info
	Ele         NullableFloat64 `xml:"ele,omitempty"`
	Timestamp   string          `xml:"time,omitempty"`
	MagVar      string          `xml:"magvar,omitempty"`
	GeoIdHeight string          `xml:"geoidheight,omitempty"`
	// Description info
	Name  string         `xml:"name,omitempty"`
	Cmt   string         `xml:"cmt,omitempty"`
	Desc  string         `xml:"desc,omitempty"`
	Src   string         `xml:"src,omitempty"`
	Links []gpx10GpxLink `xml:"link"`
	Sym   string         `xml:"sym,omitempty"`
	Type  string         `xml:"type,omitempty"`
	// Accuracy info
	Fix           string   `xml:"fix,omitempty"`
	Sat           *int     `xml:"sat,omitempty"`
	Hdop          *float64 `xml:"hdop,omitempty"`
	Vdop          *float64 `xml:"vdop,omitempty"`
	Pdop          *float64 `xml:"pdop,omitempty"`
	AgeOfDGpsData *float64 `xml:"ageofdgpsdata,omitempty"`
	DGpsId        *int     `xml:"dgpsid,omitempty"`

	// Those two values are here for simplicity, but they are available only when this is part of a track segment (not route or waypoint)!
	Course string `xml:"course,omitempty"`
	Speed  string `speed:"speed,omitempty"`
}

type gpx10GpxRte struct {
	XMLName xml.Name `xml:"rte"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links       []Link   `xml:"link"`
	Number NullableInt      `xml:"number,omitempty"`
	Type   string           `xml:"type,omitempty"`
	Points []*gpx10GpxPoint `xml:"rtept"`
}

type gpx10GpxTrkSeg struct {
	XMLName xml.Name         `xml:"trkseg"`
	Points  []*gpx10GpxPoint `xml:"trkpt"`
}

// Trk is a GPX track
type gpx10GpxTrk struct {
	XMLName xml.Name `xml:"trk"`
	Name    string   `xml:"name,omitempty"`
	Cmt     string   `xml:"cmt,omitempty"`
	Desc    string   `xml:"desc,omitempty"`
	Src     string   `xml:"src,omitempty"`
	// TODO
	//Links    []Link   `xml:"link"`
	Number   NullableInt       `xml:"number,omitempty"`
	Type     string            `xml:"type,omitempty"`
	Segments []*gpx10GpxTrkSeg `xml:"trkseg,omitempty"`
}
