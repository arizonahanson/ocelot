package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

var Base = map[string]core.Any{
	"nil":    core.Nil{},
	"true":   core.Bool(true),
	"false":  core.Bool(false),
	"nil?":   Function(Nil_Q),
	"true?":  Function(True_Q),
	"false?": Function(False_Q),
	"bool":   Function(Bool),
	"not":    Function(Not),
	"type":   Function(Type),
	"def!":   Function(Def_E),
	"let*":   Function(Let_S),
	"and":    Function(And),
	"or":     Function(Or),
}

type Function func(ast core.List, env Env) (core.Any, error)

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("'%v' wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func Or(ast core.List, env Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Bool(false), nil
	}
	if len(ast) == 2 {
		return EvalTail(ast[1], env), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg != core.Bool(false) && arg != core.Nil{}) {
			return arg, nil
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func And(ast core.List, env Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Bool(true), nil
	}
	if len(ast) == 2 {
		return EvalTail(ast[1], env), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg == core.Bool(false) || arg == core.Nil{}) {
			return arg, nil
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func Bool(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.Bool(arg1 != core.Bool(false) && arg1 != core.Nil{}), nil
}

func Not(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.Bool(arg1 == core.Bool(false) || arg1 == core.Nil{}), nil
}

func False_Q(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.Bool(arg1 == core.Bool(false)), nil
}

func True_Q(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.Bool(arg1 == core.Bool(true)), nil
}

func Nil_Q(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.Bool(arg1 == core.Nil{}), nil
}

func Let_S(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 3)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("'%v' first arg should be a List, got '%v'", ast[0], ast[1])
	case core.List:
		newEnv, err := NewEnv(&env, nil, nil)
		if err != nil {
			return nil, err
		}
		err = newEnv.SetPairs(ast[1].(core.List))
		if err != nil {
			return nil, err
		}
		return EvalTail(ast[2], *newEnv), nil
	}
}

func Def_E(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 3)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("'%v' first arg should be a Symbol, got '%v'", ast[0], ast[1])
	case core.Symbol:
		val, err := Eval(ast[2], env)
		if err != nil {
			return nil, err
		}
		env.Set(ast[1].(core.Symbol), val)
		return val, nil
	}
}

// use reflection to get value type as String
func Type(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.String(fmt.Sprintf("%T", arg1)), nil
}
