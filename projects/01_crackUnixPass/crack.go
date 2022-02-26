package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"

	"os"

	"log"

	"strings"

	crypt "github.com/amoghe/go-crypt"
)

var (
	passfile   string
	dictionary string
)

func init() {
	flag.StringVar(&passfile, "f", "", "Open shadow")
	flag.StringVar(&dictionary, "d", "", "Open pass dictionary")
}

func main() {
	// разбор флагов
	flag.Parse()

	// вывод справки
	if passfile == "" || dictionary == "" {
		println("Please " + os.Args[0] + " -h")
		os.Exit(0)
	}

	// открываем shadow
	passFile, err := os.Open(passfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer passFile.Close()

	// парольный словарь
	dictFile, err := ioutil.ReadFile(dictionary)
	if err != nil {
		log.Fatalln(err)
	}

	passDict := strings.Split(string(dictFile), "\n")

	// построчно
	scanner := bufio.NewScanner(passFile)
	for scanner.Scan() {
		j := scanner.Text()
		// строки с логин/пароль
		if strings.Contains(j, ":") {
			shadowText := strings.Split(j, ":")
			user, cryptPass := shadowText[0], shadowText[1]
			fmt.Printf("[*] Cracking Password For: %v\n", user)
			for i := 0; i < len(passDict)-1; i++ {
				if testPass(cryptPass, passDict[i]) != "" {
					println(testPass(cryptPass, passDict[i]))
					break
				}
			}
		}
	}
}

func testPass(cryptPass string, passWord string) string {
	saltSearch := strings.LastIndex(cryptPass, "$")
	salt := cryptPass[0:saltSearch]

	cryptWord, err := crypt.Crypt(passWord, salt)
	if err != nil {
		log.Fatalf("Ошибка SHA: %v", err)
	}
	// если найден !
	if cryptWord == cryptPass {
		return "[+] Found PASSWORD: " + passWord
	}
	return ""
}
