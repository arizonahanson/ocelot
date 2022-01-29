package api

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

type AlpacaClient struct {
	client alpaca.Client
}

func GetAlpacaClient(baseUrl string) *AlpacaClient {
	client := alpaca.NewClient(alpaca.ClientOpts{
		BaseURL: baseUrl,
	})
	return &AlpacaClient{client}
}

func (alpaca *AlpacaClient) GetClock() (*Clock, error) {
	clock, err := alpaca.client.GetClock()
	if err != nil {
		return nil, err
	}
	return &Clock{clock.Timestamp, clock.IsOpen, clock.NextOpen, clock.NextClose}, nil
}
