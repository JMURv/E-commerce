package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port         int    `yaml:"port"`
	ServiceName  string `yaml:"serviceName"`
	RegistryAddr string `yaml:"registryAddr"`
	RedisAddr    string `yaml:"redisAddr"`
	RedisPass    string `yaml:"redisPass"`

	Kafka struct {
		Addrs             []string `yaml:"addrs"`
		NotificationTopic string   `yaml:"notificationTopic"`
	} `yaml:"kafka"`
}

func LoadConfig() (*Config, error) {
	var conf Config

	data, err := os.ReadFile("../dev.config.yaml")
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
