/*cookie-flags takes a url and returns the cookie set.*/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input string `short:"i" long:"input" default:"" description:"MAC address"`
}

// Coordinates - struct for xml data
type Coordinates struct {
	Latitude   string `xml:"latitude,attr"`
	Longitude  string `xml:"longitude,attr"`
	NLatitude  string `xml:"nlatitude,attr"`
	NLongitude string `xml:"nlongitude,attr"`
}

type location struct {
	XMLName  xml.Name    `xml:"location"`
	Location Coordinates `xml:"coordinates"`
}

func main() {
	flags.Parse(&opts)

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		os.Exit(1)
	} else if os.Args[1] == "-h" || os.Args[1] != "-i" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		os.Exit(1)
	}

	log.Println("Start")

	// format MAC:-> bc988929e608
	opts.Input = "bc988929e608"
	url := fmt.Sprintf(`http://mobile.maps.yandex.net/cellid_location/?clid=1866854&lac=-1&cellid=-1&operatorid=null&countrycode=null&signalstrength=-1&wifinetworks=%s:-65&app=ymetro`, opts.Input)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := location{}
	err = xml.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

	// header & footer KML
	kmlheader := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
<Document>`
	kmlfooter := `</Document>
</kml>`
	kmlData := recKml(opts.Input, data.Location.Longitude, data.Location.Latitude)
	fmt.Println(kmlheader + "\n" + kmlData + kmlfooter)

	log.Println("End")
	os.Exit(0)
}

func recKml(name, long, lat string) string {
	kml := fmt.Sprintf(`<Placemark>
		<name>%s</name>
		<Point>
		<coordinates>%s,%s</coordinates>
		</Point>
		</Placemark>`, name, long, lat)
	return kml
}
