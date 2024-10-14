package models

type ApolloYAMLConfig struct {
	AppID               string `yaml:"AppID"`
	Cluster             string `yaml:"Cluster"`
	Meta                string `yaml:"Meta"`
	CommonNamespaceName string `yaml:"CommonNamespaceName"`
	NamespaceName       string `yaml:"NamespaceName"`
}
