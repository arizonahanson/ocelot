package base

import (
	"github.com/starlight/ocelot/pkg/core"
)

var Base = map[string]core.Any{
	"nil":   core.Nil{},
	"true":  core.Bool(true),
	"false": core.Bool(false),
}
