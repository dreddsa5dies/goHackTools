// Package collectlinks is useful for only one task:
//   Given a response from http.Get it will use parse the page and
//   return to you a slice of all the href links found.
//
// Usage:
//
//    package main
//    import (
//      "github.com/jackdanger/collectlinks"
//      "net/http"
//      "fmt"
//    )
//
//    func main() {
//      resp, _ := http.Get("http://motherfuckingwebsite.com")
//      links := collectlinks.All(resp.Body)
//      fmt.Println(links)
//    }
//
//
// Running that will output:
//
//    [http://twitter.com/thebarrytone http://txti.es]
//
package collectlinks
