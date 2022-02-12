package ocelot

import (
	"github.com/starlight/ocelot/internal/parser"
	"github.com/starlight/ocelot/pkg/core"
)

func Eval(in string) (core.Any, error) {
	return parser.Parse("Eval", []byte(in))
}
