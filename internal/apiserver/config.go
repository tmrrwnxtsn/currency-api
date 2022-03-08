package apiserver

type Config struct {
	BindAddr       string `toml:"bind_addr"` // server address
	LogLevel       string `toml:"log_level"`
	DatabaseURL    string `toml:"database_url"`
	CurrencyAPIKey string `toml:"currency_api_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
