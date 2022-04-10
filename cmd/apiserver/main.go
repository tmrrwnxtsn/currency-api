package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/tmrrwnxtsn/currency-api/internal/apiserver"
	"log"
	"os"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to the toml config file")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error occured while loading .env file: %s", err.Error())
	}
}

func main() {
	flag.Parse()

	cfg := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		log.Fatalf("error occured while decoding config file: %s", err.Error())
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.SetCurrencyApiUrl(os.Getenv("CURRENCY_API_KEY"))

	go ratesUpdating(cfg.UpdateInterval)

	if err = apiserver.Start(cfg); err != nil {
		log.Fatalf("error occured while starting apiserver: %s", err.Error())
	}
}

func ratesUpdating(updateInterval int) {
	for range time.Tick(time.Second * time.Duration(updateInterval)) {
		fmt.Println("Updating rates...")
	}
}
