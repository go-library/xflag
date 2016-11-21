package xflag

import (
	"testing"
)

func TestFlagSetMux(t *testing.T) {
	type Opt struct {
		verbose bool
	}

	m := &FlagSetMux{}

	fs, err := NewFlagSetFromStruct(&Opt{})
	if err != nil {
		t.Fatal(err)
	}

	m.Init(fs)

	type SubCommand struct {
		enable bool
	}

	subfs, err := NewFlagSetFromStruct(&SubCommand{})
	if err != nil {
		t.Fatal(err)
	}

	m.AddFlagSet(subfs)

	err = m.Parse([]string{
		"--verbose",
		"SubCommand",
		"--enable",
	})

	if err != nil {
		t.Fatal(err)
	}

	if m.CommandName != "SubCommand" {
		t.Error("not expected command")
	}
}
