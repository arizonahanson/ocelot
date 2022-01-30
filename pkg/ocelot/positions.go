package ocelot

import "github.com/starlight/ocelot/internal/api"

func (ocelot *Ocelot) GetPositions() ([]api.Position, error) {
	return ocelot.client.GetPositions()
}
