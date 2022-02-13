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

func numberFromInt(value int64) Number {
	d := decimal.NewFromInt(value)
	n := Number(d)
	return n
}

func (value Number) toDecimal() decimal.Decimal {
	d := decimal.Decimal(value)
	return d
}

func (value Number) String() string {
	return value.toDecimal().String()
}

func (value Number) Add(arg Number) Number {
	return Number(value.toDecimal().Add(arg.toDecimal()))
}

func (value Number) Sub(arg Number) Number {
	return Number(value.toDecimal().Sub(arg.toDecimal()))
}

func (value Number) Mul(arg Number) Number {
	return Number(value.toDecimal().Mul(arg.toDecimal()))
}

func (value Number) Quot(arg Number, precision Number) Number {
	p := precision.toDecimal().IntPart()
	q, _ := value.toDecimal().QuoRem(arg.toDecimal(), int32(p))
	return Number(q)
}

func (value Number) Quot2(arg Number) Number {
	q := value.toDecimal().Div(arg.toDecimal())
	return Number(q)
}

func (value Number) Rem(arg Number, precision Number) Number {
	p := precision.toDecimal().IntPart()
	_, r := value.toDecimal().QuoRem(arg.toDecimal(), int32(p))
	return Number(r)
}

func (value Number) Round(precision Number) Number {
	p := precision.toDecimal().IntPart()
	r := value.toDecimal().Round(int32(p))
	return Number(r)
}

var Zero Number = Number(decimal.Zero)
