package main

import (
	"bufio"
	"fmt"
	"io/ioutil"

	"os"

	"log"

	"strings"

	crypt "github.com/amoghe/go-crypt"
)

func main() {
	// открываем shadow
	passFile, err := os.Open("pass")
	if err != nil {
		log.Fatalln(err)
	}
	defer passFile.Close()

	// парольный словарь
	dictFile, err := ioutil.ReadFile("dict.txt")
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
			user := strings.Split(j, ":")[0]
			cryptPass := strings.Split(j, ":")[1]
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
