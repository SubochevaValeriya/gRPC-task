package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type globalConfig struct {
	Server serverConfig `yaml:"server"`
	Client clientConfig `yaml:"client"`
}

type serverConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	PortRest string `yaml:"portRest"`
}

type clientConfig struct {
	Rusprofile rusprofileConfig `yaml:"rusprofile"`
}

type rusprofileConfig struct {
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

func parseConfig(path string) (globalConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return globalConfig{}, fmt.Errorf("can't parse config: '%w'", err)
	}
	conf := globalConfig{}
	fmt.Printf("content='%#v'\n", string(content))
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return globalConfig{}, fmt.Errorf("can't parse config: '%w'", err)
	}
	fmt.Printf("conf='%#v'\n", conf)
	return conf, nil
}
