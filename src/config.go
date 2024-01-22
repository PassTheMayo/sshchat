package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func (c *Config) ReadFile(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
