package httpclient

import (
	"errors"
	"fmt"
	"strings"

	xml2json "github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
)

//
var (
	ErrIndexOutOfKeys   = errors.New("Index out of keys")
	ErrKeyNotFound      = errors.New("Not found for key")
	ErrNotOneXMLObject  = errors.New("Not one XML object")
	ErrNotOneJSONObject = errors.New("Not one JSON object")
)

// GetText get text
func GetText(textStr string) (string, error) {
	return textStr, nil
}

// GetXML get xml
func GetXML(xmlStr, key string) ([]string, error) {
	// TODO 不严谨
	if !strings.HasPrefix(xmlStr, "<") && !strings.HasSuffix(xmlStr, ">") {
		return nil, ErrNotOneXMLObject
	}

	data, err := xml2json.Convert(strings.NewReader(xmlStr))
	if err != nil {
		return nil, fmt.Errorf("XML to JSON failure, nest error: %v", err)
	}

	if xmlStr != "" && data.Len() == 3 {
		return nil, fmt.Errorf("XML to JSON failure, nest error: panic invalid xml string, msg: %v", xmlStr)
	}

	return GetJSON(data.String(), key)
}

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

	var data = make([]string, 0, 8)
	return getKV([]gjson.Result{result}, 0, keys, data)
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
