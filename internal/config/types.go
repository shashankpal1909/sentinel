package config

import (
	"encoding/json"
	"fmt"
)

type BalancerStrategy string

const (
	RoundRobin BalancerStrategy = "round-robin"
	Random     BalancerStrategy = "random"
)

type Config struct {
	Server   ServerConfig             `yaml:"server" json:"server"`
	Services map[string]ServiceConfig `yaml:"services" json:"services"`
	Routes   []RouteConfig            `yaml:"routes" json:"routes"`
}

type ServerConfig struct {
	Port int `yaml:"port" json:"port"`
}

type HealthCheckConfig struct {
	Path               string `yaml:"path" json:"path"`
	Interval           string `yaml:"interval" json:"interval"`
	Timeout            string `yaml:"timeout" json:"timeout"`
	HealthyThreshold   int    `yaml:"healthy_threshold" json:"healthy_threshold"`
	UnhealthyThreshold int    `yaml:"unhealthy_threshold" json:"unhealthy_threshold"`
}

type ServiceConfig struct {
	Strategy    BalancerStrategy   `yaml:"strategy" json:"strategy"`
	Backends    []string           `yaml:"backends" json:"backends"`
	HealthCheck *HealthCheckConfig `yaml:"health_check" json:"health_check"`
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
