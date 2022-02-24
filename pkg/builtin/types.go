package builtin

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

// type:nil
type Nil struct{}

func (val Nil) String() string {
	return "nil"
}

func (val Nil) GoString() string {
	return val.String()
}

func (val Nil) Equal(any core.Any) bool {
	switch any.(type) {
	default:
		return false
	case Nil:
		return true
	}
}

// type:bool
type Bool bool

func (val Bool) String() string {
	return fmt.Sprintf("%v", bool(val))
}

func (val Bool) GoString() string {
	return val.String()
}

func (val Bool) Equal(any core.Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Bool:
		return val == arg
	}
}
