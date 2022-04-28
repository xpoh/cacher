package config

import (
	"log"
)

var config *Cfg

type ConfigLoader interface {
	load() (Cfg, error)
}

type Cfg struct {
	NumberOfRequests int      `yaml:"NumberOfRequests"`
	MinTimeout       int      `yaml:"MinTimeout"`
	MaxTimeout       int      `yaml:"MaxTimeout"`
	URLs             []string `yaml:"URLs"`
}

func NewConfig() *Cfg {
	log.Println("Load yaml config file...")
	cfg, err := yamlConfig{}.load()
	if err != nil {
		log.Panicln("Error load config from yaml.")
	}
	log.Printf("%v\n", cfg)
	config = &cfg
	return config
}

func GetConfig() *Cfg {
	if config == nil {
		panic("get nil config error!!!")
	}
	return config
}
