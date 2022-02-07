package apiserver

import "github.com/tmrrwnxtsn/currency-api/internal/store"

type Config struct {
	BindAddr string `toml:"bind_addr"` // server address
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
