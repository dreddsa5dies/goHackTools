// Changing file timestamp
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s FILE\n", os.Args[0])
		fmt.Printf("Example: %s test.txt\n", os.Args[0])
		os.Exit(1)
	}

	// Change timestamp to a future time
	futureTime := time.Now().Add(50 * time.Hour).Add(15 * time.Minute)
	lastAccessTime := futureTime
	lastModifyTime := futureTime
	err := os.Chtimes(os.Args[1], lastAccessTime, lastModifyTime)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("timestamp change OK - %s", futureTime.String())
}
