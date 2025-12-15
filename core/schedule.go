package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/soxft/busuanzi/config"
	"github.com/soxft/busuanzi/process/redisutil"
	"log"
	"time"
)

// InitSchedule 初始化定时任务
func InitSchedule() {
	// 启动每日归档任务
	go dailyArchiveTask()

	// 启动月度归档任务
	go monthlyArchiveTask()
}

// dailyArchiveTask 每日归档任务，在每天零点将今日数据归档到昨日
func dailyArchiveTask() {
	for {
		// 计算距离下一个零点的时间
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := nextMidnight.Sub(now)

		if config.DEBUG {
			log.Printf("[schedule] Daily archive task will run after %v (at %s)", duration, nextMidnight.Format("2006-01-02 15:04:05"))
		}

		// 等待到零点
		time.Sleep(duration)

		// 执行归档操作
		archiveDailyData()

		// 额外等待1秒，确保不会在同一秒内重复执行
		time.Sleep(time.Second)
	}
}

// monthlyArchiveTask 月度归档任务，在每月1号零点重置月度统计
func monthlyArchiveTask() {
	for {
		// 计算距离下一个月1号零点的时间
		now := time.Now()
		var nextFirstDay time.Time
		if now.Day() == 1 && now.Hour() == 0 && now.Minute() == 0 {
			// 如果当前就是1号零点，则计算到下个月1号
			nextFirstDay = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		} else {
			// 否则计算到下个月1号零点
			nextFirstDay = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		}
		duration := nextFirstDay.Sub(now)

		if config.DEBUG {
			log.Printf("[schedule] Monthly archive task will run after %v (at %s)", duration, nextFirstDay.Format("2006-01-02 15:04:05"))
		}

		// 等待到下个月1号零点
		time.Sleep(duration)

		// 执行月度归档操作
		archiveMonthlyData()

		// 额外等待1秒
		time.Sleep(time.Second)
	}
}

// archiveDailyData 归档每日数据
func archiveDailyData() {
	ctx := context.Background()
	_redis := redisutil.RDB

	// 获取所有今日数据的 key
	todayKeyPattern := "bsz:daily_pv:*:" + time.Now().Format("2006-01-02")
	todayUvKeyPattern := "bsz:daily_uv:*:" + time.Now().Format("2006-01-02")

	// 由于我们需要归档所有站点的数据，但 Redis 的 SCAN 可能消耗较大
	// 实际上，数据已经在 Redis 中按日期分割了，不需要额外归档
	// 昨日数据会自动保留在对应的日期 key 中
	// 这里主要是为了清理和日志记录

	log.Printf("[schedule] Daily archive completed. Today keys pattern: %s, %s", todayKeyPattern, todayUvKeyPattern)

	// 可选：清理过期的历史数据（超过保留期限的）
	// 这里可以根据配置清理超过 DailyExpire 的旧数据
	cleanupOldDailyData(ctx, _redis)
}

// archiveMonthlyData 归档月度数据
func archiveMonthlyData() {
	ctx := context.Background()
	_redis := redisutil.RDB

	monthKeyPattern := "bsz:monthly_pv:*:" + time.Now().Format("2006-01")
	monthUvKeyPattern := "bsz:monthly_uv:*:" + time.Now().Format("2006-01")

	log.Printf("[schedule] Monthly archive completed. Month keys pattern: %s, %s", monthKeyPattern, monthUvKeyPattern)

	// 可选：清理过期的历史月度数据
	cleanupOldMonthlyData(ctx, _redis)
}

// cleanupOldDailyData 清理过期的每日数据
func cleanupOldDailyData(_ context.Context, _ *redis.Client) {
	// 这个功能依赖于 Redis 的 TTL 自动过期机制
	// 由于我们已经在 count.go 中设置了 DailyExpire，Redis 会自动清理过期数据
	// 这里只是一个占位符，可以在未来实现更复杂的清理逻辑
	if config.DEBUG {
		log.Printf("[schedule] Old daily data cleanup is handled by Redis TTL")
	}
}

// cleanupOldMonthlyData 清理过期的月度数据
func cleanupOldMonthlyData(_ context.Context, _ *redis.Client) {
	// 同样依赖于 Redis 的 TTL 自动过期机制
	if config.DEBUG {
		log.Printf("[schedule] Old monthly data cleanup is handled by Redis TTL")
	}
}
