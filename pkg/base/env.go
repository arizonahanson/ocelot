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
	scope, val := env.find(sym)
	if scope == nil {
		return core.Null{}, fmt.Errorf("%#v: unable to resolve symbol", sym)
	}
	return val, nil
}

func (env *Env) Set(sym core.Symbol, val core.Any) core.Any {
	switch future := val.(type) {
	default:
		break
	case Future:
		// memoize futures
		rebind := func() (core.Any, error) {
			val, err := future.Get()
			return env.Set(sym, val), err
		}
		val = Future(rebind)
	}
	env.data[sym.Val] = val
	return val
}

// cause a future binding to resolve async
func (env *Env) Async(sym core.Symbol) error {
	scope, val := env.find(sym)
	if scope == nil {
		return fmt.Errorf("%#v: unable to resolve symbol", sym)
	}
	switch future := val.(type) {
	default:
		break
	case Future:
		scope.data[sym.Val] = future.Async()
	}
	return nil
}

func (env *Env) find(sym core.Symbol) (*Env, core.Any) {
	val, ok := env.data[sym.Val]
	if !ok {
		if env.outer != nil {
			return env.outer.find(sym)
		}
		return nil, core.Null{}
	}
	return env, val
}

func (env *Env) Del(sym core.Symbol) error {
	scope, _ := env.find(sym)
	if scope == nil {
		return fmt.Errorf("%#v: unable to resolve symbol", sym)
	}
	delete(scope.data, sym.Val)
	return nil
}
