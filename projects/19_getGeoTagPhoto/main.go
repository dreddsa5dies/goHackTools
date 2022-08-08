package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Use: %v photo\n", os.Args[0])
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("can't open file: %v\n", err)
		}
		defer f.Close()

		x, err := exif.Decode(f)
		if err != nil {
			log.Printf("can't decode photo: %v\n", err)
			return
		}

		lat, long, err := x.LatLong()
		if err != nil {
			log.Printf("can't get the latitude and longitude of the photo: %v\n", err)
			return
		}

		fmt.Printf("%v\nlat:\t%v\nlong:\t%v", f.Name(), lat, long)
	}
}
