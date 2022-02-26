// Finding HTML comments in a web page
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
	}

	url := os.Args[1]
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching URL. ", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	// Look for HTML comments using a regular expressionre := regexp.MustCompile("<!--(.|\n)*?-->")
	re := regexp.MustCompile("<!--(.|\n)*?-->")
	matches := re.FindAllString(string(body), -1)
	if matches == nil {
		// Clean exit if no matches found
		fmt.Println("No HTML comments found.")
		os.Exit(0)
	}

	// Print all HTML comments found
	for _, match := range matches {
		fmt.Println(match)
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v URL\n", name)
	fmt.Printf("Search for HTML comments in a URL\n")
	os.Exit(1)
}
