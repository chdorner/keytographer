package keymap

import (
	"bytes"

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
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()

	return buf.Bytes()
}
