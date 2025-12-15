package main

import (
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/core"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/soxft/busuanzi/process/webutil"
)

func main() {
	config.Init()
	redisutil.Init()

	core.InitExpire()
	core.InitSchedule() // 初始化定时任务

	webutil.Init()
}
