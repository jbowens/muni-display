package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	server "github.com/jbowens/muni/server/core/http"
	"github.com/octavore/naga/service"
)

const (
	stopKey         = "home"
	pollingInterval = 5 * time.Second
	serverURL       = "http://djroomba.com:8000/predictions/home"
)

type Module struct {
	mu             sync.Mutex
	ticker         *time.Ticker
	serverResponse server.HandlePredictionsResponse
}

func (m *Module) Init(c *service.Config) {
	c.Start = m.start
}

func (m *Module) Response() (resp server.HandlePredictionsResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp = m.serverResponse
	return resp
}

func (m *Module) start() {
	m.ticker = time.NewTicker(pollingInterval)
	go m.poll()
}

func (m *Module) poll() {
	for _ = range m.ticker.C {
		if err := m.update(); err != nil {
			// TODO(jackson): Fix this to display the error on the screen.
			panic(err)
		}
	}
}

func (m *Module) update() error {
	resp, err := http.DefaultClient.Get(serverURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if err := json.Unmarshal(b, &m.serverResponse); err != nil {
		return err
	}
	return nil
}
