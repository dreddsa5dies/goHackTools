/**
 * WEB SCRAPPER
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getListing(listingURL string) {
	response, err := http.Get(listingURL)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		bodyText, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", bodyText)
	}

}
func main() {
	getListing("https://github.com/")
}
