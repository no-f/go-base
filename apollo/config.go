package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"go-base/apollo/info"
	"go-base/apollo/model"
	"strings"
)

var apolloClientV1 agollo.Client
var apolloClientV2 agollo.Client

var allNamespaceNameV1 string
var allNamespaceNameV2 string

// init 初始化
func init() {
	apolloInitV1()
	apolloInitV2()
}

func apolloInitV1() {
	apolloConfig, err := model.LoadConfig()
	if err != nil {
		fmt.Printf("LoadConfig失败: %v", err)
	}

	allNamespaceNameV1 = info.Namespace + apolloConfig.Apollo.CommonNamespaceName
	c := &config.AppConfig{
		AppID:          "bullyun-v2",
		Cluster:        apolloConfig.Apollo.Cluster,
		IP:             apolloConfig.Apollo.Meta,
		NamespaceName:  allNamespaceNameV1,
		IsBackupConfig: true,
	}
	apolloClientV1, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
}

// apolloInit 初始化 apollo 客户端
func apolloInitV2() {
	apolloConfig, err := model.LoadConfig()
	if err != nil {
		fmt.Printf("LoadConfig失败: %v", err)
	}

	allNamespaceNameV2 = "application," + apolloConfig.Apollo.NamespaceName
	c := &config.AppConfig{
		AppID:          apolloConfig.Apollo.AppID,
		Cluster:        apolloConfig.Apollo.Cluster,
		IP:             apolloConfig.Apollo.Meta,
		NamespaceName:  allNamespaceNameV2,
		IsBackupConfig: true,
	}
	apolloClientV2, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
}

// GetConfigValue 根据给定的 key 获取配置值
func GetConfigValue(key string) string {
	if apolloClientV1 == nil {
		return ""
	}

	for _, ns := range strings.Split(allNamespaceNameV1, ",") {
		cache := apolloClientV1.GetConfigCache(ns)

		if cache == nil {
			continue
		}
		value, err := cache.Get(key)
		if err != nil {
			continue
		}
		return toString(value)
	}

	for _, ns := range strings.Split(allNamespaceNameV2, ",") {
		cache := apolloClientV2.GetConfigCache(ns)

		if cache == nil {
			continue
		}
		value, err := cache.Get(key)
		if err != nil {
			continue
		}
		return toString(value)
	}
	return ""
}

func toString(value interface{}) string {
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}
