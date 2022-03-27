package base

import (
	"errors"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Parse(in string) (core.Any, error) {
	ast, err := parser.Parse("parse", []byte(in))
	if err != nil {
		return core.Null{}, err
	}
	return ast.(core.Any), nil
}

func EvalStr(in string, env *Env) (core.Any, error) {
	if env == nil {
		return core.Null{}, errors.New("evaluation with nil env")
	}
	ast, err := Parse(in)
	if err != nil {
		return core.Null{}, err
	}
	return Eval(ast, env)
}

// eager eval
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

// lazy eval and tail-call
func EvalFuture(ast core.Any, env *Env) Future {
	return func() (core.Any, error) {
		return evalAst(ast, env)
	}
}

// primary eval entrypoint
func evalAst(ast core.Any, env *Env) (core.Any, error) {
	switch any := ast.(type) {
	default:
		// String, Number, Bool, Null
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

// (eval list expressions)
func evalList(ast core.List, env *Env) (core.Any, error) {
	// () == nil
	if len(ast) == 0 {
		return core.Null{}, nil
	}
	// eval first item
	val, err := Eval(ast[0], env)
	if err != nil {
		return core.Null{}, err
	}
	// inspect type
	switch fn := val.(type) {
	default:
		break
	case Func:
		// function
		return fn.Future(ast, env), nil
	}
	// vector
	first := core.Vector{val}
	if len(ast) == 1 {
		return first, nil
	}
	rest, err := evalVector(core.Vector(ast[1:]), env)
	if err != nil {
		return core.Null{}, err
	}
	return append(first, rest.(core.Vector)...), nil
}

// [eval vectors]
func evalVector(ast core.Vector, env *Env) (core.Any, error) {
	res := make(core.Vector, len(ast))
	for i, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return core.Null{}, err
		}
		res[i] = val
	}
	return res, nil
}

// {:eval maps}
func evalMap(ast core.Map, env *Env) (core.Any, error) {
	res := make(core.Map, len(ast))
	for key, item := range ast {
		val, err := Eval(item, env)
		if err != nil {
			return core.Null{}, err
		}
		res[key] = val
	}
	return res, nil
}
