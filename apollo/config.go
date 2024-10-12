package apollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"go-base/apollo/info"
	"go-base/apollo/model"
	"log"
	"strings"
)

var apolloClientV1 agollo.Client
var apolloClientV2 agollo.Client

var allNamespaceNameV1 string
var allNamespaceNameV2 string

var namespacesV1 []string
var namespacesV2 []string

// init 初始化
func init() {
	apolloConfig, err := model.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig失败: %v", err)
	}

	allNamespaceNameV1 = info.Namespace + apolloConfig.Apollo.CommonNamespaceName
	allNamespaceNameV2 = "application," + apolloConfig.Apollo.NamespaceName

	apolloClientV1 = initApolloClient(allNamespaceNameV1, apolloConfig)
	apolloClientV2 = initApolloClient(allNamespaceNameV2, apolloConfig)

	namespacesV1 = strings.Split(allNamespaceNameV1, ",")
	namespacesV2 = strings.Split(allNamespaceNameV2, ",")
}

func initApolloClient(namespaceName string, apolloConfig *model.ApolloConfig) agollo.Client {
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
	value := getConfigValueFromClient(apolloClientV1, namespacesV1, key)
	if value != "" {
		return value
	}
	return getConfigValueFromClient(apolloClientV2, namespacesV2, key)
}

func getConfigValueFromClient(client agollo.Client, namespaces []string, key string) string {
	for _, ns := range namespaces {
		cache := client.GetConfigCache(ns)
		if cache == nil {
			log.Printf("Namespace cache not found for: %s", ns)
			continue
		}
		value, err := cache.Get(key)
		if err == nil {
			return toString(value)
		}
		log.Printf("Failed to get config value for key %s in namespace %s: %v", key, ns, err)
	}
	return ""
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
