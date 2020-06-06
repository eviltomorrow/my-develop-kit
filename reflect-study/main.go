package main

import (
	"fmt"
	"log"
	"reflect"
)

// Transcript t
type Transcript struct {
	English int
	Math    int
	Artical int
}

// User user
type User struct {
	ID         int
	Name       string
	Age        int
	Transcript *Transcript
}

// DoSomething do something
func (u User) DoSomething() {
	fmt.Println("do something!")
}

// DoAnything do anything
func (u User) DoAnything(count int, thing string) {
	fmt.Printf("count: %d, thing: %v\r\n", count, thing)
}

func main() {
	var num float32 = 1.23
	log.Printf("type: %v\r\n", reflect.TypeOf(num))
	log.Printf("value: %v\r\n", reflect.ValueOf(num))

	log.Printf("type: %v\r\n", reflect.TypeOf(&num))
	log.Printf("value: %v\r\n", reflect.ValueOf(&num))
	log.Printf("%v\r\n", reflect.ValueOf(&num).Interface().(*float32))
	log.Printf("%v\r\n", reflect.ValueOf(num).Interface().(float32))

	// 遍历属性和方法
	var user = User{
		ID:   1,
		Name: "shepard",
		Age:  21,
		Transcript: &Transcript{
			English: 90,
			Math:    0,
			Artical: 100,
		},
	}
	var userType = reflect.TypeOf(user)
	log.Printf("user type: %v\r\n", userType.Name())
	var userValue = reflect.ValueOf(user)
	log.Printf("user value: %v\r\n", userValue)

	for i := 0; i < userType.NumField(); i++ {
		var field = userType.Field(i)
		var value = userValue.Field(i).Interface()
		log.Printf("num field: %s: %v = %v, kind: %v\r\n", field.Name, field.Type, value, reflect.TypeOf(value).Kind())

	}

	for i := 0; i < userType.NumMethod(); i++ {
		var method = userType.Method(i)
		log.Printf("%s: %v\r\n", method.Name, method.Type)
	}

	// change valeu
	var change float64 = 1.234
	var pointer = reflect.ValueOf(&change)
	var newChange = pointer.Elem()
	log.Printf("type of newChange: %v\r\n", newChange.Type())
	log.Printf("set newChange: %v\r\n", newChange.CanSet())
	newChange.SetFloat(4.321)
	log.Printf("new change: %v\r\n", change)

	// call function
	var doSomethingMethod = userValue.MethodByName("DoSomething")
	var args = []reflect.Value{}
	doSomethingMethod.Call(args)

}
