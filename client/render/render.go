package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

const (
	dpi                      = 72
	minsFontSize             = 36
	nextTrainFontSize        = 350
	nextNextTrainFontSize    = 144
	lastUpdatedAtFontSize    = 18
	informationPopupFontSize = 100
)

var (
	foreground          = image.White
	secondaryForeground = image.NewUniform(color.RGBA{0xA6, 0xE3, 0xFA, 0xFF})
	background          = image.NewUniform(color.RGBA{0x35, 0x67, 0x99, 0xFF})
	errorBackground     = image.NewUniform(color.RGBA{0x8C, 0x35, 0x1F, 0xFF})
	loadingBackground   = image.NewUniform(color.RGBA{0x3B, 0x3B, 0x3B, 0xFF})
)

// Display encapsulates all of the data that is displayed on the screen.
type Display struct {
	Loaded               bool
	NextOK               bool
	NextNextOK           bool
	NextTrainMinutes     int
	NextNextTrainMinutes int
	UpdatedSecondsAgo    int
}

func (m *Module) Display(display Display, sz size.Event, glctx gl.Context, images *glutil.Images) {
	im := images.NewImage(sz.WidthPx, sz.HeightPx)

	if m.err != nil {
		// There was an error of some kind. We should display the error indication until the error
		// is removed from the module by a timeout.
		m.renderInformation(im.RGBA, sz.Size(), errorBackground, "T_T", "T_T")
	} else if display.Loaded {
		m.render(im.RGBA, display, sz.Size())
	} else {
		loadingStr := "Loading" + strings.Repeat(".", int(time.Now().Unix()%4))
		m.renderInformation(im.RGBA, sz.Size(), loadingBackground, loadingStr, "Loading...")
	}

	im.Upload()
	im.Draw(
		sz,
		geom.Point{},
		geom.Point{X: sz.WidthPt},
		geom.Point{Y: sz.HeightPt},
		sz.Bounds())
	im.Release()
}

func (m *Module) render(rgba *image.RGBA, display Display, dimensions image.Point) {
	// Prepare a blue background to draw on.
	draw.Draw(rgba, rgba.Bounds(), background, image.ZP, draw.Src)

	// First, render the very next train's minutes on the left half of the screen.
	d := &font.Drawer{
		Dst: rgba,
		Src: foreground,
		Face: truetype.NewFace(m.font, &truetype.Options{
			Size:    nextTrainFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	textWidth := d.MeasureString(strconv.Itoa(display.NextTrainMinutes))
	d.Dot = fixed.Point26_6{
		X: fixed.I(dimensions.X/4) - (textWidth / 2),
		Y: fixed.I(int(4 * dimensions.Y / 5)),
	}
	d.DrawString(strconv.Itoa(display.NextTrainMinutes))

	// Render the little "min" label next to the next train time.
	m.renderMin(rgba, fixed.Point26_6{
		X: d.Dot.X,
		Y: d.Dot.Y,
	})

	// Now, render the next next train's minutes on the right half of the screen.
	d = &font.Drawer{
		Dst: rgba,
		Src: foreground,
		Face: truetype.NewFace(m.font, &truetype.Options{
			Size:    nextNextTrainFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	textWidth = d.MeasureString(strconv.Itoa(display.NextNextTrainMinutes))
	d.Dot = fixed.Point26_6{
		X: fixed.I(3*dimensions.X/4) - (textWidth / 2),
		Y: fixed.I(int(2 * dimensions.Y / 3)),
	}
	d.DrawString(strconv.Itoa(display.NextNextTrainMinutes))

	// Render the little "min" label next to the next, next train time.
	m.renderMin(rgba, fixed.Point26_6{
		X: d.Dot.X,
		Y: d.Dot.Y,
	})

	// Render the text indicating the freshness of the presented data.
	d = &font.Drawer{
		Dst: rgba,
		Src: secondaryForeground,
		Face: truetype.NewFace(m.font, &truetype.Options{
			Size:    lastUpdatedAtFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	updatedAt := fmt.Sprintf("Predictions from %v seconds ago. ", display.UpdatedSecondsAgo)
	textWidth = d.MeasureString(updatedAt)
	d.Dot = fixed.Point26_6{
		X: fixed.I(dimensions.X) - textWidth,
		Y: fixed.I(dimensions.Y - lastUpdatedAtFontSize),
	}
	d.DrawString(updatedAt)
}

func (m *Module) renderMin(rgba *image.RGBA, position fixed.Point26_6) {
	d := &font.Drawer{
		Dst: rgba,
		Src: secondaryForeground,
		Face: truetype.NewFace(m.font, &truetype.Options{
			Size:    minsFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	d.Dot = position
	d.DrawString("min")
}

// renderLoadingScreen will render the Loading screen.
func (m *Module) renderInformation(rgba *image.RGBA, dimensions image.Point, background image.Image, text string, textSizing string) {
	// Prepare a dark grey background to draw on.
	draw.Draw(rgba, rgba.Bounds(), background, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst: rgba,
		Src: foreground,
		Face: truetype.NewFace(m.font, &truetype.Options{
			Size:    informationPopupFontSize,
			DPI:     dpi,
			Hinting: font.HintingNone,
		}),
	}
	dy := int(math.Ceil(informationPopupFontSize * dpi / 72))
	textWidth := d.MeasureString(textSizing)
	d.Dot = fixed.Point26_6{
		X: fixed.I(dimensions.X/2) - (textWidth / 2),
		Y: fixed.I(dimensions.Y/2 + dy/2),
	}
	d.DrawString(text)
}
