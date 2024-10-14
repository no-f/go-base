package config

import "github.com/no-f/go-base/config/models"

type YAMLConfig struct {
	Logger models.LoggerYAMLConfig `yaml:"logger"`
	Apollo models.ApolloYAMLConfig `yaml:"apollo"`
}
