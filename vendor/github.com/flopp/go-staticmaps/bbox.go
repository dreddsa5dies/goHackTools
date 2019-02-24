// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"fmt"
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// CreateBBox creates a bounding box from a north-western point
// (lat/lng in degrees) and a south-eastern point (lat/lng in degrees).
// Note that you can create a bounding box wrapping over the antimeridian at
// lng=+-/180° by nwlng > selng.
func CreateBBox(nwlat float64, nwlng float64, selat float64, selng float64) (*s2.Rect, error) {
	if nwlat < -90 || nwlat > 90 {
		return nil, fmt.Errorf("Out of range nwlat (%f) must be in [-90, 90]", nwlat)
	}
	if nwlng < -180 || nwlng > 180 {
		return nil, fmt.Errorf("Out of range nwlng (%f) must be in [-180, 180]", nwlng)
	}

	if selat < -90 || selat > 90 {
		return nil, fmt.Errorf("Out of range selat (%f) must be in [-90, 90]", selat)
	}
	if selng < -180 || selng > 180 {
		return nil, fmt.Errorf("Out of range selng (%f) must be in [-180, 180]", selng)
	}

	if nwlat == selat {
		return nil, fmt.Errorf("nwlat and selat must not be equal")
	}
	if nwlng == selng {
		return nil, fmt.Errorf("nwlng and selng must not be equal")
	}

	bbox := new(s2.Rect)
	if selat < nwlat {
		bbox.Lat.Lo = selat * math.Pi / 180.0
		bbox.Lat.Hi = nwlat * math.Pi / 180.0
	} else {
		bbox.Lat.Lo = nwlat * math.Pi / 180.0
		bbox.Lat.Hi = selat * math.Pi / 180.0
	}
	bbox.Lng = s1.IntervalFromEndpoints(nwlng*math.Pi/180.0, selng*math.Pi/180.0)

	return bbox, nil
}
