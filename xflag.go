package xflag

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

type Value interface {
	Set(string) error
	Get() interface{}
}

type Flag struct {
	Short     string
	Long      string
	MetaVar   string
	Help      string
	Value     Value
	DefValue  string
	IsSet     bool
	Completor func(args []string) (completes []string)
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

	f.PrintHelp = func() {
		fmt.Fprintf(os.Stderr, "Help of %s:\n", os.Args[0])
		f.PrintDefaults()
	}

	return f
}

func canFlagValue(v reflect.Value) (ok bool) {
	var fv *Value
	t := reflect.TypeOf(fv).Elem()
	return v.Type().Implements(t)
}

func (f *FlagSet) Bind(ifaceValue interface{}, short, long, defValue, help string) (err error) {
	// for pointer
	v := reflect.ValueOf(ifaceValue)

	if v.Kind() != reflect.Ptr {
		return Errorf(f, nil, 0, "input value is not pointer type")
	}

	if v.IsNil() {
		return Errorf(f, nil, 0, "pointer field is nil")
	}

	v = v.Elem()

	ptr := v.UnsafeAddr()

	typeName := fmt.Sprintf("%s/%s", v.Type().PkgPath(), v.Type().Name())
	kind := v.Kind()

	var value Value
	switch {
	// interface Value
	case canFlagValue(v):
		value = v.Interface().(Value)
	// time.Duration
	case typeName == "time/Duration":
		value = (*durationValue)(unsafe.Pointer(ptr))
	// bool
	case kind == reflect.Bool:
		value = (*boolValue)(unsafe.Pointer(ptr))
	// flat64
	case kind == reflect.Float64:
		value = (*float64Value)(unsafe.Pointer(ptr))
	// int
	case kind == reflect.Int:
		value = (*intValue)(unsafe.Pointer(ptr))
	// int64
	case kind == reflect.Int64:
		value = (*int64Value)(unsafe.Pointer(ptr))
	// uint
	case kind == reflect.Uint:
		value = (*uintValue)(unsafe.Pointer(ptr))
	// uint64
	case kind == reflect.Uint64:
		value = (*uint64Value)(unsafe.Pointer(ptr))
	// string
	case kind == reflect.String:
		value = (*stringValue)(unsafe.Pointer(ptr))
	// []bool
	case kind == reflect.Slice && v.Type().Elem().Kind() == reflect.Bool:
		value = (*boolSliceValue)(unsafe.Pointer(ptr))
	// []string
	case kind == reflect.Slice && v.Type().Elem().Kind() == reflect.String:
		value = (*stringSliceValue)(unsafe.Pointer(ptr))
	// error
	default:
		return Errorf(f, nil, 0, "unsupported type: %v", v.Type())
	}

	err = f.setValue(value, short, long, defValue, help)
	if err != nil {
		return err
	}

	return nil
}

// set Value as flag
func (f *FlagSet) setValue(value Value, short, long, defValue, help string) (err error) {
	var metaVar string

	if strings.HasPrefix(short, "-") {
		short = short[1:]
	}

	if len(short) > 1 {
		metaVar = strings.TrimSpace(short[1:])
		short = short[:1]
	}

	if strings.HasPrefix(long, "--") {
		long = long[2:]
	}

	if i := strings.IndexAny(long, "= "); i != -1 {
		metaVar = long[i+1:]
		long = long[:i]
	}

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
		return Errorf(f, flag, 0, "flag name undefined")
	}

	if short != "" {
		if _, has := f.shortFlags[short]; has {
			return Errorf(f, flag, 0, "short flag redefined")
		}
		f.shortFlags[short] = flag
	}

	if long != "" {
		if _, has := f.longFlags[long]; has {
			return Errorf(f, flag, 0, "long flag redefined")
		}
		f.longFlags[long] = flag
	}

	return nil
}

func (f *FlagSet) Flag(name string) *Flag {
	if strings.HasPrefix(name, "--") {
		name = name[2:]
	} else if strings.HasPrefix(name, "-") {
		name = name[1:]
	}

	if flag, ok := f.longFlags[name]; ok {
		return flag
	} else if flag, ok := f.shortFlags[name]; ok {
		return flag
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
			return Errorf(f, nil, ERR_HELP_REQUESTED, "")

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
				return Errorf(f, nil, ERR_UNDEFINED, "--%s flag is undefined", name)
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
				return Errorf(f, flag, ERR_EMPTY_VALUE, "--%s flag value was not provied", name)
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
					return Errorf(f, nil, ERR_UNDEFINED, "-%s flag is undefined", name)
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
					return Errorf(f, flag, ERR_EMPTY_VALUE, "-%s flag value was not provided", name)
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

	fmt.Fprintf(os.Stderr, format, "-h  --help", "print this message")

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
		if len(line) > 50 {
			lines = append(lines, line)
			line = ""
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func (f *FlagSet) Completions(arguments []string) (completions []string) {
	// use prev completor

	if 1 < len(arguments) {
		prev := arguments[len(arguments)-2]
		if f.Flag(prev) != nil {
			flag := f.Flag(prev)
			if boolFlag, ok := flag.Value.(boolTypeFlag); ok && boolFlag.IsBool() {
				// boolean flag
				// show all flags
			} else {
				// common flag
				if flag.Completor != nil {
					completions = append(completions, f.Flag(prev).Completor(arguments)...)
				}
				return
			}
		}
	}

	// default
	if completions == nil {
		f.Visit(func(flag *Flag) (err error) {
			if flag.Short != "" {
				completions = append(completions, fmt.Sprintf("-%s ", flag.Short))
			}

			if flag.Long != "" {
				completions = append(completions, fmt.Sprintf("--%s ", flag.Long))
			}

			return nil
		})
	}

	return
}
