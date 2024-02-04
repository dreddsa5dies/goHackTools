package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-r] <file_or_directory>\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[len(os.Args)-1]
	recursive := false
	if len(os.Args) == 3 && os.Args[1] == "-r" {
		recursive = true
	}

	if recursive {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			printOut(info)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		printOut(fileInfo)
	}
}

func printOut(info os.FileInfo) {
	fmt.Println("File name:", info.Name())
	fmt.Println("Size in bytes:", info.Size())
	fmt.Println("Permissions:", info.Mode())
	fmt.Println("Last modified:", info.ModTime())
	fmt.Println("Is Directory: ", info.IsDir())
	fmt.Printf("System interface type: %T\n", info.Sys())
	fmt.Printf("System info: %+v\n\n", info.Sys())
}
