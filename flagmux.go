package xflag

import (
	"fmt"
	"os"
	"sort"
)

// FlagSet Extented Type
type FlagSetMux struct {
	*FlagSet

	CommandName string
	Commands    map[string]*FlagSet
}

func (m *FlagSetMux) Init(fs *FlagSet) {
	m.FlagSet = fs
	m.PrintHelp = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		m.PrintDefaults()

		var cmds []string
		for cmd := range m.Commands {
			cmds = append(cmds, cmd)
		}

		sort.Sort(sort.StringSlice(cmds))

		fmt.Fprintf(os.Stderr, "\n  Commands:\n")
		for _, cmd := range cmds {
			fmt.Fprintf(os.Stderr, "    - %s\n", cmd)
		}
	}
}

func (m *FlagSetMux) AddFlagSet(fs *FlagSet) {
	if m.Commands == nil {
		m.Commands = make(map[string]*FlagSet)
	}

	m.Commands[fs.Name] = fs
}

// interface... Parser
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
			m.CommandName = ""
			return fmt.Errorf("there is no matched flagset: %v", m.CommandName)
		}

		err = m.Commands[m.CommandName].Parse(subArgs)
		if err != nil {
			return err
		}
	}
	return nil
}

// interface... CompletionPrinter
func (m *FlagSetMux) Completions(arguments []string) (completions []string) {
	m.Parse(arguments)

	if m.CommandName == "" {
		completions = m.FlagSet.Completions(arguments)
		if len(completions) > 0 {
			for cmd := range m.Commands {
				completions = append(completions, cmd)
			}
		}
	} else {
		completions = m.Commands[m.CommandName].Completions(m.FlagSet.Args()[1:])
	}

	return
}
