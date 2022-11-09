package keytographer

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/google/uuid"
)

type Renderer interface {
	RenderAllLayers(*Config) ([]*RenderedLayer, error)
	RenderLayer(*Config, string) (*RenderedLayer, error)
}

type RenderedLayer struct {
	Name string
	Svg  []byte
}

type renderer struct {
}

func NewRenderer() Renderer {
	return &renderer{}
}

func (r *renderer) RenderAllLayers(c *Config) ([]*RenderedLayer, error) {
	var layers []*RenderedLayer

	for _, layerConfig := range c.Layers {
		layer, err := r.RenderLayer(c, layerConfig.Name)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}

	return layers, nil
}

func (r *renderer) RenderLayer(c *Config, layerName string) (*RenderedLayer, error) {
	layer := c.GetLayer(layerName)
	if layer == nil {
		return nil, fmt.Errorf("layer with name %s does not exist", layerName)
	}

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	svg := r.svg(doc, c)
	r.styles(svg, c)

	for idx, key := range c.Layout.Keys {
		var label string
		if layer.Keys != nil {
			if len(layer.Keys) > idx {
				label = layer.Keys[idx].Label
			}
		}

		width := key.W * 70
		height := key.H * 70
		r.keycap(svg, label, int((key.X*75)+10), int((key.Y*75)+10), width, height)
	}

	doc.Indent(2)
	result, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	return &RenderedLayer{
		Name: layer.Name,
		Svg:  result,
	}, nil
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

func (r *renderer) keycap(svg *etree.Element, label string, x, y int, w, h float64) *etree.Element {
	g := svg.CreateElement("g")
	gid, _ := uuid.NewRandom()
	g.CreateAttr("id", gid.String())

	outer := g.CreateElement("rect")
	outer.CreateAttr("x", fmt.Sprintf(`%d`, x))
	outer.CreateAttr("y", fmt.Sprintf(`%d`, y))
	outer.CreateAttr("width", fmt.Sprintf(`%f`, w))
	outer.CreateAttr("height", fmt.Sprintf(`%f`, h))
	outer.CreateAttr("rx", "3")
	outer.CreateAttr("rx", "3")
	outer.CreateAttr("fill", "#383838")

	inx, iny, inw, inh := x+7, y+6, w-14, h-14
	inner := g.CreateElement("rect")
	inner.CreateAttr("x", fmt.Sprintf(`%d`, inx))
	inner.CreateAttr("y", fmt.Sprintf(`%d`, iny))
	inner.CreateAttr("width", fmt.Sprintf(`%f`, inw))
	inner.CreateAttr("height", fmt.Sprintf(`%f`, inh))
	inner.CreateAttr("fill", "#fff")
	inner.CreateAttr("fill-opacity", "0.1")

	text := g.CreateElement("text")
	text.CreateAttr("x", fmt.Sprintf(`%d`, inx+(int(inw/2))))
	text.CreateAttr("y", fmt.Sprintf(`%d`, iny+(int(inh/2))))
	text.CreateAttr("font-family", "Arial")
	text.CreateAttr("font-size", "16")
	text.CreateAttr("fill", "#e3e3e3")
	text.CreateAttr("text-anchor", "middle")
	text.CreateAttr("dominant-baseline", "middle")
	text.CreateText(label)

	return g
}
