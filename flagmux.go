package xflag

import (
	"flag"
	"fmt"
	"os"
)

type FlagSetMux struct {
	Command string
	*flag.FlagSet
	flags map[string]*flag.FlagSet
}

func NewFlagSetMux() (m *FlagSetMux) {
	m = new(FlagSetMux)
	m.flags = make(map[string]*flag.FlagSet)
	m.FlagSet = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	return m
}

func (m *FlagSetMux) AddFlagSet(name string, fs *flag.FlagSet) {
	m.flags[name] = fs
}

func (m *FlagSetMux) Parse(arguments []string) (err error) {
	err = m.FlagSet.Parse(arguments)
	if err != nil {
		return err
	}

	subArgs := m.FlagSet.Args()
	if len(subArgs) >= 1 {
		m.Command = subArgs[0]
		subArgs = subArgs[1:]
		if _, ok := m.flags[m.Command]; !ok {
			return fmt.Errorf("there is no matched flagset: %v", m.Command)
		}

		err = m.flags[m.Command].Parse(subArgs)
		if err != nil {
			return err
		}
	}

	return nil
}
