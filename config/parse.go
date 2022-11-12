package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Load(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func Parse(data []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	logrus.WithField("config", config).Debug("parsed config")

	return &config, nil
}
