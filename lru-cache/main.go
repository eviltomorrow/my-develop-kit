package main

import (
	"fmt"

	"github.com/eviltomorrow/my-develop-kit/lru-cache/cache"
)

func main() {
	var lru = cache.NewLRU(3)

	fmt.Println(lru)
	lru.Put("a", 1)
	fmt.Println(lru)
	lru.Put("a", 1)
	fmt.Println(lru)
	lru.Put("b", 2)
	fmt.Println(lru)
	lru.Put("c", 3)
	fmt.Println(lru)
	lru.Put("d", 4)
	fmt.Println(lru)

	var value = lru.Get("a")
	fmt.Println(value == nil)

	value = lru.Get("b")
	fmt.Println(value)
	fmt.Println(lru)
}
