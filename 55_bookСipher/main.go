// Книжный шифр (book cipher) — вид шифра, в котором каждый элемент открытого
// текста (каждая буква или слово) заменяется на указатель (например, номер страницы, строки
// и столбца) аналогичного элемента в дополнительном тексте-ключе.
// Для дешифрования необходимо иметь как закрытый текст, так и дополнительный текст-ключ.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Print("Write the message:> ")
	var startMessage string
	// scan ONE string & spaces
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		startMessage = scanner.Text()
		break
	}

	fmt.Print("Write the book-key:> ")
	var bookKey string
	_, err := fmt.Scanf("%s", &bookKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("[E]ncrypt|[D]ecrypt:> ")
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

	fmt.Println("Final message: " + encryptDecrypt(encDecBool, startMessage, bookKey))
}

func encryptDecrypt(mode bool, message, key string) string {
	final := ""

	// get book
	fileInfo, err := os.Stat(key)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Book-key does not exist.")
		}
	}
	book, err := ioutil.ReadFile(fileInfo.Name())
	if err != nil {
		log.Fatalln(err)
	}

	// get the mode & work it
	if mode {
		// encrypt
		for i := 0; i < len(message); i++ {
			var listIndexKey []int
			for k := 0; k < len(string(book)); k++ {
				if message[i] == string(book)[k] {
					listIndexKey = append(listIndexKey, k)
				}
			}
			final += strconv.Itoa(listIndexKey[i]) + "/"
		}
	} else {
		// decrypt
		for i := 0; i < len(regnumbers(message)); i++ {
			tmpNumSymbol, _ := strconv.Atoi(regnumbers(message)[i])
			for k := 0; k < len(string(book)); k++ {
				if tmpNumSymbol == k {
					final += string(string(book)[k])
				}
			}
		}
	}

	return final
}

func regnumbers(text string) []string {
	// create regexp
	regNum, _ := regexp.Compile(`[0-9*]{1,}`)
	// find number regexp
	return regNum.FindAllString(text, -1)
}
