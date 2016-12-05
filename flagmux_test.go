package xflag

import (
	"testing"
)

func TestFlagSetMux(t *testing.T) {
	type Opt struct {
		verbose bool
	}

	m := &FlagSetMux{}

	fs, err := NewFlagSetFromStruct("option", &Opt{})
	if err != nil {
		t.Fatal(err)
	}

	m.Init(fs)

	type SubCommand struct {
		enable bool
	}

	subfs, err := NewFlagSetFromStruct("command1", &SubCommand{})
	if err != nil {
		t.Fatal(err)
	}

	m.AddFlagSet(subfs)

	err = m.Parse([]string{
		"--verbose",
		"command1",
		"--enable",
	})

	if err != nil {
		t.Fatal(err)
	}

	if m.CommandName != "command1" {
		t.Error("not expected command")
	}
}
