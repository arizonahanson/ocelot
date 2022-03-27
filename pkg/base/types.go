package base

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/starlight/ocelot/pkg/core"
)

// type:func
type Func func(ast core.Expr, env *Env) (core.Any, error)

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

func (future Future) String() string {
	return "?← " // should not happen
}

func (future Future) GoString() string {
	return "?← " // should not happen
}

func (future Future) Equal(any core.Any) bool {
	return false // not comparable
}
