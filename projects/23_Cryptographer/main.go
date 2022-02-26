package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	var dir, pass string

	fmt.Print("Write the folder for encryption: ")
	fmt.Scanf("%s\n", &dir)

	if !isDir(dir) {
		log.Fatalln("It's not a directory")
	} else {
		log.Println("Directory ok")
	}

	fmt.Print("Write the password: ")
	fmt.Scanf("%s\n", &pass)

	regStr, _ := regexp.Compile(`([0-9a-zA-Z]){8,}`)
	if regStr.MatchString(pass) {
		log.Println("Pass ok")
	} else {
		log.Println("Bad password")
		log.Fatalln("Use good password")
	}

	// get the hash
	hasher := md5.New()
	hasher.Write([]byte(pass))
	hash := hex.EncodeToString(hasher.Sum(nil))

	decryptFileSave(hash, dir)

	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------")

	cryptoDir(dir, hash)

	time.Sleep(2 * time.Second)
	buildDecryptFile()
}

func cryptoDir(dir string, hash string) {
	// открываем директорию
	dh, err := os.Open(dir)
	if err != nil {
		log.Fatalln(err)
	}
	defer dh.Close()

	// считывание списка файлов
	for {
		fis, err := dh.Readdir(10)
		if err == io.EOF {
			break
		}
		// проход
		for _, fi := range fis {
			// рекурсивный проход по поддиректориям
			if fi.IsDir() {
				cryptoDir(dir+"/"+fi.Name(), hash)
			} else {
				// имя файла
				log.Printf("encrypt %v\n", fi.Name())
				file, err := ioutil.ReadFile(dir + "/" + fi.Name())
				if err != nil {
					log.Fatalln(err)
				}

				encryptFile(dir+"/"+fi.Name()+".crp", file, hash)
				os.Remove(dir + "/" + fi.Name())
			}
		}
	}
}

func encrypt(data []byte, hash string) []byte {
	block, _ := aes.NewCipher([]byte(hash))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln(err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func encryptFile(filename string, data []byte, hash string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.Write(encrypt(data, hash))
}

func isDir(name string) bool {
	var stat bool

	file, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if fi.IsDir() {
		stat = true
	} else {
		stat = false
	}
	return stat
}

func decryptFileSave(hash, dir string) {
	tmp, err := os.OpenFile("./decrypt.go", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	saveFile := fmt.Sprintf(`package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	hash := %s
	dir := %s
	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------")
	decryptoDir(dir, hash)
	os.Exit(0)
}

func decryptoDir(dir string, hash string) {
	// открываем директорию
	dh, err := os.Open(dir)
	if err != nil {
		log.Fatalln(err)
	}
	defer dh.Close()
	// считывание списка файлов
	for {
		fis, err := dh.Readdir(10)
		if err == io.EOF {
			break
		}
		// проход
		for _, fi := range fis {
			// рекурсивный проход по поддиректориям
			if fi.IsDir() {
				decryptoDir(dir+"/"+fi.Name(), hash)
			} else {
				// имя файла
				log.Println("decrypt ", fi.Name())
				decryptFile(dir+"/"+fi.Name(), hash)
				os.Remove(dir + "/" + fi.Name())
			}
		}
	}
}

func decrypt(data []byte, hash string) []byte {
	key := []byte(hash)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return plaintext
}
		
func decryptFile(filename string, hash string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create(strings.TrimSuffix(filename, ".crp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.Write(decrypt(data, hash))
}`, `"`+hash+`"`, `"`+dir+`"`)

	log.Printf("Save decrypt file >> %v\n", tmp.Name())
	tmp.WriteString(saveFile)
	tmp.Close()
}

func buildDecryptFile() {
	cmd := exec.Command("/usr/local/go/bin/go", "build", "./decrypt.go")

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	os.Remove("./decrypt.go")
}
