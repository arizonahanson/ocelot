package base

import "github.com/starlight/ocelot/pkg/core"

// lazy function call
func (fn Func) Future(ast core.List, env *Env) Future {
	return func() (core.Any, error) {
		return fn(ast, env)
	}
}

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

// resolve future synchronously and return new future
func (future Future) Sync() Future {
	val, err := future.Get()
	return func() (core.Any, error) {
		return val, err
	}
}

// resolve future asynchronously and return new future
func (future Future) Async() Future {
	tunnel := make(chan Future, 1)
	// resolve
	send := func() {
		tunnel <- future.Sync()
	}
	go send()
	// await future
	recv := func() (core.Any, error) {
		return <-tunnel, nil
	}
	return recv
}
