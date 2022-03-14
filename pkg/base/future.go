package base

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

// trampoline to resolve future values
func (future Future) Get() (val core.Any, err error) {
	val, err = future()
	for {
		if err != nil {
			return
		}
		switch future := val.(type) {
		default:
			return
		case Future:
			val, err = future()
		}
	}
}

// resolve future asynchronously and return new future
func (future Future) Async() Future {
	tunnel := make(chan Future, 1)
	// resolve
	send := func() {
		val, err := future.Get()
		tunnel <- func() (core.Any, error) {
			return val, err
		}
	}
	go send()
	// await future
	recv := func() (core.Any, error) {
		return <-tunnel, nil
	}
	return recv
}

// trace errors mapped to source ast
func (future Future) Trace(ast core.List) Future {
	return func() (val core.Any, err error) {
		val, err = future.Get()
		if err != nil {
			err = fmt.Errorf("%#v: %s", ast[0], err)
		}
		return
	}
}

// lazy function call
func (fn Func) Future(ast core.List, env *Env) Future {
	future := func() (core.Any, error) {
		return fn(ast, env)
	}
	return Future(future).Trace(ast)
}
