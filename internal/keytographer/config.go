package keytographer

type Config struct {
	Name     string
	Keyboard string
	Canvas   *CanvasConfig
}

type CanvasConfig struct {
	Width           int
	Height          int
	BackgroundColor string `yaml:"background_color"`
}
