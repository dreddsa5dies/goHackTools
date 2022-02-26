package main

// If you need run this programm, you need build this and run from root or sudo users:
// go build main.go
// sudo ./main

// Device is typically /dev/sda but may also be /dev/sdb, /dev/sdc
// Use mount, or df -h to get info on which drives are being used
// You will need sudo to access some disks at this level
import (
	"io"
	"log"
	"os"
)

func main() {
	path := "/dev/sda"
	log.Println("[+] Reading boot sector of " + path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// The file.Read() function will read a tiny file in to a large
	// byte slice, but io.ReadFull() will return an
	// error if the file is smaller than the byte slice.
	byteSlice := make([]byte, 512)
	// ReadFull Will error if 512 bytes not available to read
	numBytesRead, err := io.ReadFull(file, byteSlice)
	if err != nil {
		log.Fatal("Error reading 512 bytes from file. " + err.Error())
	}
	log.Printf("Bytes read: %d\n\n", numBytesRead)
	log.Printf("Data as decimal:\n%d\n\n", byteSlice)
	log.Printf("Data as hex:\n%x\n\n", byteSlice)
	log.Printf("Data as string:\n%s\n\n", byteSlice)
}
