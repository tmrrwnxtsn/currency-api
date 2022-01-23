package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/tmrrwnxtsn/currency-api/internal/apiserver"
	"github.com/tmrrwnxtsn/currency-api/internal/config"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to the toml config file")
}

func main() {
	flag.Parse()

	cfg := config.New()
	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		log.Fatalf("error occured while decoding config file: %s", err.Error())
	}

	srv := apiserver.New(cfg)
	if err = srv.Run(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
