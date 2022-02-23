package base

import (
	"fmt"
	"reflect"

	"github.com/starlight/ocelot/pkg/core"
)

type Nil struct{}

type Bool bool

func (val Nil) String() string {
	return "nil"
}

func (val Nil) GoString() string {
	return val.String()
}

var Base = map[string]core.Any{
	// nil / bool
	"nil":    Nil{},
	"true":   Bool(true),
	"false":  Bool(false),
	"nil?":   core.Function(nilQ),
	"true?":  core.Function(trueQ),
	"false?": core.Function(falseQ),
	"bool":   core.Function(_bool),
	"not":    core.Function(not),
	"and":    core.Function(and),
	"or":     core.Function(or),
	// numbers
	"add":   core.Function(add),
	"sub":   core.Function(sub),
	"mul":   core.Function(mul),
	"quot":  core.Function(quot),
	"rem":   core.Function(rem),
	"quot*": core.Function(quotS),
	"lt?":   core.Function(ltQ),
	"lteq?": core.Function(lteqQ),
	"gt?":   core.Function(gtQ),
	"gteq?": core.Function(gteqQ),
	// special
	"type":   core.Function(_type),
	"equal?": core.Function(equalQ),
	"def!":   core.Function(defE),
	"let":    core.Function(let),
	"do":     core.Function(do),
	"if":     core.Function(_if),
	"fn*":    core.Function(fnS),
	"prn":    core.Function(prn),
	"eval":   core.Function(eval),
	"quote":  core.Function(quote),
	// lists
	"list":   core.Function(list),
	"list?":  core.Function(listQ),
	"empty?": core.Function(emptyQ),
	"count":  core.Function(count),
}

func exactLen(ast core.List, num int) error {
	if len(ast) != num {
		return fmt.Errorf("%#v: wanted %d arg(s), got %d", ast[0], num-1, len(ast)-1)
	}
	return nil
}

func rangeLen(ast core.List, min int, max int) error {
	if len(ast) < min || len(ast) > max {
		return fmt.Errorf("%#v: wanted %d-%d args, got %d", ast[0], min-1, max-1, len(ast)-1)
	}
	return nil
}

func minLen(ast core.List, min int) error {
	if len(ast) < min {
		return fmt.Errorf("%#v: wanted at least %d args, got %d", ast[0], min-1, len(ast)-1)
	}
	return nil
}

func nilQ(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Nil{}), nil
}

func trueQ(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(true)), nil
}

func falseQ(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(false)), nil
}

func _bool(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 != Bool(false) && arg1 != Nil{}), nil
}

func not(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	return Bool(arg1 == Bool(false) || arg1 == Nil{}), nil
}

func and(ast core.List, env *core.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Bool(true), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg == Bool(false) || arg == Nil{}) {
			return arg, nil
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func or(ast core.List, env *core.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Bool(false), nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		arg, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		if (arg != Bool(false) && arg != Nil{}) {
			return arg, nil
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func evalNumber(ast core.Any, env *core.Env) (*core.Number, error) {
	arg, err := Eval(ast, env)
	if err != nil {
		return nil, err
	}
	switch val := arg.(type) {
	default:
		return nil, fmt.Errorf("called with non-number %#v", val)
	case *core.Number:
		return val, nil
	}
}

func add(ast core.List, env *core.Env) (core.Any, error) {
	res := core.Zero.Decimal()
	for _, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		res = res.Add(arg.Decimal())
	}
	return (*core.Number)(&res), nil
}

func sub(ast core.List, env *core.Env) (core.Any, error) {
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
	return (*core.Number)(&res), nil
}

func mul(ast core.List, env *core.Env) (core.Any, error) {
	res := core.One.Decimal()
	for _, item := range ast[1:] {
		arg, err := evalNumber(item, env)
		if err != nil {
			return nil, fmt.Errorf("%#v: %s", ast[0], err)
		}
		res = res.Mul(arg.Decimal())
	}
	return (*core.Number)(&res), nil
}

func quot(ast core.List, env *core.Env) (core.Any, error) {
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
	return (*core.Number)(&q), nil
}

func rem(ast core.List, env *core.Env) (core.Any, error) {
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
	return (*core.Number)(&r), nil
}

func quotS(ast core.List, env *core.Env) (core.Any, error) {
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
	return (*core.Number)(&res), nil
}

func _type(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	arg1, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	str := core.String(fmt.Sprintf("%T", arg1))
	return &str, nil
}

func defE(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-symbol %#v", ast[0], ast[1])
	case *core.Symbol:
		break
	}
	val, err := Eval(ast[2], env)
	if err != nil {
		return nil, err
	}
	env.Set(ast[1].(*core.Symbol), val)
	return val, nil
}

func let(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 3); err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("%#v: called with non-list %#v", ast[0], ast[1])
	case core.List:
		break
	}
	newEnv, err := core.NewEnv(env, nil, nil)
	if err != nil {
		return nil, err
	}
	pairs := ast[1].(core.List)
	if len(pairs)%2 != 0 || len(pairs) == 0 {
		return nil, fmt.Errorf("%#v: binding missing", ast[0])
	}
	for {
		switch pairs[0].(type) {
		default:
			return nil, fmt.Errorf("%#v: called with non-symbol %#v", ast[0], pairs[0])
		case *core.Symbol:
			break
		}
		val, err := Eval(pairs[1], env)
		if err != nil {
			return nil, err
		}
		newEnv.Set(pairs[0].(*core.Symbol), val)
		pairs = pairs[2:]
		if len(pairs) == 0 {
			break
		}
	}
	return EvalTail(ast[2], newEnv), nil
}

func do(ast core.List, env *core.Env) (core.Any, error) {
	if len(ast) == 1 {
		return Nil{}, nil
	}
	for _, item := range ast[1 : len(ast)-1] {
		_, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
	}
	return EvalTail(ast[len(ast)-1], env), nil
}

func _if(ast core.List, env *core.Env) (core.Any, error) {
	if err := rangeLen(ast, 3, 4); err != nil {
		return nil, err
	}
	cond, err := Eval(ast[1], env)
	if err != nil {
		return nil, err
	}
	if (cond != Bool(false) && cond != Nil{}) {
		return EvalTail(ast[2], env), nil
	}
	if len(ast) == 4 {
		return EvalTail(ast[3], env), nil
	}
	return Nil{}, nil
}

func fnS(ast core.List, env *core.Env) (core.Any, error) {
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
		case *core.Symbol:
			break
		}
	}
	body := ast[2]
	lambda := func(args core.List, outer *core.Env) (core.Any, error) {
		err := exactLen(args, len(binds)+1)
		if err != nil {
			return nil, err
		}
		exprs, err := list(args, outer)
		if err != nil {
			return nil, err
		}
		newEnv, err := core.NewEnv(env, binds, exprs.(core.List))
		if err != nil {
			return nil, err
		}
		return EvalTail(body, newEnv), nil
	}
	return core.Function(lambda), nil
}

func prn(ast core.List, env *core.Env) (core.Any, error) {
	vals, err := list(ast, env)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", vals)
	return Nil{}, nil
}

func list(ast core.List, env *core.Env) (core.Any, error) {
	exprs := make(core.List, len(ast)-1)
	for i, item := range ast[1:] {
		val, err := Eval(item, env)
		if err != nil {
			return nil, err
		}
		exprs[i] = val
	}
	return exprs, nil
}

func listQ(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	val, err := Eval(ast[1], env)
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

func count(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	val, err := Eval(ast[1], env)
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
	case core.List:
		cnt = len(any)
		break
	}
	return core.NewNumber(cnt), nil
}

func emptyQ(ast core.List, env *core.Env) (core.Any, error) {
	cnt, err := count(ast, env)
	if err != nil {
		return nil, err
	}
	return Bool(cnt.(*core.Number).Decimal().Equal(core.Zero.Decimal())), nil
}

func equalQ(ast core.List, env *core.Env) (core.Any, error) {
	if err := minLen(ast, 3); err != nil {
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
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func isEqual(first core.Any, item core.Any) bool {
	if reflect.TypeOf(item) != reflect.TypeOf(first) {
		return false
	}
	switch val := first.(type) {
	default:
		return item == val
	case core.Function:
		return false
	case *core.String:
		if *item.(*core.String) != *val {
			return false
		}
	case *core.Symbol:
		if item.(*core.Symbol).Val != val.Val {
			return false
		}
	case *core.Key:
		if *item.(*core.Key) != *val {
			return false
		}
	case *core.Number:
		if !item.(*core.Number).Decimal().Equal(val.Decimal()) {
			return false
		}
	case core.List:
		if len(val) != len(item.(core.List)) {
			return false
		}
		for i, a := range val {
			b := item.(core.List)[i]
			if !isEqual(a, b) {
				return false
			}
		}
	case core.Vector:
		if len(val) != len(item.(core.Vector)) {
			return false
		}
		for i, a := range val {
			b := item.(core.Vector)[i]
			if !isEqual(a, b) {
				return false
			}
		}
	}
	return true
}

func ltQ(ast core.List, env *core.Env) (core.Any, error) {
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

func lteqQ(ast core.List, env *core.Env) (core.Any, error) {
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

func gtQ(ast core.List, env *core.Env) (core.Any, error) {
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

func gteqQ(ast core.List, env *core.Env) (core.Any, error) {
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

func quote(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	return ast[1], nil
}

func eval(ast core.List, env *core.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return nil, err
	}
	// double-eval TCO'd
	return dualEvalTail(ast[1], env), nil
}

func dualEvalTail(ast core.Any, env *core.Env) Thunk {
	return func() (core.Any, error) {
		val, err := Eval(ast, env)
		if err != nil {
			return nil, err
		}
		return Eval(val, env)
	}
}
