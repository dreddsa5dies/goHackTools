package main

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v domainName\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "Ex:\t%v github.com\n", os.Args[0])
		fmt.Printf("Output format: IP\n")
		os.Exit(1)
	}

	var msg dns.Msg
	fqdn := dns.Fqdn(os.Args[1])
	msg.SetQuestion(fqdn, dns.TypeA)
	in, err := dns.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	if len(in.Answer) < 1 {
		fmt.Println("No records")
		return
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.A)
		}
	}
}
