// Extracting links from url to Maltego

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jackdanger/collectlinks"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Extracting links from url to Maltego.")
		fmt.Println("Usage: " + os.Args[0] + "<url>")
		fmt.Println("U need url: http://domain.com")
		os.Exit(1)
	}

	URL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// response url
	resp, err := http.Get(URL.String())
	if err != nil {
		log.Fatalf("error get %v", err)
	}
	// отложенное закрытие коннекта
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	// print Maltego transform response
	fmt.Println("<MaltegoMessage>")
	fmt.Println("<MaltegoTransformResponseMessage>")
	fmt.Println("  <Entities>")

	for _, link := range links {
		switch {
		case strings.HasPrefix(link, "http"):
			// if the first characters of the link are HTTP (or HTTPS), we will output it into
			// correct format as an entity Maltego
			fmt.Println("<Entity Type=\"maltego.Domain\">")
			fmt.Println("<Value>" + link + "</Value>")
		case strings.HasPrefix(link, "https"):
			fmt.Println("<Entity Type=\"maltego.Domain\">")
			fmt.Println("<Value>" + link + "</Value>")
		}
	}

	fmt.Println("  </Entities>")
	fmt.Println("</MaltegoTransformResponseMessage>")
	fmt.Println("</MaltegoMessage>")
}
