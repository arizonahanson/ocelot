package ocelot

import (
	"fmt"

	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string, env *core.Env) (core.Any, error) {
	ast, err := parser.Parse("Eval", []byte(in))
	if err != nil {
		return nil, err
	}
	if env == nil {
		base, err := core.BaseEnv()
		if err != nil {
			return nil, err
		}
		env = base
	}
	return eval_ast(ast, *env)
}

func eval_ast(ast core.Any, env core.Env) (core.Any, error) {
	switch ast.(type) {
	default:
		return ast, nil
	case core.Symbol:
		return onSymbol(ast.(core.Symbol), env)
	case core.Vector:
		return onVector(ast.(core.Vector), env)
	case core.List:
		return onList(ast.(core.List), env)
	}
}

func onSymbol(ast core.Symbol, env core.Env) (core.Any, error) {
	return env.Get(string(ast))
}

func onVector(ast core.Vector, env core.Env) (core.Vector, error) {
	res := []core.Any{}
	for _, item := range ast {
		any, err := eval_ast(item, env)
		if err != nil {
			return nil, err
		}
		res = append(res, any)
	}
	return core.Vector(res), nil
}

func defSpecial(ast core.List, env core.Env) (core.Any, error) {
	if len(ast) != 3 {
		return nil, fmt.Errorf("'def!' received wrong number of args: %d", len(ast)-1)
	}
	switch ast[1].(type) {
	default:
		return nil, fmt.Errorf("first parameter to def! must be a symbol: %s", ast[1])
	case core.Symbol:
		val, err := eval_ast(ast[2], env)
		if err != nil {
			return nil, err
		}
		env.Set(ast[1].(core.Symbol), val)
		return val, nil
	}
}

func letSpecial(ast core.List, env core.Env) (core.Any, error) {
	if len(ast) != 3 {
		return nil, fmt.Errorf("'let*' received wrong number of args: %d", len(ast)-1)
	}
	newEnv, err := core.NewEnv(&env, nil, nil)
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
	return eval_ast(ast[2], *newEnv)
}

func doSpecial(ast core.List, env core.Env) (core.Any, error) {
	var result core.Any = core.Nil{}
	for _, item := range ast[1:] {
		val, err := eval_ast(item, env)
		if err != nil {
			return nil, err
		}
		result = val
	}
	return result, nil
}

func ifSpecial(ast core.List, env core.Env) (core.Any, error) {
	if len(ast) < 3 {
		return nil, fmt.Errorf("wrong number of parameters to 'if': %d", len(ast)-1)
	}
	cond, err := eval_ast(ast[1], env)
	if err != nil {
		return nil, err
	}
	if core.IsTruthy(cond) {
		return eval_ast(ast[2], env)
	}
	if len(ast) == 4 {
		return eval_ast(ast[3], env)
	}
	return core.Nil{}, nil
}

func setPairs(pairs core.List, newEnv core.Env) error {
	if len(pairs) < 2 {
		return fmt.Errorf("missing parameter in let*")
	}
	switch pairs[0].(type) {
	default:
		return fmt.Errorf("non-symbol parameter in let*: %v", pairs[0])
	case core.Symbol:
		val, err := eval_ast(pairs[1], newEnv)
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

func fnStar(ast core.List, env core.Env) (core.Any, error) {
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
			newEnv, err := core.NewEnv(&env, binds, args)
			if err != nil {
				return nil, err
			}
			return eval_ast(body, *newEnv)
		}
		return core.Function(fn), nil
	}
}

func onList(ast core.List, env core.Env) (core.Any, error) {
	res := []core.Any{}
	if len(ast) > 0 {
		first := ast[0]
		switch first.(type) {
		default:
			break
		case core.Symbol:
			// check for special forms
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
		any, err := eval_ast(item, env)
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
		case core.Function:
			fn := first.(core.Function)
			return apply(fn, res[1:], env)
		}
	}
	// empty list
	return core.List(res), nil
}

func apply(fn core.Function, args []core.Any, env core.Env) (core.Any, error) {
	return fn(args)
}
