package ymlconfig

import (
	"github.com/no-f/go-base/config/models"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// LoadYAMLConfig load yaml config
func LoadYAMLConfig(filePath string) (configFromYAML *models.YAMLConfig, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open YAML config")
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var cfg models.YAMLConfig
	if err := decoder.Decode(&cfg); err != nil {
		log.Printf("failed to decode YAML config file")
		return nil, err
	}
	return &cfg, nil
}
