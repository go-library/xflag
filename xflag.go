package xflag

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	ErrHelp = fmt.Errorf("help requested")
)

type Value interface {
	Set(string) error
	Get() interface{}
}

type Flag struct {
	Short    string
	Long     string
	MetaVar  string
	Help     string
	Value    Value
	DefValue string
	IsSet    bool
}

type FlagSet struct {
	Name      string
	PrintHelp func()

	shortFlags map[string]*Flag
	longFlags  map[string]*Flag
	args       []string
}

func NewFlagSet(name string) (fs *FlagSet) {
	f := &FlagSet{
		Name:       name,
		shortFlags: make(map[string]*Flag),
		longFlags:  make(map[string]*Flag),
	}

	return f
}

// set Value as flag
func (f *FlagSet) Var(value Value, short, long, defValue, metaVar, help string) (err error) {
	flag := &Flag{
		Short:    short,
		Long:     long,
		MetaVar:  metaVar,
		Help:     help,
		Value:    value,
		DefValue: defValue,
		IsSet:    false,
	}

	if short == "" && long == "" {
		return fmt.Errorf("%s FlagSet name undefined", f.Name)
	}

	if short != "" {
		if _, has := f.shortFlags[short]; has {
			return fmt.Errorf("%s FlagSet redefined: %s", f.Name, short)
		}
		f.shortFlags[short] = flag
	}

	if long != "" {
		if _, has := f.longFlags[long]; has {
			return fmt.Errorf("%s FlagSet redefined: %s", f.Name, long)
		}
		f.longFlags[long] = flag
	}

	return nil
}

// visit all flag
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

// parse arguments
// -f           // only boolean
// -fvalue      // without boolean
// -f value     // without boolean
// --flag       // only boolean
// --flag=value // any type
// --flag value // without boolean
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
		case window[0] == "-h" || window[0] == "--help":
			if f.PrintHelp == nil {
				fmt.Fprintf(os.Stderr, "Help of %s:\n", os.Args[0])
				f.PrintDefaults()
			} else {
				f.PrintHelp()
			}
			return ErrHelp

		case window[0] == "--":
			// -- terminator
			window = window[1:]
			isFinished = true

		case strings.HasPrefix(window[0], "--"):
			// --* long flags
			terms := strings.SplitN(window[0][2:], "=", 2)

			// get flag name
			name = terms[0]
			if flag, has = f.longFlags[name]; !has {
				return fmt.Errorf("%s FlagSet undefined: --%s", f.Name, name)
			}

			// check boolean field
			if boolFlag, ok := flag.Value.(boolTypeFlag); ok && boolFlag.IsBool() {
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
				return fmt.Errorf("%s FlagSet value not provided: --%s", f.Name, name)
			}

			// set Value
			flag.IsSet = true
			err = flag.Value.Set(value)
			if err != nil {
				return err
			}
			window = window[shift:]

		case strings.HasPrefix(window[0], "-"):
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
					return fmt.Errorf("%s FlagSet undefined: -%s", f.Name, name)
				}

				// check boolean field
				if boolFlag, ok := flag.Value.(boolTypeFlag); ok && boolFlag.IsBool() {
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
					return fmt.Errorf("%s FlagSet value not provided: -%s", f.Name, name)
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
		if !f.IsSet && f.DefValue != "" {
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

// return remained arguments
func (f *FlagSet) Args() []string {
	return f.args
}

type byName []*Flag

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool {
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

// print default values
func (f *FlagSet) PrintDefaults() {
	const format = "  %-30s  %s\n"
	var flags []*Flag

	f.Visit(func(f *Flag) error {
		flags = append(flags, f)
		return nil
	})

	sort.Sort(byName(flags))

	fmt.Fprintf(os.Stderr, format, "-h --help", "print this message")

	var short, long, metaVar string
	for _, f := range flags {
		if boolFlag, ok := f.Value.(boolTypeFlag); ok && boolFlag.IsBool() {
			metaVar = ""
		} else {
			metaVar = f.MetaVar
		}

		// long flag name formating
		if f.Short != "" && metaVar != "" && f.Long == "" {
			short = fmt.Sprintf("-%s %s", f.Short, metaVar)
		} else if f.Short != "" {
			short = fmt.Sprintf("-%s", f.Short)
		} else {
			short = ""
		}

		// short flag name formating
		if f.Long != "" && metaVar != "" {
			long = fmt.Sprintf("--%s=%s", f.Long, metaVar)
		} else if f.Long != "" {
			long = fmt.Sprintf("--%s", f.Long)
		} else {
			long = ""
		}

		lines := splitHelp(f.Help)
		if f.DefValue != "" {
			lines = append(lines, fmt.Sprintf("(default: %s)", f.DefValue))
		}

		if len(lines) == 0 {
			lines = append(lines, "")
		}

		for i := range lines {
			var l []string

			if short != "" {
				l = append(l, short)
			}
			if long != "" {
				l = append(l, long)
			}

			fmt.Fprintf(os.Stderr, format, strings.Join(l, "  "), lines[i])

			if i == 0 {
				long = ""
				short = ""
			}
		}
	}
}

func splitHelp(help string) (lines []string) {
	var (
		line  string
		terms = strings.Split(help, " ")
	)

	for i := range terms {
		line = fmt.Sprintf("%s %s", line, terms[i])
		line = strings.TrimSpace(line)
		if len(line) > 30 {
			lines = append(lines, line)
			line = ""
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}
