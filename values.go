package xflag

import (
	"strconv"
	"time"
)

type boolTypeFlag interface {
	IsBool() bool
}

// bool

type boolValue bool

func (b *boolValue) Set(s string) (err error) {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) IsBool() bool { return true }

// []bool

type boolSliceValue []bool

func (l *boolSliceValue) Set(s string) (err error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*l = append(*l, v)
	return nil
}

func (l *boolSliceValue) Get() interface{} { return []bool(*l) }

func (b *boolSliceValue) IsBool() bool { return true }

// int

type intValue int

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} { return int(*i) }

// int64

type int64Value int64

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() interface{} { return int64(*i) }

// uint
type uintValue uint

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}
func (i *uintValue) Get() interface{} { return uint(*i) }

// uint64
type uint64Value uint64

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} { return uint64(*i) }

// string

type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

// []string

type stringSliceValue []string

func (l *stringSliceValue) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *stringSliceValue) Get() interface{} { return []string(*l) }

// flat64

type float64Value float64

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() interface{} { return float64(*f) }

// time.Duration

type durationValue time.Duration

func (d *durationValue) Set(s string) (err error) {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} { return time.Duration(*d) }
