package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

type Env struct {
	outer *Env
	data  map[core.Symbol]core.Any
}

func BaseEnv() (*Env, error) {
	env, err := NewEnv(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range Base {
		env.Set(core.Symbol(key), value)
	}
	return env, nil
}

func NewEnv(outer *Env, binds core.List, exprs core.List) (*Env, error) {
	data := make(map[core.Symbol]core.Any)
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("binds and exprs must be the same length: %d", len(binds))
	}
	for i, bind := range binds {
		switch bind.(type) {
		default:
			return nil, fmt.Errorf("binds must be symbols: %v", bind)
		case core.Symbol:
			expr := exprs[i]
			data[bind.(core.Symbol)] = expr
		}
	}
	return &Env{outer, data}, nil
}

func (env Env) Set(key core.Symbol, value core.Any) {
	env.data[key] = value
}

func (env Env) Get(key core.Symbol) (core.Any, error) {
	value, ok := env.data[key]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(key)
		}
		return nil, fmt.Errorf("not found: %s", key)
	}
	return value, nil
}

func (env Env) SetPairs(pairs core.List) error {
	if len(pairs) < 2 {
		return fmt.Errorf("missing parameter in let*")
	}
	switch pairs[0].(type) {
	default:
		return fmt.Errorf("non-symbol parameter in let*: %v", pairs[0])
	case core.Symbol:
		val, err := EvalAst(pairs[1], env)
		if err != nil {
			return err
		}
		env.Set(pairs[0].(core.Symbol), val)
		if len(pairs) > 2 {
			env.SetPairs(pairs[2:])
		}
	}
	return nil
}
