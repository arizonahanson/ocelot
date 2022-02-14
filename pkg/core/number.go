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

func (value Number) Decimal() decimal.Decimal {
	d := decimal.Decimal(value)
	return d
}

func (value Number) String() string {
	return value.Decimal().String()
}

func (value Number) Quot(arg Number, precision Number) Number {
	p := precision.Decimal().IntPart()
	q, _ := value.Decimal().QuoRem(arg.Decimal(), int32(p))
	return Number(q)
}

func (value Number) Rem(arg Number, precision Number) Number {
	p := precision.Decimal().IntPart()
	_, r := value.Decimal().QuoRem(arg.Decimal(), int32(p))
	return Number(r)
}

func (value Number) Round(precision Number) Number {
	p := precision.Decimal().IntPart()
	r := value.Decimal().Round(int32(p))
	return Number(r)
}
