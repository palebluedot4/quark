package once

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Config struct {
	Addr    string
	Timeout time.Duration
}

var Load = sync.OnceValues(parse)

var (
	once   sync.Once
	cfg    *Config
	cfgErr error
)

func LoadManual() (*Config, error) {
	once.Do(func() {
		cfg, cfgErr = parse()
	})
	return cfg, cfgErr
}

func parse() (*Config, error) {
	timeout, err := time.ParseDuration(getenv("APP_TIMEOUT", "30s"))
	if err != nil {
		return nil, fmt.Errorf("parse APP_TIMEOUT: %w", err)
	}
	return &Config{
		Addr:    getenv("APP_ADDR", "localhost:8080"),
		Timeout: timeout,
	}, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
