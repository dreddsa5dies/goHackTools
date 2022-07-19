package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	const arg = 2
	// go run goNmapScan.go IP
	// справка
	if len(os.Args) != arg {
		fmt.Fprintf(os.Stderr, "Использование: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	// установка местонахождения NMAP
	binary, err := exec.LookPath("/usr/bin/nmap")
	if err != nil {
		log.Fatalln(err)
	}

	// установка аргументов
	args := []string{"nmap", "-v", "-A", os.Args[1]}

	env := os.Environ()

	err = syscall.Exec(binary, args, env)
	if err != nil {
		log.Fatalln(err)
	}
}
