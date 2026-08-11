package atomicx

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Addr    string
	Timeout time.Duration
}

var (
	mu  sync.Mutex
	cfg atomic.Pointer[Config]
)

func Load() (*Config, error) {
	if c := cfg.Load(); c != nil {
		return c, nil
	}
	mu.Lock()
	defer mu.Unlock()
	if c := cfg.Load(); c != nil {
		return c, nil
	}
	c, err := parse()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	cfg.Store(c)
	return c, nil
}

func Reload() error {
	mu.Lock()
	defer mu.Unlock()
	c, err := parse()
	if err != nil {
		return fmt.Errorf("reload config: %w", err)
	}
	cfg.Store(c)
	return nil
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
