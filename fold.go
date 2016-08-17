package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
	in.Split(scanNRunes)
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
		return len(bs), bs, nil
	}
	return len(data), data, nil
}
