package main

import (
	"log"
	"os"
	"unsafe"
)

type user struct {
	id   int
	name string
	age  int
}

func main() {
	var v = 10
	log.Printf("bool: %v", unsafe.Sizeof(true))
	log.Printf("int32: %v", unsafe.Sizeof(int32(100)))
	log.Printf("int: %v", unsafe.Sizeof(int(1)))
	log.Printf("*T: %v", unsafe.Sizeof(&v))
	log.Printf("string: %v", unsafe.Sizeof("string"))
	log.Printf("[]int: %v", unsafe.Sizeof([]int{v}))
	log.Printf("map: %v", unsafe.Sizeof(make(map[string]string)))
	log.Printf("func: %v", unsafe.Sizeof(os.OpenFile))
	log.Printf("chan: %v", unsafe.Sizeof(make(chan string)))
	var i interface{}
	log.Printf("interface{}: %v", unsafe.Sizeof(i))

	log.Printf("alignof string: %v", unsafe.Alignof("a"))
	var u1 = user{
		// id: 10,
	}
	log.Printf("offsetof string: %v", unsafe.Offsetof(u1.age))

	var u2 = &user{}
	log.Printf("pointer: %v", unsafe.Pointer(u2))

}
