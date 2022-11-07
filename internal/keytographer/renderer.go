package keytographer

import (
	"fmt"

	"github.com/beevik/etree"
	uuid "github.com/satori/go.uuid"
)

type Renderer interface {
	Render(*Config) ([]byte, error)
}

type renderer struct {
}

func NewRenderer() Renderer {
	return &renderer{}
}

func (r *renderer) Render(c *Config) ([]byte, error) {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	svg := r.svg(doc, c)
	r.styles(svg, c)

	r.keycap(svg, "S", 10, 10)
	r.keycap(svg, "T", 85, 10)
	r.keycap(svg, "R", 160, 10)
	r.keycap(svg, "A", 235, 10)

	r.keycap(svg, "O", 10, 85)
	r.keycap(svg, "I", 85, 85)
	r.keycap(svg, "Y", 160, 85)
	r.keycap(svg, "E", 235, 85)

	doc.Indent(2)
	result, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *renderer) svg(doc *etree.Document, c *Config) *etree.Element {
	svg := doc.CreateElement("svg")
	svg.CreateAttr("xmlns", "http://www.w3.org/2000/svg")
	svg.CreateAttr("xmlns:xlink", "http://www.w3.org/1999/xlink")
	svg.CreateAttr("width", fmt.Sprintf(`%d`, c.Canvas.Width))
	svg.CreateAttr("height", fmt.Sprintf(`%d`, c.Canvas.Height))

	return svg
}

func (r *renderer) styles(svg *etree.Element, c *Config) *etree.Element {
	style := svg.CreateElement("style")

	backgroundColor := c.Canvas.BackgroundColor
	if backgroundColor == "" {
		backgroundColor = "#FFFFFF"
	}

	style.CreateCData(fmt.Sprintf(`
svg {
  background-color: %s;
}`, backgroundColor))

	return style
}

func (r *renderer) keycap(svg *etree.Element, label string, x, y int) *etree.Element {
	g := svg.CreateElement("g")
	g.CreateAttr("id", uuid.NewV4().String())

	outer := g.CreateElement("rect")
	outer.CreateAttr("x", fmt.Sprintf(`%d`, x))
	outer.CreateAttr("y", fmt.Sprintf(`%d`, y))
	outer.CreateAttr("width", "70")
	outer.CreateAttr("height", "70")
	outer.CreateAttr("rx", "3")
	outer.CreateAttr("rx", "3")
	outer.CreateAttr("fill", "#383838")

	inx, iny, inw, inh := x+7, y+6, 56, 56
	inner := g.CreateElement("rect")
	inner.CreateAttr("x", fmt.Sprintf(`%d`, inx))
	inner.CreateAttr("y", fmt.Sprintf(`%d`, iny))
	inner.CreateAttr("width", fmt.Sprintf(`%d`, inw))
	inner.CreateAttr("height", fmt.Sprintf(`%d`, inh))
	inner.CreateAttr("fill", "#fff")
	inner.CreateAttr("fill-opacity", "0.1")

	text := g.CreateElement("text")
	text.CreateAttr("x", fmt.Sprintf(`%d`, inx+(inw/2)))
	text.CreateAttr("y", fmt.Sprintf(`%d`, iny+(inh/2)))
	text.CreateAttr("font-family", "Arial")
	text.CreateAttr("font-size", "16")
	text.CreateAttr("fill", "#e3e3e3")
	text.CreateAttr("text-anchor", "middle")
	text.CreateAttr("dominant-baseline", "middle")
	text.CreateText(label)

	return g
}
