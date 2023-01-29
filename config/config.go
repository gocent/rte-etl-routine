package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	Auth    auth
	Ecowatt ecowatt
}

type auth struct {
	URI  string `env:"AUTH_URI,required"`
	Code string `env:"AUTH_CODE,required"`
}

type ecowatt struct {
	URI string `env:"ECOWATT_URI,required"`
}

var (
	once     sync.Once
	instance *Config
	err      error
)

func GetEnv() *Config {
	once.Do(func() {
		instance, err = getInstance()
	})
	if err != nil {
		log.Fatal("Unable to parse env files")
	}
	return instance
}

func getInstance() (*Config, error) {
	environment := os.Getenv("ECOWATT_ROUTINE_ENV")
	if environment == "" {
		environment = "development"
	}
	var err error

	switch environment {
	case "prod":
		err = godotenv.Load(".env")
	case "staging":
		err = godotenv.Load(".env." + environment)
	default:
		err = godotenv.Load(".env." + environment + ".local")
	}
	if err != nil {
		log.Printf("Unable to load environment file for %s", environment)
		return nil, err
	}
	cfg := Config{}
	err = env.Parse(&cfg.Auth)    // 👈 Parse environment variables into `config`
	err = env.Parse(&cfg.Ecowatt) // 👈 Parse environment variables into `config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return &cfg, err
}
