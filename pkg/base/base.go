package base

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/starlight/ocelot/pkg/core"
)

var Base = map[string]core.Any{
	"nil":    core.Nil{},
	"true":   core.Bool(true),
	"false":  core.Bool(false),
	"nil?":   Function(nil_Q),
	"true?":  Function(true_Q),
	"false?": Function(false_Q),
	"bool":   Function(bool),
	"not":    Function(not),
	"type":   Function(type_),
	"def!":   Function(def_E),
	"let*":   Function(let_S),
	"and":    Function(and),
	"or":     Function(or),
}

type Function func(ast core.List, env Env) (core.Any, error)

func (fn Function) String() string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	str := strs[len(strs)-1]
	str = strings.ReplaceAll(str, "_", "")
	str = strings.ReplaceAll(str, "E", "!")
	str = strings.ReplaceAll(str, "Q", "?")
	str = strings.ReplaceAll(str, "S", "*")
	return "&" + str
}

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("'%v' wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func or(ast core.List, env Env) (core.Any, error) {
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

func and(ast core.List, env Env) (core.Any, error) {
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

func bool(ast core.List, env Env) (core.Any, error) {
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

func not(ast core.List, env Env) (core.Any, error) {
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

func false_Q(ast core.List, env Env) (core.Any, error) {
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

func true_Q(ast core.List, env Env) (core.Any, error) {
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

func nil_Q(ast core.List, env Env) (core.Any, error) {
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

func let_S(ast core.List, env Env) (core.Any, error) {
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

func def_E(ast core.List, env Env) (core.Any, error) {
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
func type_(ast core.List, env Env) (core.Any, error) {
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
