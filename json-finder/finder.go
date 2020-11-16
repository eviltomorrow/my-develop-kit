package finder

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

//
const (
	ParentNodeType = "parent"
	ChildNodeType  = "child"
)

//
var (
	ErrKeyNotFound = errors.New("Not found the key")
)

// Key string
type Key struct {
	Feilds     []string `json:"feilds"`
	K          string   `json:"k"`
	V          string   `json:"v"`
	ParentKeys []*Key   `json:"parent_keys"`
	IsFind     bool     `json:"is_find"`
	ValueType  string   `json:"value_type"`

	Err      error  `json:"error"`
	nodeType string `json:"-"`
}

// Len len
func (k *Key) Len() int {
	return len(k.ParentKeys)
}

// Less less
func (k *Key) Less(i, j int) bool {
	return len(k.ParentKeys[i].Feilds) < len(k.ParentKeys[j].Feilds)
}

// Swap swap
func (k *Key) Swap(i, j int) {
	k.ParentKeys[i], k.ParentKeys[j] = k.ParentKeys[j], k.ParentKeys[i]
}

// E build error
func (k *Key) E(err error) {
	k.Err = err
}

// Find find
func (k *Key) Find() {
	k.IsFind = true
}

func (k *Key) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(k)
	return string(buf)
}

// BuildKey build key
func BuildKey(parentKeys []string, valueKey string) (*Key, error) {
	var pks = make([]*Key, 0, 4)
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

		var k = &Key{
			Feilds: pk,
			K:      parentKey,

			nodeType: ParentNodeType,
		}
		pks = append(pks, k)
	}

	var k = &Key{
		Feilds:     vk,
		K:          valueKey,
		ParentKeys: pks,

		nodeType: ChildNodeType,
	}
	return k, nil
}

// GetKey get value with key
func GetKey(results []gjson.Result, level int, key *Key) (*Key, error) {
	if level > len(key.Feilds) {
		return key, nil
	}
	if len(results) == 0 {
		return key, nil
	}

	var cache = make([]*Key, 0, 8)
	var k = key.Feilds[level]
	for _, result := range results {
		var newKey = key
		if key.nodeType == ChildNodeType {
			newKey = deepCloneKey(key)
		}
		var current = result.Get(k)
		if !current.Exists() {
			key.Find()
			key.E(ErrKeyNotFound)
			return key, nil
		}

		switch {
		case result.IsObject():
			if level == len(key.Feilds)-1 {
				result.String()
			} else {

			}
		case result.IsArray():
		default:
		}
	}
	return nil, nil
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

	key, err := GetKey([]gjson.Result{result}, 0, key)
	if err != nil {
		return nil, err
	}

	fmt.Println(key)
	return nil, nil
}

func deepCloneKey(key *Key) *Key {
	var feilds = make([]string, 0, len(key.Feilds))
	for _, d := range key.Feilds {
		feilds = append(feilds, d)
	}

	var pks []*Key
	if key.nodeType == ParentNodeType {
		pks = make([]*Key, 0, len(key.ParentKeys))
		for _, pk := range key.ParentKeys {
			pks = append(pks, deepCloneKey(pk))
		}
	}

	var newKey = &Key{
		Feilds:     feilds,
		K:          key.K,
		V:          key.V,
		ParentKeys: pks,
		IsFind:     key.IsFind,
		ValueType:  key.ValueType,
		Err:        key.Err,
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
