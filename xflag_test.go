package xflag

import (
	"testing"
	"time"
)

func TestXFlagParseIndirect(t *testing.T) {
	{
		type Opt struct {
			number *int
		}

		var number int

		opt := &Opt{
			number: &number,
		}
		fs, err := NewFlagSetFromStruct(opt)
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
			duration time.Duration
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--duration", "10m",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.duration != time.Minute*10 {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			duration time.Duration
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--duration=10m",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.duration != time.Minute*10 {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			ok bool
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--ok",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.ok != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			ok bool
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"--ok=TRUE",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.ok != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			v bool `xflag-short:"v"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-v",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.v != true {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			str string `xflag-short:"s"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-s", "Input",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.str != "Input" {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			str string `xflag-short:"s"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-sInput",
		})

		if err != nil {
			t.Error(err)
		}

		if opt.str != "Input" {
			t.Error("value are incorrect")
		}
	}

	{
		type Opt struct {
			v []bool `xflag-short:"v"`
		}

		opt := &Opt{}
		fs, err := NewFlagSetFromStruct(opt)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.Parse([]string{
			"-vvv",
		})

		if err != nil {
			t.Error(err)
		}

		if len(opt.v) != 3 {
			t.Error("value are incorrect")
		}
	}
}
