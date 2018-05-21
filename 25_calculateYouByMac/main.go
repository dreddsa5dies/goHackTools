/*cookie-flags takes a url and returns the cookie set.*/
package main

import (
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

	fmt.Println(string(body))

	log.Println("End")
	os.Exit(0)
}

func recKml(name string, long, lat float64) string {
	kml := fmt.Sprintf(`<Placemark>
		<name>%s</name>
		<Point>
		<coordinates>%6f,%6f</coordinates>
		</Point>
		</Placemark>`, name, long, lat)
	return kml
}
