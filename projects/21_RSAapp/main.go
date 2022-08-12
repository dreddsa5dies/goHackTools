// use https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec // так надо
	"crypto/rand"
	"encoding/hex"
	"io"
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
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) == 1 {
		log.Fatalf("Usage:\t%v -h\n", os.Args[0])
	}

	if opts.Data != "" {
		regStr := regexp.MustCompile(`([0-9a-zA-Z]){8,}`)

		if regStr.MatchString(opts.Password) {
			log.Println("Pass ok")
			encryptFile(opts.FileName, []byte(opts.Data), opts.Password)
			log.Println("Encrypt ok")
		} else {
			log.Fatalln("Bad password\nUse good password")
		}
	} else if opts.Password != "" {
		log.Printf("Decrypt:\n%v\n", string(decryptFile(opts.FileName, opts.Password)))
	}

	os.Exit(0)
}

// Hashing Passwords to Compatible Cipher Keys
func createHash(key string) string {
	hasher := md5.New() //nolint:gosec // так надо
	hasher.Write([]byte(key))

	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("EncryptFunc NewGCM: %v\n", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("EncryptFunc NonceSize: %v\n", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("DecryptFunc NewCipher: %v\n", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("DecryptFunc NewGCM: %v\n", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("PASSWORD FAIL\n%v\n", err)
	}

	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("encryptFileFunc Create: %v\n", err)
	}
	defer f.Close()

	_, err = f.Write(encrypt(data, passphrase))
	if err != nil {
		log.Panicln(err)
	}
}

func decryptFile(filename, passphrase string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("decryptFileFunc ReadFile: %v\n", err)
	}

	return decrypt(data, passphrase)
}
