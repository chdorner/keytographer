package keytographer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed config.schema.yaml
	schemaYAML string
	schema     *jsonschema.Schema
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

	return &config, nil
}

func Validate(data []byte) error {
	if schema == nil {
		err := initSchema()
		if err != nil {
			return err
		}
	}

	var v interface{}
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	err = schema.Validate(v)
	if err != nil {
		logrus.WithField("error", err).Debug("configuration is invalid")
	}
	return err
}

func initSchema() error {
	var schemaInterface interface{}
	err := yaml.Unmarshal([]byte(schemaYAML), &schemaInterface)
	if err != nil {
		return err
	}

	schemaJSON, err := json.Marshal(schemaInterface)
	if err != nil {
		return err
	}

	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020

	err = c.AddResource("schema.json", bytes.NewReader(schemaJSON))
	if err != nil {
		return err
	}

	schema, err = c.Compile("schema.json")
	return err
}
