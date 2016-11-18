package xflag

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Value interface {
	Set(string) error
	Get() interface{}
}

type Flag struct {
	Short    string
	Long     string
	Usage    string
	Value    Value
	DefValue string
	IsSet    bool
}

type FlagSet struct {
	Name string

	shortFlags map[string]*Flag
	longFlags  map[string]*Flag
	args       []string

	subCommands map[string]*FlagSet
}

func NewFlagSet(name string) (fs *FlagSet) {
	f := &FlagSet{
		Name:       name,
		shortFlags: make(map[string]*Flag),
		longFlags:  make(map[string]*Flag),
	}

	return f
}

func (f *FlagSet) Var(value Value, short, long, defValue, usage string) (err error) {
	flag := &Flag{
		Short:    short,
		Long:     long,
		Usage:    usage,
		Value:    value,
		DefValue: defValue,
		IsSet:    false,
	}

	if short == "" && long == "" {
		return fmt.Errorf("%s flag name undefined", f.Name)
	}

	if short != "" {
		if _, has := f.shortFlags[short]; has {
			return fmt.Errorf("%s flag redefined: %s", f.Name, short)
		}
		f.shortFlags[short] = flag
	}

	if long != "" {
		if _, has := f.longFlags[long]; has {
			return fmt.Errorf("%s flag redefined: %s", f.Name, long)
		}
		f.longFlags[long] = flag
	}

	return nil
}

func (f *FlagSet) Visit(fn func(*Flag) error) error {
	var err error

	for _, f := range f.shortFlags {
		err = fn(f)
		if err != nil {
			return err
		}
	}

	for _, f := range f.longFlags {
		if f.Short == "" {
			err = fn(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// -s            bool
// -s value      w/o bool
// -sValue       w/o bool
//
// --long        bool
// --long value  w/o bool
// --long=value  any
func (f *FlagSet) Parse(args []string) (err error) {
	var (
		isFinished = false
		name       string
		value      string
		isBoolFlag bool
		window     = args
		shift      int
		flag       *Flag
		has        bool
	)

	for {
		if len(window) == 0 || isFinished {
			break
		}

		switch {

		case window[0] == "--":
			// -- terminator
			window = window[1:]
			isFinished = true

		case window[0][:2] == "--":
			// --* long flags
			terms := strings.SplitN(window[0][2:], "=", 2)

			// get flag name
			name = terms[0]
			if flag, has = f.longFlags[name]; !has {
				return fmt.Errorf("%s long flag undefined: %s", f.Name, name)
			}

			// check boolean field
			if boolFlag, ok := flag.Value.(BoolTypeFlag); ok && boolFlag.IsBool() {
				isBoolFlag = true
			} else {
				isBoolFlag = false
			}

			// get value & shift window
			if len(terms) == 2 {
				// --any=value
				value = terms[1]
				shift = 1
			} else if isBoolFlag {
				// --bool
				value = "true"
				shift = 1
			} else if len(window) > 1 {
				// --wobool value
				value = window[1]
				shift = 2
			} else {
				return fmt.Errorf("%s long flag value not provided: %s", f.Name, name)
			}

			// set Value
			flag.IsSet = true
			err = flag.Value.Set(value)
			if err != nil {
				return err
			}
			window = window[shift:]

		case window[0][:1] == "-":
			shift = 1
			// -* short flags
			opt := window[0][1:]

			for {
				if len(opt) == 0 {
					break
				}
				// get flag
				name = string(opt[0])
				if flag, has = f.shortFlags[name]; !has {
					return fmt.Errorf("%s short flag undefined: %s", f.Name, name)
				}

				// check boolean field
				if boolFlag, ok := flag.Value.(BoolTypeFlag); ok && boolFlag.IsBool() {
					isBoolFlag = true
				} else {
					isBoolFlag = false
				}

				// get value & shift opt
				if isBoolFlag {
					// -b
					value = "true"
					opt = opt[1:]
				} else if len(opt) > 1 {
					// -fValue
					value = opt[1:]
					opt = ""
				} else if len(window) > 1 {
					// -f value
					value = window[1]
					opt = opt[1:]
					shift = 2
				} else {
					return fmt.Errorf("%s short flag value not provided: %s", f.Name, name)
				}

				// set value
				flag.IsSet = true
				err = flag.Value.Set(value)
				if err != nil {
					return err
				}
			} // loop
			window = window[shift:]
		default:
			// *
			isFinished = true
		}

		f.args = window
	}

	// set default values
	err = f.Visit(func(f *Flag) error {
		if !f.IsSet {
			err = f.Value.Set(f.DefValue)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *FlagSet) Args() []string {
	return f.args
}

type ByName []*Flag

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	var (
		istr, jstr string
	)

	if a[i].Long == "" {
		istr = a[i].Short
	} else {
		istr = a[i].Long
	}

	if a[j].Long == "" {
		jstr = a[j].Short
	} else {
		jstr = a[j].Long
	}

	return 0 > strings.Compare(istr, jstr)
}

func (f *FlagSet) PrintUsage() {
	var flags []*Flag

	f.Visit(func(f *Flag) error {
		flags = append(flags, f)
		return nil
	})

	sort.Sort(ByName(flags))
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	var short, long string
	for _, f := range flags {
		if f.Short != "" {
			short = "-" + f.Short
		} else {
			short = ""
		}

		if f.Long != "" {
			long = "--" + f.Long
		} else {
			long = ""
		}

		lines := splitUsage(f.Usage)
		for i := range lines {
			fmt.Fprintf(os.Stderr, " %2s  %-15s  %s\n", short, long, lines[i])
			if i == 0 {
				long = ""
				short = ""
			}
		}
	}
}

func splitUsage(usage string) (lines []string) {
	var (
		line  string
		terms = strings.Split(usage, " ")
	)

	for i := range terms {
		line = fmt.Sprintf("%s %s", line, terms[i])
		if len(line) > 20 {
			lines = append(lines, line)
			line = ""
		}
	}
	lines = append(lines, line)

	return lines
}
