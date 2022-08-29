// Finding HTML comments in a web page
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

	response, err := http.Get(u.String())
	if err != nil {
		log.Println("error fetching URL. ", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading HTTP body. ", err)
		return
	}

	// Look for HTML comments using a regular expressions := regexp.MustCompile("<!--(.|\n)*?-->")
	re := regexp.MustCompile("<!--(.|\n)*?-->")

	matches := re.FindAllString(string(body), -1)
	if matches == nil {
		// Clean exit if no matches found
		fmt.Println("No HTML comments found.")
		return
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
