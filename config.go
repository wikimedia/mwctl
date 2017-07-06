package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	DELIM = "."
)

type Service struct {
	Source string `yaml:"local_source"`
}

type Config struct {
	Services map[string]Service
}

// NewConfig returns an initialized (empty) Config struct.
func NewConfig() *Config {
	return &Config{Services: make(map[string]Service)}
}

func (c *Config) Get(k string) (string, error) {
	key := strings.Split(k, DELIM)

	switch strings.ToLower(key[0]) {
	case "service":
		if len(key) == 1 {
			// Return the YAML output of the whole service stanza
			data, err := yaml.Marshal(c.Services)
			if err != nil {
				return "", err
			}
			return string(data), nil

		} else if len(key) == 2 {
			// Return the YAML output for the specific service
			data, err := yaml.Marshal(c.Services[key[1]])
			if err != nil {
				return "", err
			}
			return string(data), nil
		}

		// Return the requested value for a service
		switch strings.ToLower(key[2]) {
		case "source":
			return c.Services[key[1]].Source, nil
		default:
			return "", &KeyError{k}
		}

	default:
		return "", &KeyError{k}
	}
}

func (c *Config) Set(k string, v string) error {
	key := strings.Split(k, DELIM)

	switch strings.ToLower(key[0]) {
	case "service":
		if len(key) != 3 {
			return &ServiceError{key[2]}
		}

		switch strings.ToLower(key[2]) {
		case "source":
			if val, ok := c.Services[key[1]]; ok {
				val.Source = v
				c.Services[key[1]] = val
			} else {
				c.Services[key[1]] = Service{Source: v}
			}
		default:
			return &KeyError{k}
		}

		return nil
	default:
		return &KeyError{k}
	}
}

type KeyError struct {
	name string
}

func (e *KeyError) Error() string {
	return fmt.Sprintf("Invalid key: %s", e.name)
}

type ServiceError struct {
	name string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("Invalid service name: %s", e.name)
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
