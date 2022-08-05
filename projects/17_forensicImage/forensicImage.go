// simple forensic image file use stlib Golang
// test gif, png, jpg

package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Printf("Use: %v fileUnknown fileUnknown fileUnknown ...\n", os.Args[0])
	} else {
		for _, file := range os.Args[1:] {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "file: %v\n", err)
				f.Close()
				continue
			}

			_, kind, err := image.Decode(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", f.Name(), err)
				f.Close()
				continue
			}
			fmt.Fprintln(os.Stdout, f.Name(), "Format =", kind)

			f.Close()
		}
	}
}
