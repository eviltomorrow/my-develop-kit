package tool

import (
	"fmt"
	"reflect"
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

// SetString set string
func SetString(key string, value string, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.String {
		return fmt.Errorf("Not a string type")
	}
	rv.SetString(value)
	return nil
}

// SetInt64 set int64
func SetInt64(key string, value int64, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Int64 {
		return fmt.Errorf("Not a int64 type")
	}
	rv.SetInt(value)
	return nil
}

// SetInt32 set int32
func SetInt32(key string, value int32, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Int32 {
		return fmt.Errorf("Not a int32 type")
	}
	rv.SetInt(int64(value))
	return nil
}

// SetInt16 set int16
func SetInt16(key string, value int32, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Int16 {
		return fmt.Errorf("Not a int16 type")
	}
	rv.SetInt(int64(value))
	return nil
}

// SetInt8 set int8
func SetInt8(key string, value int8, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Int8 {
		return fmt.Errorf("Not a int8 type")
	}
	rv.SetInt(int64(value))
	return nil
}

// SetInt set int
func SetInt(key string, value int, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Int {
		return fmt.Errorf("Not a int type")
	}
	rv.SetInt(int64(value))
	return nil
}

// SetFloat64 set float64
func SetFloat64(key string, value float64, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Float64 {
		return fmt.Errorf("Not a float64 type")
	}
	rv.SetFloat(float64(value))
	return nil
}

// SetFloat32 set float32
func SetFloat32(key string, value float32, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Float32 {
		return fmt.Errorf("Not a float32 type")
	}
	rv.SetFloat(float64(value))
	return nil
}

// SetInterface set interface
func SetInterface(key string, value interface{}, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}
	if rt.Kind() != reflect.Interface {
		return fmt.Errorf("Not a interface type")
	}
	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceString set slice string
func SetSliceString(key string, value []string, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.String {
		return fmt.Errorf("Not a slice string type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceInt64 set slice int64
func SetSliceInt64(key string, value []int64, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Int64 {
		return fmt.Errorf("Not a slice int64 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceInt32 set slice int32
func SetSliceInt32(key string, value []int32, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Int32 {
		return fmt.Errorf("Not a slice int32 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceInt16 set slice int16
func SetSliceInt16(key string, value []int16, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Int16 {
		return fmt.Errorf("Not a slice int16 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceInt8 set slice int8
func SetSliceInt8(key string, value []int8, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Int8 {
		return fmt.Errorf("Not a slice int8 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceInt set slice int
func SetSliceInt(key string, value []int, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Int {
		return fmt.Errorf("Not a slice int type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceFloat64 set slice float64
func SetSliceFloat64(key string, value []float64, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Float64 {
		return fmt.Errorf("Not a slice float64 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// SetSliceFloat32 set slice float32
func SetSliceFloat32(key string, value []float32, object interface{}) error {
	var loc = strings.Split(key, ".")
	rt, rv, err := FindObjectByTag(loc, object)
	if err != nil {
		return err
	}
	if !rv.CanSet() {
		return fmt.Errorf("Object not support modify")
	}

	if rt.Kind() != reflect.Slice && rt.Elem().Kind() != reflect.Float32 {
		return fmt.Errorf("Not a slice float32 type")
	}

	rv.Set(reflect.ValueOf(value))
	return nil
}

// // SetMap set map
// func SetMap(key string, k interface{}, v interface{}, object interface{}) error {
// 	var loc = strings.Split(key, ".")
// 	rt, rv, err := FindObjectByTag(loc, object)
// 	if err != nil {
// 		return err
// 	}
// 	if !rv.CanSet() {
// 		return fmt.Errorf("Object not support modify")
// 	}

// 	return nil
// }
