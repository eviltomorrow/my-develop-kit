package finder

import (
	"github.com/tidwall/gjson"
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
func FindKV(jsonStr string, primaryKey []string, key string) []*Value {
	return nil
}
