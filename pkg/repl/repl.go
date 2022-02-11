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
	out := ocelot.Eval(in)
	if out != core.String("") {
		fmt.Println(out)
	}
}

func Repl(prompt string) {
	p := goprompt.New(
		executor,
		completer,
		goprompt.OptionPrefix(prompt),
	)
	p.Run()
}
