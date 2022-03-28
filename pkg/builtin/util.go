package builtin

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/base"
	"github.com/starlight/ocelot/pkg/core"
)

func exactLen(ast core.Expr, num int) error {
	if len(ast) != num {
		return fmt.Errorf("wanted %d arg(s), got %d", num-1, len(ast)-1)
	}
	return nil
}

func rangeLen(ast core.Expr, min int, max int) error {
	if len(ast) < min || len(ast) > max {
		return fmt.Errorf("wanted %d-%d args, got %d", min-1, max-1, len(ast)-1)
	}
	return nil
}

func minLen(ast core.Expr, min int) error {
	if len(ast) < min {
		return fmt.Errorf("wanted at least %d args, got %d", min-1, len(ast)-1)
	}
	return nil
}

func evalNumber(ast core.Any, env *base.Env) (*core.Number, error) {
	val, err := base.Eval(ast, env)
	if err != nil {
		return nil, err
	}
	switch num := val.(type) {
	default:
		return nil, fmt.Errorf("called with non-number %#v", val)
	case core.Number:
		return &num, nil
	}
}

func oneArg(ast core.Expr, env *base.Env) (core.Any, error) {
	if err := exactLen(ast, 2); err != nil {
		return core.Null{}, err
	}
	val, err := base.Eval(ast[1], env)
	if err != nil {
		return core.Null{}, err
	}
	return val, nil
}

// eval then eval the result (lazy)
func dualFutureEval(ast core.Any, env *base.Env) base.Future {
	return func() (core.Any, error) {
		val, err := base.Eval(ast, env)
		if err != nil {
			return val, err
		}
		return base.Eval(val, env)
	}
}

func cons(first core.Any, rest core.Expr) core.Expr {
	ast := make(core.Expr, len(rest)+1)
	ast[0] = first
	for i, item := range rest {
		ast[i+1] = item
	}
	return ast
}
