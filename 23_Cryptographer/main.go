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
	"regexp"
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

	decryptFile(hash)

	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------")

	cryptoDir(dir, hash)

	os.Exit(0)
}

func cryptoDir(dir string, hash string) {
	// открываем директорию
	dh, err := os.Open(dir)
	if err != nil {
		log.Fatalln("os.Open: %v\n", err)
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
					log.Fatalln("ioutil.ReadFile: %v\n", err)
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
		log.Fatalln("EncryptFunc NewGCM: %v\n", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln("EncryptFunc NonceSize: %v\n", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func encryptFile(filename string, data []byte, hash string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalln("encryptFileFunc Create: %v\n", err)
	}
	defer f.Close()
	f.Write(encrypt(data, hash))
}

func isDir(name string) bool {
	var stat bool

	file, err := os.Open(name)
	if err != nil {
		log.Fatalln("os.Open: %v\n", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatalln("file.Stat: %v\n", err)
	}

	if fi.IsDir() {
		stat = true
	} else {
		stat = false
	}
	return stat
}

func decryptFile(hash string) {
	tmp, err := os.OpenFile("/tmp/decrypt.go", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln("os.OpenFile\n")
	}
	defer tmp.Close()

	log.Printf("Save decrypt file >> %v\n", tmp.Name())
	tmp.WriteString(hash)
}

/*
var str = `package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"io/ioutil"
	"os"
)

func main() {

}

func decrypt(data []byte, hash string) []byte {
	key := []byte(createHash(hash))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("DecryptFunc NewCipher: %v\n", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln("DecryptFunc NewGCM: %v\n", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalln("HASH FAIL\n%v\n", err)
	}
	return plaintext
}

func decryptFile(filename string, hash string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("decryptFileFunc ReadFile: %v\n", err)
	}
	return decrypt(data, hash)
}`
*/
