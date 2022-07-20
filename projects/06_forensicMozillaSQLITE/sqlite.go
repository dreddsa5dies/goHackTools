// If error - run on terminal: "CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go run sqlite.go"
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	locate string // locate sqlite file
	param  bool   // cookie or places

)

func init() {
	flag.StringVar(&locate, "l", "", "Locate sqlite file? Default - none.")
	flag.BoolVar(&param, "p", false, "Cookies or Places? Default - Cookie.")
}

func main() {
	// разбор флагов
	flag.Parse()

	if locate == "" {
		fmt.Println("Please " + os.Args[0] + " -h")
		os.Exit(0)
	}

	database, err := sql.Open("sqlite3", locate)
	if err != nil {
		log.Fatalln(err)
	}

	if !param {
		// cookies
		rows, err := database.Query("SELECT host, name, value FROM moz_cookies")
		if err != nil {
			log.Fatalln(err)
		}

		var host, name, value string

		println("[*] -- Found Cookies --")

		for rows.Next() {
			if err := rows.Scan(&host, &name, &value); err == nil {
				fmt.Printf("[+] Host:" + host + ", Cookie:" + name + ", Value:" + value + "\n")
			}
		}

		return
	}

	rows, err := database.Query("select url, datetime(visit_date/1000000, 'unixepoch') from moz_places, moz_historyvisits where visit_count > 0 and moz_places.id==moz_historyvisits.place_id;")
	if err != nil {
		log.Fatalln(err)
	}

	var url, date string

	println("[*] -- Found History --")

	for rows.Next() {
		if err := rows.Scan(&url, &date); err == nil {
			fmt.Printf("[+] " + date + " - Visited: " + url + "\n")
		}
	}
}
