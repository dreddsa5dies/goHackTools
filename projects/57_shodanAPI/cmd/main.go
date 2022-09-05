package main

import (
	"fmt"
	"log"
	"os"

	"shodan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: main <searchterm>")
	}

	// use printenv - show ENV
	// env SHODAN_API_KEY=apikey
	apiKey := os.Getenv("SHODAN_API_KEY")

	s := shodan.New(apiKey)

	info, err := s.APIInfo()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf(
		"Query Credits: %d\nScan Credits:  %d\n\n",
		info.QueryCredits,
		info.ScanCredits)

	hostSearch, err := s.HostSearch(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	for i := range hostSearch.Matches {
		fmt.Printf("%18s%8d\n", hostSearch.Matches[i].IPString, hostSearch.Matches[i].Port)
	}
}
