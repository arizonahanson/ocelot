package api

import (
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

type AlpacaClient struct {
	client alpaca.Client
}

func GetAlpacaClient(baseUrl string) ApiClient {
	client := alpaca.NewClient(alpaca.ClientOpts{
		BaseURL: baseUrl,
	})
	return ApiClient(&AlpacaClient{client})
}

func (alpaca *AlpacaClient) GetClock() (*Clock, error) {
	clock, err := alpaca.client.GetClock()
	if err != nil {
		return nil, err
	}
	return &Clock{clock.Timestamp, clock.IsOpen, clock.NextOpen, clock.NextClose}, nil
}

func (alpaca *AlpacaClient) GetPositions() ([]Position, error) {
	positions, err := alpaca.client.ListPositions()
	if err != nil {
		return nil, err
	}
	result := []Position{}
	for _, position := range positions {
		p := Position{
			position.Symbol,
			position.Qty,
			*position.MarketValue,
		}
		result = append(result, p)
	}
	return result, nil
}
