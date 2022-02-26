// Identifying alternative sites by spoofing user agents

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// check HELP
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage:\t%v URL\n", os.Args[0])
		os.Exit(1)
	}

	inputURL := os.Args[1]

	// check URL
	_, err := url.ParseRequestURI(inputURL)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Usage:\t%v URL\n", os.Args[0])
		log.Fatalln(err)
	}

	// tests user-agents
	userAgents := make(map[string]string)
	userAgents["Chrome on Windows 8.1"] = "Mozilla/5.0 (Windows NT6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.115 Safari/537.36"
	userAgents["Safari on iOS"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4"
	userAgents["IE6 on Windows XP"] = "Mozilla/5.0 (Windows; U; MSIE 6.0; WindowsNT 5.1; SV1; .NET CLR 2.0.50727)"
	userAgents["Googlebot"] = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

	// save requests from user-agents
	resultReq := make(map[string]string)

	// iterate
	for k, v := range userAgents {
		resultReq[k] = doReq(inputURL, v)
	}

	// 	The next part of the code creates an md5s array and then iterates through the responses,
	// grabbing the response.text file. From this, it generates an md5 hash of the response
	// content and stores it into the md5s array:

	md5s := make(map[string]string)
	for name, response := range resultReq {
		md5s[name] = md5EmptyHash(response)
	}

	// 	The final part of the code iterates through the md5s array and compares each item to the
	// original baseline request, in this recipe Chrome on Windows 8.1:
	for name, md5 := range md5s {
		if name != "Chrome on Windows 8.1" {
			if md5 != md5s["Chrome on Windows 8.1"] {
				fmt.Println(name, "differs from baseline")
			} else {
				fmt.Printf("No alternative site found via User-Agent spoofing: %v", md5)
			}
		}
	}

	// END
	os.Exit(0)
}

// doReq - get body request
func doReq(inputURL, userAgent string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", inputURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

// get md5 hash
func md5EmptyHash(message string) string {
	h := md5.New()
	io.WriteString(h, message)
	return fmt.Sprintf("%x", h.Sum(nil))
}
