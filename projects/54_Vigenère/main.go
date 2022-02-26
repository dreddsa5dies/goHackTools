// Шифр Виженера (фр. Chiffre de Vigenère) — метод полиалфавитного шифрования
// буквенного текста с использованием ключевого слова.
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Chiffre de Vigenère")
	fmt.Print("key > ")
	var key string
	_, err := fmt.Scanf("%s", &key)
	if err != nil {
		log.Fatalln(err)
	}

	var message string
	fmt.Print("message > ")
	_, err = fmt.Scanf("%s", &message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print("encipher > ")
	fmt.Println(encipher(message, key))
	fmt.Print("decipher > ")
	fmt.Println(decipher(encipher(message, key), key))
}

func sanitize(in string) string {
	out := []rune{}
	for _, v := range in {
		if 65 <= v && v <= 90 {
			out = append(out, v)
		} else if 97 <= v && v <= 122 {
			out = append(out, v-32)
		}
	}

	return string(out)
}

func quartets(in string) string {
	out := make([]rune, 0, len(in))
	for i, v := range in {
		if i%4 == 0 && i != 0 {
			out = append(out, rune(32))
		}
		out = append(out, v)
	}
	return string(out)
}

func encodePair(a, b rune) rune {
	return (((a - 'A') + (b - 'A')) % 26) + 'A'
}

func decodePair(a, b rune) rune {
	return (((((a - 'A') - (b - 'A')) + 26) % 26) + 'A')
}

func encipher(msg, key string) string {
	smsg, skey := sanitize(msg), sanitize(key)
	out := make([]rune, 0, len(msg))
	for i, v := range smsg {
		out = append(out, encodePair(v, rune(skey[i%len(skey)])))
	}
	return string(out)
}

func decipher(msg, key string) string {
	smsg, skey := sanitize(msg), sanitize(key)
	out := make([]rune, 0, len(msg))
	for i, v := range smsg {
		out = append(out, decodePair(v, rune(skey[i%len(skey)])))
	}
	return string(out)
}
