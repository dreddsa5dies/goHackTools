package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var regexes = []*regexp.Regexp{
	// and many others
	regexp.MustCompile(`(?i)user`),
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)kdb`),
	regexp.MustCompile(`(?i)login`),
	regexp.MustCompile(`(?i)shadow`),
}

func walkFn(path string, f os.FileInfo, err error) error {
	for _, r := range regexes {
		if r.MatchString(path) {
			fmt.Printf("[+] HIT: %s\n", path)
		}
	}
	return nil
}

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v /filepath\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "Ex:\t%v ../.\n", os.Args[0])
		os.Exit(1)
	}

	root := os.Args[1]
	if err := filepath.Walk(root, walkFn); err != nil {
		log.Panicln(err)
	}
}
