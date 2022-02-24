package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

type Env struct {
	outer *Env
	data  map[string]core.Any
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
		return nil, fmt.Errorf("%#v: unable to resolve symbol", sym)
	}
	return value, nil
}

func NewEnv(outer *Env) *Env {
	data := make(map[string]core.Any)
	return &Env{outer, data}
}
