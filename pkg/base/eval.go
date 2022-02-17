package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

func EvalAst(ast core.Any, env Env) (core.Any, error) {
	switch ast.(type) {
	default:
		return ast, nil
	case core.Symbol:
		return env.Get(ast.(core.Symbol))
	case core.List:
		return onList(ast.(core.List), env)
	case core.Vector:
		return onVector(ast.(core.Vector), env)
	}
}

type Function func(args []core.Any) (core.Any, error)

func onList(ast core.List, env Env) (core.Any, error) {
	res := []core.Any{}
	if len(ast) > 0 {
		// check for special forms
		first := ast[0]
		switch first.(type) {
		default:
			break
		case core.Symbol:
			switch first.(core.Symbol) {
			default:
				break
			case core.Symbol("def!"):
				return defSpecial(ast, env)
			case core.Symbol("let*"):
				return letSpecial(ast, env)
			case core.Symbol("do"):
				return doSpecial(ast, env)
			case core.Symbol("if"):
				return ifSpecial(ast, env)
			case core.Symbol("fn*"):
				return fnStar(ast, env)
			}
		}
	}
	// default list resolution
	for _, item := range ast {
		any, err := EvalAst(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	// check for function application
	if len(res) > 0 {
		first := res[0]
		switch first.(type) {
		default:
			return core.List(res), nil
		case Function:
			fn := first.(Function)
			return fn(res[1:])
		}
	}
	// empty list
	return core.List(res), nil
}

func onVector(ast core.Vector, env Env) (core.Vector, error) {
	res := []core.Any{}
	for _, item := range ast {
		any, err := EvalAst(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	return core.Vector(res), nil
}

func defSpecial(ast core.List, env Env) (core.Any, error) {
	if len(ast) != 3 {
		return nil, fmt.Errorf("'def!' received wrong number of args: %d", len(ast)-1)
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("first parameter to def! must be a symbol: %s", ast[1])
	case core.Symbol:
		val, err := EvalAst(ast[2], env)
		if err != nil {
			return nil, err
		}
		env.Set(ast[1].(core.Symbol), val)
		return val, nil
	}
}

func letSpecial(ast core.List, env Env) (core.Any, error) {
	if len(ast) != 3 {
		return nil, fmt.Errorf("'let*' received wrong number of args: %d", len(ast)-1)
	}
	newEnv, err := NewEnv(&env, nil, nil)
	if err != nil {
		return nil, err
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("first parameter to let* must be a list: %v", ast[1])
	case core.List:
		pairs := ast[1].(core.List)
		setPairs(pairs, *newEnv)
	}
	return EvalAst(ast[2], *newEnv)
}

func doSpecial(ast core.List, env Env) (core.Any, error) {
	var result core.Any = core.Nil{}
	for _, item := range ast[1:] {
		val, err := EvalAst(item, env)
		if err != nil {
			return nil, err
		}
		result = val
	}
	return result, nil
}

func ifSpecial(ast core.List, env Env) (core.Any, error) {
	if len(ast) < 3 {
		return nil, fmt.Errorf("wrong number of parameters to 'if': %d", len(ast)-1)
	}
	cond, err := EvalAst(ast[1], env)
	if err != nil {
		return nil, err
	}
	if (cond != core.Bool(false) && cond != core.Nil{}) {
		return EvalAst(ast[2], env)
	}
	if len(ast) == 4 {
		return EvalAst(ast[3], env)
	}
	return core.Nil{}, nil
}

func setPairs(pairs core.List, newEnv Env) error {
	if len(pairs) < 2 {
		return fmt.Errorf("missing parameter in let*")
	}
	switch pairs[0].(type) {
	default:
		return fmt.Errorf("non-symbol parameter in let*: %v", pairs[0])
	case core.Symbol:
		val, err := EvalAst(pairs[1], newEnv)
		if err != nil {
			return err
		}
		newEnv.Set(pairs[0].(core.Symbol), val)
		if len(pairs) > 2 {
			setPairs(pairs[2:], newEnv)
		}
	}
	return nil
}

func fnStar(ast core.List, env Env) (core.Any, error) {
	if len(ast) < 3 {
		return nil, fmt.Errorf("wrong number of arguments to 'fn*', got: %d", len(ast)-1)
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("wrong type for parameter list: %v", ast[1])
	case core.List:
		binds := ast[1].(core.List)
		body := ast[2]
		fn := func(args []core.Any) (core.Any, error) {
			newEnv, err := NewEnv(&env, binds, args)
			if err != nil {
				return nil, err
			}
			return EvalAst(body, *newEnv)
		}
		return Function(fn), nil
	}
}
