package ocelot

import (
	"fmt"
	"os"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/builtin"
	"github.com/starlight/ocelot/pkg/core"
	"golang.org/x/term"
)

func Repl(prompt string) error {
	env, err := builtin.BuiltinEnv()
	if err != nil {
		return err
	}
	executor := func(in string) {
		if in == "" {
			return
		}
		val, err := base.EvalStr(in, env)
		if err != nil {
			fmt.Println(err)
			return
		}
		Print(val)
	}
	completer := func(d goprompt.Document) []goprompt.Suggest {
		return []goprompt.Suggest{}
	}
	// fix ctrl-c stops working after exit
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

func Print(ast core.Any) {
	fmt.Print(color.WhiteString("→ "))
	switch any := ast.(type) {
	default:
		fmt.Printf("%#v\n", any)
	case core.Vector:
		for _, item := range any {
			fmt.Printf("%#v ", item)
		}
		fmt.Println("")
	}
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
