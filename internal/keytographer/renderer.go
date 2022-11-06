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
	s.Def()
	r.keycap(s)
	s.DefEnd()
	s.Style("text/css", r.styles(c))

	s.Use(10, 10, "#keycap")
	s.Use(90, 10, "#keycap")

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

func (r *renderer) keycap(s *svg.SVG) {
	s.Gid("keycap")
	s.Roundrect(0, 0, 70, 70, 3, 3, "fill=\"#383838\"")
	s.Roundrect(7, 6, 56, 56, 3, 3, "fill=\"#FFFFFF\"", "fill-opacity=\"0.1\"")
	s.Gend()
}
