package main

import (
	"flag"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/apiserver"
	"github.com/tmrrwnxtsn/currency-conversion-api/internal/config"
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

// @title           Simple Currency API
// @version         1.0
// @description     This is a simple currency API that allows you to create exchange rates and convert values from one currency to another.
// @termsOfService  http://swagger.io/terms/

// @license.name  The MIT License (MIT)
// @license.url   https://mit-license.org/

// @host      localhost:8080
// @BasePath  /api/v1

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
		log.Fatalf("error occured while starting API server: %s", err.Error())
	}
}
