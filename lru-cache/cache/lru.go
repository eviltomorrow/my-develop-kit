package cache

import (
	"bytes"
	"fmt"
)

// LRU cache
type LRU struct {
	size    int
	cap     int
	mapping map[string]*Node

	head *Node
	tail *Node
}

// NewLRU new lru
func NewLRU(cap int) *LRU {
	if cap <= 0 {
		panic("Panic: create lru failure, cap <= 0")
	}
	lru := &LRU{
		size:    0,
		cap:     cap,
		mapping: make(map[string]*Node, cap),

		head: &Node{},
		tail: &Node{},
	}

	lru.head.next = lru.tail
	return lru
}

// Get get key
func (l *LRU) Get(key string) interface{} {
	node, ok := l.mapping[key]
	if !ok {
		return nil
	}

	l.moveToHead(node)

	return node.Value
}

// Put put key value
func (l *LRU) Put(key string, value interface{}) {
	node, ok := l.mapping[key]
	if ok {
		l.moveToHead(node)
		node.Value = value
	} else {
		l.eliminateTail()

		var node = &Node{Key: key, Value: value}
		var temp = l.head.next
		l.head.next = node
		node.pre = l.head
		node.next = temp
		temp.pre = node
		l.size++
		l.mapping[key] = node
	}
}

func (l *LRU) moveToHead(node *Node) {
	node.pre.next = node.next
	node.next.pre = node.pre

	var temp = l.head.next
	l.head.next = node
	node.next = temp
	temp.pre = node
	node.pre = l.head
}

func (l *LRU) eliminateTail() {
	if l.size < l.cap {
		return
	}
	var last = l.tail.pre
	last.pre.next = l.tail
	l.tail.pre = last.pre
	l.size--
	delete(l.mapping, last.Key)
	last = nil
}

func (l *LRU) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("LRU cap: %d, size: %d\r\n", l.cap, l.size))
	var format = " --> loc:%2d {key: %s, value: %v}\r\n"
	if l.head.next == l.tail {
		buf.WriteString("  Empty lru cache\r\n")
	} else {
		var count int
		for cur := l.head.next; cur != l.tail; cur = cur.next {
			count++
			buf.WriteString(fmt.Sprintf(format, count, cur.Key, cur.Value))
		}
	}
	return buf.String()
}

// Node node
type Node struct {
	Key   string
	Value interface{}

	pre  *Node
	next *Node
}
