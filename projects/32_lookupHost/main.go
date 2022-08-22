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
	fmt.Println("Looking up IP addresses for hostname: " + arg)

	ips, err := net.LookupHost(arg)
	if err != nil {
		log.Fatal(err)
	}

	for i := range ips {
		fmt.Println(ips[i])
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v hostname\n", name)
	fmt.Printf("Looking up IP addresses for hostname\n")
	os.Exit(1)
}
