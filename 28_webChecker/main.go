// if you need check web resource - use it

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v url\n", os.Args[0])
		fmt.Printf("Output format: write to url.log\n")
		os.Exit(1)
	}

	// validate url
	u, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// write to .log file
	pwdDir, _ := os.Getwd()
	fLog, err := os.OpenFile(pwdDir+`/`+u.Hostname()+`.log`, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Fatalln(err)
	}
	defer fLog.Close()
	log.SetOutput(fLog)

	// checker
	for {
		resp, err := http.Get(u.String())
		if err != nil {
			log.Printf("Error connection %s\n", err)
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("Error. %s http-status: %d\n", u.String(), resp.StatusCode)
			return
		}

		log.Printf("Online. %s http-status: %d\n", u.String(), resp.StatusCode)

		time.Sleep(3 * time.Second)
	}
}
