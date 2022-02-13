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

func Print(ast core.Any) {
	switch ast.(type) {
	default:
		fmt.Println(ast)
	case core.List:
		for _, node := range ast.(core.List) {
			fmt.Printf("%v ", node)
		}
		fmt.Println("")
	}
}

func Repl(prompt string) error {
	env, err := core.BaseEnv()
	if err != nil {
		return err
	}
	executor := func(in string) {
		if in == "" {
			return
		}
		out, err := ocelot.Eval(in, env)
		if err != nil {
			fmt.Println(err)
			return
		}
		// print
		fmt.Print("↪ ")
		Print(out)
	}
	// prompt
	p := goprompt.New(
		executor,
		completer,
		goprompt.OptionPrefix(prompt),
	)
	p.Run()
	return nil
}
