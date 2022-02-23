package core

import (
	"fmt"
)

type Env struct {
	outer *Env
	data  map[string]Any
}

func NewEnv(outer *Env, binds Vector, exprs List) (*Env, error) {
	data := make(map[string]Any)
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("binds and exprs must be the same length: %d", len(binds))
	}
	for i, bind := range binds {
		switch bind.(type) {
		default:
			return nil, fmt.Errorf("binds must be symbols: %v", bind)
		case *Symbol:
			break
		}
		expr := exprs[i]
		data[bind.(*Symbol).Val] = expr
	}
	return &Env{outer, data}, nil
}

func (env *Env) Set(sym *Symbol, value Any) {
	env.data[sym.Val] = value
}

func (env *Env) Get(sym *Symbol) (Any, error) {
	value, ok := env.data[sym.Val]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(sym)
		}
		return nil, fmt.Errorf("%#v: not found", sym)
	}
	return value, nil
}
