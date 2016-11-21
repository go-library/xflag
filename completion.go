package xflag

import (
	"fmt"
	"strings"
)

func PrintCompletions(iface interface{}, args []string) (err error) {
	var fs *FlagSet
	switch i := iface.(type) {
	case *FlagSet:
		fs = i
	case *FlagSetMux:
		fs = i.FlagSet
	}

	err = fs.Parse(args)

	if err != nil {
		return err
	}

	var complSet []string
	fs.Visit(func(f *Flag) (err error) {
		var short = ""
		var long = ""

		if f.Short != "" {
			short = "-" + f.Short
		}
		if f.Long != "" {
			long = "--" + f.Long
		}

		compl := fmt.Sprintf("%s %s", short, long)
		compl = strings.TrimSpace(compl)
		complSet = append(complSet, compl)
		return nil
	})

	for i := range complSet {
		fmt.Printf("%s ", complSet[i])
	}
	fmt.Println()

	return nil
}
