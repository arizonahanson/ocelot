package repl

import (
	"fmt"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/pkg/core"
	"github.com/starlight/ocelot/pkg/ocelot"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func executor(in string) {
	if in == "" {
		return
	}
	env := make(map[string]interface{})
	env["true"] = core.Bool(true)
	env["false"] = core.Bool(false)
	env["nil"] = core.Nil{}
	out, err := ocelot.Eval(in, env)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast := out.(core.List)
	fmt.Print("↪ ")
	for _, node := range ast {
		fmt.Printf("%v ", node)
	}
	fmt.Println("")
}

func Repl(prompt string) {
	p := goprompt.New(
		executor,
		completer,
		goprompt.OptionPrefix(prompt),
	)
	p.Run()
}
