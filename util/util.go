package util

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	// 方法1: 检查X-Forwarded-For (代理/负载均衡器)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if ip != "" && !isLocalIP(ip) {
			return ip
		}
	}

	// 方法2: 检查X-Real-IP (Nginx代理)
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		if !isLocalIP(xri) {
			return xri
		}
	}

	// 方法3: 使用Gin内置方法
	if ip := c.ClientIP(); ip != "" {
		if !isLocalIP(ip) {
			return ip
		}
	}

	// 方法4: 使用RemoteIP作为后备
	if ip := c.RemoteIP(); ip != "" {
		if !isLocalIP(ip) {
			return ip
		}
		if ip == "::1" || ip == "127.0.0.1" {
			return "本地访问(localhost)"
		}
	}

	return "未知IP"
}

// 检查是否为本地IP
func isLocalIP(ip string) bool {
	localIPs := []string{
		"127.0.0.1", "::1", "localhost",
		"0.0.0.0", "::", "fe80::",
	}

	for _, localIP := range localIPs {
		if strings.HasPrefix(ip, localIP) {
			return true
		}
	}

	// 检查私有网络IP段
	if strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.") {
		return false // 这些是私有IP，但不是localhost
	}

	return false
}
