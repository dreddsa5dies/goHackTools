// Finding named hosts on a network

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Using: %s ip-addr\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "First three octets. Ex: 192.168.0\n")
		os.Exit(1)
	}

	subnetToScan := os.Args[1]

	activeThreads := 0
	doneChannel := make(chan bool)

	for ip := 0; ip <= 255; ip++ {
		fullIP := subnetToScan + "." + strconv.Itoa(ip)
		go resolve(fullIP, doneChannel)
		activeThreads++
	}

	// Wait for all threads to finish
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
}

func resolve(ip string, doneChannel chan bool) {
	addresses, err := net.LookupAddr(ip)
	if err == nil {
		fmt.Printf("%s - %s\n", ip, strings.Join(addresses, ", "))
	}
	doneChannel <- true
}
