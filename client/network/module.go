package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/jbowens/muni-display/client/render"
	server "github.com/jbowens/muni-display/server/core/http"
	"github.com/octavore/naga/service"
)

const (
	stopKey         = "home"
	pollingInterval = 5 * time.Second
	serverURL       = "http://djroomba.com:8000/predictions/home"
)

type Module struct {
	Render *render.Module

	mu             sync.Mutex
	ticker         *time.Ticker
	serverResponse server.HandlePredictionsResponse
	updated        chan struct{}
}

func (m *Module) Init(c *service.Config) {
	c.Setup = m.setup
}

func (m *Module) Response() (resp server.HandlePredictionsResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp = m.serverResponse
	return resp
}

func (m *Module) Updated() <-chan struct{} {
	return m.updated
}

func (m *Module) Update() {
	m.update()
}

func (m *Module) setup() error {
	m.updated = make(chan struct{}, 1)
	m.ticker = time.NewTicker(pollingInterval)
	go m.poll()
	return nil
}

func (m *Module) poll() {
	for _ = range m.ticker.C {
		m.update()
	}
}

func (m *Module) update() {
	resp, err := http.DefaultClient.Get(serverURL)
	if err != nil {
		m.err(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.err(err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if err := json.Unmarshal(b, &m.serverResponse); err != nil {
		m.err(err)
		return
	}
	m.updated <- struct{}{}
}

func (m *Module) err(err error) {
	m.Render.DisplayError(err)
}
