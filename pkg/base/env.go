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
	val, ok := env.data[sym.Val]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(sym)
		}
		return core.Nil{}, fmt.Errorf("%#v: unable to resolve symbol", sym)
	}
	return val, nil
}

func (env *Env) Set(sym core.Symbol, val core.Any) core.Any {
	switch future := val.(type) {
	default:
		// not a future
		break
	case Future:
		// wrap future in future that updates binding
		rebind := func() (core.Any, error) {
			val, err := future.Get()
			return env.Set(sym, val), err
		}
		val = Future(rebind)
	}
	env.data[sym.Val] = val
	return val
}
