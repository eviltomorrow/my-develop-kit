package finder

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

//
var (
	ErrNotJSONObject = fmt.Errorf("Not JSON Object")
)

// KeyNode key node
type KeyNode struct {
	Name        string      `json:"name"`         // 键名
	IsRoot      bool        `json:"is_root"`      // 是否顶结点
	Level       int         `json:"level"`        // 层级
	Data        interface{} `json:"data"`         // 值
	ChildNode   *KeyNode    `json:"child_node"`   // 子结点
	BrotherNode *KeyNode    `json:"brother_node"` // 兄弟结点
}

func (k *KeyNode) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}

// BuildKeyTree build tree
func BuildKeyTree(keys []string) (*KeyNode, error) {
	var rootNode = &KeyNode{
		IsRoot: true,
		Level:  0,
		Name:   "root",
	}

	for _, key := range keys {
		var cache = make([]string, 0, 16)
		var begin int
		var flag bool
		for i := 0; i < len(key); i++ {
			if key[i] == '.' && !flag {
				cache = append(cache, key[begin:i])
				begin = i + 1
			}
			if key[i] == '\\' {
				flag = true
			} else {
				flag = false
			}
		}

		if begin != len(key) {
			cache = append(cache, key[begin:])
		}

		rootNode = buildNode(rootNode, 0, cache)
	}
	return rootNode, nil
}

func buildNode(rootNode *KeyNode, index int, keys []string) *KeyNode {
	var currentNode = rootNode
	for i, key := range keys {
		if currentNode.ChildNode == nil {
			currentNode.ChildNode = &KeyNode{
				Name:  key,
				Level: i,
			}
			currentNode = currentNode.ChildNode
		} else {
			currentNode = currentNode.ChildNode
			if currentNode.Name == key {
				continue
			}

			for currentNode.BrotherNode != nil {
				currentNode = currentNode.BrotherNode
				if currentNode.Name == key {
					continue
				}
			}
			currentNode.BrotherNode = &KeyNode{
				Name:  key,
				Level: i,
			}
			currentNode = currentNode.BrotherNode
		}
	}
	return rootNode
}
