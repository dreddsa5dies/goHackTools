package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"log"

	"github.com/shavac/gexpect"
)

func main() {
	// где находится ssh
	ssh, err := exec.LookPath("ssh")
	if err != nil {
		log.Println(err)
	}

	// новое подключение
	child, _ := gexpect.NewSubProcess(ssh, "user@127.0.0.1")
	if err := child.Start(); err != nil {
		fmt.Println(err)
	}
	// закрытие соединения
	defer child.Close()
	// ввод пароля
	if idx, _ := child.ExpectTimeout(0*time.Second, regexp.MustCompile("password:")); idx >= 0 {
		child.SendLine("pass")
	}

	// ввод команды
	child.SendLine("sudo cat /etc/shadow | grep root")

	// время ожидания
	child.InteractTimeout(3 * time.Second)
}
