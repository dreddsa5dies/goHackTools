package main

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	printRecord("173.255.226.98")

	printRecord("81.2.69.142")

	printRecord("35.184.160.12")
}

func printRecord(tgt string) {
	if tgt == "" {
		log.Fatalln("Error IP")
	}

	absPath, _ := filepath.Abs("GeoLite2-City.mmdb")
	db, err := geoip2.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(tgt)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[*] Target: %v Geo-located.\n", tgt)
	fmt.Printf("[+] %v, %v, %v\n", record.City.Names["ru"], record.Subdivisions[0].Names["ru"], record.Country.Names["ru"])
	fmt.Printf("[+] ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("[+] Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("[+] Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	println()
}
