package keytographer

type Config struct {
	Name     string
	Keyboard string `yaml:"keyboard,omitempty"`
	Canvas   CanvasConfig
	Layers   []Layer `yaml:"layers,omitempty"`
	Layout   LayoutConfig
}

type CanvasConfig struct {
	Width           int    `yaml:"width"`
	Height          int    `yaml:"height"`
	BackgroundColor string `yaml:"background_color,omitempty"`
}

type Layer struct {
	Name string
	Keys []LayerKey `yaml:"keys,omitempty"`
}

type LayerKey struct {
	Code   string `yaml:"code"`
	Label  string `yaml:"label"`
	Shift  string `yaml:"shift"`
	Hold   string `yaml:"hold"`
	Active bool   `yaml:"active"`
}

type LayoutConfig struct {
	Keys []LayoutKeyConfig `yaml:"keys"`
}

type LayoutKeyConfig struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
	W float64 `yaml:"w"`
	H float64 `yaml:"h"`
}

func (c *Config) GetLayer(name string) *Layer {
	for _, layer := range c.Layers {
		if layer.Name == name {
			return &layer
		}
	}
	return nil
}
