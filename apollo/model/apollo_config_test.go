package model

import (
	"testing"
)

// TestLoadConfig 测试LoadConfig函数
func TestLoadConfig(t *testing.T) {
	// 调用待测试的函数
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig失败: %v", err)
	}

	// 验证结果
	if config.Apollo.AppID != "xx-service" {
		t.Errorf("AppID不匹配，期望得到'xx-service'，实际得到'%s'", config.Apollo.AppID)
	}

	if config.Apollo.Cluster != "PRO" {
		t.Errorf("Cluster不匹配，期望得到'PRO'，实际得到'%s'", config.Apollo.Cluster)
	}

}
