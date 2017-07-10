package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/alexmullins/zip"
)

var (
	zipfile    string
	dictionary string
)

func init() {
	flag.StringVar(&zipfile, "f", "", "Open zipfile")
	flag.StringVar(&dictionary, "d", "", "Open pass dictionary")
}

func main() {
	// разбор флагов
	flag.Parse()

	// вывод справки
	if zipfile == "" || dictionary == "" {
		println("Please " + os.Args[0] + " -h")
		os.Exit(0)
	}

	// запароленный архив
	zipr, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}

	// парольный словарь
	// ../00_addMaterials/dict.txt
	dictFile, err := os.Open(dictionary)
	if err != nil {
		log.Fatalln(err)
	}
	defer dictFile.Close()

	scanner := bufio.NewScanner(dictFile)
	for scanner.Scan() {
		pass := scanner.Text()
		for _, z := range zipr.File {
			z.SetPassword(pass)
			_, err := z.Open()
			// если все ок
			if err == nil {
				println("[+] Found password")
				println("[+] Password = " + pass)
				os.Exit(0)
			}
		}
	}
}
