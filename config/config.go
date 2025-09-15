package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host        string `env:"HOST,notEmpty"`
	User_DB     string `env:"USER_DB,notEmpty"`
	Password_DB string `env:"PASSWORD_DB,notEmpty"`
	Db_name     string `env:"DB_NAME,notEmpty"`
	Port        string `env:"PORT,notEmpty"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)

	return cfg, err
}
