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
	X float32 `yaml:"x"`
	Y float32 `yaml:"y"`
	W float32 `yaml:"w"`
	H float32 `yaml:"h"`
}
