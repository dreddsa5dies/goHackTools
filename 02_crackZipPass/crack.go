package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/alexmullins/zip"
)

var (
	zipfile     string
	dictionary  string
	concurrency int
)

func init() {
	flag.StringVar(&zipfile, "f", "", "Open zipfile")
	flag.StringVar(&dictionary, "d", "", "Open pass dictionary")
	flag.IntVar(&concurrency, "c", runtime.NumCPU(), "Number of workers to use")
}

func main() {
	// разбор флагов
	flag.Parse()

	// вывод справки
	if zipfile == "" || dictionary == "" {
		println("Please " + os.Args[0] + " -h")
		os.Exit(0)
	}

	word := make(chan string, 0)
	found := make(chan string, 0)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go zipCrackWorker(word, found, &wg)
	}

	// парольный словарь
	// ../00_addMaterials/dict.txt
	dictFile, err := os.Open(dictionary)
	if err != nil {
		log.Fatalln(err)
	}
	defer dictFile.Close()

	scanner := bufio.NewScanner(dictFile)
	go func() {
		for scanner.Scan() {
			pass := scanner.Text()
			word <- pass
		}
		close(word)
	}()

	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case f := <-found:
		println("[+] Found password")
		println("[+] Password =", f)
		return
	case <-done:
		println("[+] Password not found")
		return
	}

}

func zipCrackWorker(word <-chan string, found chan<- string, wg *sync.WaitGroup) {
	// запароленный архив
	zipr, err := zip.OpenReader(zipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer zipr.Close()
	defer wg.Done()
	for w := range word {
		for _, z := range zipr.File {
			z.SetPassword(w)
			_, err := z.Open()
			// если все ок
			if err == nil {
				found <- w
			}
		}
	}
}
