package core

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

type Any interface{}

type Nil struct{}

type List []Any

type Vector []Any

type Symbol string

type String string

type Bool bool

type Number decimal.Decimal

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

func (val Number) String() string {
	return decimal.Decimal(val).String()
}

func (val Nil) String() string {
	return "nil"
}
