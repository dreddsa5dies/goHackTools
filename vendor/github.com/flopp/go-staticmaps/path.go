// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tkrajina/gpxgo/gpx"
)

// Path represents a path or area on the map
type Path struct {
	MapObject
	Positions []s2.LatLng
	Color     color.Color
	Weight    float64
}

// NewPath creates a new Path
func NewPath(positions []s2.LatLng, col color.Color, weight float64) *Path {
	p := new(Path)
	p.Positions = positions
	p.Color = col
	p.Weight = weight

	return p
}

// ParsePathString parses a string and returns a path
func ParsePathString(s string) ([]*Path, error) {
	paths := make([]*Path, 0, 0)
	currentPath := new(Path)
	currentPath.Color = color.RGBA{0xff, 0, 0, 0xff}
	currentPath.Weight = 5.0

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "color:"); ok {
			var err error
			if currentPath.Color, err = ParseColorString(suffix); err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "weight:"); ok {
			var err error
			if currentPath.Weight, err = strconv.ParseFloat(suffix, 64); err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "gpx:"); ok {
			gpxData, err := gpx.ParseFile(suffix)
			if err != nil {
				return nil, err
			}
			for _, trk := range gpxData.Tracks {
				for _, seg := range trk.Segments {
					p := new(Path)
					p.Color = currentPath.Color
					p.Weight = currentPath.Weight
					for _, pt := range seg.Points {
						p.Positions = append(p.Positions, s2.LatLngFromDegrees(pt.GetLatitude(), pt.GetLongitude()))
					}
					if len(p.Positions) > 0 {
						paths = append(paths, p)
					}
				}
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			currentPath.Positions = append(currentPath.Positions, s2.LatLngFromDegrees(lat, lng))
		}
	}
	if len(currentPath.Positions) > 0 {
		paths = append(paths, currentPath)
	}
	return paths, nil
}

func (p *Path) extraMarginPixels() float64 {
	return 0.5 * p.Weight
}

func (p *Path) bounds() s2.Rect {
	r := s2.EmptyRect()
	for _, ll := range p.Positions {
		r = r.AddPoint(ll)
	}
	return r
}

func (p *Path) draw(gc *gg.Context, trans *transformer) {
	if len(p.Positions) <= 1 {
		return
	}

	gc.ClearPath()
	gc.SetLineWidth(p.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	for _, ll := range p.Positions {
		gc.LineTo(trans.ll2p(ll))
	}
	gc.SetColor(p.Color)
	gc.Stroke()
}
