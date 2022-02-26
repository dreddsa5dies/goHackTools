// Detecting a ZIP archive in a JPEG image & extract it

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func usage(name string) {
	fmt.Fprintf(os.Stdout, "Usage:\t%v file.jpg\n", name)
	fmt.Printf("Detecting a ZIP archive in a JPEG image\n")
	fmt.Printf("For unarchiving you need 7z\n")
	os.Exit(1)
}

func main() {
	// usage
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		usage(os.Args[0])
	}

	fileJpg := os.Args[1]

	if !strings.HasSuffix(fileJpg, "jpg") {
		usage(os.Args[0])
	}

	// Zip signature is "\x50\x4b\x03\x04"
	file, err := os.Open(fileJpg)
	if err != nil {
		log.Fatal(err)
	}
	bufferedReader := bufio.NewReader(file)
	fileStat, _ := file.Stat()
	// 0 is being cast to an int64 to force i to be initialized as
	// int64 because filestat.Size() returns an int64 and must be
	// compared against the same type
	for i := int64(0); i < fileStat.Size(); i++ {
		myByte, err := bufferedReader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if myByte == '\x50' {
			// First byte match. Check the next 3 bytes
			byteSlice := make([]byte, 3)
			// Get bytes without advancing pointer with Peek
			byteSlice, err = bufferedReader.Peek(3)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Equal(byteSlice, []byte{'\x4b', '\x03', '\x04'}) {
				log.Printf("Found zip signature at byte %d.", i)
			}
		}
	}

	// Unzip it
	for {
		var unz string
		fmt.Print("Unzip it? (Y/N) > ")
		_, err := fmt.Scanf("%s", &unz)
		if err != nil {
			fmt.Println("Wrong data")
			continue
		}
		switch {
		case unz == "Y":
			fmt.Println("OK")

			// where 7z
			binary, err := exec.LookPath("/usr/bin/7z")
			if err != nil {
				log.Fatalln(err)
			}

			// args
			args := []string{"7z", "e", os.Args[1]}

			env := os.Environ()

			err = syscall.Exec(binary, args, env)
			if err != nil {
				log.Fatalln(err)
			}

		case unz == "N":
			fmt.Println("OK")
			fmt.Println("Exit")
			os.Exit(0)
		default:
			fmt.Println("Wrong data")
			continue
		}
	}

}
