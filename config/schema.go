package config

import (
	"bytes"
	_ "embed"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed config.schema.yaml
	schemaYAML string
	schema     *jsonschema.Schema
)

func Schema() *jsonschema.Schema {
	if schema == nil {
		err := initSchema()
		if err != nil {
			panic(err)
		}
	}
	return schema
}

func Validate(data []byte) error {
	var v interface{}
	err := yaml.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	err = Schema().Validate(v)
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
