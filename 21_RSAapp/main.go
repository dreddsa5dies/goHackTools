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
	"regexp"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Data     string `short:"d" long:"data" default:"" description:"Data for encrypt"`
	Password string `short:"p" long:"pass" default:"" description:"Password"`
	FileName string `short:"f" long:"filename" default:"encrypt.txt" description:"Save|Open filename"`
	// Open     bool   `short:"o" long:"open" default:"false" description:"Open filename"`
}

func main() {
	flags.Parse(&opts)

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage:\t%v -h\n", os.Args[0])
		os.Exit(1)
	}

	if opts.Data != "" {
		regStr, _ := regexp.Compile(`([0-9a-zA-Z]){8,}`)
		if regStr.MatchString(opts.Password) {
			log.Println("Pass ok")
			encryptFile(opts.FileName, []byte(opts.Data), opts.Password)
			log.Println("Encrypt ok")
		} else {
			fmt.Fprintf(os.Stdout, "Bad password\n")
			fmt.Fprintf(os.Stdout, "Use good password\n")
			os.Exit(1)
		}
	} else if opts.Password != "" {
		fmt.Fprintf(os.Stdout, "Decrypt:\n%v\n", string(decryptFile(opts.FileName, opts.Password)))
	}

	os.Exit(0)
}

// Hashing Passwords to Compatible Cipher Keys
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
