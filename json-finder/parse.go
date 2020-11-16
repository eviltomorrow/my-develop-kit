package finder

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// GetValue get value with key
func GetValue(jsonStr string, key string) (string, error) {
	var result = gjson.Parse(jsonStr)
	if !result.IsObject() {
		return "", fmt.Errorf("Test")
	}

	return result.Get(key).String(), nil
}
