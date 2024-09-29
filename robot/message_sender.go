package robot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       100, // TODO 可根据需要调整
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true, // 禁用压缩以减少CPU负载
		},
	}
	once sync.Once
)

type WeChatMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// AsyncSendMessage 异步发送消息到企业微信群机器人，不等待结果，且不因异常而阻塞
func AsyncSendMessage(msg string) {
	go func(msg string) {
		message := WeChatMessage{
			MsgType: "markdown",
			Markdown: struct {
				Content string `json:"content"`
			}{
				Content: generateMarkdownWarning(msg),
			},
		}

		requestBody, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Failed to create HTTP request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to send message, status code: %d", resp.StatusCode)
		}
	}(msg)
}

// generateMarkdownWarning 生成Markdown格式的警告信息
func generateMarkdownWarning(msg string) string {
	service := ""
	cluster := ""
	level := "ERROR"
	time := getCurrentSystemTime()
	exception := msg

	return fmt.Sprintf(`# 警告 %s
> 当前集群: %s
> 异常级别: %s
> 触发时间: %s
> 异常: %s`, service, cluster, level, time, exception)
}

// getCurrentSystemTime 返回当前系统时间的RFC3339格式字符串
func getCurrentSystemTime() string {
	return time.Now().Format(time.RFC3339)
}
