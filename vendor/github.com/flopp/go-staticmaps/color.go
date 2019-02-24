// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"fmt"
	"image/color"
	"regexp"
	"strings"
)

// ParseColorString parses hex color strings (i.e. `0xRRGGBB`, `#RRGGBB`, `0xRRGGBBAA`, `#RRGGBBAA`), and names colors (e.g. 'black', 'blue', ...)
func ParseColorString(s string) (color.Color, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	re := regexp.MustCompile(`^(0x|#)([A-Fa-f0-9]{6})$`)
	matches := re.FindStringSubmatch(s)
	if matches != nil {
		var r, g, b int
		fmt.Sscanf(matches[2], "%2x%2x%2x", &r, &g, &b)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}, nil
	}

	re = regexp.MustCompile(`^(0x|#)([A-Fa-f0-9]{8})$`)
	matches = re.FindStringSubmatch(s)
	if matches != nil {
		var r, g, b, a int
		fmt.Sscanf(matches[2], "%2x%2x%2x%2x", &r, &g, &b, &a)
		rr := float64(r) * float64(a) / 256.0
		gg := float64(g) * float64(a) / 256.0
		bb := float64(b) * float64(a) / 256.0
		return color.RGBA{uint8(rr), uint8(gg), uint8(bb), uint8(a)}, nil
	}

	switch s {
	case "black":
		return color.RGBA{0x00, 0x00, 0x00, 0xff}, nil
	case "blue":
		return color.RGBA{0x00, 0x00, 0xff, 0xff}, nil
	case "brown":
		return color.RGBA{0x96, 0x4b, 0x00, 0xff}, nil
	case "green":
		return color.RGBA{0x00, 0xff, 0x00, 0xff}, nil
	case "orange":
		return color.RGBA{0xff, 0x7f, 0x00, 0xff}, nil
	case "purple":
		return color.RGBA{0x7f, 0x00, 0x7f, 0xff}, nil
	case "red":
		return color.RGBA{0xff, 0x00, 0x00, 0xff}, nil
	case "yellow":
		return color.RGBA{0xff, 0xff, 0x00, 0xff}, nil
	case "white":
		return color.RGBA{0xff, 0xff, 0xff, 0xff}, nil
	case "transparent":
		return color.RGBA{0x00, 0x00, 0x00, 0x00}, nil
	}
	return color.Transparent, fmt.Errorf("Cannot parse color string: %s", s)
}

// Luminance computes the luminance (~ brightness) of the given color. Range: 0.0 for black to 1.0 for white.
func Luminance(col color.Color) float64 {
	r, g, b, _ := col.RGBA()
	return (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / float64(0xffff)
}
