package export

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/chdorner/keytographer/internal/keytographer"
	"github.com/stretchr/testify/assert"
)

func TestCExportExistingFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		fixture    string
		configFile string
		expected   string
	}{
		{"paintbrush", "keytographer.yaml", "keymap.c"},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			fixtureDir := filepath.Join(cwd, "..", "..", "test", "fixtures", c.fixture)

			outDir, err := os.MkdirTemp("", "")
			assert.NoError(t, err)
			outPath := filepath.Join(outDir, "actual.c")

			expected, err := os.ReadFile(filepath.Join(fixtureDir, c.expected))
			assert.NoError(t, err)

			err = os.WriteFile(outPath, expected, 0644)
			assert.NoError(t, err)

			exporter := NewCExporter(context.TODO())
			config := loadAndParse(t, filepath.Join(fixtureDir, c.configFile))

			err = exporter.Export(config, outPath)
			assert.NoError(t, err)

			actual, err := os.ReadFile(outPath)
			assert.NoError(t, err)
			assert.Equal(t, string(expected), string(actual))
		})
	}
}

func TestCExportNewFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		fixture    string
		configFile string
	}{
		{"paintbrush", "keytographer.yaml"},
	}

	for _, c := range cases {
		t.Run(c.fixture, func(t *testing.T) {
			fixtureDir := filepath.Join(cwd, "..", "..", "test", "fixtures", c.fixture)

			outDir, err := os.MkdirTemp("", "")
			assert.NoError(t, err)
			outPath := filepath.Join(outDir, "actual.c")

			exporter := NewCExporter(context.TODO())
			config := loadAndParse(t, filepath.Join(fixtureDir, c.configFile))

			err = exporter.Export(config, outPath)
			assert.NoError(t, err)

			_, err = os.ReadFile(outPath)
			assert.NoError(t, err)
		})
	}
}

func loadAndParse(t *testing.T, configFile string) *keytographer.Config {
	data, err := keytographer.Load(configFile)
	if err != nil {
		t.Fatal(err)
	}

	config, err := keytographer.Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	return config
}
