package core

import (
	"fmt"
)

type Env struct {
	outer *Env
	data  map[Symbol]Any
}

func NewEnv(outer *Env) Env {
	data := make(map[Symbol]Any)
	return Env{outer, data}
}

func (env Env) Set(key Symbol, value Any) {
	env.data[key] = value
}

func (env Env) Get(key string) (Any, error) {
	value, ok := env.data[Symbol(key)]
	if !ok {
		if env.outer != nil {
			return env.outer.Get(key)
		}
		return nil, fmt.Errorf("not found: %s", key)
	}
	return value, nil
}
