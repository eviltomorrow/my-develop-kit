package main

import "fmt"

func main() {
	var ch = make(chan struct{})
	var str = "abcdefghijklmn"

	go func() {
		for i := 0; i < len(str); i++ {
			<-ch
			if i%2 == 0 {
				fmt.Printf("go-2: %c\r\n", str[i])
			}
		}
	}()

	go func() {
		for i := 0; i < len(str); i++ {
			ch <- struct{}{}
			if i%2 == 1 {
				fmt.Printf("go-1: %c\r\n", str[i])
			}
		}
	}()
	for {
	}
}
