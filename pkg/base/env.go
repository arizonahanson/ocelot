package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

type Env struct {
	outer *Env
	data  map[string]core.Any
}

func NewEnv(outer *Env) *Env {
	data := make(map[string]core.Any)
	return &Env{outer, data}
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

func (env *Env) Set(sym core.Symbol, value core.Any) {
	switch future := value.(type) {
	default:
		break
	case Future:
		// bind symbol in env to future value evaluated once
		once := func() (val core.Any, err error) {
			val, err = future.Resolve()
			if err != nil {
				return
			}
			env.Set(sym, val)
			return
		}
		value = Future(once)
	}
	env.data[sym.Val] = value
}
