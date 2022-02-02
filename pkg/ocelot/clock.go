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
	TTL    time.Duration
}

func (ocelot *Ocelot) GetClock(maxTTL time.Duration) (*Clock, error) {
	request_time := time.Now()
	clock, err := ocelot.client.GetClock()
	response_time := time.Now()
	if err != nil {
		return nil, err
	}
	rtt := response_time.Sub(request_time).Round(time.Millisecond)
	owd := clock.Timestamp.Sub(request_time).Round(time.Millisecond)
	ttl := rtt - owd
	if maxTTL > 0 && ttl > maxTTL || ttl < -maxTTL {
		return nil, fmt.Errorf("TTL of %s exceeded: %s", maxTTL, ttl)
	}
	return &Clock{clock, rtt, owd, ttl}, nil
}
