package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	args                  = 2
	seconds time.Duration = 10
)

func main() {
	// справка
	if len(os.Args) != args {
		fmt.Fprintf(os.Stderr, "Using: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	target := os.Args[1]

	activeThreads := 0
	doneChannel := make(chan bool)

	for port := 0; port <= 65535; port++ {
		go testTCPConnection(target, port, doneChannel)
		activeThreads++
	}
	// Wait for all threads to finish
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
}

func testTCPConnection(ip string, port int, doneChannel chan bool) {
	_, err := net.DialTimeout(
		"tcp",
		ip+":"+strconv.Itoa(port),
		time.Second*seconds,
	)
	if err == nil {
		fmt.Printf("Port %d: Open\n", port)
	}
	doneChannel <- true
}
