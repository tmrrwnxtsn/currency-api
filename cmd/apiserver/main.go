package main

import (
	"flag"
	"github.com/tmrrwnxtsn/currency-api/internal/apiserver"
	"github.com/tmrrwnxtsn/currency-api/internal/config"
	"log"
)

var (
	tomlConfigPath string
	envConfigPath  string
)

func init() {
	flag.StringVar(&tomlConfigPath, "toml-config-path", "configs/apiserver.toml", "path to the toml config file")
	flag.StringVar(&envConfigPath, "env-config-path", "configs/.env", "path to the .env config file")
}

func main() {
	flag.Parse()

	cfg := config.New()
	if err := cfg.LoadTomlConfig(tomlConfigPath); err != nil {
		log.Fatalf("error occured while loading toml config file: %s", err.Error())
	}

	if err := cfg.LoadEnvConfig(envConfigPath); err != nil {
		log.Fatalf("error occured while loading env config file: %s", err.Error())
	}

	if err := apiserver.Start(cfg); err != nil {
		log.Fatalf("error occured while starting apiserver: %s", err.Error())
	}
}
