package repl

import (
	"fmt"
	"reflect"

	goprompt "github.com/c-bata/go-prompt"
	"github.com/starlight/ocelot/pkg/core"
	"github.com/starlight/ocelot/pkg/ocelot"
)

func completer(d goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func executor(in string) {
	if in == "" {
		return
	}
	env := make(map[string]interface{})
	env["true"] = core.Bool(true)
	env["false"] = core.Bool(false)
	env["nil"] = core.Nil{}
	env["type"] = core.Function(func(args ...core.Any) (core.Any, error) {
		if len(args) == 0 {
			return core.Nil{}, nil
		}
		if len(args) > 1 {
			return nil, fmt.Errorf("too many args for type: %d", len(args))
		}
		typeStr := reflect.TypeOf(args[0]).String()
		return core.String(typeStr), nil
	})
	out, err := ocelot.Eval(in, env)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast := out.(core.List)
	fmt.Print("↪ ")
	for _, node := range ast {
		fmt.Printf("%v ", node)
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
