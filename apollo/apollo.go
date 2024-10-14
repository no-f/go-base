package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/no-f/go-base/apollo/info"
	"github.com/no-f/go-base/config"
	"github.com/no-f/go-base/config/models"
	"github.com/no-f/go-base/logger"
	"log"
	"strings"
	"sync"
)

var apolloClient agollo.Client
var namespaceNames []string

// Initialize 初始化 Apollo 客户端
func Initialize(apolloConfig *models.ApolloYAMLConfig) {

	logger.Info("开始初始化 APOLLO 客户端")
	allNamespaceNames := []string{
		info.BullyunV2Namespace + apolloConfig.CommonNamespaceName,
		info.ServiceNamespace + apolloConfig.NamespaceName,
	}
	namespaceName := strings.Join(allNamespaceNames, ",")

	c := &config.AppConfig{
		AppID:          apolloConfig.AppID,
		Cluster:        apolloConfig.Cluster,
		IP:             apolloConfig.Meta,
		NamespaceName:  namespaceName,
		IsBackupConfig: false,
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		log.Fatalf("StartWithConfig失败: %v", err)
	}

	apolloClient = client
	namespaceNames = strings.Split(strings.Join(allNamespaceNames, ","), ",")
}

// GetConfigValue 根据给定的 key 获取配置值
func GetConfigValue(key string) string {
	// 使用并发方式从命名空间获取配置值
	return fetchConfigValueFromClient(apolloClient, namespaceNames, key)
}

// getConfigValueFromClient 从指定命名空间列表中获取配置值
func fetchConfigValueFromClient(client agollo.Client, namespaces []string, key string) string {
	results := make(chan string, len(namespaces))
	var wg sync.WaitGroup

	// 并发获取每个命名空间的配置
	for _, ns := range namespaces {
		wg.Add(1)
		go func(ns string) {
			defer wg.Done()
			cache := client.GetConfigCache(ns)
			if cache != nil {
				value, err := cache.Get(key)
				if err == nil {
					if valueStr, ok := value.(string); ok {
						results <- valueStr
						return
					}
				}
			} else {
				log.Printf("Namespace cache not found for: %s", ns)
			}
			results <- ""
		}(ns)
	}

	// 等待所有并发获取完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 获取第一个非空的配置值
	for result := range results {
		if result != "" {
			return result
		}
	}
	return ""
}

func init() {
	configFromYAML, err := ymlconfig.LoadYAMLConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML config: %v", err)
	}

	//init logger
	logger.Initialize(&configFromYAML.Logger)

	//init apollo
	Initialize(&configFromYAML.Apollo)
}
