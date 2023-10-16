/*cookie-flags takes a url and returns the cookie set.*/
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Input string `short:"i" long:"input" default:"" description:"URL"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		os.Exit(1)
	}

	log.Println("Start")

	_, err = url.ParseRequestURI(opts.Input)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		log.Fatalln(err)
	}

	input := opts.Input
	timeout := flag.Int("timeout", 1000, "timeout for requests")

	cookies := doReq(context.Background(), input, *timeout)
	if len(cookies) == 0 {
		log.Println("Not cookies")
	} else {
		for _, c := range cookies {
			fmt.Printf("%s: %s\n", input, c)
		}
	}

	log.Println("End")

	os.Exit(0)
}

func doReq(ctx context.Context, location string, timeout int) []string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, location, http.NoBody)
	if err != nil {
		log.Fatalln(err)
	}

	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Duration(timeout)*time.Millisecond)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // так надо
	}

	res, err := tr.RoundTrip(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	cookies := make([]string, 0, len(res.Cookies()))
	for _, c := range res.Cookies() {
		cookies = append(cookies, c.Raw)
	}

	return cookies
}
