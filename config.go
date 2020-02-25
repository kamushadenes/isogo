package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigISO struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Discover    bool   `yaml:"discover"`
	Regex       string `yaml:"regex"`
	Destination string `yaml:"destination"`
}

type ConfigKeep struct {
	Directory string `yaml:"directory"`
	Regex     string `yaml:"regex"`
	Last      int    `yaml:"last"`
}

type Config struct {
	ISOs []*ConfigISO  `yaml:"isos"`
	Keep []*ConfigKeep `yaml:"keep"`
}

func readConfig(fname string) (*Config, error) {
	body, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	var cnf Config

	err = yaml.Unmarshal(body, &cnf)
	if err != nil {
		return nil, err
	}

	return &cnf, nil
}
