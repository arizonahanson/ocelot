package api

import (
	"time"

	"github.com/shopspring/decimal"
)

type ApiClient interface {
	GetClock() (*Clock, error)
	GetPositions() ([]Position, error)
}

type Clock struct {
	Timestamp time.Time
	IsOpen    bool
	NextOpen  time.Time
	NextClose time.Time
}

type Position struct {
	Symbol   string
	Quantity decimal.Decimal
	Value    decimal.Decimal
}
