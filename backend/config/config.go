package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Local    string `yaml:"local"`
	Port     int    `yaml:"port"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type RedisConfig struct {
	Local        string `yaml:"local"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	Db           int    `yaml:"db"`
	Poolsize     int    `yaml:"poolsize"`
	Minidleconns int    `yaml:"minidleconns"`
	Maxretries   int    `yaml:"maxretries"`
}

var (
	instance *Config
	once     sync.Once
)

var ConfigPath = "config/config.yaml"

func InitConf() *Config {
	once.Do(func() {
		instance = &Config{}
		data, err := ioutil.ReadFile(ConfigPath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, instance)
		if err != nil {
			panic(err)
		}
	})
	return instance
}
