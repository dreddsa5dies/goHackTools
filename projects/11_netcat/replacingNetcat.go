package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func showLocalAddrs() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

// Listen port
func Listen(port int) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Accept error:", err)
		}
		log.Println("accept:", conn.RemoteAddr())

		go func(c net.Conn) {
			io.Copy(os.Stdout, c)
			log.Println("closed:", conn.RemoteAddr())
			defer c.Close()
		}(conn)
	}
}

// Dial port
func Dial(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// To retrieve text through | (pipe)
		buffer, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = io.Copy(conn, bytes.NewReader(buffer))
	} else {
		_, err = io.Copy(conn, os.Stdin)
	}

	return err
}

func main() {
	port := flag.Int("p", 0, "local port number")

	flag.Usage = func() {
		fmt.Println(strings.Replace(
			`connect:	$name HOSTNAME PORT
listen:		$name -p PORT
	-p	listen port number`,
			"$name", filepath.Base(os.Args[0]), -1))
	}

	flag.Parse()

	if *port > 0 {
		log.Fatal(Listen(*port))
	}

	if flag.NArg() != 2 {
		flag.Usage()
		return
	}

	dialPort := 0
	fmt.Sscanf(flag.Arg(1), "%d", &dialPort)
	Dial(flag.Arg(0), dialPort)
}
