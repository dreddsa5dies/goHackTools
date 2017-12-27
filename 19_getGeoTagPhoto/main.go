package main

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Use: %v photo\n", os.Args[0])
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "file: %v\n", err)
			os.Exit(1)
		}

		x, err := exif.Decode(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Decode: %v\n", err)
			os.Exit(1)
		}

		lat, long, err := x.LatLong()
		if err != nil {
			fmt.Fprintf(os.Stderr, "LatLong: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stdout, f.Name())
			fmt.Fprintln(os.Stdout, fmt.Sprintf("lat:\t%v\nlong:\t%v", lat, long))
			os.Exit(0)
		}
	}
}
