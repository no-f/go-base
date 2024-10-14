package models

type YAMLConfig struct {
	Logger LoggerYAMLConfig `yaml:"logger"`
	Apollo ApolloYAMLConfig `yaml:"apollo"`
}
