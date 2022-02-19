package base

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/starlight/ocelot/pkg/core"
)

var Base = map[string]core.Any{
	// nil / bool
	"nil":    core.Nil{},
	"true":   core.Bool(true),
	"false":  core.Bool(false),
	"nil?":   core.Function(nil_Q),
	"true?":  core.Function(true_Q),
	"false?": core.Function(false_Q),
	"bool":   core.Function(bool),
	"not":    core.Function(not),
	"and":    core.Function(and),
	"or":     core.Function(or),
	// numbers
	"add":   core.Function(add),
	"sub":   core.Function(sub),
	"mul":   core.Function(mul),
	"quot":  core.Function(quot),
	"rem":   core.Function(rem),
	"quot*": core.Function(quot_S),
	// special
	"type": core.Function(type_),
	"def!": core.Function(def_E),
	"let*": core.Function(let_S),
}

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("'%v' wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func nil_Q(ast core.List, env core.Env) (core.Any, error) {
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

func true_Q(ast core.List, env core.Env) (core.Any, error) {
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

func false_Q(ast core.List, env core.Env) (core.Any, error) {
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

func bool(ast core.List, env core.Env) (core.Any, error) {
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

func not(ast core.List, env core.Env) (core.Any, error) {
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

func and(ast core.List, env core.Env) (core.Any, error) {
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

func or(ast core.List, env core.Env) (core.Any, error) {
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

func add(ast core.List, env core.Env) (core.Any, error) {
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
			break
		}
		res = res.Add(arg.(core.Number).Decimal())
	}
	return core.Number(res), nil
}

func sub(ast core.List, env core.Env) (core.Any, error) {
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
			break
		}
		num := arg.(core.Number).Decimal()
		if i == 0 && len(ast) > 2 {
			res = num
		} else {
			res = res.Sub(num)
		}
	}
	return core.Number(res), nil
}

func mul(ast core.List, env core.Env) (core.Any, error) {
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
			break
		}
		res = res.Mul(arg.(core.Number).Decimal())
	}
	return core.Number(res), nil
}

func quot(ast core.List, env core.Env) (core.Any, error) {
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

func rem(ast core.List, env core.Env) (core.Any, error) {
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

func quot_S(ast core.List, env core.Env) (core.Any, error) {
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
			break
		}
		num := arg.(core.Number).Decimal()
		if i == 0 && len(ast) > 2 {
			res = num
		} else {
			res = res.Div(num)
		}
	}
	return core.Number(res), nil
}

func type_(ast core.List, env core.Env) (core.Any, error) {
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

func def_E(ast core.List, env core.Env) (core.Any, error) {
	err := exactLen(ast, 3)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("'%v' first arg should be a Symbol, got '%v'", ast[0], ast[1])
	case core.Symbol:
		break
	}
	val, err := Eval(ast[2], env)
	if err != nil {
		return nil, err
	}
	env.Set(ast[1].(core.Symbol), val)
	return val, nil
}

func let_S(ast core.List, env core.Env) (core.Any, error) {
	err := exactLen(ast, 3)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("'%v' first arg should be a List, got '%v'", ast[0], ast[1])
	case core.List:
		break
	}
	newEnv, err := core.NewEnv(&env, nil, nil)
	if err != nil {
		return nil, err
	}
	pairs := ast[1].(core.List)
	if len(pairs)%2 != 0 || len(pairs) == 0 {
		return nil, fmt.Errorf("'%v' first arg should be an even List, has length %d", ast[0], len(pairs))
	}
	for {
		switch pairs[0].(type) {
		default:
			return nil, fmt.Errorf("'%v' called with non-symbol '%v'", ast[0], pairs[0])
		case core.Symbol:
			break
		}
		val, err := Eval(pairs[1], env)
		if err != nil {
			return nil, err
		}
		env.Set(pairs[0].(core.Symbol), val)
		pairs = pairs[2:]
		if len(pairs) == 0 {
			break
		}
	}
	return EvalTail(ast[2], *newEnv), nil
}
