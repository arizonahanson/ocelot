package repl

import (
	goprompt "github.com/c-bata/go-prompt"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func Read(prompt string) string {
	return goprompt.Input(prompt, completer)
}

func Repl() {
	for {
		input := Read(">> ")
		if input == "(exit)" {
			break
		}
		// TODO: Eval
		println(input)
	}
}
