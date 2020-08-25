package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	demo1()
}

func demo1() {
	var text = `this is commander shepard!`
	var scanner = bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer([]byte{1, 2, 3}, 2)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
