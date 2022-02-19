package base

import (
	"github.com/starlight/ocelot/pkg/core"
)

type Thunk func() (core.Any, error)

// trampoline eval for making non-tail calls
func Eval(ast core.Any, env core.Env) (core.Any, error) {
	value, err := evalAny(ast, env)
	if err != nil {
		return nil, err
	}
	for {
		switch value.(type) {
		default:
			return value, nil
		case Thunk:
			break
		}
		thunk := value.(Thunk)
		next, err := thunk()
		if err != nil {
			return nil, err
		}
		value = next
	}
}

// thunked eval for tail-calls
func EvalTail(ast core.Any, env core.Env) Thunk {
	return func() (core.Any, error) {
		return evalAny(ast, env)
	}
}

// thunked function call
func FnTail(fn core.Function, ast core.List, env core.Env) Thunk {
	return func() (core.Any, error) {
		return fn(ast, env)
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
	case core.Map:
		return evalMap(ast.(core.Map), env)
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
					// tail-call function (unevaluated ast)
					return FnTail(val.(core.Function), ast, env), nil
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

func evalMap(ast core.Map, env core.Env) (core.Map, error) {
	res := make(map[core.Key]core.Any)
	for key, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[key] = val
	}
	return core.Map(res), nil
}

func BaseEnv() (*core.Env, error) {
	env, err := core.NewEnv(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for sym, value := range Base {
		env.Set(core.Symbol{Val: sym, Pos: nil}, value)
	}
	return env, nil
}
