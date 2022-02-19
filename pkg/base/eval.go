package base

import (
	"github.com/starlight/ocelot/pkg/core"
)

// trampoline eval for non tail-calls
func Eval(ast core.Any, env core.Env) (core.Any, error) {
	value, err := evalAny(ast, env)
	if err != nil {
		return nil, err
	}
	for {
		switch value.(type) {
		default:
			return value, nil
		case ThunkType:
			break
		}
		thunk := value.(ThunkType)
		next, err := thunk()
		if err != nil {
			return nil, err
		}
		value = next
	}
}

type ThunkType func() (core.Any, error)

// thunked eval for tail-calls
func EvalTail(ast core.Any, env core.Env) ThunkType {
	return func() (core.Any, error) {
		return evalAny(ast, env)
	}
}

// eval impl
func evalAny(ast core.Any, env core.Env) (core.Any, error) {
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

func evalList(ast core.List, env core.Env) (core.Any, error) {
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
				case core.Function:
					// call function (unevaluated ast) & return
					fn := val.(core.Function)
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

func evalVector(ast core.Vector, env core.Env) (core.Vector, error) {
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

func BaseEnv() (*core.Env, error) {
	env, err := core.NewEnv(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range Base {
		env.Set(core.Symbol(key), value)
	}
	return env, nil
}
