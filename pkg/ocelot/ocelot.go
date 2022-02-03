package ocelot

import (
	"github.com/starlight/ocelot/internal/api"
)

type Ocelot struct {
	client api.ApiClient
}

func GetOcelot(alpacaUrl string) *Ocelot {
	// uses Alpaca; but behind an interface
	client := api.GetAlpacaClient(alpacaUrl)
	return &Ocelot{client}
}
