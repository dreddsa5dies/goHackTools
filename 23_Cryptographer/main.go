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

	hashSaveFile, err := os.OpenFile("save_hash", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.OpenFile\n")
		os.Exit(1)
	}
	defer hashSaveFile.Close()

	fmt.Fprintf(os.Stdout, "Save IT >> %v\n", hashSaveFile.Name())
	hashSaveFile.WriteString(hash)

	fmt.Fprintln(os.Stdout, "-----------------------------------------------------------")

	cryptoDir(dir, hash)

	os.Exit(0)
}

func cryptoDir(dir string, hash string) {
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
			// рекурсивный проход по поддиректориям
			if fi.IsDir() {
				cryptoDir(dir+"/"+fi.Name(), hash)
			} else {
				// имя файла
				fmt.Fprintf(os.Stdout, "encrypt %v\n", fi.Name())
				file, err := ioutil.ReadFile(dir + "/" + fi.Name())
				if err != nil {
					fmt.Fprintf(os.Stderr, "ioutil.ReadFile: %v\n", err)
					os.Exit(1)
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

func encryptFile(filename string, data []byte, hash string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encryptFileFunc Create: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	f.Write(encrypt(data, hash))
}

/*
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

func decryptFile(filename string, passphrase string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decryptFileFunc ReadFile: %v\n", err)
		os.Exit(1)
	}
	return decrypt(data, passphrase)
}
*/

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
