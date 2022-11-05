package renderer

import (
	"bytes"

	svg "github.com/ajstarks/svgo"
)

type RenderConfig struct{}

type Renderer interface {
	Render() []byte
}

type renderer struct {
	counter int
}

func NewRenderer(config *RenderConfig) Renderer {
	return &renderer{
		counter: 0,
	}
}

func (r *renderer) Render() []byte {
	buf := bytes.NewBuffer([]byte{})

	s := svg.New(buf)
	s.Start(500, 500)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()

	return buf.Bytes()
}
