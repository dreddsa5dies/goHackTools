// Base64 encoding data

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		usage(os.Args[0])
	}

	switch fl := os.Args[1]; fl {
	case "encode":
		// Encode bytes to base64 encoded string.
		encodedString := base64.StdEncoding.EncodeToString([]byte(os.Args[2]))
		fmt.Printf("%s\n", encodedString)
	case "decode":
		// Decode base64 encoded string to bytes.
		decodedData, err := base64.StdEncoding.DecodeString(os.Args[2])
		if err != nil {
			log.Fatal("Error decoding data. ", err)
		}
		fmt.Printf("%s\n", decodedData)
	default:
		usage(os.Args[0])
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stderr, "Using: %s encode DATA\n", name)
	fmt.Fprintf(os.Stderr, "OR %s decode DATA\n", name)
	os.Exit(1)
}
