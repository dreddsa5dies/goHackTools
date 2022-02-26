// Finding writable files
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Recursively look for files with the write bit set for everyone.")
		fmt.Println("Usage: " + os.Args[0] + "DIR")
		fmt.Println("Example: " + os.Args[0] + " /var/log")
		os.Exit(1)
	}

	dirPath := os.Args[1]

	err := filepath.Walk(dirPath, checkFilePermissions)
	if err != nil {
		log.Fatal(err)
	}
}

func checkFilePermissions(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		// no need to see
		return nil
	}

	// Bitwise operators to isolate specific bit groups
	maskedPermissions := fileInfo.Mode().Perm() & 0002
	if maskedPermissions == 0002 {
		fmt.Println(fileInfo.Mode().Perm().String() + " " + path)
	}

	return nil
}
