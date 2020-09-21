package httpclient

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

//
var (
	ErrIndexOutOfKeys   = errors.New("Index out of keys")
	ErrKeyNotFound      = errors.New("Not found for key")
	ErrNotOneJSONObject = errors.New("Not one JSON object")
)

// GetJSON get json slice
func GetJSON(jsonStr, key string) ([]string, error) {
	var result = gjson.Parse(jsonStr)
	if !result.IsObject() {
		return nil, ErrNotOneJSONObject
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

	return getKV([]gjson.Result{result}, 0, keys, nil)
}

func getKV(results []gjson.Result, index int, keys []string, data []string) ([]string, error) {
	if index > len(keys)-1 {
		return data, nil
	}

	if len(results) == 0 {
		return data, nil
	}

	var key = keys[index]
	for _, result := range results {
		var result = result.Get(key)
		if !result.Exists() {
			return nil, ErrKeyNotFound
		}

		switch {
		case result.IsObject():
			index++
			cache, err := getKV([]gjson.Result{result}, index, keys, data)
			if err != nil {
				return nil, err
			}
			data = append(data, cache...)

		case result.IsArray():
			if index == len(keys)-1 {
				data = append(data, result.String())
			} else {
				index++
				var count = result.Get("#").Int()
				var i int64
				var collections = make([]gjson.Result, 0, count)
				for ; i < count; i++ {
					collections = append(collections, result.Get(fmt.Sprintf("%d", i)))
				}

				cache, err := getKV(collections, index, keys, data)
				if err != nil {
					return nil, err
				}
				data = append(data, cache...)
			}
		default:
			if data == nil {
				data = make([]string, 0, 128)
			}

			if index == len(keys)-1 {
				data = append(data, result.String())
			} else {
				index++
				cache, err := getKV([]gjson.Result{result}, index, keys, data)
				if err != nil {
					return nil, err
				}
				data = append(data, cache...)
			}
		}
	}
	return data, nil
}
