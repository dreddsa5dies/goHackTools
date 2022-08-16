/*
API not work today, so sorry
cookie-flags takes a url and returns the cookie set.
*/
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

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
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v MAC-addr...\n", os.Args[0])
		fmt.Printf("Output format: MAC Longitude Latitude\n")
		os.Exit(1)
	}

	regMac := regexp.MustCompile(`([([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)

	for k := 1; k < len(os.Args); k++ {
		if !regMac.MatchString(os.Args[k]) {
			fmt.Fprintf(os.Stderr, "%s not found? please get MAC format 00:00:00:00:00:00 or 00-00-00-00-00-00\n", os.Args[k])
			continue
		}

		data, err := getLocation(macFormat(os.Args[k]))
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(os.Args[k], data.Location.Longitude, data.Location.Latitude)
	}
}

func macFormat(mac string) string {
	mac = strings.ToLower(mac)
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")

	return mac
}

func getLocation(mac string) (*location, error) {
	params := url.Values{}

	params.Set("wifinetworks", mac+":-65")

	urlPath := "http://mobile.maps.yandex.net/cellid_location/?" + params.Encode()

	resp, err := http.Get(urlPath) //nolint:gosec // так надо
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := location{}

	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
