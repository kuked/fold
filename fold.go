package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]
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
	in.Split(scan10Bytes)
	for in.Scan() {
		fmt.Println(in.Text())
	}
	return in.Err()
}

func scan10Bytes(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) >= 10 {
		return 10, data[0:10], nil
	}
	return len(data), data, nil
}
