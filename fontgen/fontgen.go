package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"os/exec"
)

var thresh int = 200

var font = flag.String("font", "Monospace", "Font name")
var size = flag.String("size", "30", "Font size")

func main() {
	flag.Parse()
	fname := flag.Arg(0)
	if fname == "" {
		usage()
		os.Exit(1)
	}
	fm := make(map[rune][]byte)

	for c := '0'; c <= '9'; c++ {
		fm[c] = genChar(c)
	}
	for c := 'a'; c <= 'z'; c++ {
		fm[c] = genChar(c)
	}
	for c := 'A'; c <= 'Z'; c++ {
		fm[c] = genChar(c)
	}

	file, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(fm)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
}

func convertImage(file io.Reader) []byte {
	m, _, err := image.Decode(file)

	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	xp := bounds.Max.X - bounds.Min.X
	yp := bounds.Max.Y - bounds.Min.Y
	a := make([]byte, yp*xp)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := m.At(x, y)
			r, g, b, _ := col.RGBA()
			gr := (299*r + 587*g + 114*b + 500) / 1000
			gr >>= 8

			px := 0
			if int(gr) <= thresh {
				px = 1
			}
			idx := xp*(y-bounds.Min.Y) + x - bounds.Min.X
			a[idx] = byte(px)
		}
	}
	return a
}

func genChar(char rune) []byte {
	c := fmt.Sprintf("%c", char)
	cmd := genImage(c, *size, *font)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	a := convertImage(stdout)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return a
}

func genImage(letter, size, font string) *exec.Cmd {
	return exec.Command("convert", "-size", "15x25", "-pointsize", size, "-background", "white", "-fill", "black", "-filter", "Point", "-font", font, "-gravity", "center", "label:"+letter, "gif:-")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: fontgen [flags] filename\n")
	flag.PrintDefaults()
}
