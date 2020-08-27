package xshell

import (
	"fmt"
	"testing"
)

func TestOprationDB(t *testing.T) {
	var db = &db{
		Path:      "/tmp/xshell-resource.db",
		Resources: make([]*resource, 0, 20),
	}
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")
	var r1 = &resource{
		Username: "god",
		Password: "god",
		Host:     "127.0.0.1",
		Port:     22,
	}
	db.insert(r1)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	var r2 = &resource{
		Username: "shepard",
		Password: "shepard",
		Host:     "localhost",
		Port:     3306,
	}
	db.insert(r2)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	var r3 = &resource{
		Username: "root",
		Password: "root",
		Host:     "192.168.0.1",
		Port:     1521,
	}
	db.insert(r3)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	resources := db.query()
	fmt.Println(resources)

	resources = db.match("127.0.")
	fmt.Println(resources)
	fmt.Println("------------------------------------------------------")

	resources = db.match("127.0.1")
	fmt.Println(resources)
	fmt.Println("------------------------------------------------------")

	db.delete(1)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	db.delete(0)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	db.insert(r1)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	db.update(1, r2)
	fmt.Println(db.String())
	fmt.Println("------------------------------------------------------")

	err := db.dump()
	fmt.Println(err)
}
