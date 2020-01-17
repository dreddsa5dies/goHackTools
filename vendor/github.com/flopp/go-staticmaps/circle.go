package sm

import (
	"image/color"
	"log"
	"math"
	"strings"

	"strconv"

	"github.com/flopp/go-coordsparser"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Circle represents a circle on the map
type Circle struct {
	MapObject
	Position s2.LatLng
	Color    color.Color
	Fill     color.Color
	Weight   float64
	Radius   float64 // in m.
}

// NewCircle creates a new circle
func NewCircle(pos s2.LatLng, col, fill color.Color, radius, weight float64) *Circle {
	return &Circle{
		Position: pos,
		Color:    col,
		Fill:     fill,
		Weight:   weight,
		Radius:   radius,
	}
}

// ParseCircleString parses a string and returns an array of circles
func ParseCircleString(s string) (circles []*Circle, err error) {
	circles = make([]*Circle, 0, 0)

	var col color.Color = color.RGBA{0xff, 0, 0, 0xff}
	var fill color.Color = color.Transparent
	radius := 100.0
	weight := 5.0

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "color:"); ok {
			col, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "fill:"); ok {
			fill, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "radius:"); ok {
			if radius, err = strconv.ParseFloat(suffix, 64); err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "weight:"); ok {
			if weight, err = strconv.ParseFloat(suffix, 64); err != nil {
				return nil, err
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			c := NewCircle(s2.LatLngFromDegrees(lat, lng), col, fill, radius, weight)
			circles = append(circles, c)
		}
	}
	return circles, nil
}

func (m *Circle) getLatLng(plus bool) s2.LatLng {
	const (
		R = 6371000.0
	)
	th := m.Radius / R
	br := 0 / float64(s1.Degree)
	if !plus {
		th *= -1
	}
	lat := m.Position.Lat.Radians()
	lat1 := math.Asin(math.Sin(lat)*math.Cos(th) + math.Cos(lat)*math.Sin(th)*math.Cos(br))
	lng1 := m.Position.Lng.Radians() +
		math.Atan2(math.Sin(br)*math.Sin(th)*math.Cos(lat),
			math.Cos(th)-math.Sin(lat)*math.Sin(lat1))
	return s2.LatLng{
		Lat: s1.Angle(lat1),
		Lng: s1.Angle(lng1),
	}
}

func (m *Circle) extraMarginPixels() float64 {
	return 0.5 * m.Weight
}

func (m *Circle) bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(m.getLatLng(false))
	r = r.AddPoint(m.getLatLng(true))
	return r
}

func (m *Circle) draw(gc *gg.Context, trans *transformer) {
	if !CanDisplay(m.Position) {
		log.Printf("Circle coordinates not displayable: %f/%f", m.Position.Lat.Degrees(), m.Position.Lng.Degrees())
		return
	}

	ll := m.getLatLng(true)
	x, y := trans.ll2p(m.Position)
	x1, y1 := trans.ll2p(ll)
	radius := math.Sqrt(math.Pow(x1-x, 2) + math.Pow(y1-y, 2))
	gc.ClearPath()
	gc.SetLineWidth(m.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	gc.DrawCircle(x, y, radius)
	gc.SetColor(m.Fill)
	gc.FillPreserve()
	gc.SetColor(m.Color)
	gc.Stroke()
}
