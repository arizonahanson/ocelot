package base

import (
	"fmt"
	"reflect"

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
	"bool":   core.Function(bool_),
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
	"type":   core.Function(type_),
	"def!":   core.Function(def_E),
	"let*":   core.Function(let_S),
	"do":     core.Function(do),
	"if":     core.Function(if_),
	"fn*":    core.Function(fn_S),
	"prn":    core.Function(prn),
	"list":   core.Function(list),
	"list?":  core.Function(list_Q),
	"empty?": core.Function(empty_Q),
	"count":  core.Function(count),
	"equal?": core.Function(equal_Q),
}

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("%v wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func rangeLen(ast core.List, min int, max int) error {
	if len(ast) < min || len(ast) > max {
		return fmt.Errorf("%v wanted %d-%d args, got %d", ast[0], min-1, max-1, len(ast)-1)
	}
	return nil
}

func minLen(ast core.List, min int) error {
	if len(ast) < min {
		return fmt.Errorf("%v wanted at least %d args, got %d", ast[0], min-1, len(ast)-1)
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

func bool_(ast core.List, env core.Env) (core.Any, error) {
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
	res := core.Zero.Decimal()
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
		case core.Number:
			break
		}
		res = res.Add(arg.(core.Number).Decimal())
	}
	return core.Number(res), nil
}

func sub(ast core.List, env core.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for i, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
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
	res := core.One.Decimal()
	for _, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
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
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
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
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
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
	res := core.One.Decimal()
	for i, item := range ast[1:] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		switch arg.(type) {
		default:
			return nil, fmt.Errorf("%v called with non-number '%v'", ast[0], arg)
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
		return nil, fmt.Errorf("%v first arg should be a Symbol, got '%v'", ast[0], ast[1])
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
		return nil, fmt.Errorf("%v first arg should be a List, got '%v'", ast[0], ast[1])
	case core.List:
		break
	}
	newEnv, err := core.NewEnv(&env, nil, nil)
	if err != nil {
		return nil, err
	}
	pairs := ast[1].(core.List)
	if len(pairs)%2 != 0 || len(pairs) == 0 {
		return nil, fmt.Errorf("%v first arg should be an even List, has length %d", ast[0], len(pairs))
	}
	for {
		switch pairs[0].(type) {
		default:
			return nil, fmt.Errorf("%v called with non-symbol '%v'", ast[0], pairs[0])
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

func do(ast core.List, env core.Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Nil{}, nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		_, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func if_(ast core.List, env core.Env) (core.Any, error) {
	err := rangeLen(ast, 3, 4)
	if err != nil {
		return nil, err
	}
	cond, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	if (cond != core.Bool(false) && cond != core.Nil{}) {
		return EvalTail(ast[2], env), nil
	}
	if len(ast) == 4 {
		return EvalTail(ast[3], env), nil
	}
	return core.Nil{}, nil
}

func fn_S(ast core.List, env core.Env) (core.Any, error) {
	err := exactLen(ast, 3)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%v called with non-list for first arg: '%v'", ast[0], ast[1])
	case core.List:
		break
	}
	binds := ast[1].(core.List)
	for _, item := range binds {
		switch item.(type) {
		default:
			return nil, fmt.Errorf("%v bind expression contained non-symbol: '%v'", ast[0], item)
		case core.Symbol:
			break
		}
	}
	body := ast[2]
	lambda := func(args core.List, outer core.Env) (core.Any, error) {
		err := exactLen(args, len(binds)+1)
		if err != nil {
			return nil, err
		}
		exprs, err := list(args, outer)
		if err != nil {
			return nil, err
		}
		newEnv, err := core.NewEnv(&env, binds, exprs.(core.List))
		if err != nil {
			return nil, err
		}
		return EvalTail(body, *newEnv), nil
	}
	return core.Function(lambda), nil
}

func prn(ast core.List, env core.Env) (core.Any, error) {
	vals, err := list(ast, env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", vals)
	return core.Nil{}, nil
}

func list(ast core.List, env core.Env) (core.Any, error) {
	exprs := []core.Any{}
	for _, item := range ast[1:] {
		val, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, val)
	}
	return core.List(exprs), nil
}

func list_Q(ast core.List, env core.Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	val, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.List:
		return core.Bool(true), nil
	}
}

func count(ast core.List, env core.Env) (core.Any, error) {
	err := exactLen(ast, 2)
	if err != nil {
		return nil, err
	}
	val, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	switch val.(type) {
	default:
		return nil, fmt.Errorf("%v arg should be a List, got '%v'", ast[0], val)
	case core.List:
		break
	}
	return core.NewNumber(len(val.(core.List))), nil
}

func empty_Q(ast core.List, env core.Env) (core.Any, error) {
	cnt, err := count(ast, env)
	if err != nil {
		return nil, err
	}
	return core.Bool(cnt.(core.Number).Decimal().Equal(core.Zero.Decimal())), nil
}

func equal_Q(ast core.List, env core.Env) (core.Any, error) {
	err := minLen(ast, 3)
	if err != nil {
		return nil, err
	}
	first, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	for _, item := range ast[2:] {
		value, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		if !isEqual(first, value) {
			return core.Bool(false), nil
		}
	}
	return core.Bool(true), nil
}

func isEqual(first core.Any, item core.Any) bool {
	if reflect.TypeOf(item) != reflect.TypeOf(first) {
		return false
	}
	switch first.(type) {
	default:
		if item != first {
			return false
		}
	case core.Function:
		return false
	case core.Symbol:
		if item.(core.Symbol).Val != first.(core.Symbol).Val {
			return false
		}
	case core.Number:
		if !item.(core.Number).Decimal().Equal(first.(core.Number).Decimal()) {
			return false
		}
	case core.List:
		if len(first.(core.List)) != len(item.(core.List)) {
			return false
		}
		for i, a := range first.(core.List) {
			b := item.(core.List)[i]
			if !isEqual(a, b) {
				return false
			}
		}
	case core.Vector:
		if len(first.(core.Vector)) != len(item.(core.Vector)) {
			return false
		}
		for i, a := range first.(core.Vector) {
			b := item.(core.Vector)[i]
			if !isEqual(a, b) {
				return false
			}
		}
	}
	return true
}
