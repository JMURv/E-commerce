package conf

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
	MongoAddr    string `yaml:"mongoAddr"`

	Kafka struct {
		Addrs     []string `yaml:"addrs"`
		Topics    []string `yaml:"topics"`
		TopicName struct {
			Create string `yaml:"create"`
			Update string `yaml:"update"`
			Delete string `yaml:"delete"`
		} `yaml:"topicName"`
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
