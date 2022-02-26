// Perform an HTTP HEAD request on a URL and print out headers

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
	}

	url := os.Args[1]

	// Perform HTTP HEAD
	response, err := http.Head(url)
	if err != nil {
		log.Fatal("Error fetching URL. ", err)
	}

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
