package main

import (
	"fmt"
)

func main() {
	var x = 10
	var y = 20

	if x != y {
		x ^= y
		y ^= x
		x ^= y
	}

	fmt.Printf("x: %d, y: %d\r\n", x, y)

	x = 10
	y = 20
	x, y = y, x
	fmt.Printf("x: %d, y: %d\r\n", x, y)

	x = x + y
	y = x - y
	x = x - y
	fmt.Printf("x: %d, y: %d\r\n", x, y)

}
