package finder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

//
var (
	ErrKeyNotFound   = errors.New("Not Found Key")
	ErrNotJSONObject = errors.New("Not JSON Object")
)

// Value value
type Value struct {
	ExpExpression string
	PrimaryKey    []string
	PrimaryValue  []string
	Key           string
	Value         string

	jsonResult *gjson.Result
}

// FindKV find kv
func FindKV(jsonStr string, primaryKey []string, key string) ([]*Value, error) {
	var result = gjson.Parse(jsonStr)
	if !result.IsObject() {
		return nil, ErrNotJSONObject
	}

	for _, pk := range primaryKey {
		if strings.LastIndex(pk, ".") == -1 {
			if !strings.HasPrefix(key, pk) {
				return nil, fmt.Errorf("Key[%s] not has prifix with primaryKey[%s]", key, pk)
			}
		} else {
			var s = pk[:strings.LastIndex(pk, ".")]
			if !strings.HasPrefix(key, pk) {
				return nil, fmt.Errorf("Key[%s] not has prifix with primaryKey[%s]", key, s)
			}
		}
	}

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

	return nil, nil
}

func getKV(results []gjson.Result, index int, keys []string) (map[string]string, error) {

	return nil, nil
}
