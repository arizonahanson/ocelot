package ocelot

import (
	"fmt"
	"os"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/pkg/base"
	"golang.org/x/term"
)

func Repl(prompt string) error {
	env, err := base.BaseEnv()
	if err != nil {
		return err
	}
	executor := func(in string) {
		if in == "" {
			return
		}
		out, err := base.EvalStr(in, env)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(out)
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
