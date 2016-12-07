package xflag

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

func doCompletion(f *FlagSet) {

	envs := os.Environ()
	for i := range envs {
		if strings.HasPrefix(envs[i], "XFLAG_COMPLETION=1") {
			fmt.Println(strings.Join(genComplWords(f, os.Args[1:]), " "))
			os.Exit(0)
		}
		if strings.HasPrefix(envs[i], "XFLAG_COMPLETION_SCRIPT=1") {
			printBashScript()
			os.Exit(0)
		}
	}
}

func printBashScript() {
	templ := `
CMD="{{.ori}}"

function _{{.base}}(){
  local cur prev cword words
	_get_comp_words_by_ref -n =: cur words
  OPTS=($(XFLAG_COMPLETION=1 "${words[@]}"))
  COMPREPLY=($(compgen -W "${OPTS[*]}" -- ${cur}))
	__ltrim_colon_completions "$cur"
  return 0
}

complete -F _{{.base}} $CMD
`
	data := map[string]interface{}{
		"ori":  os.Args[0],
		"base": path.Base(os.Args[0]),
	}

	t := template.Must(template.New("script").Parse(templ))
	t.Execute(os.Stdout, data)
}

func genComplWords(f *FlagSet, arguments []string) (completions []string) {
	cmdComplete := func(args []string) (compl []string) {
		// use prev-flag completor
		if 1 < len(args) {
			prev := args[len(args)-2]
			if flag := f.Flag(prev); flag != nil {
				// not bool type flag
				if boolFlag, ok := flag.Value.(boolTypeFlag); !ok || !boolFlag.IsBool() {
					if flag.Completor != nil {
						compl = append(compl, f.Flag(prev).Completor(f.Args())...)
					}
					return
				}
			}
		}

		// append flags
		f.Visit(func(flag *Flag) (err error) {
			if flag.Short != "" {
				compl = append(compl, fmt.Sprintf("-%s ", flag.Short))
			}

			if flag.Long != "" {
				compl = append(compl, fmt.Sprintf("--%s ", flag.Long))
			}

			return
		})

		// flagset completor
		if f.Completor != nil {
			compl = append(compl, f.Completor(f.Args())...)
		}

		return
	}

	//f.Parse(arguments)

	// self completion
	if f.cmdName == "" {
		completions = cmdComplete(arguments)
		if len(completions) > 0 {
			for cmd := range f.cmdSet {
				completions = append(completions, cmd)
			}
		}
		// sub command completion
	} else {
		if arguments[len(arguments)-1] == f.cmdName {
			completions = []string{f.cmdName}
		} else {
			completions = genComplWords(f.cmdSet[f.cmdName], f.Args()[1:])
		}
	}

	return

}
