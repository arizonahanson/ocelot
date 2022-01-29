package api

import "time"

type ApiClient interface {
	GetClock() (*Clock, error)
}

type Clock struct {
	Timestamp time.Time
	IsOpen    bool
	NextOpen  time.Time
	NextClose time.Time
}
