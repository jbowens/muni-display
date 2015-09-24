package http

import (
	"fmt"
	"net/http"

	"github.com/jbowens/muni/server/core/config"
	"github.com/jbowens/muni/server/core/predictions"
	"github.com/octavore/naga/service"
)

// Module implements naga/service.Module and encapsulates the MUNI http server.
type Module struct {
	Config      *config.Module
	Predictions *predictions.Module
	config      httpConfig
	mux         *http.ServeMux
}

type httpConfig struct {
	BindAddress string `json:"bind_address"`
}

func (m *Module) Init(c *service.Config) {
	c.Start = m.start
	c.Setup = m.setup
}

func (m *Module) setup() error {
	if err := m.Config.Load("config.json", &m.config); err != nil {
		return err
	}

	m.mux = http.NewServeMux()
	m.mux.HandleFunc("/predictions/", m.handlePredictions)
	return nil
}

func (m *Module) start() {
	if err := http.ListenAndServe(m.config.BindAddress, m); err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %s...\n", m.config.BindAddress)
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(rw, req)
}
