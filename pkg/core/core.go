package core

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// parse quoted, escaped string
func (val String) Unquote() (String, error) {
	str, err := strconv.Unquote(val.Val)
	return String{str}, err
}

// parse number
func (val String) Number() (Number, error) {
	dec, err := decimal.NewFromString(val.Val)
	return Number(dec), err
}

// useful numbers
var Zero = NewNumber(0)
var One = NewNumber(1)

// number from int
func NewNumber(num int) Number {
	return Number(decimal.NewFromInt(int64(num)))
}

// number to decimal impl
func (val Number) Decimal() decimal.Decimal {
	return decimal.Decimal(val)
}
