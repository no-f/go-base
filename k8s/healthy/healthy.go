package healthy

import (
	"log"
	"net/http"
	"strings"
)

func ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getRealIP(r)
		log.Printf("k8s healthy check: ReadinessHandler - ClientIP: %s\n", clientIP)
		w.WriteHeader(http.StatusOK)
	}
}

func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getRealIP(r)
		log.Printf("k8s healthy check: LivenessHandler - ClientIP: %s\n", clientIP)
		w.WriteHeader(http.StatusOK)
	}
}

// getRealIP 从请求中提取客户端的真实 IP 地址
func getRealIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}
	return r.RemoteAddr
}
