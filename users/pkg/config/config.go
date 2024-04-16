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

	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`

	Redis struct {
		Addr string `yaml:"addr"`
		Pass string `yaml:"pass"`
	} `yaml:"redis"`

	Kafka struct {
		Addrs             []string `yaml:"addrs"`
		NotificationTopic string   `yaml:"notificationTopic"`
	} `yaml:"kafka"`

	Jaeger struct {
		Sampler struct {
			Type  string `yaml:"type"`
			Param int    `yaml:"param"`
		} `yaml:"sampler"`
		Reporter struct {
			LogSpans           bool   `yaml:"LogSpans"`
			LocalAgentHostPort string `yaml:"LocalAgentHostPort"`
		} `yaml:"reporter"`
	} `yaml:"jaeger"`
}

func LoadConfig(configName string) (*Config, error) {
	var conf Config

	data, err := os.ReadFile(fmt.Sprintf("../%v.yaml", configName))
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
