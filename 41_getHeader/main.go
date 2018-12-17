// HTTP response headers
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
	}

	resp, err := http.Get(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Print(v)
		fmt.Println()
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v URL\n", name)
	fmt.Printf("HTTP response headers\n")
	os.Exit(1)
}
