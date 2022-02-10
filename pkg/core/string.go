package core

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type String string

func (value String) ToNumber() (*Number, error) {
	d, err := decimal.NewFromString(string(value))
	if err != nil {
		return nil, err
	}
	n := Number(d)
	return &n, nil
}

func (value String) Unquote() (*String, error) {
	unquoted, err := strconv.Unquote(string(value))
	if err != nil {
		return nil, err
	}
	s := String(unquoted)
	return &s, nil
}
