package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	"rsc.io/pdf"
)

func main() {
	// справка
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Использование: %s FILE.pdf\n", os.Args[0])
		os.Exit(1)
	}

	printMeta(os.Args[1])
}

func printMeta(filename string) {
	pdfFile, err := pdf.Open(filename)
	if err != nil {
		log.Println(err)
	}

	docInfo := pdfFile.Trailer().Key("Info").String()

	docInfo = strings.TrimLeft(docInfo, "<<")
	docInfo = strings.TrimRight(docInfo, ">>")
	info := strings.Split(docInfo, "/")
	for _, v := range info {
		println(v)
	}
}
