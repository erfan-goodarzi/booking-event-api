package config

import (
	"fmt"
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type AppConfig struct {
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

var localConfig *Config = nil

func Load() *Config {
	if localConfig != nil {
		return localConfig
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalln(err)
	}

	localConfig = &config

	return &config
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
	)
}
