package main

import (
	"fmt"
	"time"

	"github.com/jbowens/muni/client/network"
	"github.com/jbowens/muni/client/render"
	"github.com/octavore/naga/service"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type Module struct {
	Network *network.Module
	Render  *render.Module
	loaded  bool
}

func (m *Module) Init(c *service.Config) {
	c.Setup = m.setup
}

func (m *Module) setup() error {
	// Register this module's main function as the main function of the
	// application.
	app.Main(m.main)
	return nil
}

func (m *Module) draw(glctx gl.Context, sz size.Event, images *glutil.Images) {
	serverResponse := m.Network.Response()

	glctx.ClearColor(1, 1, 1, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	var display render.Display
	display.Loaded = m.loaded
	if len(serverResponse.Predictions) > 0 {
		display.NextOK = true
		display.NextTrainMinutes = serverResponse.Predictions[0].Minutes
		display.TransitRouteName = fmt.Sprintf("%s (%s)", serverResponse.Stop.Route, serverResponse.Stop.Direction)
		if len(serverResponse.Predictions) > 1 {
			display.NextNextOK = true
			display.NextNextTrainMinutes = serverResponse.Predictions[1].Minutes
		}
		display.UpdatedSecondsAgo = int(time.Now().Sub(serverResponse.LastRefresh).Seconds())
	}
	m.Render.Display(display, sz, glctx, images)
}

// main is the entry point of the application. This is function gets registered
// as the main function of the application.
func (m *Module) main(a app.App) {
	var images *glutil.Images
	var glctx gl.Context
	sz := size.Event{}

	m.Render.InitApp(a)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			a.Send(paint.Event{})

		case <-m.Network.Updated():
			m.loaded = true

		case e := <-a.Events():
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				glctx, _ = e.DrawContext.(gl.Context)
				if glctx != nil {
					glctx = e.DrawContext.(gl.Context)
					images = glutil.NewImages(glctx)
				}
			case size.Event:
				sz = e
			case paint.Event:
				m.draw(glctx, sz, images)
				a.Publish()
			}
		}
	}
}
