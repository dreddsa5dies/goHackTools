package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	pcapFile string
	handle   *pcap.Handle
	err      error
)

func main() {
	// справка
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Использование: %s PCAP_FILE\n", os.Args[0])
		os.Exit(1)
	}

	pcapFile = os.Args[1]

	// открытие файла
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// header & footer KML
	kmlheader := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
<Document>`
	kmlfooter := `</Document>
</kml>`

	// data structs
	allIP := make(map[string]int)
	allKml := ""

	// просмотр всех пакетов
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// от кого кому
		if net := packet.NetworkLayer(); net != nil {
			src, dst := net.NetworkFlow().Endpoints()
			allIP[src.String()]++
			allIP[dst.String()]++
		}
	}

	// все IP
	for k, _ := range allIP {
		kml, err := recKml(k)
		if err != nil {
			log.Fatal(err)
		}
		allKml += kml + "\n"
	}

	fmt.Println(kmlheader + "\n" + allKml + kmlfooter)
}

func recKml(ip string) (string, error) {
	if ip == "" {
		return "", errors.New("Not IP")
	}

	absPath, _ := filepath.Abs("GeoLite2-City.mmdb")
	db, err := geoip2.Open(absPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	ipNew := net.ParseIP(ip)
	record, err := db.City(ipNew)
	if err != nil {
		return "", err
	}

	kml := fmt.Sprintf(`<Placemark>
<name>%s</name>
<Point>
<coordinates>%6f,%6f</coordinates>
</Point>
</Placemark>`, ip, record.Location.Longitude, record.Location.Latitude)
	return kml, nil
}
