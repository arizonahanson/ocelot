package core

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type Any interface{}

type Nil struct{}

type List []Any

type Vector []Any

type Symbol struct {
	Val string
	Pos *Position
}

type Position struct {
	Line, Col, Offset int
}

type String string

type Key string

type Map map[Key]Any

type Bool bool

type Number decimal.Decimal

type Function func(ast List, env Env) (Any, error)

var Zero = Number(decimal.Zero)
var One = Number(decimal.NewFromInt(1))

// parse quoted, escaped string
func (val String) Unquote() (String, error) {
	str, err := strconv.Unquote(string(val))
	if err != nil {
		return String(""), err
	}
	return String(str), nil
}

// parse number
func (val String) Number() (Number, error) {
	dec, err := decimal.NewFromString(string(val))
	if err != nil {
		return Number(decimal.Zero), err
	}
	return Number(dec), nil
}

func (val Symbol) String() string {
	return fmt.Sprintf("ocelot:%d:%d (%d): '%s'", val.Pos.Line, val.Pos.Col, val.Pos.Offset, val.Val)
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

func (val Map) String() string {
	res := []string{}
	for key, value := range val {
		res = append(res, fmt.Sprintf("%v", key))
		switch value.(type) {
		default:
			res = append(res, fmt.Sprintf("%v", value))
			break
		case List:
			res = append(res, fmt.Sprintf("(%v)", value))
			break
		}
	}
	return "{" + strings.Join(res, " ") + "}"
}

func (val Number) String() string {
	return val.Decimal().String()
}

func (val Number) Decimal() decimal.Decimal {
	return decimal.Decimal(val)
}

func (val Nil) String() string {
	return "nil"
}

func (fn Function) String() string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	str := strs[len(strs)-1]
	str = strings.ReplaceAll(str, "_", "")
	str = strings.ReplaceAll(str, "E", "!")
	str = strings.ReplaceAll(str, "Q", "?")
	str = strings.ReplaceAll(str, "S", "*")
	return "&" + str
}
