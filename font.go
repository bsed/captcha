package captcha

import (
	"encoding/gob"
	"os"
)

const (
	fontWidth  = 15
	fontHeight = 25
	blackChar  = 1
)

type Font struct {
	data map[rune][]byte
}

var fonts = make(map[string]*Font)
var selectedFont string

func LoadFontFromFile(fname string) *Font {
	f := make(map[rune][]byte)

	file, err := os.Open(fname)
	if err != nil {
		return nil
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&f)
	return &Font{f}
}

func AddFont(name string, f *Font) {
	fonts[name] = f
	if len(fonts) == 1 {
		SelectFont(name)
	}
}

func SelectFont(name string) {
	_, ok := fonts[name]
	if ok {
		selectedFont = name
	}
}

func Digit2Rune(d byte) rune {
	switch {
	case 0 <= d && d <= 9:
		return rune(d) + '0'
	case 10 <= d && d <= 10+byte('Z'-'A'):
		return rune(d) + 'A' - 10
	case 11+byte('Z'-'A') <= d && d <= 11+byte('Z'-'A')+byte('z'-'a'):
		return rune(d) - 'Z' + 'A' + 'a' - 11
	}
	return 0
}

func Rune2Digit(c rune) byte {
	switch {
	case '0' <= c && c <= '9':
		return byte(c - '0')
	case 'A' <= c && c <= 'Z':
		return byte(c - 'A' + 10)
	case 'a' <= c && c <= 'z':
		return byte(c - 'a' + 'Z' - 'A' + 11)
	}
	return 0
}

func getChar(d byte) []byte {
	if selectedFont == "" {
		panic("No font selected")
	}
	r := Digit2Rune(d)
	return fonts[selectedFont].data[r]
}
