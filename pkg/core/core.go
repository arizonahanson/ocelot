package core

import (
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// parse quoted, escaped string
func (val String) Unquote() (String, error) {
	// Go Unquote doesn't support solidus-escape: "\/"
	str := strings.Replace(val.Val, "\\/", "/", -1)
	str, err := strconv.Unquote(str)
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

func NewSymbol(sym string, pos *Position) Symbol {
	return Symbol{Val: sym, Pos: pos}
}
