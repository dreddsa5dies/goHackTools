// CSPRNG - Cryptographically secure pseudo-random number generator

package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

func main() {
	fmt.Println("Random value:")

	// Generate a random int
	limit := int64(math.MaxInt64) // Highest random number allowed
	randInt, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Int:\t", randInt)

	// Alternatively, you could generate the random bytes
	// and turn them into the specific data type needed.
	// binary.Read() will only read enough bytes to fill the data type
	var number uint32
	err = binary.Read(rand.Reader, binary.BigEndian, &number)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Uint32:\t", number)

	// Or just generate a random byte slice
	numBytes := 4
	randomBytes := make([]byte, numBytes)
	rand.Read(randomBytes)
	fmt.Println("Byte:\t", randomBytes)
}
