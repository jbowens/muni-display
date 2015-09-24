package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/octavore/naga/service"
)

// Module implements naga/service.Module and encapsulates logic surrounding
// loading configuration files.
type Module struct{}

// Init implements the service.Module interface.
func (m *Module) Init(c *service.Config) {}

// Load loads the given config file into the provided struct.
func (m *Module) Load(filename string, dst interface{}) error {
	b, err := ioutil.ReadFile("./config/" + filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dst)
}
