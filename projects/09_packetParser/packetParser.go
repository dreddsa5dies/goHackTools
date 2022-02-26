package main

import (
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	device       string = "wls1"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
)

func main() {
	// открытие интерфейса
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// просмотр всех пакетов
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// от кого кому
		if net := packet.NetworkLayer(); net != nil {
			src, dst := net.NetworkFlow().Endpoints()
			fmt.Printf("[+] Src: %v, --> Dst: %v \n", src, dst)
			printRecord(src.String(), dst.String())
		}
	}
}

func printRecord(src string, dst string) {
	if src == "" || dst == "" {
		log.Fatalln("Error IP")
	}

	absPath, _ := filepath.Abs("GeoLite2-City.mmdb")
	db, err := geoip2.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ipSRC := net.ParseIP(src)
	recordSRC, err := db.City(ipSRC)
	if err != nil {
		log.Fatal(err)
	}

	ipDST := net.ParseIP(dst)
	recordDST, err := db.City(ipDST)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] SRC: %v, %v\n", recordSRC.City.Names["ru"], recordSRC.Country.Names["ru"])
	fmt.Printf("[+] DST: %v, %v\n", recordDST.City.Names["ru"], recordDST.Country.Names["ru"])
}
