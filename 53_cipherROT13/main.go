// Шифр ROT13 – часть шифра Цезаря с позицией 13. Особенность ROT13
// заключена в принципе инволюции при которой не нужен переключатель режимов
// шифрования/расшифрования.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Println(`ROT13 ("rotate by 13 places", sometimes hyphenated ROT-13)`)
		fmt.Println("Usage: " + os.Args[0] + " <string>")
		os.Exit(1)
	}

	input := os.Args[1]

	mapped := strings.Map(rot13, input)
	fmt.Println(input)
	fmt.Println(mapped)
}

func rot13(r rune) rune {
	if r >= 'a' && r <= 'z' {
		// Rotate lowercase letters 13 places.
		if r >= 'm' {
			return r - 13
		} else {
			return r + 13
		}
	} else if r >= 'A' && r <= 'Z' {
		// Rotate uppercase letters 13 places.
		if r >= 'M' {
			return r - 13
		} else {
			return r + 13
		}
	}
	// Do nothing.
	return r
}
