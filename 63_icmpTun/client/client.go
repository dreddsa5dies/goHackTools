package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:12000")
	if err != nil {
		log.Fatal(err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		fmt.Print("Input lowercase sentence > ")
		var msg string
		_, err := fmt.Scanf("%s", &msg)
		if err != nil {
			log.Fatal(err)
		}
		buf := []byte(msg)
		_, err = conn.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 1)
	}
}
