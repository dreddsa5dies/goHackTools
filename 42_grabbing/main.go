// Grabbing a banner from a service

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Using: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	ipToScan := os.Args[1]

	activeThreads := 0
	doneChannel := make(chan bool)

	for port := 0; port <= 65535; port++ {
		go grabBanner(ipToScan, port, doneChannel)
		activeThreads++
	}

	// Wait for all threads to finish
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
}

func grabBanner(ip string, port int, doneChannel chan bool) {
	connection, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port), time.Second*10)
	if err != nil {
		doneChannel <- true
		return
	}

	// See if server offers anything to read
	buffer := make([]byte, 4096)
	connection.SetReadDeadline(time.Now().Add(time.Second * 5))

	// Set timeout
	numBytesRead, err := connection.Read(buffer)
	if err != nil {
		doneChannel <- true
		return
	}

	fmt.Printf("Banner from port %d\n%s\n", port, buffer[0:numBytesRead])
	doneChannel <- true
}
