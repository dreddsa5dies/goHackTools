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
				continue
			}
			defer f.Close()

			_, kind, err := image.Decode(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v\n", f.Name(), err)
				continue
			}
			fmt.Fprintln(os.Stderr, f.Name(), "Format =", kind)
		}
	}
}
