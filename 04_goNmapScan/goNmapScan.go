package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// go run goNmapScan.go IP

	// установка местонахождения NMAP
	binary, lookErr := exec.LookPath("/usr/bin/nmap")
	if lookErr != nil {
		panic(lookErr)
	}

	// установка аргументов
	args := []string{"nmap", "-v", "-A", os.Args[1]}

	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
