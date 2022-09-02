// Extracting emails from url to Maltego

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
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// write body to []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// create email regexp
	regMail := regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)

	// variable to emails
	var mailAddr []string

	// search emails
	if regMail.MatchString(string(body)) {
		mailAddr = regMail.FindAllString(string(body), -1)
	}

	// print Maltego transform response
	fmt.Println("<MaltegoMessage>")
	fmt.Println("<MaltegoTransformResponseMessage>")
	fmt.Println("  <Entities>")

	for _, mail := range mailAddr {
		fmt.Println("    <Entity Type=\"maltego.EmailAddress\">")
		fmt.Println("      <Value>" + mail + "</Value>")
		fmt.Println("    </Entity>")
	}

	fmt.Println("  </Entities>")
	fmt.Println("</MaltegoTransformResponseMessage>")
	fmt.Println("</MaltegoMessage>")
}
