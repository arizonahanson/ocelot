package core

import (
	"github.com/shopspring/decimal"
)

type Any interface{}

type Bool bool

type Nil struct{}

type Number decimal.Decimal

type Symbol string

type List []Any
