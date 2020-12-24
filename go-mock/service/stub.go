package service

import "fmt"

// Say say
func Say() {
	fmt.Println("I am Tom")
}

// WapperFunc wapper func
func WapperFunc() {
	Say()
}
