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
	"nil?":   base.Func(nilQ),
	"true?":  base.Func(trueQ),
	"false?": base.Func(falseQ),
	"bool":   base.Func(_bool),
	"not":    base.Func(not),
	"and":    base.Func(and),
	"or":     base.Func(or),
	// numbers
	"add":   base.Func(add),
	"sub":   base.Func(sub),
	"mul":   base.Func(mul),
	"quot":  base.Func(quot),
	"rem":   base.Func(rem),
	"quot*": base.Func(quotS),
	"lt?":   base.Func(ltQ),
	"lteq?": base.Func(lteqQ),
	"gt?":   base.Func(gtQ),
	"gteq?": base.Func(gteqQ),
	// special
	"type":   base.Func(_type),
	"equal?": base.Func(equalQ),
	"def!":   base.Func(defE),
	"let":    base.Func(let),
	"do":     base.Func(do),
	"if":     base.Func(_if),
	"fn*":    base.Func(fnS),
	"prn":    base.Func(prn),
	"eval":   base.Func(eval),
	"quote":  base.Func(quote),
	// lists
	"list":   base.Func(list),
	"list?":  base.Func(listQ),
	"empty?": base.Func(emptyQ),
	"count":  base.Func(count),
	"map":    base.Func(_map),
}

func nilQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Nil{}), nil
}

func trueQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(true)), nil
}

func falseQ(ast core.List, env *base.Env) (core.Any, error) {
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

func not(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := base.Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(false) || arg1 == Nil{}), nil
}

func and(ast core.List, env *base.Env) (core.Any, error) {
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
	return base.EvalLazy(ast[len(ast)-1], env), nil
}

func or(ast core.List, env *base.Env) (core.Any, error) {
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
	return base.EvalLazy(ast[len(ast)-1], env), nil
}

func add(ast core.List, env *base.Env) (core.Any, error) {
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

func sub(ast core.List, env *base.Env) (core.Any, error) {
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

func mul(ast core.List, env *base.Env) (core.Any, error) {
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

func quot(ast core.List, env *base.Env) (core.Any, error) {
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

func rem(ast core.List, env *base.Env) (core.Any, error) {
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

func quotS(ast core.List, env *base.Env) (core.Any, error) {
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

func defE(ast core.List, env *base.Env) (core.Any, error) {
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

func let(ast core.List, env *base.Env) (core.Any, error) {
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
	return base.EvalLazy(ast[2], newEnv), nil
}

func do(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Nil{}, nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		_, err := base.Eval(item, env)
		if err != nil {
			return nil, err
		}
	}
	return base.EvalLazy(ast[len(ast)-1], env), nil
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
		return base.EvalLazy(ast[2], env), nil
	}
	if len(ast) == 4 {
		return base.EvalLazy(ast[3], env), nil
	}
	return Nil{}, nil
}

func fnS(ast core.List, env *base.Env) (core.Any, error) {
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
	for _, item := range binds {
		switch item.(type) {
		default:
			return nil, fmt.Errorf("%#v: bind expression contained non-symbol %#v", ast[0], item)
		case core.Symbol:
			break
		}
	}
	body := ast[2]
	lambda := func(args core.List, outer *base.Env) (core.Any, error) {
		err := exactLen(args, len(binds)+1)
		if err != nil {
			return nil, err
		}
		// binds in fn* scope, args eval in outer scope lazily
		inner := base.NewEnv(env)
		for i, bind := range binds {
			expr := args[i+1]
			// in inner env, bind sym to expr, lazy eval in outer env
			inner.SetLazy(bind.(core.Symbol), base.EvalLazy(expr, outer))
		}
		// done with these scopes
		return base.EvalLazy(body, inner), nil
	}
	return base.Func(lambda), nil
}

func prn(ast core.List, env *base.Env) (core.Any, error) {
	vals, err := list(ast, env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", vals)
	return Nil{}, nil
}

func list(ast core.List, env *base.Env) (core.Any, error) {
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

func listQ(ast core.List, env *base.Env) (core.Any, error) {
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

func count(ast core.List, env *base.Env) (core.Any, error) {
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

func emptyQ(ast core.List, env *base.Env) (core.Any, error) {
	cnt, err := count(ast, env)
	if err != nil {
		return nil, err
	}
	return Bool(cnt.Equal(core.Zero)), nil
}

func equalQ(ast core.List, env *base.Env) (core.Any, error) {
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

func ltQ(ast core.List, env *base.Env) (core.Any, error) {
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

func lteqQ(ast core.List, env *base.Env) (core.Any, error) {
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

func gtQ(ast core.List, env *base.Env) (core.Any, error) {
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

func gteqQ(ast core.List, env *base.Env) (core.Any, error) {
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

func quote(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	return ast[1], nil
}

func eval(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	// double-eval TCO'd
	return dualEvalLazy(ast[1], env), nil
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
