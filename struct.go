package xflag

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func canFlagValue(v reflect.Value) (ok bool) {
	var fv *Value
	t := reflect.TypeOf(fv).Elem()
	return v.Type().Implements(t)
}

/*
  struct tags:
  @ xflag-name
  @ xflag-default
  @ xflag-usage
*/
func NewFlagSetFromStruct(opt interface{}) (fs *FlagSet, err error) {

	optValue := reflect.ValueOf(opt)
	if optValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("input agrument opt is not pointer of struct")
	}

	optValue = optValue.Elem()
	if optValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input agrument opt is not struct")
	}

	t := optValue.Type()

	fs = NewFlagSet(t.Name())

	for i := 0; i < t.NumField(); i++ {
		var (
			field      = t.Field(i)
			fieldValue = optValue.Field(i)
			short      = ""
			long       = ""
			usage      = ""
			defValue   = ""
		)

		if v, ok := field.Tag.Lookup("xflag-short"); ok {
			short = strings.TrimSpace(v)
		}

		if v, ok := field.Tag.Lookup("xflag-long"); ok {
			long = strings.TrimSpace(v)
		}

		if v, ok := field.Tag.Lookup("xflag-default"); ok {
			defValue = strings.TrimSpace(v)
		}

		if v, ok := field.Tag.Lookup("xflag-usage"); ok {
			usage = strings.TrimSpace(v)
		}

		if short == "" && long == "" {
			long = strings.ToLower(field.Name)
		}

		// for pointer

		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				return nil, fmt.Errorf("pointer field is nil")
			}
			fieldValue = fieldValue.Elem()
		}

		ptr := fieldValue.UnsafeAddr()

		typeName := fmt.Sprintf("%s/%s", field.Type.PkgPath(), field.Type.Name())

		kind := fieldValue.Kind()

		var v Value
		switch {
		// interface Value
		case canFlagValue(fieldValue):
			v = fieldValue.Interface().(Value)
		// time.Duration
		case typeName == "time/Duration":
			v = (*durationValue)(unsafe.Pointer(ptr))
		// bool
		case kind == reflect.Bool:
			v = (*boolValue)(unsafe.Pointer(ptr))
		// flat64
		case kind == reflect.Float64:
			v = (*float64Value)(unsafe.Pointer(ptr))
		// int
		case kind == reflect.Int:
			v = (*intValue)(unsafe.Pointer(ptr))
		// int64
		case kind == reflect.Int64:
			v = (*int64Value)(unsafe.Pointer(ptr))
		// uint
		case kind == reflect.Uint:
			v = (*uintValue)(unsafe.Pointer(ptr))
		// uint64
		case kind == reflect.Uint64:
			v = (*uint64Value)(unsafe.Pointer(ptr))
		// string
		case kind == reflect.String:
			v = (*stringValue)(unsafe.Pointer(ptr))
		// []bool
		case kind == reflect.Slice && field.Type.Elem().Kind() == reflect.Bool:
			v = (*boolSliceValue)(unsafe.Pointer(ptr))
		// []string
		case kind == reflect.Slice && field.Type.Elem().Kind() == reflect.String:
			v = (*stringSliceValue)(unsafe.Pointer(ptr))
		// error
		default:
			return nil, fmt.Errorf("unsupported type: %v", field.Type)
		}
		fs.Var(v, short, long, defValue, usage)
	}

	return fs, nil

}
