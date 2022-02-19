package base

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/shopspring/decimal"
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
	"add":    Function(add),
	"sub":    Function(sub),
	"mul":    Function(mul),
	"quot":   Function(quot),
	"rem":    Function(rem),
	"quot*":  Function(quot_S),
}

type Function func(ast core.List, env Env) (core.Any, error)

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("'%v' wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func add(ast core.List, env Env) (core.Any, error) {
	res := decimal.Zero
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			res = res.Add(arg.(core.Number).Decimal())
		}
	}
	return core.Number(res), nil
}

func sub(ast core.List, env Env) (core.Any, error) {
	res := decimal.Zero
	for i, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			num := arg.(core.Number).Decimal()
			if i == 0 && len(ast) > 2 {
				res = num
			} else {
				res = res.Sub(num)
			}
		}
	}
	return core.Number(res), nil
}

func mul(ast core.List, env Env) (core.Any, error) {
	res := decimal.NewFromInt32(1)
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			res = res.Mul(arg.(core.Number).Decimal())
		}
	}
	return core.Number(res), nil
}

func quot(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 4)
	if err != nil {
		return nil, err
	}
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			continue
		}
	}
	arg1 := ast[1].(core.Number).Decimal()
	arg2 := ast[2].(core.Number).Decimal()
	prec := ast[3].(core.Number).Decimal().IntPart()
	q, _ := arg1.QuoRem(arg2, int32(prec))
	return core.Number(q), nil
}

func rem(ast core.List, env Env) (core.Any, error) {
	err := exactLen(ast, 4)
	if err != nil {
		return nil, err
	}
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			continue
		}
	}
	arg1 := ast[1].(core.Number).Decimal()
	arg2 := ast[2].(core.Number).Decimal()
	prec := ast[3].(core.Number).Decimal().IntPart()
	_, r := arg1.QuoRem(arg2, int32(prec))
	return core.Number(r), nil
}

func quot_S(ast core.List, env Env) (core.Any, error) {
	res := decimal.NewFromInt32(1)
	for i, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-number '%v'", ast[0], arg)
		case core.Number:
			num := arg.(core.Number).Decimal()
			if i == 0 && len(ast) > 2 {
				res = num
			} else {
				res = res.Div(num)
			}
		}
	}
	return core.Number(res), nil
}

func or(ast core.List, env Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Bool(false), nil
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

func (fn Function) String() string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	str := strs[len(strs)-1]
	str = strings.ReplaceAll(str, "_", "")
	str = strings.ReplaceAll(str, "E", "!")
	str = strings.ReplaceAll(str, "Q", "?")
	str = strings.ReplaceAll(str, "S", "*")
	return "&" + str
}
