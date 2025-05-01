package config

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB     DatabaseConfig     `yaml:"db"`
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
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Name     string `yaml:"name"`
	MinConns int    `yaml:"min_conns"`
	MaxConns int    `yaml:"max_conns"`
}

func New(path string) (*Config, error) {
	cfg := &Config{
		DB: DatabaseConfig{
			Host:     "postgres",
			Port:     5432,
			User:     "user",
			Pass:     "secret",
			Name:     "postgres_db",
			MinConns: 5,
			MaxConns: 10,
		},
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
func (d *DatabaseConfig) URL() string {

	name := strings.TrimSpace(d.Name)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_min_conns=%d&pool_max_conns=%d",
		d.User, d.Pass,
		d.Host, d.Port,
		name,
		d.MinConns, d.MaxConns,
	)
}
