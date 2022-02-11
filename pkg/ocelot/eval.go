package ocelot

import "github.com/starlight/ocelot/pkg/core"

func Eval(in string) core.Any {
	return core.String(in)
}
