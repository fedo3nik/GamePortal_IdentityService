package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host              string `envconfig:"IDENTITY_SERVICE_HOST"`
	Port              string `envconfig:"IDENTITY_SERVICE_PORT"`
	ConnURI           string `envconfig:"IDENTITY_SERVICE_CONN_URI"`
	DB                string `envconfig:"IDENTITY_SERVICE_DB"`
  GrpcPort          string `envconfig:"GRPC_PORT"`
	RefreshPrivateKey string `envconfig:"IDENTITY_SERVICE_REFRESH_PRIVATE_KEY"`
	AccessPrivateKey  string `envconfig:"IDENTITY_SERVICE_ACCESS_PRIVATE_KEY"`
  RefreshPublicKey  string `envconfig:"IDENTITYSERVICE_REFRESH_PUBLIC_KEY"`
	AccessPublicKey   string `envconfig:"IDENTITYSERVICE_ACCESS_PUBLIC_KEY"`
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
