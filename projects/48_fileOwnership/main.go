// Changing file ownership

package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 4 {
		fmt.Println("Change the owner of a file.")
		fmt.Println("Usage: " + os.Args[0] + " <user> <group> <filepath>")
		fmt.Println("Example: " + os.Args[0] + " dano dano test.txt")
		fmt.Println("Example: sudo " + os.Args[0] + " root root test.txt")
		os.Exit(1)
	}

	username := os.Args[1]
	groupname := os.Args[2]
	filePath := os.Args[3]

	// Look up user based on name and get ID
	userInfo, err := user.Lookup(username)
	if err != nil {
		log.Fatal("Error looking up user "+username+". ", err)
	}

	uid, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		log.Fatal("Error converting "+userInfo.Uid+" to integer. ", err)
	}

	// Look up group name and get group ID
	group, err := user.LookupGroup(groupname)
	if err != nil {
		log.Fatal("Error looking up group "+groupname+". ", err)
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		log.Fatal("Error converting "+group.Gid+" to integer. ", err)
	}

	fmt.Printf("Changing owner of %s to %s(%d):%s(%d).\n", filePath, username, uid, groupname, gid)
	os.Chown(filePath, uid, gid)
}
