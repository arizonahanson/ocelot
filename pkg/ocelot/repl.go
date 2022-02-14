package ocelot

import (
	"fmt"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/pkg/core"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
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
		out, err := Eval(in, env)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(out)
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
