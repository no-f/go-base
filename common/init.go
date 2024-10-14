package common

import (
	"github.com/no-f/go-base/apollo"
	"github.com/no-f/go-base/config"
	"github.com/no-f/go-base/logger"
	"log"
)

// init 初始化
func init() {
	configFromYAML, err := ymlconfig.LoadYAMLConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML config: %v", err)
	}

	//init logger
	logger.Initialize(&configFromYAML.Logger)

	//init apollo
	apollo.Initialize(&configFromYAML.Apollo)
}
