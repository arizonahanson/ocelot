package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type Any interface{}

type Number decimal.Decimal

type String string

type Symbol struct {
	Val string
	Pos *Position
}

type Position struct {
	Line, Col, Offset int
}

type Key string

type Map map[Key]Any

type Vector []Any

type List []Any

var Zero = NewNumber(0)
var One = NewNumber(1)

func NewNumber(num int) Number {
	return Number(decimal.NewFromInt(int64(num)))
}

func (val Number) String() string {
	return val.Decimal().String()
}

func (val Number) GoString() string {
	return val.String()
}

func (val Number) Decimal() decimal.Decimal {
	return decimal.Decimal(val)
}

// parse quoted, escaped string
func (val String) Unquote() (String, error) {
	str, err := strconv.Unquote(string(val))
	return String(str), err
}

// parse number
func (val String) Number() (Number, error) {
	dec, err := decimal.NewFromString(string(val))
	return Number(dec), err
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

func (key Key) GoString() string {
	return fmt.Sprintf("%s", key)
}

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

func (val List) GoString() string {
	str := "("
	for i, item := range val {
		if i != 0 {
			str += " "
		}
		if i >= 2 {
			str += "..."
			break
		}
		str += fmt.Sprintf("%#v", item)
	}
	return str + ")"
}

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
		if i >= 2 {
			str += " ..."
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
