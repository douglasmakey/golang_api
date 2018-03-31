package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Host         string `json:"host"`
		Port         string `json:"port"`
		IsProduction bool   `json:"is_production"`
		PasswordSalt []byte `json:"password_salt"`
		JwtSecret    []byte `json:"jwt_secret"`
		Debug        bool   `json:"debug"`
	} `json:"server"`
	Postgres struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DB       string `json:"db"`
	} `json:"postgres"`
}

var cfg Config

func FromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}
