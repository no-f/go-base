package model

import (
	"gopkg.in/yaml.v3"
	"os"
)

// ApolloConfig 定义了Apollo配置的结构
type ApolloConfig struct {
	Apollo struct {
		AppID               string `yaml:"AppID"`
		Cluster             string `yaml:"Cluster"`
		Meta                string `yaml:"Meta"`
		CommonNamespaceName string `yaml:"CommonNamespaceName"`
		NamespaceName       string `yaml:"NamespaceName"`
	} `yaml:"apollo"`
}

// LoadConfig 从YAML文件中读取配置并返回配置对象
func LoadConfig() (*ApolloConfig, error) {
	var apolloConfig ApolloConfig

	// 从文件中读取YAML配置
	yamlFile, err := os.ReadFile("E:\\no_f\\code\\github-x\\go\\go-base\\config.yaml")
	if err != nil {
		return nil, err
	}

	// 解析YAML配置到ApolloConfig结构体
	err = yaml.Unmarshal(yamlFile, &apolloConfig)
	if err != nil {
		return nil, err
	}
	setApolloConfigFromEnv(apolloConfig)
	return &apolloConfig, nil
}

// 从环境变量中读取Apollo配置 k8s 编排文件部署
func setApolloConfigFromEnv(apolloConfig ApolloConfig) {
	apolloMeta := os.Getenv("APOLLO_META")
	if apolloMeta != "" {
		apolloConfig.Apollo.Meta = apolloMeta
	}

	apolloCluster := os.Getenv("APOLLO_CLUSTER")
	if apolloMeta != "" {
		apolloConfig.Apollo.Cluster = apolloCluster
	}
}
