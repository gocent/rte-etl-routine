package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	Auth     auth
	Influxdb influxdb
	Ecowatt  ecowatt
}

type auth struct {
	URI  string `env:"AUTH_URI,required"`
	Code string `env:"AUTH_CODE,required"`
}

type ecowatt struct {
	URI string `env:"ECOWATT_URI,required"`
}

type influxdb struct {
	Org   string `env:"INFLUXDB_ORG,required"`
	Token string `env:"INFLUXDB_TOKEN,required"`
	Host  string `env:"INFLUXDB_HOST,required"`
	Port  string `env:"INFLUXDB_PORT,required"`
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
	environment := os.Getenv("RTE_ETL_ROUTINE_ENV")
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
	err = env.Parse(&cfg.Auth)
	err = env.Parse(&cfg.Ecowatt)
	err = env.Parse(&cfg.Influxdb)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return &cfg, err
}
