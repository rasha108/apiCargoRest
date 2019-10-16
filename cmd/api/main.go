package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/rasha108/apiCargoRest.git/internal/app/api"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Start(config); err != nil {
		log.Fatal(err)
	}
}
