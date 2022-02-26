// Extracting emails from url to Maltego

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
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Println("Extracting links from url to Maltego.")
		fmt.Println("Usage: " + os.Args[0] + "<url>")
		fmt.Println("U need url: http://domain.com")
		os.Exit(1)
	}

	url := os.Args[1]

	// response url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка запроса %v", err)
	}
	defer resp.Body.Close()

	// write body to []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка записи ответа %v", err)
	}

	// create email regexp
	regMail, _ := regexp.Compile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)

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
