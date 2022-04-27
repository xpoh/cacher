package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type yamlConfig struct{}

func (y yamlConfig) load() (Cfg, error) {
	cfg := Cfg{}
	yamlFile, err := os.ReadFile("configs/config.yaml")

	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)

	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
