package core

import (
	"fmt"
)

type Env struct {
	outer *Env
	data  map[Symbol]Any
}

func NewEnv(outer *Env, binds List, exprs List) (*Env, error) {
	data := make(map[Symbol]Any)
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("binds and exprs must be the same length: %d", len(binds))
	}
	for i, bind := range binds {
		switch bind.(type) {
		default:
			return nil, fmt.Errorf("binds must be symbols: %v", bind)
		case Symbol:
			expr := exprs[i]
			data[bind.(Symbol)] = expr
		}
	}
	return &Env{outer, data}, nil
}

func (env Env) Set(key Symbol, value Any) {
	env.data[key] = value
}

func (env Env) Get(key Symbol) (Any, error) {
	value, ok := env.data[key]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(key)
		}
		return nil, fmt.Errorf("not found: %s", key)
	}
	return value, nil
}
