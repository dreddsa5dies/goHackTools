// Creating a steganographic image archive

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stdout, "Usage:\t%v file.jpg file.zip\n", os.Args[0])
		fmt.Printf("Create a hidden archive into new jpg\n")
		os.Exit(1)
	}

	fileJpg := os.Args[1]
	fileZip := os.Args[2]

	if !strings.HasSuffix(fileJpg, "jpg") || !strings.HasSuffix(fileZip, "zip") {
		fmt.Fprintf(os.Stdout, "Usage:\t%v file.jpg file.zip\n", os.Args[0])
		fmt.Printf("Create a hidden archive into new jpg\n")
		os.Exit(1)
	}

	// Open original file
	firstFile, err := os.Open(fileJpg)
	if err != nil {
		log.Fatal(err)
	}
	defer firstFile.Close()
	// Second file
	secondFile, err := os.Open(fileZip)
	if err != nil {
		log.Fatal(err)
	}
	defer secondFile.Close()

	// New file for output
	newFile, err := os.Create("stego_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Copy the bytes to destination from source_, err = io.Copy(newFile, firstFile)
	_, err = io.Copy(newFile, firstFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(newFile, secondFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
}
