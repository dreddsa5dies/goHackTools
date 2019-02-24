// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sm (~ static maps) renders static map images from OSM tiles with markers, paths, and filled areas.
package sm

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// Context holds all information about the map image that is to be rendered
type Context struct {
	width  int
	height int

	hasZoom bool
	zoom    int

	hasCenter bool
	center    s2.LatLng

	hasBoundingBox bool
	boundingBox    s2.Rect

	background color.Color

	markers  []*Marker
	paths    []*Path
	areas    []*Area
	circles  []*Circle
	overlays []*TileProvider

	userAgent    string
	tileProvider *TileProvider

	overrideAttribution *string
}

// NewContext creates a new instance of Context
func NewContext() *Context {
	t := new(Context)
	t.width = 512
	t.height = 512
	t.hasZoom = false
	t.hasCenter = false
	t.hasBoundingBox = false
	t.background = nil
	t.userAgent = ""
	t.tileProvider = NewTileProviderOpenStreetMaps()
	return t
}

// SetTileProvider sets the TileProvider to be used
func (m *Context) SetTileProvider(t *TileProvider) {
	m.tileProvider = t
}

// SetUserAgent sets the HTTP user agent string used when downloading map tiles
func (m *Context) SetUserAgent(a string) {
	m.userAgent = a
}

// SetSize sets the size of the generated image
func (m *Context) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetZoom sets the zoom level
func (m *Context) SetZoom(zoom int) {
	m.zoom = zoom
	m.hasZoom = true
}

// SetCenter sets the center coordinates
func (m *Context) SetCenter(center s2.LatLng) {
	m.center = center
	m.hasCenter = true
}

// SetBoundingBox sets the bounding box
func (m *Context) SetBoundingBox(bbox s2.Rect) {
	m.boundingBox = bbox
	m.hasBoundingBox = true
}

// SetBackground sets the background color (used as a fallback for areas without map tiles)
func (m *Context) SetBackground(col color.Color) {
	m.background = col
}

// AddMarker adds a marker to the Context
func (m *Context) AddMarker(marker *Marker) {
	m.markers = append(m.markers, marker)
}

// ClearMarkers removes all markers from the Context
func (m *Context) ClearMarkers() {
	m.markers = nil
}

// AddPath adds a path to the Context
func (m *Context) AddPath(path *Path) {
	m.paths = append(m.paths, path)
}

// ClearPaths removes all paths from the Context
func (m *Context) ClearPaths() {
	m.paths = nil
}

// AddArea adds an area to the Context
func (m *Context) AddArea(area *Area) {
	m.areas = append(m.areas, area)
}

// ClearAreas removes all areas from the Context
func (m *Context) ClearAreas() {
	m.areas = nil
}

// AddCircle adds an circle to the Context
func (m *Context) AddCircle(circle *Circle) {
	m.circles = append(m.circles, circle)
}

// ClearCircles removes all circles from the Context
func (m *Context) ClearCircles() {
	m.circles = nil
}

// AddOverlay adds an overlay to the Context
func (m *Context) AddOverlay(overlay *TileProvider) {
	m.overlays = append(m.overlays, overlay)
}

// ClearOverlays removes all overlays from the Context
func (m *Context) ClearOverlays() {
	m.overlays = nil
}

// OverrideAttribution sets a custom attribution string (or none if empty)
//
// Pay attention you might be violating the terms of usage for the
// selected map provider - only use the function if you are aware of this!
func (m *Context) OverrideAttribution(attribution string) {
	m.overrideAttribution = &attribution
}

func (m *Context) determineBounds() s2.Rect {
	r := s2.EmptyRect()
	for _, marker := range m.markers {
		r = r.Union(marker.bounds())
	}
	for _, path := range m.paths {
		r = r.Union(path.bounds())
	}
	for _, area := range m.areas {
		r = r.Union(area.bounds())
	}
	for _, circle := range m.circles {
		r = r.Union(circle.bounds())
	}
	return r
}

func (m *Context) determineExtraMarginPixels() float64 {
	p := 0.0
	for _, marker := range m.markers {
		if pp := marker.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	for _, path := range m.paths {
		if pp := path.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	for _, area := range m.areas {
		if pp := area.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	for _, circle := range m.circles {
		if pp := circle.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	return p
}

func (m *Context) determineZoom(bounds s2.Rect, center s2.LatLng) int {
	b := bounds.AddPoint(center)
	if b.IsEmpty() || b.IsPoint() {
		return 15
	}

	tileSize := m.tileProvider.TileSize
	margin := 4.0 + m.determineExtraMarginPixels()
	w := (float64(m.width) - 2.0*margin) / float64(tileSize)
	h := (float64(m.height) - 2.0*margin) / float64(tileSize)
	minX := (b.Lo().Lng.Degrees() + 180.0) / 360.0
	maxX := (b.Hi().Lng.Degrees() + 180.0) / 360.0
	minY := (1.0 - math.Log(math.Tan(b.Lo().Lat.Radians())+(1.0/math.Cos(b.Lo().Lat.Radians())))/math.Pi) / 2.0
	maxY := (1.0 - math.Log(math.Tan(b.Hi().Lat.Radians())+(1.0/math.Cos(b.Hi().Lat.Radians())))/math.Pi) / 2.0

	dx := maxX - minX
	for dx < 0 {
		dx = dx + 1
	}
	for dx > 1 {
		dx = dx - 1
	}
	dy := math.Abs(maxY - minY)

	zoom := 1
	for zoom < 30 {
		tiles := float64(uint(1) << uint(zoom))
		if dx*tiles > w || dy*tiles > h {
			return zoom - 1
		}
		zoom = zoom + 1
	}

	return 15
}

func (m *Context) determineZoomCenter() (int, s2.LatLng, error) {
	bounds := m.determineBounds()
	if m.hasBoundingBox && !m.boundingBox.IsEmpty() {
		center := m.boundingBox.Center()
		return m.determineZoom(m.boundingBox, center), center, nil
	} else if m.hasCenter {
		if m.hasZoom {
			return m.zoom, m.center, nil
		}
		return m.determineZoom(bounds, m.center), m.center, nil
	} else if !bounds.IsEmpty() {
		center := bounds.Center()
		if m.hasZoom {
			return m.zoom, center, nil
		}
		return m.determineZoom(bounds, center), center, nil
	}

	return 0, s2.LatLngFromDegrees(0, 0), errors.New("Cannot determine map extent: no center coordinates given, no bounding box given, no content (markers, paths, areas) given")
}

type transformer struct {
	zoom               int
	numTiles           float64 // number of tiles per dimension at this zoom level
	tileSize           int     // tile size in pixels from this provider
	pWidth, pHeight    int     // pixel size of returned set of tiles
	pCenterX, pCenterY int     // pixel location of requested center in set of tiles
	tCountX, tCountY   int     // download area in tile units
	tCenterX, tCenterY float64 // tile index to requested center
	tOriginX, tOriginY int     // bottom left tile to download
	pMinX, pMaxX       int
}

func newTransformer(width int, height int, zoom int, llCenter s2.LatLng, tileSize int) *transformer {
	t := new(transformer)

	t.zoom = zoom
	t.numTiles = math.Exp2(float64(t.zoom))
	t.tileSize = tileSize

	// fractional tile index to center of requested area
	t.tCenterX, t.tCenterY = t.ll2t(llCenter)

	ww := float64(width) / float64(tileSize)
	hh := float64(height) / float64(tileSize)

	// origin tile to fulfill request
	t.tOriginX = int(math.Floor(t.tCenterX - 0.5*ww))
	t.tOriginY = int(math.Floor(t.tCenterY - 0.5*hh))

	// tiles in each axis to fulfill request
	t.tCountX = 1 + int(math.Floor(t.tCenterX+0.5*ww)) - t.tOriginX
	t.tCountY = 1 + int(math.Floor(t.tCenterY+0.5*hh)) - t.tOriginY

	// final pixel dimensions of area returned
	t.pWidth = t.tCountX * tileSize
	t.pHeight = t.tCountY * tileSize

	// Pixel location in returned image for center of requested area
	t.pCenterX = int((t.tCenterX - float64(t.tOriginX)) * float64(tileSize))
	t.pCenterY = int((t.tCenterY - float64(t.tOriginY)) * float64(tileSize))

	t.pMinX = t.pCenterX - width/2
	t.pMaxX = t.pMinX + width

	return t
}

// ll2t returns fractional tile index for a lat/lng points
func (t *transformer) ll2t(ll s2.LatLng) (float64, float64) {
	x := t.numTiles * (ll.Lng.Degrees() + 180.0) / 360.0
	y := t.numTiles * (1 - math.Log(math.Tan(ll.Lat.Radians())+(1.0/math.Cos(ll.Lat.Radians())))/math.Pi) / 2.0
	return x, y
}

func (t *transformer) ll2p(ll s2.LatLng) (float64, float64) {
	x, y := t.ll2t(ll)
	x = float64(t.pCenterX) + (x-t.tCenterX)*float64(t.tileSize)
	y = float64(t.pCenterY) + (y-t.tCenterY)*float64(t.tileSize)

	offset := t.numTiles * float64(t.tileSize)
	if x < float64(t.pMinX) {
		for x < float64(t.pMinX) {
			x = x + offset
		}
	} else if x >= float64(t.pMaxX) {
		for x >= float64(t.pMaxX) {
			x = x - offset
		}
	}
	return x, y
}

// Rect returns an s2.Rect bounding box around the set of tiles described by transformer
func (t *transformer) Rect() (bbox s2.Rect) {
	// transform from https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Go
	invNumTiles := 1.0 / t.numTiles
	// Get latitude bounds
	n := math.Pi - 2.0*math.Pi*float64(t.tOriginY)*invNumTiles
	bbox.Lat.Hi = math.Atan(0.5 * (math.Exp(n) - math.Exp(-n)))
	n = math.Pi - 2.0*math.Pi*float64(t.tOriginY+t.tCountY)*invNumTiles
	bbox.Lat.Lo = math.Atan(0.5 * (math.Exp(n) - math.Exp(-n)))
	// Get longtitude bounds, much easier
	bbox.Lng.Lo = float64(t.tOriginX)*invNumTiles*2.0*math.Pi - math.Pi
	bbox.Lng.Hi = float64(t.tOriginX+t.tCountX)*invNumTiles*2.0*math.Pi - math.Pi
	return bbox
}

// Render actually renders the map image including all map objects (markers, paths, areas)
func (m *Context) Render() (image.Image, error) {
	zoom, center, err := m.determineZoomCenter()
	if err != nil {
		return nil, err
	}

	tileSize := m.tileProvider.TileSize
	trans := newTransformer(m.width, m.height, zoom, center, tileSize)
	img := image.NewRGBA(image.Rect(0, 0, trans.pWidth, trans.pHeight))
	gc := gg.NewContextForRGBA(img)
	if m.background != nil {
		draw.Draw(img, img.Bounds(), &image.Uniform{m.background}, image.ZP, draw.Src)
	}

	// fetch and draw tiles to img
	layers := []*TileProvider{m.tileProvider}
	if m.overlays != nil {
		layers = append(layers, m.overlays...)
	}

	for _, layer := range layers {
		if err := m.renderLayer(gc, zoom, trans, tileSize, layer); err != nil {
			return nil, err
		}
	}

	// draw map objects
	for _, area := range m.areas {
		area.draw(gc, trans)
	}
	for _, path := range m.paths {
		path.draw(gc, trans)
	}
	for _, marker := range m.markers {
		marker.draw(gc, trans)
	}
	for _, circle := range m.circles {
		circle.draw(gc, trans)
	}

	// crop image
	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{trans.pCenterX - int(m.width)/2, trans.pCenterY - int(m.height)/2},
		draw.Src)

	attribution := m.tileProvider.Attribution
	if m.overrideAttribution != nil {
		attribution = *m.overrideAttribution
	}

	// draw attribution
	if attribution == "" {
		return croppedImg, nil
	}
	_, textHeight := gc.MeasureString(attribution)
	boxHeight := textHeight + 4.0
	gc = gg.NewContextForRGBA(croppedImg)
	gc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	gc.DrawRectangle(0.0, float64(m.height)-boxHeight, float64(m.width), boxHeight)
	gc.Fill()
	gc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	gc.DrawString(attribution, 4.0, float64(m.height)-4.0)

	return croppedImg, nil
}

// RenderWithBounds actually renders the map image including all map objects (markers, paths, areas).
// The returned image covers requested area as well as any tiles necessary to cover that area, which may
// be larger than the request.
//
// Specific bounding box of returned image is provided to support image registration with other data
func (m *Context) RenderWithBounds() (image.Image, s2.Rect, error) {
	zoom, center, err := m.determineZoomCenter()
	if err != nil {
		return nil, s2.Rect{}, err
	}

	tileSize := m.tileProvider.TileSize
	trans := newTransformer(m.width, m.height, zoom, center, tileSize)
	img := image.NewRGBA(image.Rect(0, 0, trans.pWidth, trans.pHeight))
	gc := gg.NewContextForRGBA(img)
	if m.background != nil {
		draw.Draw(img, img.Bounds(), &image.Uniform{m.background}, image.ZP, draw.Src)
	}

	// fetch and draw tiles to img
	layers := []*TileProvider{m.tileProvider}
	if m.overlays != nil {
		layers = append(layers, m.overlays...)
	}

	for _, layer := range layers {
		if err := m.renderLayer(gc, zoom, trans, tileSize, layer); err != nil {
			return nil, s2.Rect{}, err
		}
	}

	// draw map objects
	for _, area := range m.areas {
		area.draw(gc, trans)
	}
	for _, path := range m.paths {
		path.draw(gc, trans)
	}
	for _, circle := range m.circles {
		circle.draw(gc, trans)
	}
	for _, marker := range m.markers {
		marker.draw(gc, trans)
	}

	// draw attribution
	if m.tileProvider.Attribution == "" {
		return img, trans.Rect(), nil
	}
	_, textHeight := gc.MeasureString(m.tileProvider.Attribution)
	boxHeight := textHeight + 4.0
	gc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	gc.DrawRectangle(0.0, float64(trans.pHeight)-boxHeight, float64(trans.pWidth), boxHeight)
	gc.Fill()
	gc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	gc.DrawString(m.tileProvider.Attribution, 4.0, float64(m.height)-4.0)

	return img, trans.Rect(), nil
}

func (m *Context) renderLayer(gc *gg.Context, zoom int, trans *transformer, tileSize int, provider *TileProvider) error {
	t := NewTileFetcher(provider)
	if m.userAgent != "" {
		t.SetUserAgent(m.userAgent)
	}

	tiles := (1 << uint(zoom))
	for xx := 0; xx < trans.tCountX; xx++ {
		x := trans.tOriginX + xx
		if x < 0 {
			x = x + tiles
		} else if x >= tiles {
			x = x - tiles
		}
		for yy := 0; yy < trans.tCountY; yy++ {
			y := trans.tOriginY + yy
			if y < 0 || y >= tiles {
				log.Printf("Skipping out of bounds tile %d/%d", x, y)
			} else {
				if tileImg, err := t.Fetch(zoom, x, y); err == nil {
					gc.DrawImage(tileImg, xx*tileSize, yy*tileSize)
				} else {
					log.Printf("Error downloading tile file: %s", err)
				}
			}
		}
	}

	return nil
}
