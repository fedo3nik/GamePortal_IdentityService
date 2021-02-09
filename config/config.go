package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host    string `envconfig:"IDENTITYSERVICE_HOST"`
	Port    string `envconfig:"IDENTITYSERVICE_PORT"`
	ConnURI string `envconfig:"IDENTITYSERVICE_CONN_URI"`
	DB      string `envconfig:"IDENTITYSERVICE_DB"`
}

func NewConfig() (*Config, error) {
	var c Config

	err := envconfig.Process("identityservice", &c)
	if err != nil {
		log.Printf("Process config file error: %v\n", err)
		return nil, err
	}

	return &c, nil
}