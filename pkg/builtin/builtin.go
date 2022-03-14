package builtin

import (
	"fmt"
	"time"

	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
)

func BuiltinEnv() (*base.Env, error) {
	env := base.NewEnv(nil)
	for sym, val := range Builtin {
		env.Set(core.NewSymbol(sym, nil), val)
	}
	return env, nil
}

var Builtin = map[string]core.Any{
	// nil / bool
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
	"div":   base.Func(_div),
	"rem":   base.Func(_rem),
	"div*":  base.Func(_divS),
	"lt?":   base.Func(_ltQ),
	"lteq?": base.Func(_lteqQ),
	"gt?":   base.Func(_gtQ),
	"gteq?": base.Func(_gteqQ),
	// special
	"equal?": base.Func(_equalQ),
	"def!":   base.Func(_defE),
	"deffn!": base.Func(_deffnE),
	"do":     base.Func(_do),
	"func":   base.Func(_func),
	"let":    base.Func(_let),
	"wait":   base.Func(_wait),
	"async":  base.Func(_async),
	"if":     base.Func(_if),
	"prn":    base.Func(_prn),
	"eval":   base.Func(_eval),
	"quote":  base.Func(_quote),
	"map":    base.Func(_map),
	"apply":  base.Func(_apply),
	"throw":  base.Func(_throw),
	"try":    base.Func(_try),
	"catch":  base.Func(_func), // alias
	// type check
	"type":    base.Func(_type),
	"bool?":   base.Func(_boolQ),
	"number?": base.Func(_numberQ),
	"string?": base.Func(_stringQ),
	"symbol?": base.Func(_symbolQ),
	"key?":    base.Func(_keyQ),
	"list?":   base.Func(_listQ),
	"vector?": base.Func(_vectorQ),
	"map?":    base.Func(_mapQ),
	// sequences
	"empty?": base.Func(_emptyQ),
	"count":  base.Func(_count),
}

func _nilQ(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val == core.Nil{}), nil
}

func _trueQ(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val == core.Bool(true)), nil
}

func _falseQ(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val == core.Bool(false)), nil
}

func _bool(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val != core.Bool(false) && val != core.Nil{}), nil
}

func _not(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val == core.Bool(false) || val == core.Nil{}), nil
}

func _and(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Bool(true), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		val, err := base.Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		if (val == core.Bool(false) || val == core.Nil{}) {
			return val, nil
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _or(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Bool(false), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		val, err := base.Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		if (val != core.Bool(false) && val != core.Nil{}) {
			return val, nil
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _numberQ(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Number:
		return core.Bool(true), nil
	}
}

func _add(ast core.List, env *base.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for _, item := range ast[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		res = res.Add(val.Decimal())
	}
	return core.Number(res), nil
}

func _sub(ast core.List, env *base.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for i, item := range ast[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		if i == 0 && len(ast) > 2 {
			res = val.Decimal()
		} else {
			res = res.Sub(val.Decimal())
		}
	}
	return core.Number(res), nil
}

func _mul(ast core.List, env *base.Env) (core.Any, error) {
	res := core.One.Decimal()
	for _, item := range ast[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		res = res.Mul(val.Decimal())
	}
	return core.Number(res), nil
}

func _div(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 4); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	val3, err := evalNumber(ast[3], env)
	if err != nil {
		return core.Nil{}, err
	}
	q, _ := val1.Decimal().QuoRem(val2.Decimal(), int32(val3.Decimal().IntPart()))
	return core.Number(q), nil
}

func _rem(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 4); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	val3, err := evalNumber(ast[3], env)
	if err != nil {
		return core.Nil{}, err
	}
	_, r := val1.Decimal().QuoRem(val2.Decimal(), int32(val3.Decimal().IntPart()))
	return core.Number(r), nil
}

func _divS(ast core.List, env *base.Env) (core.Any, error) {
	res := core.One.Decimal()
	for i, item := range ast[1:] {
		val, err := evalNumber(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		if i == 0 && len(ast) > 2 {
			res = val.Decimal()
		} else {
			res = res.Div(val.Decimal())
		}
	}
	return core.Number(res), nil
}

func _type(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	if (val == core.Nil{}) {
		return core.Nil{}, nil
	}
	str := core.String{Val: fmt.Sprintf("%T", val)}
	return str, nil
}

func _defE(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	switch sym := ast[1].(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-symbol %#v", ast[1])
	case core.Symbol:
		env.Set(sym, base.EvalFuture(ast[2], env))
		return core.Nil{}, nil
	}
}

func _deffnE(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 4); err != nil {
		return core.Nil{}, err
	}
	switch ast[2].(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-vector %#v", ast[2])
	case core.Vector:
		break
	}
	fn := base.Func(_func).Future(cons(ast[0], ast[2:]), env)
	return _defE(append(ast[:2], fn), env)
}

func _let(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	var pairs core.Vector
	switch arg1 := ast[1].(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-sequence %#v", ast[1])
	case core.Vector:
		pairs = arg1
		break
	case core.List:
		pairs = core.Vector(arg1)
		break
	}
	if len(pairs)%2 != 0 {
		return core.Nil{}, fmt.Errorf("binding missing")
	}
	newEnv := base.NewEnv(env)
	for {
		if len(pairs) == 0 {
			break
		}
		switch sym := pairs[0].(type) {
		default:
			return core.Nil{}, fmt.Errorf("called with non-symbol %#v", pairs[0])
		case core.Symbol:
			newEnv.Set(sym, base.EvalFuture(pairs[1], newEnv))
			pairs = pairs[2:]
		}
	}
	return base.EvalFuture(ast[2], newEnv), nil
}

func _do(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) == 1 {
		return core.Nil{}, nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		_, err := base.Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
	}
	return base.EvalFuture(ast[len(ast)-1], env), nil
}

func _if(ast core.List, env *base.Env) (core.Any, error) {
	if err := rangeLen(ast, 3, 4); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	if (val != core.Bool(false) && val != core.Nil{}) {
		return base.EvalFuture(ast[2], env), nil
	}
	if len(ast) == 4 {
		return base.EvalFuture(ast[3], env), nil
	}
	return core.Nil{}, nil
}

func _func(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	switch ast[1].(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-vector %#v", ast[1])
	case core.Vector:
		break
	}
	binds := ast[1].(core.Vector)
	symbols := make([]core.Symbol, len(binds))
	for i, item := range binds {
		switch sym := item.(type) {
		default:
			return core.Nil{}, fmt.Errorf("bind expression contained non-symbol %#v", item)
		case core.Symbol:
			symbols[i] = sym
			break
		}
	}
	body := ast[2]
	fn := func(args core.List, outer *base.Env) (core.Any, error) {
		err := exactLen(args, len(symbols)+1)
		if err != nil {
			return core.Nil{}, err
		}
		local := base.NewEnv(env)
		for i, sym := range symbols {
			// bind sym to arg in local, but lazy eval arg in outer
			local.Set(sym, base.EvalFuture(args[i+1], outer))
		}
		// future that places breaks in error trace
		future := func() (val core.Any, err error) {
			val, err = base.Eval(body, local)
			if err != nil {
				err = fmt.Errorf("error\n  %v", err)
			}
			return
		}
		return base.Future(future), nil
	}
	return base.Func(fn), nil
}

func _prn(ast core.List, env *base.Env) (core.Any, error) {
	if len(ast) > 0 {
		var str string
		for i, arg := range ast[1:] {
			if i != 0 {
				str += " "
			}
			val, err := base.Eval(arg, env)
			if err != nil {
				return core.Nil{}, err
			}
			str += fmt.Sprintf("%v", val)
		}
		fmt.Println(str)
	}
	return core.Nil{}, nil
}

func _listQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.List:
		return core.Bool(true), nil
	}
}

func _vectorQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Vector:
		return core.Bool(true), nil
	}
}

func _symbolQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Symbol:
		return core.Bool(true), nil
	}
}

func _boolQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Bool:
		return core.Bool(true), nil
	}
}

func _keyQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Key:
		return core.Bool(true), nil
	}
}

func _stringQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.String:
		return core.Bool(true), nil
	}
}

func _mapQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch val.(type) {
	default:
		return core.Bool(false), nil
	case core.Map:
		return core.Bool(true), nil
	}
}

func _count(ast core.List, env *base.Env) (core.Any, error) {
	val, err := oneLen(ast, env)
	if err != nil {
		return core.Nil{}, err
	}
	var cnt int
	switch any := val.(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-collection %#v", any)
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
		return core.Nil{}, err
	}
	return core.Bool(cnt.Equal(core.Zero)), nil
}

func _equalQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := minLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	for _, item := range ast[2:] {
		val2, err := base.Eval(item, env)
		if err != nil {
			return core.Nil{}, err
		}
		if !val1.Equal(val2) {
			return core.Bool(false), nil
		}
	}
	return core.Bool(true), nil
}

func _ltQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val1.Decimal().LessThan(val2.Decimal())), nil
}

func _lteqQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val1.Decimal().LessThanOrEqual(val2.Decimal())), nil
}

func _gtQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val1.Decimal().GreaterThan(val2.Decimal())), nil
}

func _gteqQ(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	val2, err := evalNumber(ast[2], env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Bool(val1.Decimal().GreaterThanOrEqual(val2.Decimal())), nil
}

func _quote(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	return ast[1], nil
}

func _eval(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	// double-eval TCO'd
	return dualEvalFuture(ast[1], env), nil
}

// converts lazy-futures to async-futures
func _async(ast core.List, env *base.Env) (res core.Any, err error) {
	res = core.Nil{}
	if err = exactLen(ast, 2); err != nil {
		return
	}
	switch arg := ast[1].(type) {
	default:
		err = fmt.Errorf("called with non-symbol %#v", ast[1])
	case core.Symbol:
		err = env.Async(arg)
	case core.Vector:
	loop:
		for _, item := range arg {
			switch sym := item.(type) {
			default:
				err = fmt.Errorf("called with non-symbol %#v", item)
				break loop
			case core.Symbol:
				if err = env.Async(sym); err != nil {
					break loop
				}
			}
		}
	}
	return
}

func _apply(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch fn := val1.(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-function: %#v", val1)
	case base.Func:
		val2, err := base.Eval(ast[2], env)
		if err != nil {
			return core.Nil{}, err
		}
		switch vec := val2.(type) {
		default:
			return core.Nil{}, fmt.Errorf("called with non-vector: %#v", val2)
		case core.Vector:
			return fn.Future(cons(ast[1], core.List(vec)), env), nil
		}
	}
}

func _map(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	val1, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	switch fn := val1.(type) {
	default:
		return core.Nil{}, fmt.Errorf("called with non-function: %#v", val1)
	case base.Func:
		val2, err := base.Eval(ast[2], env)
		if err != nil {
			return core.Nil{}, err
		}
		switch val2.(type) {
		default:
			return core.Nil{}, fmt.Errorf("called with non-vector: %#v", val2)
		case core.Vector:
			break
		}
		res := make(core.Vector, len(val2.(core.Vector)))
		for i, item := range val2.(core.Vector) {
			res[i] = fn.Future(cons(ast[1], core.List{item}), env)
		}
		return base.EvalFuture(res, env), nil
	}
}

func _wait(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	arg, err := evalNumber(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%ss", arg))
	if err != nil {
		return core.Nil{}, err
	}
	time.Sleep(dur)
	return core.Nil{}, nil
}

func _throw(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Nil{}, err
	}
	arg, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Nil{}, err
	}
	return core.Nil{}, fmt.Errorf("%s", arg)
}

func _try(ast core.List, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return core.Nil{}, err
	}
	res, err := base.Eval(ast[1], env)
	if err != nil {
		res2, err2 := base.Eval(ast[2], env)
		if err2 != nil {
			return core.Nil{}, err
		}
		switch catch := res2.(type) {
		default:
			return core.Nil{}, fmt.Errorf("called with non-function catch %#v", res2)
		case base.Func:
			// lazy-call catch-function
			return catch.Future(core.List{ast[2], core.String{Val: err.Error()}}, env), nil
		}
	}
	return res, nil
}
