// Extracting links from url to Maltego

package main

import (
	"fmt"
	"log"
	"net/http"
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

	url := os.Args[1]

	// response url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка запроса %v", err)
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
		case strings.HasPrefix(link, "https"):
			// If the first character is a / , which indicates that the link is a relative link, then we’ll output it to
			// the correct format after we have prepended the target URL to the link. While this recipe shows
			// how to deal with one example of a relative link, it is important to note that there are other
			// types of relative links, such as just a filename (example.php), a directory, and also a relative
			// path dot notation (../../example.php), as shown here:
			fmt.Println("<Entity Type=\"maltego.Domain\">")
			fmt.Println("<Value>" + url + link + "</Value>")
		}
	}

	fmt.Println("  </Entities>")
	fmt.Println("</MaltegoTransformResponseMessage>")
	fmt.Println("</MaltegoMessage>")
}
