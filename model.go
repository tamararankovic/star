package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Star struct {
	Conf Config `yaml:"star"`
}

type Config struct {
	Version        string            `yaml:"version"`
	NodeId         string            `yaml:"nodeid"`
	RTopic         string            `yaml:"rtopic"`
	STopic         string            `yaml:"stopic"`
	ErrTopic       string            `yaml:"errtopic"`
	Flusher        string            `yaml:"flusher"`
	InstrumentConf map[string]string `yaml:"instrument"`
}

func ConfigFile(n ...string) (*Config, error) {
	path := "config.yml"
	if len(n) > 0 {
		path = n[0]
	}

	yamlFile, err := ioutil.ReadFile(path)
	check(err)

	var conf Star
	err = yaml.Unmarshal(yamlFile, &conf)
	check(err)

	return &conf.Conf, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
