package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Data    string `short:"d" long:"data" default:"" description:"Data for decrypt/encrypt"`
	Key     int    `short:"k" long:"key" default:"" description:"Key value for Cipher of Caesar"`
	Verbose bool   `short:"v" long:"verbose" description:"Decrypt string"`
}

func main() {
	flags.Parse(&opts)

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		os.Exit(1)
	}

	if opts.Verbose {
		// Decrypt
		fmt.Fprintf(os.Stdout, "Decrypt: %s\n", decrypt(opts.Data, opts.Key))
	} else {
		// Encrypt
		fmt.Fprintf(os.Stdout, "Encrypt: %s\n", encrypt(strings.ToLower(opts.Data), opts.Key))
	}

	os.Exit(0)
}

func encrypt(data string, key int) string {
	result := strings.Map(func(r rune) rune {
		return caesar(r, -key)
	}, data)
	return result
}

func decrypt(data string, key int) string {
	result := strings.Map(func(r rune) rune {
		return caesar(r, +key)
	}, data)
	return result
}

func caesar(r rune, shift int) rune {
	s := int(r) + shift
	if s > 'z' {
		return rune(s - 26)
	} else if s < 'a' {
		return rune(s + 26)
	}
	return rune(s)
}
