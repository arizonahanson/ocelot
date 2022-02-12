package ocelot

import (
	"fmt"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string, env *core.Env) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	var e core.Env
	if env == nil {
		e = core.BaseEnv()
	} else {
		e = *env
	}
	any, err := eval_ast(ast, e)
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
	if len(ast) > 0 {
		first := ast[0]
		switch first.(type) {
		default:
			break
		case core.Symbol:
			// check for special forms
			switch ast[0].(core.Symbol) {
			default:
				break
			case core.Symbol("def!"):
				if len(ast) != 3 {
					return nil, fmt.Errorf("'def!' received wrong number of args: %d", len(ast)-1)
				}
				switch ast[1].(type) {
				default:
					return nil, fmt.Errorf("first parameter to def! must be a symbol: %s", ast[1])
				case core.Symbol:
					val, err := eval_ast(ast[2], env)
					if err != nil {
						return nil, err
					}
					env.Set(ast[1].(core.Symbol), val)
					return val, nil
				}
			}
		}
	}
	// default list resolution
	for _, item := range ast {
		any, err := eval_ast(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	// check for function application
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
	// empty list
	return core.List(res), nil
}

func apply(fn core.Function, args []core.Any, env core.Env) (core.Any, error) {
	// TODO env?
	return fn(args...)
}
