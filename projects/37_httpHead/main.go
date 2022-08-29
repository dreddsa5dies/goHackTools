// Perform an HTTP HEAD request on a URL and print out headers

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
	}

	urlArg := os.Args[1]

	u, err := url.ParseRequestURI(urlArg)
	if err != nil {
		log.Fatal("error check URL. ", err)
	}

	response, err := http.Head(u.String())
	if err != nil {
		log.Println("error fetching URL. ", err)
		return
	}
	defer response.Body.Close()

	// Print out each header key and value pair
	for key, value := range response.Header {
		fmt.Printf("%s: %s\n", key, value[0])
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v URL\n", name)
	fmt.Printf("Perform an HTTP HEAD request to a URL\n")
	os.Exit(1)
}
