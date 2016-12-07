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
	var (
		words = arguments
		cword = len(words) - 1
		cur   string
		prev  string
	)

	if cword >= 0 {
		cur = words[cword]
	}
	if cword >= 1 {
		prev = words[cword-1]
	}

	cmdComplete := func() (compl []string) {
		// comple parameter
		if prev != "" {
			if flag := f.Flag(prev); flag != nil {
				if boolFlag, ok := flag.Value.(boolTypeFlag); !ok || !boolFlag.IsBool() {
					if flag.Completor != nil {
						compl = append(compl, f.Flag(prev).Completor(words)...)
					}
					return
				}
			}
		}

		if strings.HasPrefix(cur, "-") && len(f.Args()) < 2 {
			// complete flag
			f.Visit(func(flag *Flag) (err error) {
				if flag.Short != "" {
					compl = append(compl, fmt.Sprintf("-%s ", flag.Short))
				}

				if flag.Long != "" {
					compl = append(compl, fmt.Sprintf("--%s ", flag.Long))
				}

				return
			})
		} else {
			// complete argument
			for cmd := range f.cmdSet {
				compl = append(compl, cmd)
			}

			if f.Completor != nil {
				compl = append(compl, f.Completor(words)...)
			}
		}

		return
	}

	// self completion
	if f.SubCommandName() == "" {
		completions = cmdComplete()
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
