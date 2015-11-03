package core

import (
	"fmt"

	"github.com/jbowens/muni-display/server/core/http"
	"github.com/octavore/naga/service"
)

// Module implements naga/service.Module and encapsulates the entire muni
// application server
type Module struct {
	HTTP *http.Module
}

func (m *Module) Init(c *service.Config) {
	c.Start = m.start
}

func (m *Module) start() {
	fmt.Println("Starting app...")
}
