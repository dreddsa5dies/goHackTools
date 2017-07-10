package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

func main() {
	// какая подсеть
	netIP := myNet()

	// поиск хостов
	// TODO: без fastpinga надо делать
	allHosts := searchHosts(netIP)
	println(allHosts)

	// сканирование
	for _, host := range allHosts {

		wg := sync.WaitGroup{}

		c := func(ports int) {
			conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(ports), time.Duration(1)*time.Second)
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
				fmt.Println(ports, " open")
			}
			wg.Done()
		}

		wg.Add(65534)
		for i := 0; i < 65534; i++ {
			go c(i)
		}
		wg.Wait()
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
				result = append(result, network.String())
			}
		}
	}

	return result
}

// определение доступности хоста
func pingTarget(ipTarget string) bool {
	p := fastping.NewPinger()

	err := p.AddIP(ipTarget)
	if err != nil {
		log.Println(err)
	}

	target := false

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		target = true
	}

	err = p.Run()
	if err != nil {
		log.Println(err)
	}

	return target
}

// поиск хостов
func searchHosts(netIP []string) []string {
	allHosts := make([]string, 0)

	for _, a := range netIP {
		addressHosts := strings.TrimRight(a, "0")

		wg := sync.WaitGroup{}

		p := func(hostPing string) {
			if pingTarget(hostPing) {
				allHosts = append(allHosts, hostPing)
				wg.Done()
			}
		}

		wg.Add(253)
		for i := 1; i < 254; i++ {
			go p(addressHosts + strconv.Itoa(i))
		}
		wg.Wait()
	}

	return allHosts
}
