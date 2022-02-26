package builtin

import (
	"fmt"

	"github.com/starlight/ocelot/pkg/core"
)

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
