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
	"os"
	"regexp"
)

func main() {
	var dir, pass string

	fmt.Print("Write the folder for encryption: ")
	fmt.Scanf("%s\n", &dir)

	if !isDir(dir) {
		fmt.Fprintf(os.Stderr, "It's not a directory\n")
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "Directory ok\n")
	}

	fmt.Print("Write the password: ")
	fmt.Scanf("%s\n", &pass)

	regStr, _ := regexp.Compile(`([0-9a-zA-Z]){8,}`)
	if regStr.MatchString(pass) {
		fmt.Fprintf(os.Stdout, "Pass ok\n")
	} else {
		fmt.Fprintf(os.Stderr, "Bad password\n")
		fmt.Fprintf(os.Stderr, "Use good password\n")
		os.Exit(1)
	}

	// get the hash
	hasher := md5.New()
	hasher.Write([]byte(pass))
	hash := hex.EncodeToString(hasher.Sum(nil))

	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------")
	os.Exit(0)
}

func readdir(dir string) {
	// открываем директорию
	dh, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Open: %v\n", err)
		os.Exit(1)
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
			// имя файла
			fmt.Println(fi.Name())
			// рекурсивный проход по поддиректориям
			if fi.IsDir() {
				readdir(dir + "/" + fi.Name())
			}
		}
	}
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EncryptFunc NewGCM: %v\n", err)
		os.Exit(1)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Fprintf(os.Stderr, "EncryptFunc NonceSize: %v\n", err)
		os.Exit(1)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DecryptFunc NewCipher: %v\n", err)
		os.Exit(1)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DecryptFunc NewGCM: %v\n", err)
		os.Exit(1)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PASSWORD FAIL\n%v\n", err)
		os.Exit(1)
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encryptFileFunc Create: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decryptFileFunc ReadFile: %v\n", err)
		os.Exit(1)
	}
	return decrypt(data, passphrase)
}

func isDir(name string) bool {
	var stat bool

	file, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.Open: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "file.Stat: %v\n", err)
		os.Exit(1)
	}

	if fi.IsDir() {
		stat = true
	} else {
		stat = false
	}
	return stat
}
