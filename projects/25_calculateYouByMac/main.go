/*cookie-flags takes a url and returns the cookie set.*/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	regMac, err := regexp.Compile(`([([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)
	if err != nil {
		log.Fatalln(err)
	}

	for k := 1; k < len(os.Args); k++ {
		if !regMac.MatchString(os.Args[k]) {
			fmt.Fprintf(os.Stderr, "%s not found? please get MAC format 00:00:00:00:00:00 or 00-00-00-00-00-00\n", os.Args[k])
			continue
		}

		macAdress := macFormat(os.Args[k])

		url := fmt.Sprintf(`http://mobile.maps.yandex.net/cellid_location/?clid=1866854&lac=-1&cellid=-1&operatorid=null&countrycode=null&signalstrength=-1&wifinetworks=%s:-65&app=ymetro`, macAdress)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("http.Get: %s", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ioutil.ReadAll: %s", err)
		}

		data := location{}
		err = xml.Unmarshal(body, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s not found\n", os.Args[k])
			continue
		}

		fmt.Println(os.Args[k], data.Location.Longitude, data.Location.Latitude)
	}
	os.Exit(0)
}

func macFormat(mac string) string {
	mac = strings.ToLower(mac)
	mac = strings.Replace(mac, ":", "", -1)
	mac = strings.Replace(mac, "-", "", -1)
	return mac
}
