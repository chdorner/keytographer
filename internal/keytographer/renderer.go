package keytographer

import (
	"bytes"
	"fmt"

	svg "github.com/ajstarks/svgo"
)

type Renderer interface {
	Render(*Config) []byte
}

type renderer struct {
}

func NewRenderer() Renderer {
	return &renderer{}
}

func (r *renderer) Render(c *Config) []byte {
	buf := bytes.NewBuffer([]byte{})

	s := svg.New(buf)
	s.Start(c.Canvas.Width, c.Canvas.Height)
	s.Style("text/css", r.styles(c))
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()

	return buf.Bytes()
}

func (r *renderer) styles(c *Config) string {
	backgroundColor := c.Canvas.BackgroundColor
	if backgroundColor == "" {
		backgroundColor = "#FFFFFF"
	}

	styles := fmt.Sprintf(`
svg {
  background-color: %s;
}`, backgroundColor)
	return styles
}
