package ocelot

import (
	"github.com/starlight/ocelot/internal/api"
)

type Ocelot struct {
	client api.ApiClient
}

func GetOcelot() *Ocelot {
	// uses Alpaca; but behind an interface
	alpaca := api.GetAlpacaClient("https://paper-api.alpaca.markets")
	client := api.ApiClient(alpaca)
	return &Ocelot{client}
}
