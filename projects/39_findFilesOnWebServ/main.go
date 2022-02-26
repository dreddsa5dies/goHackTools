// Finding unlisted files on a webserver
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	// Load command line arguments
	if len(os.Args) != 4 {
		fmt.Println(os.Args[0] + " - Perform an HTTP HEAD request to a URL")
		fmt.Println("Usage: " + os.Args[0] +
			" <wordlist_file> <URL> <maxThreads>")
		fmt.Println("Example: " + os.Args[0] +
			" wordlist.txt https://www.devdungeon.com 10")
		os.Exit(1)
	}

	wordlistFilename := os.Args[1]
	baseURL := os.Args[2]
	maxThreads, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal("Error converting maxThread value to integer. ", err)
	}

	// Track how many threads are active to avoid
	// flooding a web server
	activeThreads := 0
	doneChannel := make(chan bool)

	// Open word list file for reading
	wordlistFile, err := os.Open(wordlistFilename)
	if err != nil {
		log.Fatal("Error opening wordlist file. ", err)
	}

	// Read each line and do an HTTP HEAD
	scanner := bufio.NewScanner(wordlistFile)
	for scanner.Scan() {
		go checkIfURLExists(baseURL, scanner.Text(), doneChannel)
		activeThreads++
		// Wait until a done signal before next if max threads reached
		if activeThreads >= maxThreads {
			<-doneChannel
			activeThreads--
		}
	}

	// Wait for all threads before repeating and fetching a new batch
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}

	// Scanner errors must be checked manually
	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading wordlist file. ", err)
	}
}

// Given a base URL (protocol+hostname) and a filepath (relative URL)
// perform an HTTP HEAD and see if the path exists.
// If the path returns a 200 OK print out the path
func checkIfURLExists(baseURL, filePath string, doneChannel chan bool) {
	// Create URL object from raw string
	targetURL, err := url.Parse(baseURL)
	if err != nil {
		log.Println("Error parsing base URL. ", err)
	}

	// Set the part of the URL after the host name
	targetURL.Path = filePath

	// Perform a HEAD only, checking status without
	// downloading the entire file
	response, err := http.Head(targetURL.String())
	if err != nil {
		log.Println("Error fetching ", targetURL.String())
	}

	// If server returns 200 OK file can be downloaded
	if response.StatusCode == 200 {
		log.Println(targetURL.String())
	}

	// Signal completion so next thread can start
	doneChannel <- true
}
