package http

import (
	"fmt"
	"net/http"

	"github.com/jbowens/muni/server/core/predictions"
	"github.com/octavore/naga/service"
)

const (
	bindAddress = "localhost:8080"
)

// Module implements naga/service.Module and encapsulates the MUNI http server.
type Module struct {
	Predictions *predictions.Module
	mux         *http.ServeMux
}

func (m *Module) Init(c *service.Config) {
	c.Start = m.start
	c.Setup = m.setup
}

func (m *Module) setup() error {
	m.mux = http.NewServeMux()
	m.mux.HandleFunc("/predictions/", m.handlePredictions)
	return nil
}

func (m *Module) start() {
	if err := http.ListenAndServe(bindAddress, m); err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %s...\n", bindAddress)
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(rw, req)
}
