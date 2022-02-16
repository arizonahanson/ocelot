package ocelot

import (
	"fmt"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func Eval(in string, env *base.Env) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	if env == nil {
		base, err := base.BaseEnv()
		if err != nil {
			return nil, err
		}
		env = base
	}
	return base.EvalAst(ast, *env)
}

func Repl(prompt string) error {
	env, err := base.BaseEnv()
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
