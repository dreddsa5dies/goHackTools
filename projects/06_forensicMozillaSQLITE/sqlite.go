package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	/*
		// cookies
		database, _ := sql.Open("sqlite3", "./cookies.sqlite")
		rows, _ := database.Query("SELECT host, name, value FROM moz_cookies")
		var host, name, value string

		println("[*] -- Found Cookies --")

		for rows.Next() {
			rows.Scan(&host, &name, &value)
			fmt.Printf("[+] Host:" + host + ", Cookie:" + name + ", Value:" + value + "\n")
		}
	*/

	// история посещений
	database, _ := sql.Open("sqlite3", "./places.sqlite")
	rows, _ := database.Query("select url, datetime(visit_date/1000000, 'unixepoch') from moz_places, moz_historyvisits where visit_count > 0 and moz_places.id==moz_historyvisits.place_id;")
	var url, date string

	println("[*] -- Found History --")

	for rows.Next() {
		rows.Scan(&url, &date)
		fmt.Printf("[+] " + date + " - Visited: " + url + "\n")
	}
}
