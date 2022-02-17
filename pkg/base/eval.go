package base

import (
	"github.com/starlight/ocelot/pkg/core"
)

func EvalAst(ast core.Any, env Env) (core.Any, error) {
	switch ast.(type) {
	default:
		// String, Number
		return ast, nil
	case core.Symbol:
		return env.Get(ast.(core.Symbol))
	case core.List:
		return onList(ast.(core.List), env)
	case core.Vector:
		return onVector(ast.(core.Vector), env)
	}
}

type Function func(args []core.Any, env Env) (core.Any, error)

func onList(ast core.List, env Env) (core.Any, error) {
	res := []core.Any{}
	for i, item := range ast {
		if i == 0 {
			// check for function symbols.
			switch item.(type) {
			default:
				// first isn't symbol, eval
				any, err := EvalAst(item, env)
				if err != nil {
					return nil, err
				}
				res = append(res, any)
				break
			case core.Symbol:
				// first is a symbol, get
				val, err := env.Get(item.(core.Symbol))
				if err != nil {
					return nil, err
				}
				switch val.(type) {
				default:
					// not a function, append
					res = append(res, val)
					break
				case Function:
					// call function (unevaluated params) & return
					fn := val.(Function)
					return fn(ast[1:], env)
				}
			}
		} else {
			// default list resolution for rest
			any, err := EvalAst(item, env)
			if err != nil {
				return nil, err
			}
			res = append(res, any)
		}
	}
	// empty list
	return core.List(res), nil
}

func onVector(ast core.Vector, env Env) (core.Vector, error) {
	res := []core.Any{}
	for _, item := range ast {
		any, err := EvalAst(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	return core.Vector(res), nil
}
