package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Host string
	Port int
}

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path"`
	DbConn      string     `yaml:"db_conn"`
	HttpServer  HttpServer `yaml:"http_server"`
	JwtSecret   string     `yaml:"jwt_secret"`
}

func NewConfig() (*Config, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		return nil, fmt.Errorf("config file path is required")
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
