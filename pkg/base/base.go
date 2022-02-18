package base

import (
	"fmt"
	"reflect"

	"github.com/starlight/ocelot/pkg/core"
)

var Base = map[string]core.Any{
	"nil":   core.Nil{},
	"true":  core.Bool(true),
	"false": core.Bool(false),
	"type":  Function(Type),
	"def!":  Function(Def_E),
	"let*":  Function(Let_S),
}

type Function func(ast core.List, env Env) (core.Any, error)

func Let_S(ast core.List, env Env) (core.Any, error) {
	if len(ast) != 3 {
		return nil, fmt.Errorf("'%v' wanted 2 args, got %d", ast[0], len(ast)-1)
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
	if len(ast) != 3 {
		return nil, fmt.Errorf("'%v' wanted 2 args, got %d", ast[0], len(ast)-1)
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
	if len(ast) != 2 {
		return nil, fmt.Errorf("'%v' wanted 1 arg, got %d", ast[0], len(ast)-1)
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return core.String(fmt.Sprintf("%v", reflect.TypeOf(arg1))), nil
}
