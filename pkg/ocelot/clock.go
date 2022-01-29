package ocelot

import (
	"fmt"
	"time"

	"github.com/starlight/ocelot/internal/api"
)

type Clock struct {
	Market *api.Clock
	RTT    time.Duration
	OWD    time.Duration
	LAG    time.Duration
}

func (ocelot *Ocelot) GetClock() (*Clock, error) {
	request_time := time.Now()
	clock, err := ocelot.client.GetClock()
	response_time := time.Now()
	if err != nil {
		return nil, err
	}
	rtt := response_time.Sub(request_time).Round(time.Millisecond)
	owd := clock.Timestamp.Sub(request_time).Round(time.Millisecond)
	if owd > time.Second || owd < -time.Second {
		return nil, fmt.Errorf("DELAY 1s exceeded: %s", owd)
	}
	lag := rtt - owd
	clock.Timestamp = clock.Timestamp.Round(time.Millisecond)
	return &Clock{clock, rtt, owd, lag}, nil
}
