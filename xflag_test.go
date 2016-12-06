package xflag

import (
	"testing"
	"time"
)

func TestXFlagParseIndirect(t *testing.T) {
	{
		type Opt struct {
			Number *int
		}

		var number int

		opt := &Opt{
			Number: &number,
		}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--number", "10",
		})

		if err != nil {
			t.Error(err)
		}

		if number != 10 {
			t.Error("unexpected value")
		}
	}
}

func TestXFlagParse(t *testing.T) {
	{
		type Opt struct {
			Duration time.Duration
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--duration", "10m",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.Duration != time.Minute*10 {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			Duration time.Duration
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--duration=10m",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.Duration != time.Minute*10 {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			OK bool
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--ok",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.OK != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			OK bool
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--ok=TRUE",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.OK != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			V bool `xflag:"v,verbose"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-v",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.V != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			Str string `xflag:"s"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-s", "Input",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.Str != "Input" {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			Str string `xflag:"s"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-sInput",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.Str != "Input" {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			V []bool `xflag:"v"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct("opt", opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-vvv",
		})

		if err != nil {
			t.Error(err)
		}

		if len(opt.V) != 3 {
			t.Error("value are incorrect")
		}
	}
}
