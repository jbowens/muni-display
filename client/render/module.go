package render

import (
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/octavore/naga/service"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/exp/font"
)

type Module struct {
	font *truetype.Font
	app  app.App
	err  error
}

var _ service.Module = &Module{}

func (m *Module) Init(c *service.Config) {
	c.Setup = m.setup
}

func (m *Module) setup() (err error) {
	// Retrieve the default system font, encoded as a TTF.
	ttfBytes := font.Default()

	m.font, err = truetype.Parse(ttfBytes)
	return err
}

func (m *Module) DisplayError(err error) {
	m.err = err
	go func() {
		time.Sleep(5 * time.Second)
		m.err = nil
	}()
}

func (m *Module) InitApp(app app.App) {
	m.app = app

}
