package controller

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// 获取当前请求的域名和协议
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// 支持代理情况下的协议检测
	if c.GetHeader("X-Forwarded-Proto") != "" {
		scheme = c.GetHeader("X-Forwarded-Proto")
	}

	host := c.Request.Host
	currentDomain := scheme + "://" + host

	c.HTML(200, "index.html", gin.H{
		"Domain": currentDomain,
	})
}
