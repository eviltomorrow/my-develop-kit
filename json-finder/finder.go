package finder

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

//
const (
	ParentNodeType  = "parent"
	ChildNodeType   = "child"
	BrotherNodeType = "brother"
)

//
var (
	ErrKeyNotFound = errors.New("Not found the key")
)

// Key string
type Key struct {
	Fields      []string `json:"fields"`
	K           string   `json:"k"`
	V           string   `json:"v"`
	IsFind      bool     `json:"is_find"`
	ParentKeys  []*Key   `json:"parent_keys"`
	BrotherKeys []*Key   `json:"brother_keys"`

	Err      error `json:"error"`
	nodeType string
}

// Len len
func (k *Key) Len() int {
	return len(k.ParentKeys)
}

// Less less
func (k *Key) Less(i, j int) bool {
	return len(k.ParentKeys[i].Fields) < len(k.ParentKeys[j].Fields)
}

// Swap swap
func (k *Key) Swap(i, j int) {
	k.ParentKeys[i], k.ParentKeys[j] = k.ParentKeys[j], k.ParentKeys[i]
}

// Val v
func (k *Key) Val(v string) {
	k.V = v
}

// E build error
func (k *Key) E(err error) {
	k.Err = err
}

// Find find
func (k *Key) Find() {
	k.IsFind = true
}

// PrintPtr print ptr
func (k *Key) PrintPtr() {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Current key: %p\r\n", k))
	buffer.WriteString(fmt.Sprintf("\tFields: %p\r\n", k.Fields))
	buffer.WriteString("\tParentKeys: \r\n")

	for _, pk := range k.ParentKeys {
		buffer.WriteString(fmt.Sprintf("\t\tPK: %p\r\n", pk))
	}

	buffer.WriteString("\tBrotherKeys: \r\n")

	for _, bk := range k.BrotherKeys {
		buffer.WriteString(fmt.Sprintf("\t\tBK: %p\r\n", bk))
	}

	fmt.Printf("%s\r\n", buffer.String())
}

func (k *Key) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}

// BuildKey build key
func BuildKey(parentKeys []string, brotherKeys []string, valueKey string) (*Key, error) {
	var pks = make([]*Key, 0, len(parentKeys))
	var vk = parseKey(valueKey)
	for _, parentKey := range parentKeys {
		var pk = parseKey(parentKey)
		if len(pk) > len(vk) {
			return nil, fmt.Errorf("The ParentKey[%s] has a higher depth than ValueKey[%s]", parentKey, valueKey)
		}

		for i := 0; i < len(pk); i++ {
			if i == len(pk)-1 {
				break
			}

			if pk[i] != vk[i] {
				return nil, fmt.Errorf("The PrimaryKey[%s] and ValueKey[%s] not belonging to the same tree", parentKey, valueKey)
			}
		}

		pks = append(pks, &Key{
			Fields: pk,
			K:      parentKey,

			nodeType: ParentNodeType,
		})
	}

	var bks = make([]*Key, 0, len(brotherKeys))
	for _, brotherKey := range brotherKeys {
		var bk = parseKey(brotherKey)
		if len(bk) != len(vk) {
			return nil, fmt.Errorf("The BrotherKey[%s] has a different depth with ValueKey[%s]", brotherKey, valueKey)
		}

		for i, k := range vk {
			if i == len(vk)-1 {
				break
			}
			if bk[i] != k {
				return nil, fmt.Errorf("The BrotherKey[%s] and ValueKey[%s] not belonging to the same tree", brotherKey, valueKey)
			}
		}

		bks = append(bks, &Key{
			Fields: bk,
			K:      brotherKey,

			nodeType: BrotherNodeType,
		})
	}

	var k = &Key{
		Fields:      vk,
		K:           valueKey,
		ParentKeys:  pks,
		BrotherKeys: bks,

		nodeType: ChildNodeType,
	}

	sort.Sort(k)
	return k, nil
}

// GetKey get value with key
func GetKey(results []gjson.Result, level int, key *Key) ([]*Key, error) {
	if level > len(key.Fields) {
		return []*Key{}, nil
	}
	if len(results) == 0 {
		return []*Key{}, nil
	}

	var cache = make([]*Key, 0, 8)
	var k = key.Fields[level]
	for _, result := range results {
		var newKey = key
		if key.nodeType == ChildNodeType {
			newKey = deepCloneKey(key)
		}

		var current = result.Get(k)

		switch {
		case !current.Exists():
			newKey.E(ErrKeyNotFound)
			newKey.Find()

		case level == len(newKey.Fields)-1:
			newKey.Val(current.String())
			newKey.Find()

		default:
			// 计算基础值
			switch {
			case current.IsObject():
				level++
				data, err := GetKey([]gjson.Result{current}, level, newKey)
				if err != nil {
					return nil, err
				}
				level--
				cache = append(cache, data...)

			case current.IsArray():
				var count = current.Get("#").Int()
				var i int64
				var collections = make([]gjson.Result, 0, count)
				for ; i < count; i++ {
					collections = append(collections, current.Get(fmt.Sprintf("%d", i)))
				}

				level++
				data, err := GetKey(collections, level, newKey)
				if err != nil {
					return nil, err
				}
				level--
				cache = append(cache, data...)

			default:
				if level == len(newKey.Fields)-1 {
					newKey.Val(current.String())
				} else {
					newKey.E(ErrKeyNotFound)
				}
				newKey.Find()
			}
		}

		if newKey.IsFind {
			cache = append(cache, newKey)
		}

		for _, k := range cache {
			for _, pk := range k.ParentKeys {
				if !pk.IsFind && level >= 0 && level == len(pk.Fields)-1 {
					pk.Find()
					var current = result.Get(pk.Fields[level])
					if !current.Exists() {
						pk.E(ErrKeyNotFound)
					} else {
						pk.Val(current.String())
					}
				}
			}

			for _, bk := range k.BrotherKeys {
				if !bk.IsFind && level >= 0 && level == len(bk.Fields)-1 {
					bk.Find()
					var current = result.Get(bk.Fields[level])
					if !current.Exists() {
						bk.E(ErrKeyNotFound)
					} else {
						bk.Val(current.String())
					}
				}
			}
		}
	}

	return cache, nil
}

// FindKey find value with key
func FindKey(jsonStr string, key *Key) ([]string, error) {
	if key == nil {
		return nil, fmt.Errorf("The key is nil")
	}

	if !gjson.Valid(jsonStr) {
		return nil, fmt.Errorf("Invalid JSON")
	}

	var result = gjson.Parse(jsonStr)
	if !result.IsObject() {
		return nil, fmt.Errorf("Not a JSON Object")
	}

	cache, err := GetKey([]gjson.Result{result}, 0, key)
	if err != nil {
		return nil, err
	}

	for _, k := range cache {
		fmt.Printf("k: %v\r\n", k.String())
	}
	return []string{fmt.Sprintf("%v", cache)}, nil
}

func deepCloneKey(key *Key) *Key {
	var fields = make([]string, 0, len(key.Fields))
	for _, d := range key.Fields {
		fields = append(fields, d)
	}

	var pks []*Key
	var bks []*Key
	if key.nodeType == ChildNodeType {
		pks = make([]*Key, 0, len(key.ParentKeys))
		for _, pk := range key.ParentKeys {
			pks = append(pks, deepCloneKey(pk))
		}

		bks = make([]*Key, 0, len(key.BrotherKeys))
		for _, bk := range key.BrotherKeys {
			bks = append(bks, deepCloneKey(bk))
		}
	}

	var newKey = &Key{
		Fields:      fields,
		K:           key.K,
		V:           key.V,
		ParentKeys:  pks,
		BrotherKeys: bks,
		IsFind:      key.IsFind,
		nodeType:    key.nodeType,
		Err:         key.Err,
	}
	return newKey
}
func parseKey(key string) []string {
	var keys = make([]string, 0, 16)
	var begin int
	var flag bool

	for i := 0; i < len(key); i++ {
		if key[i] == '.' && !flag {
			keys = append(keys, key[begin:i])
			begin = i + 1
		}
		if key[i] == '\\' {
			flag = true
		} else {
			flag = false
		}
	}

	if begin != len(key) {
		keys = append(keys, key[begin:])
	}
	return keys
}
