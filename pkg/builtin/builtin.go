package builtin

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
)

func BuiltinEnv() (*base.Env, error) {
	env := base.NewEnv(nil)
	for sym, value := range Builtin {
		env.Set(core.Symbol{Val: sym, Pos: nil}, value)
	}
	return env, nil
}

var Builtin = map[string]core.Any{
	// nil / bool
	"nil":    Nil{},
	"true":   Bool(true),
	"false":  Bool(false),
	"nil?":   base.Func(_nilQ),
	"true?":  base.Func(_trueQ),
	"false?": base.Func(_falseQ),
	"bool":   base.Func(_bool),
	"not":    base.Func(_not),
	"and":    base.Func(_and),
	"or":     base.Func(_or),
	// numbers
	"add":   base.Func(_add),
	"sub":   base.Func(_sub),
	"mul":   base.Func(_mul),
	"quot":  base.Func(_quot),
	"rem":   base.Func(_rem),
	"quot*": base.Func(_quotS),
	"lt?":   base.Func(_ltQ),
	"lteq?": base.Func(_lteqQ),
	"gt?":   base.Func(_gtQ),
	"gteq?": base.Func(_gteqQ),
	// special
	"func":   base.Func(_func),
	"type":   base.Func(_type),
	"equal?": base.Func(_equalQ),
	"def!":   base.Func(_defE),
	"let":    base.Func(_let),
	"do":     base.Func(_do),
	"if":     base.Func(_if),
	"prn":    base.Func(_prn),
	"eval":   base.Func(_eval),
	"quote":  base.Func(_quote),
	// lists
	"list":   base.Func(_list),
	"list?":  base.Func(_listQ),
	"empty?": base.Func(_emptyQ),
	"count":  base.Func(_count),
	"map":    base.Func(_map),
}

func _nilQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Nil{}), nil
}

func _trueQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(true)), nil
}

func _falseQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(false)), nil
}

func _bool(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 != Bool(false) && arg1 != Nil{}), nil
}

func _not(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(false) || arg1 == Nil{}), nil
}

func _and(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Bool(true), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg == Bool(false) || arg == Nil{}) {
			return arg, nil
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _or(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Bool(false), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg != Bool(false) && arg != Nil{}) {
			return arg, nil
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _add(ast core.List, env *base.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for _, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		res = res.Add(arg.Decimal())
	}
	return core.Number(res), nil
}

func _sub(ast core.List, env *base.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for i, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		if i == 0 && len(ast) > 2 {
			res = arg.Decimal()
		} else {
			res = res.Sub(arg.Decimal())
		}
	}
	return core.Number(res), nil
}

func _mul(ast core.List, env *base.Env) (core.Any, error) {
	res := core.One.Decimal()
	for _, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		res = res.Mul(arg.Decimal())
	}
	return core.Number(res), nil
}

func _quot(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 4); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg3, err := evalNumber(ast[3], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	q, _ := arg1.Decimal().QuoRem(arg2.Decimal(), int32(arg3.Decimal().IntPart()))
	return core.Number(q), nil
}

func _rem(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 4); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg3, err := evalNumber(ast[3], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	_, r := arg1.Decimal().QuoRem(arg2.Decimal(), int32(arg3.Decimal().IntPart()))
	return core.Number(r), nil
}

func _quotS(ast core.List, env *base.Env) (core.Any, error) {
	res := core.One.Decimal()
	for i, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		if i == 0 && len(ast) > 2 {
			res = arg.Decimal()
		} else {
			res = res.Div(arg.Decimal())
		}
	}
	return core.Number(res), nil
}

func _type(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	str := core.String{Val: fmt.Sprintf("%T", arg1)}
	return str, nil
}

func _defE(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-symbol %#v", ast[0], ast[1])
	case core.Symbol:
		break
	}
	val, err := base.Eval(ast[2], env)
	if err != nil {
		return nil, err
	}
	env.Set(ast[1].(core.Symbol), val)
	return val, nil
}

func _let(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-list %#v", ast[0], ast[1])
	case core.List:
		break
	}
	newEnv := base.NewEnv(env)
	pairs := ast[1].(core.List)
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("%#v: binding missing", ast[0])
	}
	for {
		if len(pairs) == 0 {
			break
		}
		switch sym := pairs[0].(type) {
		default:
			return nil, fmt.Errorf("%#v: called with non-symbol %#v", ast[0], pairs[0])
		case core.Symbol:
			val, err := base.Eval(pairs[1], env)
			if err != nil {
				return nil, err
			}
			newEnv.Set(sym, val)
			pairs = pairs[2:]
		}
	}
	return base.EvalFuture(ast[2], newEnv), nil
}

func _do(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Nil{}, nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		_, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _if(ast core.List, env *base.Env) (core.Any, error) {
	if err := rangeLen(ast, 3, 4); err != nil {
		return nil, err
	}
	cond, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	if (cond != Bool(false) && cond != Nil{}) {
		return base.EvalFuture(ast[2], env), nil
	}
	if len(ast) == 4 {
		return base.EvalFuture(ast[3], env), nil
	}
	return Nil{}, nil
}

func _func(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-vector %#v", ast[0], ast[1])
	case core.Vector:
		break
	}
	binds := ast[1].(core.Vector)
	symbols := make([]core.Symbol, len(binds))
	for i, item := range binds {
		switch sym := item.(type) {
		default:
			return nil, fmt.Errorf("%#v: bind expression contained non-symbol %#v", ast[0], item)
		case core.Symbol:
			symbols[i] = sym
			break
		}
	}
	body := ast[2]
	lambda := func(args core.List, outer *base.Env) (core.Any, error) {
		err := exactLen(args, len(symbols)+1)
		if err != nil {
			return nil, err
		}
		local := base.NewEnv(env)
		for i, symbol := range symbols {
			// bind sym to arg in local, but lazy eval arg in outer
			local.SetFuture(symbol, base.EvalFuture(args[i+1], outer))
		}
		// lazy eval body in local
		return base.EvalFuture(body, local), nil
	}
	return base.Func(lambda), nil
}

func _prn(ast core.List, env *base.Env) (core.Any, error) {
	vals, err := _list(ast, env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", vals)
	return Nil{}, nil
}

func _list(ast core.List, env *base.Env) (core.Any, error) {
	exprs := make(core.List, len(ast)-1)
	for i, item := range ast[1:] {
		val, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
		exprs[i] = val
	}
	return exprs, nil
}

func _listQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	switch val.(type) {
	default:
		return Bool(false), nil
	case core.List:
		return Bool(true), nil
	}
}

func _count(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	var cnt int
	switch any := val.(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-collection %#v", ast[0], any)
	case core.Vector:
		cnt = len(any)
		break
	case core.Map:
		cnt = len(any)
		break
	case core.List:
		cnt = len(any)
		break
	}
	return core.NewNumber(cnt), nil
}

func _emptyQ(ast core.List, env *base.Env) (core.Any, error) {
	cnt, err := _count(ast, env)
	if err != nil {
		return nil, err
	}
	return Bool(cnt.Equal(core.Zero)), nil
}

func _equalQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := minLen(ast, 3); err != nil {
		return nil, err
	}
	first, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	for _, item := range ast[2:] {
		value, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
		if !first.Equal(value) {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func _ltQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	return Bool(arg1.Decimal().LessThan(arg2.Decimal())), nil
}

func _lteqQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	return Bool(arg1.Decimal().LessThanOrEqual(arg2.Decimal())), nil
}

func _gtQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	return Bool(arg1.Decimal().GreaterThan(arg2.Decimal())), nil
}

func _gteqQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	arg1, err := evalNumber(ast[1], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	arg2, err := evalNumber(ast[2], env)
	if err != nil {
		return nil, fmt.Errorf("%#v: %s", ast[0], err)
	}
	return Bool(arg1.Decimal().GreaterThanOrEqual(arg2.Decimal())), nil
}

func _quote(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	return ast[1], nil
}

func _eval(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	// double-eval TCO'd
	return dualEvalFuture(ast[1], env), nil
}

func _map(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	switch fn := arg1.(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-function: %#v", ast[0], arg1)
	case base.Func:
		arg2, err := base.Eval(ast[2], env)
		if err != nil {
			return nil, err
		}
		switch list := arg2.(type) {
		default:
			return nil, fmt.Errorf("%#v: called with non-list: %#v", ast[0], arg2)
		case core.List:
			res := make(core.List, len(list))
			for i, item := range list {
				val, err := base.Eval(core.List{fn, item}, env)
				if err != nil {
					return nil, err
				}
				res[i] = val
			}
			return res, nil
		}
	}
}
