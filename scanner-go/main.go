package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// demo1()
	// demo2()
	// demo3()
	// demo4()
	demo5()
}

func demo1() {
	var text = `t h is is commander shepard!`
	var scanner = bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func demo2() {
	var text = "abcdefghijklmn"
	var scanner = bufio.NewScanner(strings.NewReader(text))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%s %t\r\n", data, atEOF)
		return 0, nil, nil
	}
	scanner.Buffer(make([]byte, 0, 2), bufio.MaxScanTokenSize)
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Printf("text: %s\r\n", scanner.Text())
	}
}

func demo3() {
	var text = "abcdefghijklmn"
	var scanner = bufio.NewScanner(strings.NewReader(text))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%s %t\r\n", data, atEOF)
		return 0, nil, nil
	}
	scanner.Buffer(make([]byte, 2), 40)
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Printf("text: %s\r\n", scanner.Text())
	}
}

func demo4() {
	var text = "abcdefghijklmn"
	var scanner = bufio.NewScanner(strings.NewReader(text))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%s %t\r\n", data, atEOF)
		return 0, nil, bufio.ErrFinalToken
	}
	scanner.Buffer(make([]byte, 2), 2)
	scanner.Split(split)
	for scanner.Scan() {
		fmt.Printf("text: %s\r\n", scanner.Text())
	}
}

func demo5() {
	input := "foo|bar"
	scanner := bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '|'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data[:len(data)], nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	for scanner.Scan() {
		// if scanner.Text() != "" {
		fmt.Println(scanner.Text())
		// }
	}
}
