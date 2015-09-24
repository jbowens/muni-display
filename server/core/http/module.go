package http

import (
	"fmt"
	"net/http"

	"github.com/octavore/naga/service"
)

const (
	bindAddress = "localhost:80"
)

// Module implements naga/service.Module and encapsulates the MUNI http server.
type Module struct{}

func (m *Module) Init(c *service.Config) {
	c.Start = m.Start
}

func (m *Module) Start() {
	if err := http.ListenAndServe(bindAddress, m); err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %s...\n", bindAddress)
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("sup"))
}
