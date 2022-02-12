package core

import "github.com/shopspring/decimal"

type Number decimal.Decimal

func numberFromString(value String) (Number, error) {
	d, err := decimal.NewFromString(string(value))
	if err != nil {
		return Number(decimal.Zero), err
	}
	n := Number(d)
	return n, nil
}

func (value Number) toDecimal() decimal.Decimal {
	d := decimal.Decimal(value)
	return d
}

func (value Number) String() string {
	return value.toDecimal().String()
}
