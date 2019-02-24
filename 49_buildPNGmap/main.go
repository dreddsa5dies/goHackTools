package main

import (
	"image/color"

	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(600, 600)

	// get geotag from https://geocode-maps.yandex.ru )) - Москва, Тверская, 6
	ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(55.757926, 37.607242), color.RGBA{0xff, 0, 0, 0xff}, 16.0))
	// and Home
	ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(55.919609, 37.742699), color.RGBA{0xff, 0, 0, 0xff}, 10.0))
	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		panic(err)
	}
}
