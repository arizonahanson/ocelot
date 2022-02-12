package ocelot

import (
	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	env := core.BaseEnv()
	any, err := eval_ast(ast, env)
	return any, err
}

func eval_ast(ast core.Any, env core.Env) (core.Any, error) {
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

func onSymbol(ast core.Symbol, env core.Env) (core.Any, error) {
	return env.Get(string(ast))
}

func onVector(ast core.Vector, env core.Env) (core.Vector, error) {
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

func onList(ast core.List, env core.Env) (core.Any, error) {
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

func apply(fn core.Function, args []core.Any, env core.Env) (core.Any, error) {
	// TODO env?
	return fn(args...)
}
