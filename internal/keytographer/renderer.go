package keytographer

import (
	"bytes"
	"fmt"

	svg "github.com/ajstarks/svgo"
	uuid "github.com/satori/go.uuid"
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

	r.keycap(s, keycapOptions{10, 10, "S"})
	r.keycap(s, keycapOptions{85, 10, "T"})
	r.keycap(s, keycapOptions{160, 10, "R"})
	r.keycap(s, keycapOptions{235, 10, "A"})

	r.keycap(s, keycapOptions{10, 85, "O"})
	r.keycap(s, keycapOptions{85, 85, "I"})
	r.keycap(s, keycapOptions{160, 85, "Y"})
	r.keycap(s, keycapOptions{235, 85, "E"})

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

type keycapOptions struct {
	x     int
	y     int
	label string
}

func (r *renderer) keycap(s *svg.SVG, opts keycapOptions) {
	inx, iny, inw, inh := opts.x+7, opts.y+6, 56, 56
	fontSize := 16
	s.Gid(fmt.Sprintf("keycap-%s", uuid.NewV4().String()))
	s.Roundrect(opts.x, opts.y, 70, 70, 3, 3, "fill=\"#383838\"")
	s.Roundrect(inx, iny, inw, inh, 2, 2, "fill=\"#FFFFFF\"", "fill-opacity=\"0.1\"")
	s.Text(
		inx+(inw/2),
		iny+(inh/2)+(fontSize/3),
		opts.label,
		"font-family=\"Arial\"",
		fmt.Sprintf("font-size=\"%d\"", fontSize),
		"fill=\"#e3e3e3\"",
		"text-anchor=\"middle\"",
	)
	s.Gend()
}
