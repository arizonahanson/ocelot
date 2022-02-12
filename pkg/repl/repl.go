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
	out, err := ocelot.Eval(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast := out.(core.List)
	for _, node := range ast {
		fmt.Printf("%s ", node)
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
