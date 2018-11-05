// This program will find nameservers associated with a given hostname.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		usage(os.Args[0])
	}

	arg := os.Args[1]
	nameservers, err := net.LookupNS(arg)
	if err != nil {
		log.Fatal(err)
	}
	for _, nameserver := range nameservers {
		fmt.Println(nameserver.Host)
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v hostname\n", name)
	fmt.Printf("Looking up nameservers\n")
	os.Exit(1)
}
