package ocelot

import (
	"fmt"
	"os"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
	"golang.org/x/term"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func EvalStr(in string, env *core.Env) (core.Any, error) {
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
	return base.Eval(ast, *env)
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
		out, err := EvalStr(in, env)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(out)
	}
	saveTermState()
	defer restoreTermState()
	// prompt
	p := goprompt.New(
		executor,
		completer,
		goprompt.OptionPrefix(prompt),
	)
	p.Run()
	return nil
}

var termState *term.State

func saveTermState() {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	termState = oldState
}

func restoreTermState() {
	if termState != nil {
		term.Restore(int(os.Stdin.Fd()), termState)
	}
}
