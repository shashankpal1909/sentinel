package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Server   ServerConfig             `yaml:"server" json:"server"`
	Services map[string]ServiceConfig `yaml:"services" json:"services"`
	Routes   []RouteConfig            `yaml:"routes" json:"routes"`
}

type ServerConfig struct {
	Port int `yaml:"port" json:"port"`
}

type ServiceConfig struct {
	Strategy string   `yaml:"strategy" json:"strategy"`
	Backends []string `yaml:"backends" json:"backends"`
}

type RouteConfig struct {
	Path    string `yaml:"path" json:"path"`
	Service string `yaml:"service" json:"service"`
}

func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", *c)
	}
	return string(data)
}
