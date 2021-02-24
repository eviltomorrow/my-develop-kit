package main

import (
	"fmt"
	"reflect"
)

func main() {
	demo1()
	demo2()
	demo3()
}

func demo1() {
	var num interface{} = 100
	var numType = reflect.TypeOf(num)
	fmt.Printf("Num=100, kind: %v\r\n", numType.Kind())
	fmt.Printf("Num=100, Name: %v\r\n", numType.Name())
	fmt.Printf("Num=100, PkgPath: %v\r\n", numType.PkgPath())
	fmt.Printf("Num=100, Size: %v\r\n", numType.Size())
	fmt.Printf("Num=100, Bits: %d\r\n", numType.Bits())
}

func demo2() {
	var p1 = &Person{
		Name:    "Shepaprd",
		Age:     20,
		Ability: map[string]string{"Run": "Yes"},
		Child:   nil,
	}

	var p1Type = reflect.TypeOf(p1)
	fmt.Printf("Person = p1, kind: %v\r\n", p1Type.Kind())
	fmt.Printf("Person = p1, name: %v\r\n", p1Type.Name())
	fmt.Printf("Person = p1, elem: %v\r\n", p1Type.Elem())

	var p1Elem = p1Type.Elem()
	fmt.Printf("Person = p1, NumFeild: %v\r\n", p1Elem.NumField())
	fmt.Printf("Person = p1, NumMethod: %v\r\n", p1Type.NumMethod())

	var p *Person = nil
	fmt.Printf("P == nil, %t\r\n", p == nil)

	p = (*Person)(nil)
	fmt.Printf("P == nil, %t\r\n", p == nil)
}

func demo3() {
	var text = "This is shepard"
	var textValue = reflect.ValueOf(text)
	fmt.Printf("Text: %v\r\n", textValue)
	fmt.Printf("Text String(): %v\r\n", textValue.String())

	var num = 2
	var numValue = reflect.ValueOf(num)
	fmt.Printf("Num: %v\r\n", numValue)
	fmt.Printf("Num String(): %v\r\n", numValue.String())
	fmt.Printf("Num Interface(): %v\r\n", numValue.Interface())

}

// Person p
type Person struct {
	Name    string
	Age     int
	Ability map[string]string
	Child   []string
}

// CanRun can run
func (p *Person) CanRun() bool {
	return true
}

// CanFly can fly
func (p Person) CanFly() bool {
	return false
}
