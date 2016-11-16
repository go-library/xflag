package xflag

import (
	"flag"
	"fmt"
	//"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func canFlagValue(v reflect.Value) (ok bool) {
	var fv *flag.Value
	t := reflect.TypeOf(fv).Elem()
	return v.Type().Implements(t)
}

/*
  struct tags:
  @ xflag-name
  @ xflag-default
  @ xflag-usage
*/
func NewFlagSetFromStruct(opt interface{}) (fs *flag.FlagSet, err error) {

	optValue := reflect.ValueOf(opt)
	if optValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("input agrument opt is not pointer of struct")
	}

	optValue = optValue.Elem()
	if optValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input agrument opt is not struct")
	}

	t := optValue.Type()

	fs = flag.NewFlagSet(t.Name(), flag.ContinueOnError)

	for i := 0; i < t.NumField(); i++ {
		var (
			field      = t.Field(i)
			fieldValue = optValue.Field(i)
			name       = strings.ToLower(field.Name)
			usage      = ""
			defValue   = ""
		)

		if v, ok := field.Tag.Lookup("xflag-name"); ok {
			name = v
		}

		if v, ok := field.Tag.Lookup("xflag-default"); ok {
			defValue = v
		}

		if v, ok := field.Tag.Lookup("xflag-usage"); ok {
			usage = v
		}

		ptr := fieldValue.UnsafeAddr()

		typeName := fmt.Sprintf("%s/%s", field.Type.PkgPath(), field.Type.Name())

		kind := fieldValue.Kind()

		switch {
		case canFlagValue(fieldValue):
			v := fieldValue.Interface().(flag.Value)
			fs.Var(v, name, usage)
		case typeName == "time/Duration":
			var v time.Duration
			if defValue != "" {
				v, err = time.ParseDuration(defValue)
				if err != nil {
					return nil, err
				}
			}
			fs.DurationVar((*time.Duration)(unsafe.Pointer(ptr)), name, v, usage)
		case kind == reflect.Bool:
			var v bool
			if defValue != "" {
				defValue = strings.ToLower(defValue)
				switch defValue {
				case "true", "t":
					v = true
				case "false", "f":
					v = false
				}
			}
			fs.BoolVar((*bool)(unsafe.Pointer(ptr)), name, v, usage)
		case kind == reflect.Float64:
			var v float64
			if defValue != "" {
				v, err = strconv.ParseFloat(defValue, 64)
				if err != nil {
					return nil, err
				}
			}
			fs.Float64Var((*float64)(unsafe.Pointer(ptr)), name, v, usage)
		case kind == reflect.Int:
			var v int64
			if defValue != "" {
				v, err = strconv.ParseInt(defValue, 10, 32)
				if err != nil {
					return nil, err
				}
			}
			fs.IntVar((*int)(unsafe.Pointer(ptr)), name, int(v), usage)
		case kind == reflect.Int64:
			var v int64
			if defValue != "" {
				v, err = strconv.ParseInt(defValue, 10, 64)
				if err != nil {
					return nil, err
				}
			}
			fs.Int64Var((*int64)(unsafe.Pointer(ptr)), name, v, usage)
		case kind == reflect.String:
			fs.StringVar((*string)(unsafe.Pointer(ptr)), name, defValue, usage)
		case kind == reflect.Slice && field.Type.Elem().Kind() == reflect.String: // []string
			fs.Var((*StringList)(unsafe.Pointer(ptr)), name, usage)
		default:
			return nil, fmt.Errorf("unsupported type: %v", field.Type)
		}
	}

	return fs, nil

}
