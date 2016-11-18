package xflag

import (
	"strconv"
	"time"
)

type BoolTypeFlag interface {
	IsBool() bool
}

// bool

type BoolValue bool

func (b *BoolValue) Set(s string) (err error) {
	v, err := strconv.ParseBool(s)
	*b = BoolValue(v)
	return err
}

func (b *BoolValue) Get() interface{} { return bool(*b) }

func (b *BoolValue) IsBool() bool { return true }

// []bool

type BoolSliceValue []bool

func (l *BoolSliceValue) Set(s string) (err error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*l = append(*l, v)
	return nil
}

func (l *BoolSliceValue) Get() interface{} { return []bool(*l) }

func (b *BoolSliceValue) IsBool() bool { return true }

// int

type IntValue int

func (i *IntValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = IntValue(v)
	return err
}

func (i *IntValue) Get() interface{} { return int(*i) }

// int64

type Int64Value int64

func (i *Int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = Int64Value(v)
	return err
}

func (i *Int64Value) Get() interface{} { return int64(*i) }

// uint
type UintValue uint

func (i *UintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = UintValue(v)
	return err
}
func (i *UintValue) Get() interface{} { return uint(*i) }

// uint64
type Uint64Value uint64

func (i *Uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = Uint64Value(v)
	return err
}

func (i *Uint64Value) Get() interface{} { return uint64(*i) }

// string

type StringValue string

func (s *StringValue) Set(val string) error {
	*s = StringValue(val)
	return nil
}

func (s *StringValue) Get() interface{} { return string(*s) }

// []string

type StringSliceValue []string

func (l *StringSliceValue) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *StringSliceValue) Get() interface{} { return []string(*l) }

// flat64

type Float64Value float64

func (f *Float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = Float64Value(v)
	return err
}

func (f *Float64Value) Get() interface{} { return float64(*f) }

// time.Duration

type DurationValue time.Duration

func (d *DurationValue) Set(s string) (err error) {
	v, err := time.ParseDuration(s)
	*d = DurationValue(v)
	return err
}

func (d *DurationValue) Get() interface{} { return time.Duration(*d) }
