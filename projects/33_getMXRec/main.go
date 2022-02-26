// This program will take a domain name and return the MX records.
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
	mxRecords, err := net.LookupMX(arg)
	if err != nil {
		log.Fatal(err)
	}
	for _, mxRecord := range mxRecords {
		fmt.Printf("Host: %s\tPreference: %d\n", mxRecord.Host, mxRecord.Pref)
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v hostname\n", name)
	fmt.Printf("Looking up MX records\n")
	os.Exit(1)
}
