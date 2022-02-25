package base

import (
	"errors"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Parse(in string) (core.Any, error) {
	ast, err := parser.Parse("parse", []byte(in))
	if err != nil {
		return nil, err
	}
	return ast.(core.Any), nil
}

func EvalStr(in string, env *Env) (core.Any, error) {
	if env == nil {
		return nil, errors.New("evaluation with nil env")
	}
	ast, err := Parse(in)
	if err != nil {
		return nil, err
	}
	return Eval(ast, env)
}

// eval that resolves future values
func Eval(ast core.Any, env *Env) (value core.Any, err error) {
	value, err = evalAst(ast, env)
	if err != nil {
		return
	}
	switch future := value.(type) {
	default:
		return
	case Future:
		return future.Resolve()
	}
}

// trampoline to resolve future values
func (future Future) Resolve() (value core.Any, err error) {
	value, err = future()
	for {
		if err != nil {
			return
		}
		switch next := value.(type) {
		default:
			return
		case Future:
			value, err = next()
		}
	}
}

// thunked eval for tail-calls (lazy)
func EvalFuture(ast core.Any, env *Env) Future {
	return func() (core.Any, error) {
		return evalAst(ast, env)
	}
}

// thunked function call (always lazy)
func FnFuture(fn Func, ast core.List, env *Env) Future {
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
	case core.Symbol:
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
			// check for function
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
			case Func:
				// tail-call function (unevaluated ast)
				return FnFuture(first, ast, env), nil
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
