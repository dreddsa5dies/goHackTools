package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	tmp := []string{
		"173.255.226.98",
		"81.2.69.142",
		"35.184.160.12",
	}

	db, err := getDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := range tmp {
		if err = printRecord(tmp[i], db); err != nil {
			log.Println(err)
		}
	}
}

func printRecord(tgt string, db *geoip2.Reader) error {
	if tgt == "" {
		return errors.New("error IP")
	}

	ip := net.ParseIP(tgt)

	record, err := db.City(ip)
	if err != nil {
		return err
	}

	fmt.Printf("[*] Target: %v Geo-located.\n", tgt)
	fmt.Printf("[+] %v, %v, %v\n", record.City.Names["ru"], record.Subdivisions[0].Names["ru"], record.Country.Names["ru"])
	fmt.Printf("[+] ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("[+] Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("[+] Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)

	return nil
}

func getDb() (*geoip2.Reader, error) {
	absPath, err := filepath.Abs("GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}

	db, err := geoip2.Open(absPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
