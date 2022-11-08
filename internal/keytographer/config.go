package keytographer

type Config struct {
	Name     string
	Keyboard string `yaml:"keyboard,omitempty"`
	Canvas   CanvasConfig
	Layout   LayoutConfig
}

type CanvasConfig struct {
	Width           int    `yaml:"width"`
	Height          int    `yaml:"height"`
	BackgroundColor string `yaml:"background_color,omitempty"`
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
