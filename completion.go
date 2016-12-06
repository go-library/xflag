package xflag

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

type Completer interface {
	Completions(arguments []string) (completions []string)
}

func HelpCompletion(c Completer) {

	envs := os.Environ()
	for i := range envs {
		if strings.HasPrefix(envs[i], "XFLAG_COMPLETION=1") {
			fmt.Println(strings.Join(c.Completions(os.Args[1:]), " "))
			os.Exit(0)
		}
		if strings.HasPrefix(envs[i], "XFLAG_COMPLETION_SCRIPT=1") {
			PrintBashScript()
			os.Exit(0)
		}
	}
}

func PrintBashScript() {
	templ := `
CMD="{{.ori}}"

function _{{.base}}(){
  local cur prev opts
  COMPREPLY=()
  CURL="${COMP_WORDS[COMP_CWORD]}"
  PREV="${COMP_WORDS[COMP_CWORD-1]}"
  OPTS=$(XFLAG_COMPLETION=1 "${COMP_WORDS[@]}")
  COMPREPLY=( $(compgen -W "${OPTS}" -- ${CURL}) )
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
