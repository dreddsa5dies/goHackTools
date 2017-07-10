package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

func main() {
	// какая подсеть
	netIP := myNet()

	for _, a := range netIP {
		addressHosts := strings.TrimRight(a, "0")
		for i := 1; i < 254; i++ {
			if pingTarget(addressHosts + strconv.Itoa(i)) {
				println(addressHosts + strconv.Itoa(i))
			}
		}
	}

	// поиск хостов

	// сканирование на SSH

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
