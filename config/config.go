package config

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Listen string             `yaml:"listen"`
	LB     LoadBalancerConfig `yaml:"lb" envPrefix:"LB_"`
	Rate   Rate               `yaml:"rate"`
	Health Health             `yaml:"health"`
}

type LoadBalancerConfig struct {
	Alg      string   `yaml:"Alg"`
	Backends []string `yaml:"backends" env:"BACKENDS" envSeparator:","`
}

type Rate struct {
	Capacity int `yaml:"capacity" env:"CAPACITY"`
	RPS      int `yaml:"rps" env:"RPS"`
}
type Health struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

func New(path string) (*Config, error) {
	cfg := &Config{
		Listen: ":8080",
		LB:     LoadBalancerConfig{},
		Rate: Rate{
			Capacity: 100,
			RPS:      10,
		},
	}
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	}

	return cfg, nil
}
