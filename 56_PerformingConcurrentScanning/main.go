package main

import (
	"fmt"
	"net"
	"os"
	"sort"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		usage(os.Args[0])
	}

	target := os.Args[1]
	ports := make(chan int, 100)
	results := make(chan int)

	var openports []int
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, target)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

func worker(ports, results chan int, target string) {
	for p := range ports {
		address := fmt.Sprintf(target+":%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v scanme.nmap.org\n", name)
	fmt.Printf("Scanning scanme.nmap.org\n")
	os.Exit(1)
}
