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

func EvalStr(in string, env *Env) (core.Any, error) {
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
	return Eval(ast, env)
}

// trampoline eval for making non-tail calls (eager)
func Eval(ast core.Any, env *Env) (value core.Any, err error) {
	value, err = evalAst(ast, env)
	for {
		if err != nil {
			return
		}
		switch thunk := value.(type) {
		default:
			return
		case Thunk:
			value, err = thunk()
		}
	}
}

// thunked eval for tail-calls (lazy)
func EvalTail(ast core.Any, env *Env) Thunk {
	return func() (core.Any, error) {
		return evalAst(ast, env)
	}
}

// thunked function call (always lazy)
func FnTail(fn Function, ast core.List, env *Env) Thunk {
	return func() (core.Any, error) {
		return fn(ast, env)
	}
}

// eval impl
func evalAst(ast core.Any, env *Env) (core.Any, error) {
	switch any := ast.(type) {
	default:
		// String, Number, Key
		return any, nil
	case *core.Symbol:
		return env.Get(any)
	case core.List:
		return evalList(any, env)
	case core.Vector:
		return evalVector(any, env)
	case core.Map:
		return evalMap(any, env)
	}
}

func evalList(ast core.List, env *Env) (core.Any, error) {
	var res core.List
	for i, item := range ast {
		if i == 0 {
			// check for function symbols.
			any, err := Eval(item, env)
			if err != nil {
				return nil, err
			}
			switch first := any.(type) {
			default:
				// not a function
				res = make(core.List, len(ast))
				res[0] = first
				continue
			case Function:
				// tail-call function (unevaluated ast)
				return FnTail(first, ast, env), nil
			}
		}
		// default list resolution for rest
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[i] = any
	}
	return res, nil
}

func evalVector(ast core.Vector, env *Env) (core.Vector, error) {
	res := make(core.Vector, len(ast))
	for i, item := range ast {
		any, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[i] = any
	}
	return res, nil
}

func evalMap(ast core.Map, env *Env) (core.Map, error) {
	res := make(core.Map, len(ast))
	for key, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		res[key] = val
	}
	return res, nil
}

func BaseEnv() (*Env, error) {
	env, err := NewEnv(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for sym, value := range Base {
		env.Set(&core.Symbol{Val: sym, Pos: nil}, value)
	}
	return env, nil
}
