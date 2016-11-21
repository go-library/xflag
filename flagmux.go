package xflag

import (
	"fmt"
)

type FlagSetMux struct {
	*FlagSet

	CommandName string
	Commands    map[string]*FlagSet
}

func (m *FlagSetMux) Init(fs *FlagSet) {
	m.FlagSet = fs
}

func (m *FlagSetMux) AddFlagSet(fs *FlagSet) {
	if m.Commands == nil {
		m.Commands = make(map[string]*FlagSet)
	}

	m.Commands[fs.Name] = fs
}

func (m *FlagSetMux) Parse(arguments []string) (err error) {
	err = m.FlagSet.Parse(arguments)
	if err != nil {
		return err
	}

	subArgs := m.FlagSet.Args()
	if len(subArgs) >= 1 {
		m.CommandName = subArgs[0]
		subArgs = subArgs[1:]
		if _, ok := m.Commands[m.CommandName]; !ok {
			return fmt.Errorf("there is no matched flagset: %v", m.CommandName)
		}

		err = m.Commands[m.CommandName].Parse(subArgs)
		if err != nil {
			return err
		}
	}

	return nil
}
