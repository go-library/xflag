package xflag

import (
	"flag"
	"testing"
)

func TestRepeat(t *testing.T) {
	type Opt struct {
		InputList []string
	}

	opt := &Opt{}

	fs, err := NewFlagSetFromStruct(opt)
	if err != nil {
		t.Fatal(err)
	}

	err = fs.Parse([]string{
		"-inputlist", "1",
		"-inputlist", "2",
		"-inputlist", "3",
	})

	if err != nil {
		t.Error(err)
	}

	fs.VisitAll(func(f *flag.Flag) {
		t.Logf("- %+v\n", f)
	})

}
