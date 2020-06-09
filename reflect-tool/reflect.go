package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//
var (
	TagKey = "json"
)

// FindObjectByTag find object
func FindObjectByTag(tag []string, object interface{}) (rt reflect.Type, rv reflect.Value, err error) {
	defer func() {
		if e1 := recover(); e1 != nil {
			err = fmt.Errorf("%v", e1)
			return
		}
	}()

	var otype = reflect.TypeOf(object)
	var ovalue = reflect.ValueOf(object)
	return find(tag, otype, ovalue)
}

func find(tag []string, pt reflect.Type, pv reflect.Value) (rt reflect.Type, rv reflect.Value, err error) {
	defer func() {
		if e1 := recover(); e1 != nil {
			err = fmt.Errorf("Panic: Unknown error: %v", e1)
			return
		}
	}()

	switch pt.Kind() {
	case reflect.Ptr:
		var elem = pt.Elem()
		for i := 0; i < elem.NumField(); i++ {
			var field = elem.Field(i)
			if tag[0] == field.Tag.Get(TagKey) {
				if len(tag) == 1 {
					return field.Type, pv.Elem().Field(i), nil
				}
				return find(tag[1:], field.Type, pv.Elem().Field(i))
			}
		}
		return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])
	case reflect.Struct:
		for i := 0; i < pt.NumField(); i++ {
			var field = pt.Field(i)
			if tag[0] == field.Tag.Get(TagKey) {
				if len(tag) == 1 {
					return field.Type, pv.Field(i), nil
				}
				return find(tag[1:], pv.Type(), pv.Field(i))
			}
		}
		return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])
	}
	return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])
}

// Set set
func Set(key string, value interface{}, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object value not support modified")
	}

	if rt.Kind() != reflect.ValueOf(value).Kind() {
		return fmt.Errorf("Type not match [expect: %v, actual: %v]", rt.Kind().String(), reflect.ValueOf(value).Kind())
	}

	if !reflect.TypeOf(value).AssignableTo(rt) {
		return fmt.Errorf("Assignable failure [expect: %v, actual: %v]", rt, reflect.TypeOf(value))
	}
	rv.Set(reflect.ValueOf(value))
	return nil
}

func funcInt(s string) (interface{}, error) {
	result, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcInt32(s string) (interface{}, error) {
	result, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, err
	}
	return int32(result), nil
}

func funcInt64(s string) (interface{}, error) {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcFloat32(s string) (interface{}, error) {
	result, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcFloat64(s string) (interface{}, error) {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcString(s string) (interface{}, error) {
	return s, nil
}

func funcBool(s string) (interface{}, error) {
	result, err := strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcMapStringAndString(s string) (interface{}, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func funcSliceString(s string) (interface{}, error) {
	var result []string
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
