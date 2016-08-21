package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
)

var width = flag.Int("w", 80, "")

func main() {
	flag.Parse()
	files := flag.Args()[0:]
	if len(files) == 0 {
		if err := view(os.Stdin); err != nil {
			panic(err)
		}
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fold: %v\n", err)
				continue
			}
			if err := view(f); err != nil {
				panic(err)
			}
			f.Close()
		}
	}
}

func view(f *os.File) error {
	in := bufio.NewScanner(f)
	in.Split(scanNColumns)
	for in.Scan() {
		fmt.Println(in.Text())
	}
	return in.Err()
}

func scanNBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) >= *width {
		// +1: '\n'
		if i := bytes.IndexByte(data[0:*width+1], '\n'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		return *width, data[0:*width], nil
	}
	return len(data), data, nil
}

func scanNRunes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	s := string(data)
	if utf8.RuneCountInString(s) >= *width {
		rs := []rune(s)
		bs := []byte(string(rs[0:*width]))
		// XXX
		if i := bytes.IndexByte(bs, '\n'); i >= 0 {
			return i + 1, bs[0:i], nil
		}
		return len(bs), bs, nil
	}
	return len(data), data, nil
}

var halfWidth = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0xff61, 0xffbe, 1},
		{0xffc2, 0xffc7, 1},
		{0xffca, 0xffcf, 1},
		{0xffd2, 0xffd7, 1},
		{0xffda, 0xffdc, 1},
		{0xffe9, 0xffee, 1},
	},
}

func scanNColumns(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// full ASCII characters.
	s := string(data)
	if len(data) == utf8.RuneCountInString(s) {
		return scanNBytes(data, atEOF)
	}

	if *width >= utf8.RuneCountInString(s) {
		return len(data), data, nil
	}

	// ASCII characters and non ASCII characters.
	rs := ([]rune(s))[0:*width]
	full := 0
	half := 0
	index := 0
	for _, r := range rs {
		bs := []byte(string(r))
		if bs[0] < utf8.RuneSelf || unicode.In(r, halfWidth) {
			half++
			index++
			if (full*2 + half) >= *width {
				break
			}
		} else {
			full++
			index++
			if *width == 1 {
				break
			}
			if (full*2 + half) == *width {
				break
			}
			if (full*2 + half) > *width {
				index--
				break
			}
		}
	}

	bs := []byte(string(rs[0:index]))
	if i := bytes.IndexByte(bs, '\n'); i >= 0 {
		return i + 1, bs[0:i], nil
	}
	return len(bs), bs, nil
}
