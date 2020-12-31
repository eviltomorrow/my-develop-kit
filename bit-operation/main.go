package main

import "fmt"

//
const (
	Create = 1 << iota
	Update
	Delete
)

func main() {
	fmt.Println("Create: ", Create)
	fmt.Println("Update: ", Update)
	fmt.Println("Delete: ", Delete)

	fmt.Println("Create and Update: ", Create|Update)
	fmt.Println("Create and Delete: ", Create|Delete)
	fmt.Println("Delete and Update: ", Delete|Update)
	fmt.Println("Create, Update, Delete: ", Create|Update|Delete)
}
