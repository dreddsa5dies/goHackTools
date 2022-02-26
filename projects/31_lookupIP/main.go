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
	// Parse the IP for validation
	ip := net.ParseIP(arg)
	if ip == nil {
		log.Fatal("Valid IP not detected. Value provided: " + arg)
	}
	fmt.Print(arg + "\t")
	hostnames, err := net.LookupAddr(ip.String())
	if err != nil {
		log.Fatal(err)
	}
	for _, hostnames := range hostnames {
		fmt.Print(hostnames)
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v ip\n", name)
	fmt.Printf("Looking up hostnames for IP address\n")
	os.Exit(1)
}
