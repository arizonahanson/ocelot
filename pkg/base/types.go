package base

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/starlight/ocelot/pkg/core"
)

// type:func
type Func func(ast core.List, env *Env) (core.Any, error)

func (fn Func) String() string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(name, ".")
	short := parts[len(parts)-1]
	return "&" + short
}

func (fn Func) GoString() string {
	return fn.String()
}

func (fn Func) Equal(any core.Any) bool {
	return false // not comparable
}

// type:future
type Future func() (core.Any, error)

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

// get future value async
func (future Future) Async() Future {
	// await future
	channel := make(chan Future, 1)
	await := func() (core.Any, error) {
		return <-channel, nil
	}
	// resolve future
	async := func() {
		val, err := future.Get()
		channel <- func() (core.Any, error) {
			return val, err
		}
	}
	go async()
	return Future(await)
}

func (future Future) String() string {
	return "?← " // should not happen
}

func (future Future) GoString() string {
	return "?← " // should not happen
}

func (future Future) Equal(any core.Any) bool {
	return false // not comparable
}
