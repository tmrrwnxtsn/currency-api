package config

import (
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BindAddr       string `toml:"bind_addr"`       // server address
	UpdateInterval int    `toml:"update_interval"` // in minutes

	CurrencyAPIKey string
	DatabaseURL    string
}

func New() *Config {
	return &Config{
		BindAddr:       ":8080",
		UpdateInterval: 10,
	}
}

func (c *Config) LoadTomlConfig(configPath string) error {
	_, err := toml.DecodeFile(configPath, c)
	return err
}

func (c *Config) LoadEnvConfig(configPath string) error {
	if err := godotenv.Load(configPath); err != nil {
		return err
	}

	c.DatabaseURL = os.Getenv("DATABASE_URL")
	c.CurrencyAPIKey = os.Getenv("CURRENCY_API_KEY")

	return nil
}
