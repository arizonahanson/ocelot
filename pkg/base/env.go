package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

type Env struct {
	outer *Env
	data  map[string]core.Any
}

func NewEnv(outer *Env, binds core.Vector, exprs core.List) (*Env, error) {
	data := make(map[string]core.Any)
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("binds and exprs must be the same length: %d", len(binds))
	}
	for i, bind := range binds {
		switch bind.(type) {
		default:
			return nil, fmt.Errorf("binds must be symbols: %v", bind)
		case core.Symbol:
			break
		}
		expr := exprs[i]
		data[bind.(core.Symbol).Val] = expr
	}
	return &Env{outer, data}, nil
}

func (env *Env) Set(sym core.Symbol, value core.Any) {
	env.data[sym.Val] = value
}

func (env *Env) Get(sym core.Symbol) (core.Any, error) {
	value, ok := env.data[sym.Val]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(sym)
		}
		return nil, fmt.Errorf("%#v: not found", sym)
	}
	return value, nil
}
