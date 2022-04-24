package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	IsDev bool
	Port  string `env:"PORT"`
}

const (
	devEnvFilePath = "./dev.env"
)

var (
	do  sync.Once
	cfg AppConfig
)

func GetConfig() AppConfig {

	do.Do(
		func() {
			var isDev bool
			flag.BoolVar(&isDev, "DEV", false, "")
			flag.Parse()

			cfg.IsDev = isDev

			if isDev {
				err := godotenv.Load(devEnvFilePath)
				if err != nil {
					panic(
						fmt.Sprintf("Failed to load env file, by path: %s\nErr: %s", devEnvFilePath, err),
					)

				}
			}

			err := env.Parse(&cfg)
			if err != nil {
				panic(
					fmt.Sprintf("Failed to parse env to cfg\nErr: %s", err),
				)
			}
		},
	)

	return cfg
}
