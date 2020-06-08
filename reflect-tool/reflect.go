package tool

import (
	"fmt"
	"reflect"
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

	fmt.Println(pt.Kind())
	switch pt.Kind() {
	case reflect.Ptr:
		var elem = pt.Elem()
		for i := 0; i < elem.NumField(); i++ {
			var field = elem.Field(i)
			if tag[0] == field.Tag.Get("json") {
				if len(tag) == 1 {
					return field.Type, pv.Elem().Field(i), nil
				}
				return find(tag[1:], pv.Type(), pv.Elem().Field(i))
			}
		}
		return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])
	case reflect.Struct:
		for i := 0; i < pt.NumField(); i++ {
			var field = pt.Field(i)
			if tag[0] == field.Tag.Get("json") {
				if len(tag) == 1 {
					return field.Type, pv.Field(i), nil
				}
				return find(tag[1:], pv.Type(), pv.Field(i))
			}
		}
		return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])

	default:

	}
	return rt, rv, fmt.Errorf("Not found specified object with tag[%v]", tag[0])
}
