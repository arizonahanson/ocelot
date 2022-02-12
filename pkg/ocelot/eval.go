package ocelot

import (
	"fmt"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string, env map[string]interface{}) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	any, err := eval_ast(ast, env)
	return any, err
}

func eval_ast(ast core.Any, env map[string]interface{}) (core.Any, error) {
	switch ast.(type) {
	default:
		return ast, nil
	case core.Symbol:
		val := env[string(ast.(core.Symbol))]
		if val == nil {
			return nil, fmt.Errorf("not found: %s", ast)
		}
		return val, nil
	case core.Vector:
		res := []core.Any{}
		for _, item := range ast.(core.Vector) {
			any, err := eval_ast(item, env)
			if err != nil {
				return nil, err
			}
			res = append(res, any)
		}
		return core.Vector(res), nil
	case core.List:
		res := []core.Any{}
		for _, item := range ast.(core.List) {
			any, err := eval_ast(item, env)
			if err != nil {
				return nil, err
			}
			res = append(res, any)
		}
		if len(res) > 0 {
			first := res[0]
			switch first.(type) {
			default:
				return core.List(res), nil
			case core.Function:
				fn := first.(core.Function)
				return fn(res[1:]...)
			}
		}
		return core.List(res), nil
	}
}
