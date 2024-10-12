package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/no-f/go-base/apollo/info"
	"github.com/no-f/go-base/apollo/model"
	"log"
	"strings"
	"sync"
)

var client agollo.Client
var namespaceNames []string

// init 初始化
func init() {
	apolloConfig, err := model.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig失败: %v", err)
	}

	// 构造命名空间列表
	allNamespaceNames := []string{
		info.BullyunV2Namespace + apolloConfig.Apollo.CommonNamespaceName,
		info.ServiceNamespace + apolloConfig.Apollo.NamespaceName,
	}

	// 初始化 Apollo 客户端
	client = initializeApolloClient(strings.Join(allNamespaceNames, ","), apolloConfig)

	// 将命名空间拆分为 slice
	namespaceNames = strings.Split(strings.Join(allNamespaceNames, ","), ",")
}

// initApolloClient 初始化 Apollo 客户端
func initializeApolloClient(namespaceName string, apolloConfig *model.ApolloConfig) agollo.Client {
	c := &config.AppConfig{
		AppID:          apolloConfig.Apollo.AppID,
		Cluster:        apolloConfig.Apollo.Cluster,
		IP:             apolloConfig.Apollo.Meta,
		NamespaceName:  namespaceName,
		IsBackupConfig: false,
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		log.Fatalf("StartWithConfig失败: %v", err)
	}
	return client
}

// GetConfigValue 根据给定的 key 获取配置值
func GetConfigValue(key string) string {
	// 使用并发方式从命名空间获取配置值
	return fetchConfigValueFromClient(client, namespaceNames, key)
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
				} else {
					log.Printf("Failed to get config value for key %s in namespace %s: %v", key, ns, err)
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

// toString 将 interface{} 转换为 string
func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
