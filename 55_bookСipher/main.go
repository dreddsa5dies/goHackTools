// Книжный шифр (book cipher) — Книжный шифр — вид шифра, в котором каждый элемент открытого
// текста (каждая буква или слово) заменяется на указатель (например, номер страницы, строки
// и столбца) аналогичного элемента в дополнительном тексте-ключе.
// Для дешифрования необходимо иметь как закрытый текст, так и дополнительный текст-ключ.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Print("Write the message\t:> ")
	var startMessage string
	_, err := fmt.Scanf("%s", &startMessage)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Write the book-key\t:> ")
	var bookKey string
	_, err = fmt.Scanf("%s", &bookKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("[E]ncrypt|[D]ecrypt\t:> ")
	var key string
	_, err = fmt.Scanf("%s", &key)
	if err != nil {
		log.Fatalln(err)
	}

	var encDecBool bool
	switch {
	case strings.ToLower(key) == "e":
		encDecBool = true
	case strings.ToLower(key) == "d":
		encDecBool = false
	default:
		log.Fatalln("No such key!")
	}

	fmt.Println(encryptDecrypt(encDecBool, startMessage, bookKey))
}

func encryptDecrypt(mode bool, message, key string) string {

	if mode {
		// encrypt
	} else {
		// decrypt
	}

	// get random symbol from book
	fileInfo, err := os.Stat(key)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Book-key does not exist.")
		}
	}
	dat, err := ioutil.ReadFile(fileInfo.Name())
	if err != nil {
		log.Fatalln(err)
	}
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	codingMess := string(dat[rand.Intn(len(dat))])

	return codingMess
}
