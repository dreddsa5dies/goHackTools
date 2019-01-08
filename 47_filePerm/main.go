// Changing file permissions

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Change the permissions of a file")
		fmt.Println("Usage: " + os.Args[0] + " <mode> <filepath>")
		fmt.Println("Example: " + os.Args[0] + " 777 FILE")
		fmt.Println("Example: " + os.Args[0] + " 0644 FILE")
		os.Exit(1)
	}

	mode := os.Args[1]
	filePath := os.Args[2]

	// Convert the mode value from string to uin32 to os.FileMode
	fileModeValue, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		log.Fatal("Error converting permission string to octal value. ",
			err)
	}

	fileMode := os.FileMode(fileModeValue)
	err = os.Chmod(filePath, fileMode)
	if err != nil {
		log.Fatal("Error changing permissions. ", err)
	}

	fmt.Println("Permissions changed for " + filePath)
}
