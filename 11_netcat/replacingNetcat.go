package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	// define some global variables
	listen            bool
	command           bool
	upload            = false
	execute           string
	target            string
	uploadDestination string
	port              int
)

func init() {
	flag.BoolVar(&listen, "l", false, "Listen on [host]:[port] for incoming connections")
	flag.BoolVar(&command, "c", false, "Initialize a command shell")

	flag.StringVar(&execute, "e", "", "Execute the given file upon receiving a connection")
	flag.StringVar(&target, "t", "", "Target")
	flag.StringVar(&uploadDestination, "u", "", "Upon receiving connection upload a file and write to [destination]")

	flag.IntVar(&port, "p", 0, "Port")
}

func main() {
	flag.Parse()

	if !listen && len(target) > 0 && port > 0 {
		// read in the buffer from the commandline
		// this will block, so send CTRL-D if not sending input
		// to stdin
		buffer, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}

		// send data off
		clientSender(buffer)
	}

	// we are going to listen and potentially
	// upload things, execute commands, and drop a shell back
	// depending on our command line options above
	if listen {
		serverLoop()
	}
}
