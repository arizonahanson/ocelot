package base

import (
	"errors"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Parse(in string) (core.Any, error) {
	ast, err := parser.Parse("parse", []byte(in))
	if err != nil {
		return core.Nil{}, err
	}
	return ast.(core.Any), nil
}

func EvalStr(in string, env *Env) (core.Any, error) {
	if env == nil {
		return core.Nil{}, errors.New("evaluation with nil env")
	}
	ast, err := Parse(in)
	if err != nil {
		return core.Nil{}, err
	}
	return Eval(ast, env)
}

// eval that resolves future values
func Eval(ast core.Any, env *Env) (val core.Any, err error) {
	val, err = evalAst(ast, env)
	if err != nil {
		return
	}
	switch future := val.(type) {
	default:
		return
	case Future:
		return future.Get()
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
		// String, Number, Key, Bool, Nil
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
	if len(ast) == 0 {
		return core.Nil{}, nil
	}
	// eval to vector or function-call
	val, err := Eval(ast[0], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch fn := val.(type) {
	default:
		break
	case Func:
		// tail-call function (unevaluated ast)
		return FnFuture(fn, ast, env), nil
	}
	// vector
	first := core.Vector{val}
	if len(ast) == 1 {
		return first, nil
	}
	rest, err := evalVector(core.Vector(ast[1:]), env)
	if err != nil {
		return core.Nil{}, err
	}
	return append(first, rest.(core.Vector)...), nil
}

func evalVector(ast core.Vector, env *Env) (core.Any, error) {
	res := make(core.Vector, len(ast))
	for i, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		res[i] = val
	}
	return res, nil
}

func evalMap(ast core.Map, env *Env) (core.Any, error) {
	res := make(core.Map, len(ast))
	for key, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		res[key] = val
	}
	return res, nil
}
