package main

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Sources map[string]string `yaml:"sources"`
}

func NewConfig() *Config {
	return &Config{Sources: make(map[string]string)}
}

func ReadConfigFile(cf string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(cf)
	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}

func WriteConfigFile(cf string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cf, data, 0600)
}

// GetConfigPath returns the full path to the user's configuration file.
func GetConfigPath() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return path.Join(dir, ".mwctl")
}

// GetConfig returns a Config struct corresponding to the user's configuration,
// or an empty (but initialized) struct if the file is not present.
func GetConfig(config string) (*Config, error) {
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return NewConfig(), nil
	}

	return ReadConfigFile(config)
}
