package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Api struct {
	PORT string `yaml:"port"`
	HOST string `yaml:"host"`
	Auth `yaml:"auth"`
}

type Auth struct {
	EnableAuth bool          `yaml:"enable_auth"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
	SignKey    string        `yaml:"sign_key"`
}

func (api *Api) GetAddr() string {
	return fmt.Sprintf("%s:%s", api.HOST, api.PORT)
}

type PostgresDSN string

func (p PostgresDSN) String() string {
	return string(p)
}

type ServiceConfiguration struct {
	Api         `yaml:"api"`
	PostgresDSN `yaml:"postgres_dsn"`
}

func Load() ServiceConfiguration {
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var cfg ServiceConfiguration
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func NewConfig() ServiceConfiguration {
	return Load()
}
