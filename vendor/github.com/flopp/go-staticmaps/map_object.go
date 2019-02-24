// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// MapObject is the interface for all objects on the map
type MapObject interface {
	bounds() s2.Rect
	extraMarginPixels() float64
	draw(dc *gg.Context, trans *transformer)
}

// CanDisplay checks if pos is generally displayable (i.e. its latitude is in [-85,85])
func CanDisplay(pos s2.LatLng) bool {
	const minLatitude float64 = -85.0
	const maxLatitude float64 = 85.0
	return (minLatitude <= pos.Lat.Degrees()) && (pos.Lat.Degrees() <= maxLatitude)
}
