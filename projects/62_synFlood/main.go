package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func capture(iface, target string, results map[string]int) {
	var (
		snaplen = int32(320)
		timeout = pcap.BlockForever
		/**
		TCP Flags and Their Byte Positions
		Bit 7 6 5 4 3 2 1 0
		Flag CWR ECE URG ACK PSH RST SYN FIN

		ACK and FIN: 00010001 (0x11)
		ACK: 00010000 (0x10)
		ACK and PSH: 00011000 (0x18)

		tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18
		**/
		filter = "tcp[13] == 0x11 or tcp[13] == 0x10 or tcp[13] == 0x18"
	)

	handle, err := pcap.OpenLive(iface, snaplen, true, timeout)
	if err != nil {
		log.Println(err)
		return
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Println(err)
		return
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("Capturing packets")
	for packet := range source.Packets() {
		networkLayer := packet.NetworkLayer()
		if networkLayer == nil {
			continue
		}

		transportLayer := packet.TransportLayer()
		if transportLayer == nil {
			continue
		}

		srcHost := networkLayer.NetworkFlow().Src().String()
		srcPort := transportLayer.TransportFlow().Src().String()

		if srcHost != target {
			continue
		}
		results[srcPort]++
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: sudo ./main <capture_iface> <target_ip> <port1,port2,port3>")
		os.Exit(1)
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	var devFound bool
	iface := os.Args[1]
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
			break
		}
	}

	if !devFound {
		log.Printf("device named '%s' does not exist", iface)
		return
	}

	ip := os.Args[2]
	results := make(map[string]int, 0)
	go capture(iface, ip, results)

	<-time.After(1 * time.Second)

	for _, port := range explode(os.Args[3]) {
		target := fmt.Sprintf("%s:%s", ip, port)
		fmt.Println("Trying", target)
		c, err := net.DialTimeout("tcp", target, 1000*time.Millisecond)
		if err != nil {
			continue
		}
		c.Close()
	}
	<-time.After(2 * time.Second)

	for port, confidence := range results {
		if confidence >= 1 {
			fmt.Printf("Port %s open (confidence: %d)\n", port, confidence)
		}
	}
}

func explode(portString string) []string {
	ret := make([]string, 0)

	ports := strings.Split(portString, ",")
	for _, port := range ports {
		port = strings.TrimSpace(port)
		ret = append(ret, port)
	}

	return ret
}
