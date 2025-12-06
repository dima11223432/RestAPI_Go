package main

import (
	apiserver "RestApi/internal/app/apiServer"
	"flag"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiServer.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		panic(err)
	}

	if err := apiserver.Start(config); err != nil {
		panic(err)
	}
}
