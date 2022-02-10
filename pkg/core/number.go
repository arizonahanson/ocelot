package core

import "github.com/shopspring/decimal"

type Number decimal.Decimal

func numberFromString(value String) (*Number, error) {
	d, err := decimal.NewFromString(string(value))
	if err != nil {
		return nil, err
	}
	n := Number(d)
	return &n, nil
}

func (value Number) toDecimal() *decimal.Decimal {
	d := decimal.Decimal(value)
	return &d
}

func (value Number) ToString() *String {
	s := String(value.toDecimal().String())
	return &s
}
