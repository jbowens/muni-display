package core

import (
	"fmt"

	"github.com/jbowens/muni/server/core/predictions"
	"github.com/octavore/naga/service"
)

// Module implements naga/service.Module and encapsulates the entire muni
// application server
type Module struct {
	Predictions *predictions.Module
}

func (m *Module) Init(c *service.Config) {
	c.Start = m.Start
}

func (m *Module) Start() {
	fmt.Println("Starting app...")
}
