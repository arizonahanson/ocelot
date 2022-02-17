package base

import (
	"github.com/starlight/ocelot/pkg/core"
)

// trampoline eval for non tail-calls
func Eval(ast core.Any, env Env) (core.Any, error) {
	return EvalType(evalAny).Trampoline(ast, env)
}

// thunked eval for tail-calls
func EvalTail(ast core.Any, env Env) core.Any {
	return EvalType(evalAny).Thunk(ast, env)
}

// eval impl
func evalAny(ast core.Any, env Env) (core.Any, error) {
	switch ast.(type) {
	default:
		// String, Number
		return ast, nil
	case core.Symbol:
		return env.Get(ast.(core.Symbol))
	case core.List:
		return evalList(ast.(core.List), env)
	case core.Vector:
		return evalVector(ast.(core.Vector), env)
	}
}

func evalList(ast core.List, env Env) (core.Any, error) {
	res := []core.Any{}
	for i, item := range ast {
		if i == 0 {
			// check for function symbols.
			switch item.(type) {
			default:
				// first isn't symbol
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
					continue
				case Function:
					// call function (unevaluated ast) & return
					fn := val.(Function)
					return fn(ast, env)
				}
			}
		}
		// default list resolution for rest
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	// empty list
	return core.List(res), nil
}

func evalVector(ast core.Vector, env Env) (core.Vector, error) {
	res := []core.Any{}
	for _, item := range ast {
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	return core.Vector(res), nil
}
