package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

//env variable name prefix
const envPrefix = "WEBCRAWLER"

// Config struct
type Config struct {
	HTTPPort int `envconfig:"default=8001"`
	Workers  int `envconfig:"default=10"`
}

// InitConfig initializes configuration struct from environment variables or .env file
func InitConfig() (config Config, err error) {

	//load environment variables from the .env file
	if err = godotenv.Load(); err != nil {
		fmt.Printf("Use env variables with prefix %v_ or .env file to overload configuration\n", envPrefix)
	}

	return config, envconfig.InitWithPrefix(&config, envPrefix)
}
