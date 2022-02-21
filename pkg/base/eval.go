package base

import (
	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

type Thunk func() (core.Any, error)

func Parse(in string) (core.Any, error) {
	ast, err := parser.Parse("parse", []byte(in))
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func EvalStr(in string, env *core.Env) (core.Any, error) {
	ast, err := Parse(in)
	if err != nil {
		return nil, err
	}
	if env == nil {
		base, err := BaseEnv()
		if err != nil {
			return nil, err
		}
		env = base
	}
	return Eval(ast, *env)
}

// trampoline eval for making non-tail calls (eager)
func Eval(ast core.Any, env core.Env) (core.Any, error) {
	value, err := evalAst(ast, env)
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

// thunked eval for tail-calls (lazy)
func EvalTail(ast core.Any, env core.Env) Thunk {
	return func() (core.Any, error) {
		return evalAst(ast, env)
	}
}

// thunked function call (always lazy)
func FnTail(fn core.Function, ast core.List, env core.Env) Thunk {
	return func() (core.Any, error) {
		return fn(ast, env)
	}
}

// eval impl
func evalAst(ast core.Any, env core.Env) (core.Any, error) {
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
	res := make([]core.Any, len(ast))
	for i, item := range ast {
		if i == 0 {
			// check for function symbols.
			first, err := Eval(item, env)
			if err != nil {
				return nil, err
			}
			switch first.(type) {
			default:
				// not a function
				res[0] = first
				continue
			case core.Function:
				// tail-call function (unevaluated ast)
				return FnTail(first.(core.Function), ast, env), nil
			}
		}
		// default list resolution for rest
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[i] = any
	}
	// empty list
	return core.List(res), nil
}

func evalVector(ast core.Vector, env core.Env) (core.Vector, error) {
	res := make([]core.Any, len(ast))
	for i, item := range ast {
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[i] = any
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
