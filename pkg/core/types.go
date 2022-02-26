package core

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// interface:any
type Any interface {
	fmt.Stringer
	fmt.GoStringer
	Equal(any Any) bool
}

// type:nil
type Nil struct{}

func (val Nil) String() string {
	return "nil"
}

func (val Nil) GoString() string {
	return val.String()
}

func (val Nil) Equal(any Any) bool {
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

func (val Bool) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Bool:
		return val == arg
	}
}

// type:number
type Number decimal.Decimal

func (val Number) String() string {
	return val.Decimal().String()
}

func (val Number) GoString() string {
	return val.String()
}

func (val Number) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Number:
		return val.Decimal().Equal(arg.Decimal())
	}
}

// type:string
type String struct {
	Val string
}

func (val String) String() string {
	return fmt.Sprintf("%s", val.Val)
}

func (val String) GoString() string {
	return fmt.Sprintf("%#v", val.Val)
}

func (val String) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case String:
		return val == arg
	}
}

// type:symbol
type Symbol struct {
	Val string
	Pos *Position
}

type Position struct {
	Line, Col, Offset int
}

func (val Symbol) String() string {
	return fmt.Sprintf("%s", val.Val)
}

func (val Symbol) GoString() string {
	if val.Pos != nil {
		return fmt.Sprintf("%s<%d,%d;%d>", val.Val, val.Pos.Line, val.Pos.Col, val.Pos.Offset)
	}
	return fmt.Sprintf("%s<?>", val.Val)
}

func (val Symbol) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Symbol:
		return val.Val == arg.Val
	}
}

// type:key
type Key string

func (key Key) String() string {
	return fmt.Sprintf("%s", string(key))
}

func (key Key) GoString() string {
	return key.String()
}

func (val Key) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Key:
		return val == arg
	}
}

// type:list
type List []Any

func (val List) String() string {
	str := ""
	for i, item := range val {
		if i != 0 {
			str += " "
		}
		switch item.(type) {
		default:
			str += fmt.Sprintf("%v", item)
			break
		case List:
			str += fmt.Sprintf("(%v)", item)
			break
		}
	}
	return str
}

func (val List) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case List:
		if len(val) != len(arg) {
			return false
		}
		for i, a := range val {
			b := arg[i]
			if !a.Equal(b) {
				return false
			}
		}
		return true
	}
}

func (val List) GoString() string {
	str := "("
	for i, item := range val {
		if i != 0 {
			str += " "
		}
		if i == 1 && len(val) > 2 {
			str += "..."
			break
		}
		str += fmt.Sprintf("%#v", item)
	}
	return str + ")"
}

// type:vector
type Vector []Any

func (val Vector) String() string {
	str := ""
	for i, item := range val {
		if i != 0 {
			str += " "
		}
		switch item.(type) {
		default:
			str += fmt.Sprintf("%v", item)
			break
		case List:
			str += fmt.Sprintf("(%v)", item)
			break
		}
	}
	return "[" + str + "]"
}

func (val Vector) GoString() string {
	str := "["
	for i, item := range val {
		if i != 0 {
			str += " "
		}
		if i == 1 && len(val) > 2 {
			str += "..."
			break
		}
		switch item.(type) {
		default:
			str += fmt.Sprintf("%#v", item)
			break
		case List:
			str += fmt.Sprintf("(%#v)", item)
			break
		}
	}
	return str + "]"
}

func (val Vector) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Vector:
		if len(val) != len(arg) {
			return false
		}
		for i, a := range val {
			b := arg[i]
			if !a.Equal(b) {
				return false
			}
		}
		return true
	}
}

// type:map
type Map map[Key]Any

func (val Map) String() string {
	res := make([]string, len(val)*2)
	i := 0
	for key, value := range val {
		res[i] = fmt.Sprintf("%v", key)
		switch value.(type) {
		default:
			res[i+1] = fmt.Sprintf("%v", value)
			break
		case List:
			res[i+1] = fmt.Sprintf("(%v)", value)
			break
		}
		i += 2
	}
	return "{" + strings.Join(res, " ") + "}"
}

func (val Map) GoString() string {
	return val.String()
}

func (val Map) Equal(any Any) bool {
	switch arg := any.(type) {
	default:
		return false
	case Map:
		if len(val) != len(arg) {
			return false
		}
		for key, item := range val {
			item2, ok := arg[key]
			if !ok || !item.Equal(item2) {
				return false
			}
		}
		return true
	}
}
