package apiserver

type Config struct {
	BindAddr       string `toml:"bind_addr"` // server address
	LogLevel       string `toml:"log_level"`
	DatabaseURL    string `toml:"database_url"`
	UpdateInterval int    `toml:"update_interval"` // in minutes
	CurrencyAPIUrl string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:       ":8080",
		LogLevel:       "debug",
		UpdateInterval: 10,
	}
}

func (c *Config) SetCurrencyApiUrl(currencyApiKey string) {
	c.CurrencyAPIUrl = "https://freecurrencyapi.net/api/v2/latest?apikey=" + currencyApiKey + "&base_currency=%s"
}
