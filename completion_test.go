package xflag

import (
	"testing"
)

func TestPrintCompletions(t *testing.T) {
	{
		type MainOpt struct {
			Verbose bool `xflag-short:"v"`
			Addr    string
			TTL     int
		}

		mainOpt := &MainOpt{}
		m := &FlagSetMux{}
		fs, err := NewFlagSetFromStruct(mainOpt)
		if err != nil {
			t.Fatal(err)
		}

		fs.Name = "main"
		m.Init(fs)

		err = PrintCompletions(fs, []string{
			"-v",
			"--addr",
		})
		if err != nil {
			t.Fatal(err)
		}

	}

}
