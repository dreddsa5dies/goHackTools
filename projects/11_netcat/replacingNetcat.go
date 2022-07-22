package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("p", 0, "local port number")

	flag.Usage = func() {
		fmt.Println(strings.ReplaceAll(
			`connect:	$name HOSTNAME PORT
listen:		$name -p PORT
	-p	listen port number`,
			"$name", filepath.Base(os.Args[0])))
	}

	flag.Parse()

	if *port > 0 {
		log.Fatal(Listen(*port))
	}

	two := 2
	if flag.NArg() != two {
		flag.Usage()
		return
	}

	dialPort := 0
	fmt.Sscanf(flag.Arg(1), "%d", &dialPort)

	err := Dial(flag.Arg(0), dialPort)
	if err != nil {
		log.Fatal(err)
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
			log.Println("accept error:", err)
		}

		log.Println("accept:", conn.RemoteAddr())

		go func(c net.Conn) {
			_, err = io.Copy(os.Stdout, c)
			if err != nil {
				log.Println(err)
			}

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

	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	var src io.Reader

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// To retrieve text through | (pipe)
		buffer, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		src = bytes.NewReader(buffer)
	} else {
		src = os.Stdin
	}

	_, err = io.Copy(conn, src)

	return err
}
