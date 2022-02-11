package repl

import (
	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func Read() string {
	return prompt.Input(">> ", completer)
}
