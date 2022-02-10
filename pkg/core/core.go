package core

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type Any interface{}

type Bool bool

type Nil struct{}

type Number decimal.Decimal

type String string

type Symbol string

type List []Any

func ParseNumber(value string) (*Number, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return nil, err
	}
	n := Number(d)
	return &n, nil
}

func ParseString(value string) (*String, error) {
	unquoted, err := strconv.Unquote(value)
	if err != nil {
		return nil, err
	}
	s := String(unquoted)
	return &s, nil
}
