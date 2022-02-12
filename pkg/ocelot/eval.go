package ocelot

import (
	"fmt"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	env := core.GetEnv(make(map[string]interface{}))
	any, err := eval_ast(ast, env)
	return any, err
}

func eval_ast(ast core.Any, env core.Environment) (core.Any, error) {
	switch ast.(type) {
	default:
		return ast, nil
	case core.Symbol:
		return onSymbol(ast.(core.Symbol), env)
	case core.Vector:
		return onVector(ast.(core.Vector), env)
	case core.List:
		return onList(ast.(core.List), env)
	}
}

func onVector(ast core.Vector, env core.Environment) (core.Vector, error) {
	res := []core.Any{}
	for _, item := range ast {
		any, err := eval_ast(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	return core.Vector(res), nil
}

func onList(ast core.List, env core.Environment) (core.Any, error) {
	res := []core.Any{}
	for _, item := range ast {
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
			return apply(fn, res[1:], env)
		}
	}
	return core.List(res), nil
}

func apply(fn core.Function, args []core.Any, env core.Environment) (core.Any, error) {
	// TODO env?
	return fn(args...)
}

func onSymbol(ast core.Symbol, env core.Environment) (core.Any, error) {
	val := env[ast]
	if val == nil {
		return nil, fmt.Errorf("not found: %s", ast)
	}
	return val, nil
}
