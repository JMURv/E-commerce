package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port         int    `yaml:"port"`
	ServiceName  string `yaml:"serviceName"`
	RegistryAddr string `yaml:"registryAddr"`
	RedisAddr    string `yaml:"redisAddr"`
	RedisPass    string `yaml:"redisPass"`

	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`

	Kafka struct {
		Addrs             []string `yaml:"addrs"`
		NotificationTopic string   `yaml:"notificationTopic"`
	} `yaml:"kafka"`
}

func LoadConfig(state string) (*Config, error) {
	var conf Config

	data, err := os.ReadFile(fmt.Sprintf("../%v.config.yaml", state))
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
