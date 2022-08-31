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
	var flag bool

	switch {
	case r >= 'a' && r <= 'z':
		if r >= 'm' {
			flag = true
		}
	case r >= 'A' && r <= 'Z':
		if r >= 'M' {
			flag = true
		}
	default:
		return r
	}

	if flag {
		return r - 13
	}

	return r + 13
}
