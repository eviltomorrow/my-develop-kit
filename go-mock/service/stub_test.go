package service

import (
	"fmt"
	"testing"

	"github.com/prashantv/gostub"
)

func TestStubSay(t *testing.T) {
	var stubSay = func() {
		fmt.Println("I am John")
	}

	var say = Say
	stubs := gostub.Stub(&say, stubSay)
	defer stubs.Reset()

	WapperFunc()
}
