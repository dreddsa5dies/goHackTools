// Changing the user agent of a request

package main

import (
	"log"
	"net/http"
)

func main() {
	// Create the request for use later
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://www.whoishostingthis.com/tools/user-agent/", nil)
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	// Override the user agent
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")

	_, err = client.Do(request)
	if err != nil {
		log.Fatal("Error making request. ", err)
	}
}
