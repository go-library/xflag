package xflag

import (
	"reflect"
	"strings"
)

func NewFlagSetFromStruct(name string, opt interface{}) (fs *FlagSet, err error) {

	optValue := reflect.ValueOf(opt)
	if optValue.Kind() != reflect.Ptr {
		return nil, Errorf(nil, nil, 0, "input agrument opt is not pointer of struct")
	}

	optValue = optValue.Elem()
	if optValue.Kind() != reflect.Struct {
		return nil, Errorf(nil, nil, 0, "input agrument opt is not struct")
	}

	t := optValue.Type()

	fs = NewFlagSet(name)

	for i := 0; i < t.NumField(); i++ {
		var (
			field      = t.Field(i)
			fieldValue = optValue.Field(i)
			short      = ""
			long       = ""
			help       = ""
			defValue   = ""
		)

		// struct tag parsing
		if v, ok := field.Tag.Lookup("xflag"); ok {
			terms := strings.SplitN(v, ",", 4)
			if len(terms) > 0 {
				short = strings.TrimSpace(terms[0])
			}

			if len(terms) > 1 {
				long = strings.TrimSpace(terms[1])
			}

			if len(terms) > 2 {
				defValue = strings.TrimSpace(terms[2])
			}

			if len(terms) > 3 {
				help = strings.TrimSpace(terms[3])
			}
		}

		if short == "" && long == "" {
			long = strings.ToLower(field.Name)
		}

		// bind
		if fieldValue.Kind() != reflect.Ptr {
			fieldValue = fieldValue.Addr()
		}

		err = fs.Bind(fieldValue.Interface(), short, long, defValue, help)
		if err != nil {
			return nil, err
		}
	}

	return fs, nil

}
