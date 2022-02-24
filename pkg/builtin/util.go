package builtin

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
)

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

func evalNumber(ast core.Any, env *base.Env) (*core.Number, error) {
	arg, err := base.Eval(ast, env)
	if err != nil {
		return nil, err
	}
	switch val := arg.(type) {
	default:
		return nil, fmt.Errorf("called with non-number %#v", val)
	case core.Number:
		return &val, nil
	}
}

func bindLazy(inner *base.Env, sym core.Symbol, outer *base.Env, expr core.Any) {
	lazy := func() (core.Any, error) {
		val, err := base.Eval(expr, outer)
		if err != nil {
			return nil, err
		}
		// memoize
		inner.Set(sym, val)
		return val, nil
	}
	inner.Set(sym, base.Lazy(lazy))
}

// eval then lazy-eval the result
func dualEvalLazy(ast core.Any, env *base.Env) base.Lazy {
	return func() (core.Any, error) {
		val, err := base.Eval(ast, env)
		if err != nil {
			return nil, err
		}
		return base.EvalLazy(val, env), nil
	}
}
