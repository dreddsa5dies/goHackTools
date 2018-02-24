// use https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
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
)

func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Printf("Use: %v pass data_string\n", os.Args[0])
		os.Exit(1)
	} else {
		log.Println("Starting the application...")
		ciphertext := encrypt([]byte(os.Args[2]), os.Args[1])
		log.Printf("Encrypted: %x\n", ciphertext)

		plaintext := decrypt(ciphertext, os.Args[1])
		log.Printf("Decrypted: %s\n", plaintext)

		encryptFile("encrypt.txt", []byte(os.Args[2]), os.Args[1])
		log.Println(string(decryptFile("encrypt.txt", os.Args[1])))
		os.Exit(0)
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
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
		fmt.Fprintf(os.Stderr, "DecryptFunc gcm.Open: %v\n", err)
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
